package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/service/storage_loc_service"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/service/excel_service"
)

type ErpStorageLoc struct {
	*revel.Controller
}

func (c ErpStorageLoc) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	return c.RenderTemplate("erp/sppz/storage_loc/storage_loc_index.html")
}

type ErpStorageLocViewItem struct {
	sppz_models.StorageLocInfo `xorm:"extends" json:"l"`
	sppz_models.StockInfo      `xorm:"extends" json:"s"`
}

func (c ErpStorageLoc) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT l.*, s.*
		FROM erp_storage_loc l
		    INNER JOIN erp_stock s ON l.stock_id = s.id
	*/

	query := db_session.
		Select("l.*, s.*").
		Table("erp_storage_loc").Alias("l").
		Join("INNER", []string{"erp_stock", "s"}, "l.stock_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("l.id")
	}

	storage_loc_list := make([]ErpStorageLocViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&storage_loc_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpStorageLocViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  storage_loc_list,
		Total: count,
	})
}

type StorageLocDetailForm struct {
	StorageLoc sppz_models.StorageLocInfo `json:"storage_loc"`
	StockList  []sppz_models.StockInfo    `json:"stock_list"`
}

func (f *StorageLocDetailForm) IsCreate() bool {
	return f.StorageLoc.Id == 0
}

func (f *StorageLocDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.StorageLoc.StockId).Message("仓库不能为空！")

	validation.Required(f.StorageLoc.Name).Message("名称不能为空！")
	if f.StorageLoc.Name != "" {
		validation.MinSize(f.StorageLoc.Name, 2).Message("名称长度不能小于2！")
	}
	if f.StorageLoc.Name != "" {
		validation.MaxSize(f.StorageLoc.Name, 100).Message("名称长度不能大于100！")
	}

	if f.StorageLoc.Desc.Valid == true {
		validation.MaxSize(f.StorageLoc.Desc.String, 500).Message("描述长度不能大于500！")
	}

	// validation.Required(f.StorageLoc.XAxes).Message("X轴坐标不能为空！")
	if f.StorageLoc.XAxes != "" {
		validation.MaxSize(f.StorageLoc.XAxes, 100).Message("X轴坐标长度不能大于100！")
	}

	// validation.Required(f.StorageLoc.YAxes).Message("Y轴坐标不能为空！")
	if f.StorageLoc.YAxes != "" {
		validation.MaxSize(f.StorageLoc.YAxes, 100).Message("Y轴坐标长度不能大于100！")
	}

	return validation.HasErrors() == false
}

func (c ErpStorageLoc) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	db_session := c.DbSession

	var storage_loc_id int64
	c.Params.Bind(&storage_loc_id, "id")

	var storage_loc sppz_models.StorageLocInfo
	if storage_loc_id != 0 {
		has, err := db_session.Id(storage_loc_id).Get(&storage_loc)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的仓位不存在！"})
		}
	}

	var stock_list = make([]sppz_models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	var detail_form StorageLocDetailForm
	detail_form.StorageLoc = storage_loc
	detail_form.StockList = stock_list

	return c.RenderJson(detail_form)
}

