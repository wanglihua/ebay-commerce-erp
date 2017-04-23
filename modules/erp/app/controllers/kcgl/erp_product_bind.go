package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/product_bind_log_service"
	"matrix/modules/erp/service/product_bind_service"
	"matrix/modules/erp/service/status_service"
)

type ErpProductBind struct {
	*revel.Controller
}

func (c ErpProductBind) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_bind/product_bind_index.html")
}

type ErpProductBindViewItem struct {
	models.ProductBindInfo     `xorm:"extends" json:"b"`
	models.StatusInfo          `xorm:"extends" json:"s"`
}

func (c ErpProductBind) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT b.*, m.*, s.*
		FROM erp_product_bind b
		    INNER JOIN erp_package_material m ON b.pkg_mtrl_id = m.id
		    INNER JOIN erp_status s ON b.status_id = s.id
	*/

	query := db_session.
	Select("b.*, s.*").
		Table("erp_product_bind").Alias("b").
		Join("INNER", []string{"erp_status", "s"}, "b.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_bind_list := make([]ErpProductBindViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_bind_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpProductBindViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_bind_list,
		Total: count,
	})
}

type ProductBindDetailForm struct {
	ProductBind models.ProductBindInfo       `json:"product_bind"`
	StatusList  []models.StatusInfo          `json:"status_list"`
}

func (f *ProductBindDetailForm) IsCreate() bool {
	return f.ProductBind.Id == 0
}

func (f *ProductBindDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductBind.Code).Message("商品编号不能为空！")
	if f.ProductBind.Code != "" {
		validation.MinSize(f.ProductBind.Code, 2).Message("商品编号长度不能小于2！")
	}
	if f.ProductBind.Code != "" {
		validation.MaxSize(f.ProductBind.Code, 100).Message("商品编号长度不能大于100！")
	}

	validation.Required(f.ProductBind.Name).Message("名称不能为空！")
	if f.ProductBind.Name != "" {
		validation.MinSize(f.ProductBind.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ProductBind.Name != "" {
		validation.MaxSize(f.ProductBind.Name, 200).Message("名称长度不能大于200！")
	}

	validation.Required(f.ProductBind.BindList).Message("捆绑商品列表不能为空！")
	if f.ProductBind.BindList != "" {
		validation.MinSize(f.ProductBind.BindList, 1).Message("捆绑商品列表长度不能小于1！")
	}
	if f.ProductBind.BindList != "" {
		validation.MaxSize(f.ProductBind.BindList, 500).Message("捆绑商品列表长度不能大于500！")
	}

	validation.Required(f.ProductBind.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductBind) DetailData() revel.Result {
	db_session := c.DbSession

	var product_bind_id int64
	c.Params.Bind(&product_bind_id, "id")

	var product_bind models.ProductBindInfo
	if product_bind_id != 0 {
		has, err := db_session.Id(product_bind_id).Get(&product_bind)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的绑定商品不存在！"})
		}
	} else {
		product_bind.StatusId = models.PRODUCT_BIND_STATUS_NORMAL
	}

	status_list := status_service.Get_status_list(db_session, models.STATUS_TYPE_PRODUCT_BIND)

	var detail_form ProductBindDetailForm
	detail_form.ProductBind = product_bind
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpProductBind) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductBindDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_bind_from_ui := detail_form.ProductBind

	valid, err := product_bind_service.Check_product_bind_valid(product_bind_from_ui.BindList, db_session)
	if valid == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	var affected int64
	if detail_form.IsCreate() {
		code_count, err := db_session.Where("code = ?", product_bind_from_ui.Code).Count(&product_bind_from_ui)
		core.HandleError(err)
		if code_count != 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "保存失败， 编码已存在！",
			})
		}

		name_count, err := db_session.Where("name = ?", product_bind_from_ui.Name).Count(&product_bind_from_ui)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "保存失败， 名称已存在！",
			})
		}

		affected, err = db_session.Insert(&product_bind_from_ui)
		core.HandleError(err)

		product_bind_log_service.Insert_product_bind_log(nil, &product_bind_from_ui, core.Db_operation_types.Insert, db_session, c.Session)
	} else {
		code_count, err := db_session.Where("id <> ? AND code = ?", product_bind_from_ui.Id, product_bind_from_ui.Code).Count(&product_bind_from_ui)
		core.HandleError(err)
		if code_count != 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "保存失败， 编码已存在！",
			})
		}

		name_count, err := db_session.Where("id <> ? AND name = ?", product_bind_from_ui.Id, product_bind_from_ui.Name).Count(&product_bind_from_ui)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "保存失败， 名称已存在！",
			})
		}

		var product_bind_in_db models.ProductBindInfo
		_, err = db_session.Id(product_bind_from_ui.Id).Get(&product_bind_in_db)
		core.HandleError(err)

		update_cols := []string{"code", "name", "bind_list", "desc", "status_id"}
		affected, err = db_session.Id(product_bind_from_ui.Id).Cols(update_cols...).Update(&product_bind_from_ui)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}

		product_bind_log_service.Insert_product_bind_log(&product_bind_from_ui, &product_bind_in_db, core.Db_operation_types.Update, db_session, c.Session)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductBind) Delete() revel.Result {
	db_session := c.DbSession

	product_bind_id_list := make([]int64, 0)
	c.Params.Bind(&product_bind_id_list, "id_list")

	for _, product_bind_id := range product_bind_id_list {
		var product_bind models.ProductBindInfo
		has, err := db_session.Id(product_bind_id).Get(&product_bind)
		core.HandleError(err)

		if has {
			product_bind_log_service.Insert_product_bind_log(nil, &product_bind, core.Db_operation_types.Delete, db_session, c.Session)
		}
	}

	affected, err := db_session.In("id", product_bind_id_list).Delete(new(models.ProductBindInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpProductBind) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_BIND, c.Session)

	db_session := c.DbSession

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type product_bind_import_t struct {
		OrdNum   string // 序号
		BindCode string // 编号
		BindName string // 名称
		BindList string // 绑定列表
	}

	product_bind_import_list := make([]product_bind_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_bind_import product_bind_import_t

			col_index := 0

			if len(row.Cells) < 4 { // 四列数据
				break
			}

			// 序号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			product_bind_import.OrdNum = cell_val
			col_index++

			// 编号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_bind_import.BindCode = cell_val
			col_index++

			// 名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_bind_import.BindName = cell_val
			col_index++

			// 绑定列表
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_bind_import.BindList = cell_val
			col_index++

			product_bind_import_list = append(product_bind_import_list, product_bind_import)
		}
	}

	if len(product_bind_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入物流方式列表为空！",
		})
	}

	for _, product_bind_import := range product_bind_import_list {
		var product_bind models.ProductBindInfo

		// 检查 绑定代码是否已存在
		if product_bind_service.Is_product_bind_sku(product_bind_import.BindCode, db_session) == true {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("序号 %s 的绑定编号 %s 已存在，请核查！", product_bind_import.OrdNum, product_bind_import.BindCode),
			})
		}

		// 检查 绑定列表是否合法有效
		bind_list_valid, err := product_bind_service.Check_product_bind_valid(product_bind.BindList, db_session)
		core.HandleError(err)
		if bind_list_valid == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("序号 %s 的绑定列表无效，请核查！", product_bind.BindList),
			})
		}

		product_bind.Code = product_bind_import.BindCode
		product_bind.Name = product_bind_import.BindName
		product_bind.BindList = product_bind_import.BindList
		product_bind.StatusId = models.PRODUCT_BIND_STATUS_NORMAL

		_, err = db_session.Insert(&product_bind)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}
