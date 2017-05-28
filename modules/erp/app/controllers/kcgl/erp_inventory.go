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
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/product_service"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/service/storage_loc_service"
)

type ErpInventory struct {
	*revel.Controller
}

func (c ErpInventory) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	return c.RenderTemplate("erp/kcgl/inventory/inventory_index.html")
}

func (c ErpInventory) PageData() revel.Result {
	var db_session = c.DbSession
	var stock_list = make([]models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"stock_list":         stock_list,
		"stock_id_query":     0,
		"product_code_query": "",
	})
}

type ErpInventoryViewItem struct {
	Inventory  models.InventoryInfo  `xorm:"extends" json:"ivt"`
	Stock      models.StockInfo      `xorm:"extends" json:"st"`
	StorageLoc models.StorageLocInfo `xorm:"extends" json:"loc"`
	Product    models.ProductInfo    `xorm:"extends" json:"prd"`
}

func (c ErpInventory) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	/*
		SELECT ivt.*, st.*, loc.*, prd.*
		FROM erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	query := db_session.
	Select("ivt.*, st.*, loc.*, prd.*").
		Table("erp_inventory").Alias("ivt").
		Join("INNER", []string{"erp_stock", "st"}, "ivt.stock_id = st.id").
		Join("LEFT", []string{"erp_storage_loc", "loc"}, "ivt.storage_loc_id = loc.id").
		Join("INNER", []string{"erp_product", "prd"}, "ivt.product_id = prd.id")

	//query extra filter here

	if stock_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("st.id = %d", stock_id_query))
	}

	if product_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_code_query))
	}

	query = query.Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("ivt.id")
	}

	inventory_list := make([]ErpInventoryViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&inventory_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpInventoryViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_list,
		Total: count,
	})
}

type InventoryDetailForm struct {
	Inventory      models.InventoryInfo     `json:"inventory"`
	StockList      []models.StockInfo       `json:"stock_list"`
	StorageLocList []models.StorageLocInfo  `json:"storage_loc_list"`
	ProductSku     string                   `json:"product_sku"`
}

func (f *InventoryDetailForm) IsCreate() bool {
	return f.Inventory.Id == 0
}

func (f *InventoryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Inventory.StockId).Message("仓库不能为空！")

	validation.Required(f.Inventory.StorageLocId).Message("库位不能为空！")

	validation.Required(f.Inventory.ProductId).Message("商品不能为空！")

	validation.Required(f.Inventory.Qty).Message("数量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpInventory) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	var inventory_id int64
	c.Params.Bind(&inventory_id, "id")

	var stock_list = make([]models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	var detail_form InventoryDetailForm
	detail_form.StockList = stock_list

	var inventory models.InventoryInfo
	if inventory_id != 0 {
		has, err := db_session.Id(inventory_id).Get(&inventory)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品库存不存在！"})
		}

		var storage_loc_list = make([]models.StorageLocInfo, 0)
		err = db_session.Where(fmt.Sprintf("stock_id = %d", inventory.StockId)).Find(&storage_loc_list)
		core.HandleError(err)

		detail_form.StorageLocList = storage_loc_list

		var product models.ProductInfo
		_, err = db_session.Id(inventory.ProductId).Get(&product)
		core.HandleError(err)

		detail_form.ProductSku = product.Sku
	}

	detail_form.Inventory = inventory

	return c.RenderJson(detail_form)
}

func (c ErpInventory) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	var detail_form InventoryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	inventory := detail_form.Inventory
	var product_sku = detail_form.ProductSku

	var product_id = product_service.Get_product_id_by_sku(product_sku, db_session)
	if product_id == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "保存失败，商品SKU不存在！",
		})
	}

	inventory.ProductId = product_id

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&inventory)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(inventory.Id).Update(&inventory)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpInventory) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	inventory_id_list := make([]int64, 0)
	c.Params.Bind(&inventory_id_list, "id_list")

	affected, err := db_session.In("id", inventory_id_list).Delete(new(models.InventoryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpInventory) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

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

	//序号 商品编号	仓库	仓位	库存

	type inventory_import_t struct {
		OrdNum         string  `json:"ord_num"`          // 序号
		ProductSku     string  `json:"product_sku"`      // 商品编号
		StockName      string  `json:"stock_name"`       // 仓库名称
		StorageLocName string  `json:"storage_loc_name"` // 仓位名称
		Qty            int64   `json:"qty"`              // 库存
		CostPrice      float64 `json:"cost_price"`       // 成本价

		// 下面这些字段是从数据库里面查询后保存在这里面的，这样过程中，就不用二次查询了
		stock_id       int64
		storage_loc_id int64
		product_id     int64
		ErrorMessage   string `json:"error_message"`
	}

	inventory_import_list := make([]inventory_import_t, 0)

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

			var inventory_import inventory_import_t

			col_index := 0

			var data_col_count = 6

			if len(row.Cells) < data_col_count-1 { // 6 列数据， 最后一列可能为空
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
			inventory_import.OrdNum = cell_val
			col_index++

			// 商品编号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.ProductSku = cell_val
			col_index++

			// 仓库名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StockName = cell_val
			col_index++

			// 仓位名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StorageLocName = cell_val
			col_index++

			// 库存
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			qty, _ := strconv.ParseInt(cell_val, 10, 64)
			inventory_import.Qty = qty
			col_index++

			if len(row.Cells) > data_col_count-1 { // 最后一列有数据
				// 成本价
				cell = row.Cells[col_index]
				cell_val, err = cell.String()
				core.HandleError(err)
				cell_val = strings.TrimSpace(cell_val)
				cost_price, _ := strconv.ParseFloat(cell_val, 64)
				inventory_import.CostPrice = cost_price
				col_index++
			}

			inventory_import_list = append(inventory_import_list, inventory_import)
		}
	}

	if len(inventory_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品库存列表为空！",
		})
	}

	var inventory_import_error_list = make([]inventory_import_t, 0)
	// 先检查一遍是不是有错误，
	for index, inventory_import := range inventory_import_list {

		// 检查 商品 sku
		var product_id int64 = product_service.Get_product_id_by_sku(inventory_import.ProductSku, db_session)

		if product_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的商品SKU %s 不存在！", inventory_import.OrdNum, inventory_import.ProductSku)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		}

		// 检查 仓库名称
		var stock_id int64 = stock_service.Get_stock_id_by_name(inventory_import.StockName, db_session)
		var storage_loc_id int64

		if stock_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StockName)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		} else { //仓库名称存在
			var storage_loc_name = strings.TrimSpace(inventory_import.StorageLocName)
			if storage_loc_name != "" {
				// 如果提供了仓位， 检查仓位名称
				storage_loc_id = storage_loc_service.Get_storage_loc_id_by_name(storage_loc_name, stock_id, db_session)
				if storage_loc_id == 0 {
					inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓位名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StorageLocName)
					inventory_import_error_list = append(inventory_import_error_list, inventory_import)
				}
			}
		}

		inventory_import_list[index].product_id = product_id
		inventory_import_list[index].stock_id = stock_id
		inventory_import_list[index].storage_loc_id = storage_loc_id
	}

	if len(inventory_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": inventory_import_error_list,
		})
	}

	for _, inventory_import := range inventory_import_list {

		var stock_id = inventory_import.stock_id
		var storage_loc_id = inventory_import.storage_loc_id
		var product_id = inventory_import.product_id

		var inventory models.InventoryInfo
		has, err := db_session.Where(fmt.Sprintf("stock_id = %d AND storage_loc_id = %d AND product_id = %d", stock_id, storage_loc_id, product_id)).Get(&inventory)
		core.HandleError(err)

		inventory.StockId = stock_id
		inventory.StorageLocId = storage_loc_id
		inventory.ProductId = product_id
		// inventory.Qty += inventory_import.Qty // 如果是插入的话， 开始时 inventory.Qty
		inventory.Qty = inventory_import.Qty // 这里采用覆盖原来库存的方式
		inventory.CostPrice = inventory_import.CostPrice

		if has == false {
			// 插入
			_, err = db_session.Insert(&inventory)
			core.HandleError(err)
		} else {
			// 更新
			_, err := db_session.Id(inventory.Id).Cols("stock_id", "storage_loc_id", "product_id", "qty", "cost_price").Update(&inventory)
			core.HandleError(err)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpInventory) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	xlsx_file := c.get_export_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "库存列表" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpInventory) get_export_excel(db_session *xorm.Session) *xlsx.File {
	filter, order, _, _ := core.GetGridRequestParam(c.Request)

	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	/*
		SELECT ivt.*, st.*, loc.*, prd.*
		FROM erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	query := db_session.
	Select("ivt.*, st.*, loc.*, prd.*").
		Table("erp_inventory").Alias("ivt").
		Join("INNER", []string{"erp_stock", "st"}, "ivt.stock_id = st.id").
		Join("LEFT", []string{"erp_storage_loc", "loc"}, "ivt.storage_loc_id = loc.id").
		Join("INNER", []string{"erp_product", "prd"}, "ivt.product_id = prd.id")

	//query extra filter here

	if stock_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("st.id = %d", stock_id_query))
	}

	if product_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_code_query))
	}

	query = query.Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("ivt.id")
	}

	inventory_list := make([]ErpInventoryViewItem, 0)
	err := query.Find(&inventory_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("库存列表")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = horizontal_align
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	row_height := 0.6

	// 序号	商品编号	仓库	仓位	库存	成本价（RMB）
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
		sheet.SetColWidth(col_index, col_index, 25)
		set_head_cell_style(cell)
		cell.Value = "商品编号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 22)
		set_head_cell_style(cell)
		cell.Value = "仓库"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "仓位"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "库存"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "成本价"
	}

	for index, inventory_view := range inventory_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Product.Sku

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Stock.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.StorageLoc.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(inventory_view.Inventory.Qty)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(inventory_view.Inventory.CostPrice)
	}

	return xlsx_file
}