func (c ErpStorageLoc) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	db_session := c.DbSession

	var detail_form StorageLocDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	storage_loc := detail_form.StorageLoc

	var affected int64
	if detail_form.IsCreate() {
		stockId_count, err := db_session.Where("stock_id = ? and name = ?", storage_loc.StockId, storage_loc.Name).Count(new(sppz_models.StorageLocInfo))
		core.HandleError(err)
		if stockId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，仓库、名称已存在！"})
		}

		affected, err = db_session.Insert(&storage_loc)
		core.HandleError(err)
	} else {
		stockId_count, err := db_session.Where("id <> ? and stock_id = ? and name = ?", storage_loc.Id, storage_loc.StockId, storage_loc.Name).Count(new(sppz_models.StorageLocInfo))
		core.HandleError(err)
		if stockId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，仓库、名称已存在！"})
		}

		affected, err = db_session.Id(storage_loc.Id).Update(&storage_loc)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpStorageLoc) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	db_session := c.DbSession

	storage_loc_id_list := make([]int64, 0)
	c.Params.Bind(&storage_loc_id_list, "id_list")

	affected, err := db_session.In("id", storage_loc_id_list).Delete(new(sppz_models.StorageLocInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpStorageLoc) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	var db_session = c.DbSession

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type storage_loc_import_t struct {
		OrdNum         string `json:"ord_num"`          // 序号
		StockName      string `json:"stock_name"`       // 仓库名称
		StorageLocName string `json:"storage_loc_name"` // 仓位名称
		StorageLocDesc string `json:"storage_loc_desc"` // 仓位描述

		// 下面这些字段是从数据库里面查询后保存在这里面的，这样过程中，就不用二次查询了
		stock_id     int64
		ErrorMessage string `json:"error_message"`
	}

	storage_loc_import_list := make([]storage_loc_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow

		revel.TRACE.Println("max_row", max_row)
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var storage_loc_import storage_loc_import_t

			col_index := 0

			// 序号	仓库名称	仓位名称	仓位描述

			// 序号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			storage_loc_import.OrdNum = cell_val
			col_index++

			// 仓库名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			storage_loc_import.StockName = cell_val
			col_index++

			// 仓位名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			storage_loc_import.StorageLocName = cell_val
			col_index++

			if len(row.Cells) > col_index { // 最后一列有数据
				// 成本价
				// cell = row.Cells[col_index]
				cell = sheet.Cell(row_index, col_index)
				cell_val, err = cell.String()
				core.HandleError(err)
				cell_val = strings.TrimSpace(cell_val)
				storage_loc_import.StorageLocDesc = cell_val
				col_index++
			}

			storage_loc_import_list = append(storage_loc_import_list, storage_loc_import)
		}
	}

	if len(storage_loc_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入仓库库位列表为空！",
		})
	}

	var storage_loc_import_error_list = make([]storage_loc_import_t, 0)
	// 先检查一遍是不是有错误，
	for index, storage_loc_import := range storage_loc_import_list {

		// 检查 仓库名称
		stock_id, _ := stock_service.Get_stock_id_type_by_name(storage_loc_import.StockName, db_session)
		var storage_loc_id int64

		if stock_id == 0 {
			storage_loc_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库名称 %s 不存在！", storage_loc_import.OrdNum, storage_loc_import.StockName)
			storage_loc_import_error_list = append(storage_loc_import_error_list, storage_loc_import)
		} else { //仓库名称存在
			var storage_loc_name = strings.TrimSpace(storage_loc_import.StorageLocName)
			if storage_loc_name != "" {
				// 检查仓位是否已经存在
				storage_loc_id = storage_loc_service.Get_storage_loc_id_by_name(storage_loc_name, stock_id, db_session)
				if storage_loc_id != 0 {
					storage_loc_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓位名称 %s 已经存在！", storage_loc_import.OrdNum, storage_loc_import.StorageLocName)
					storage_loc_import_error_list = append(storage_loc_import_error_list, storage_loc_import)
				}
			}
		}

		revel.TRACE.Println("stock_id", stock_id)
		storage_loc_import_list[index].stock_id = stock_id
	}

	if len(storage_loc_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": storage_loc_import_error_list,
		})
	}

	for _, storage_loc_import := range storage_loc_import_list {
		revel.TRACE.Println("storage_loc_import.stock_id", storage_loc_import.stock_id)
		var storage_loc sppz_models.StorageLocInfo
		storage_loc.StockId = storage_loc_import.stock_id
		storage_loc.Name = storage_loc_import.StorageLocName
		storage_loc.Desc = core.MakeNullString(storage_loc_import.StorageLocDesc, true)

		// 插入
		_, err = db_session.Insert(&storage_loc)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpStorageLoc) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	var db_session = c.DbSession

	xlsx_file := c.get_export_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "仓位" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpStorageLoc) get_export_excel(db_session *xorm.Session) *xlsx.File {

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = ""

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("sl.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	type storage_loc_export_item_t struct {
		StorageLoc sppz_models.StorageLocInfo `xorm:"extends" json:"sl"`
		Stock      sppz_models.StockInfo      `xorm:"extends" json:"s"`
	}

	/*
SELECT sl.*, s.*
FROM erp_storage_loc sl
    INNER JOIN erp_stock s ON s.id = sl.stock_id
	 */

	var sql_get_storage_loc_list = `SELECT sl.*, s.*
FROM erp_storage_loc sl
    INNER JOIN erp_stock s ON s.id = sl.stock_id`

	if filter != "" {
		sql_get_storage_loc_list = sql_get_storage_loc_list + " WHERE " + filter
	}

	var storage_loc_list = make([]storage_loc_export_item_t, 0)
	err := db_session.Sql(sql_get_storage_loc_list).Find(&storage_loc_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("仓位")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		excel_service.Set_head_cell_style(cell)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		excel_service.Set_text_cell_style(cell, horizontal_align)
	}

	row_height := 0.6

	// 序号	仓库名称	仓位名称	仓位描述
	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "序号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "仓库名称"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "仓位名称"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "仓位描述"
	}

	for index, storage_loc := range storage_loc_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 序号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		// 仓库名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = storage_loc.Stock.Name

		// 仓位名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = storage_loc.StorageLoc.Name

		// 仓位描述
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		if storage_loc.StorageLoc.Desc.Valid == true {
			cell.Value = storage_loc.StorageLoc.Desc.String
		} else {
			cell.Value = ""
		}
	}

	return xlsx_file
}

func (c ErpStorageLoc) GetStorageLocList() revel.Result {
	// 这个不用权限限制，别的地方都会调用

	var db_session = c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var storage_loc_list = make([]sppz_models.StorageLocInfo, 0)
	err := db_session.Where(fmt.Sprintf("stock_id = %d", stock_id)).Find(&storage_loc_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success":          true,
		"message":          "获取库位列表成功！",
		"storage_loc_list": storage_loc_list,
	})
}
