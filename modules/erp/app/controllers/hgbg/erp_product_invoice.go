package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"strconv"
)

type ErpProductInvoice struct {
	*revel.Controller
}

func (c ErpProductInvoice) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	return c.RenderTemplate("erp/hgbg/product_invoice/product_invoice_index.html")
}

func (c ErpProductInvoice) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	db_session := c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_invoice_view_item_t struct {
		Product kcgl_models.ProductInfo `xorm:"extends" json:"p"`
		Status  models.StatusInfo       `xorm:"extends" json:"s"`
	}

	/*
		SELECT p.*, s.*
		FROM erp_product p
		    INNER JOIN erp_status s ON s.id = p.status_id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("p.*, s.*")
	sb.Count("count(p.id) AS product_count")
	sb.From(`erp_product p
		    INNER JOIN erp_status s ON s.id = p.status_id`)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("p.id ASC")
	}

	sb.Limit(offset, limit)

	var product_list = make([]product_invoice_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_list)
	core.HandleError(err)

	type product_count_t struct {
		Count int64 `xorm:"'product_count'"`
	}
	var product_count product_count_t
	_, err = db_session.Sql(sb.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: product_count.Count,
	})
}

type product_weight_detail_form_t struct {
	Product kcgl_models.ProductInfo `json:"product"`
}

func (c ErpProductInvoice) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var detail_form product_weight_detail_form_t
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProductInvoice) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	db_session := c.DbSession

	var detail_form product_weight_detail_form_t
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	var product_id = detail_form.Product.Id

	var product_in_db kcgl_models.ProductInfo
	has, err := db_session.Id(product_id).Get(&product_in_db)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品不存在！",
		})
	}

	product_in_db.NeedInvoice = detail_form.Product.NeedInvoice
	product_in_db.TaxPoint = detail_form.Product.TaxPoint
	product_in_db.DrawbackPoint = detail_form.Product.DrawbackPoint

	updata_col_List := []string{"need_invoice", "tax_point", "drawback_point"}
	_, err = db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}

func (c ErpProductInvoice) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

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

	type product_invoice_t struct {
		sku            string
		need_invoice   int64
		tax_point      float64
		drawback_point float64
	}

	product_invoice_list := make([]product_invoice_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow

		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_invoice product_invoice_t

			col_index := 0 //从第二列开始导入，略去第一列的编号

			//SKU
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}

			//revel.TRACE.Println(cell_val)
			product_invoice.sku = cell_val
			col_index++

			// 是否开票
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)

			if strings.TrimSpace(cell_val) == "是" {
				product_invoice.need_invoice = 1
			} else {
				product_invoice.need_invoice = 0
			}

			col_index++

			// 收取税点
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)

			if cell_val != "" {
				tax_point, err := strconv.ParseFloat(cell_val, 64)
				core.HandleError(err)
				product_invoice.tax_point = tax_point
			} else {
				product_invoice.tax_point = 0
			}

			col_index++

			// 退税点
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)

			if cell_val != "" {
				drawback_point, err := strconv.ParseFloat(cell_val, 64)
				core.HandleError(err)
				product_invoice.drawback_point = drawback_point
			} else {
				product_invoice.tax_point = 0
			}

			col_index++

			product_invoice_list = append(product_invoice_list, product_invoice)
		}
	}

	if len(product_invoice_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	for _, product_invoice := range product_invoice_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_invoice.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_invoice.sku
		product.NeedInvoice = product_invoice.need_invoice
		product.TaxPoint = product_invoice.tax_point
		product.DrawbackPoint = product_invoice.drawback_point

		if has == true {
			_, err = db_session.Id(product.Id).Cols("need_invoice", "tax_point", "drawback_point").Update(&product)
			core.HandleError(err)
		} else {
			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_invoice.sku))
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpProductInvoice) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("id IN (%s)", strings.Join(id_str_list, ",")))
	}

	product_list := make([]kcgl_models.ProductInfo, 0)
	err := db_session.Where(filter).Find(&product_list)
	core.HandleError(err)

	xlsx_file := c.export_product_list_to_excel(db_session, product_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品是否开票导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpProductInvoice) export_product_list_to_excel(db_session *xorm.Session, product_list []kcgl_models.ProductInfo) *xlsx.File {

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := *xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00CCFFCC", "00000000") //CCFFCC
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

	//统一设置内容单元格格式的函数
	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := *xlsx.NewFont(9, "微软雅黑")
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

	head_row_height := 0.8
	common_row_height := 0.6

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	//添加工作表
	sheet, err := xlsx_file.AddSheet("商品长宽高重量列表")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	var product_status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT)).Find(&product_status_list)
	core.HandleError(err)

	var find_product_status_name = func(status_id int64) string {
		for _, status := range product_status_list {
			if status.Id == status_id {
				return status.Name
			}
		}

		return ""
	}

	//序号	SKU	中文名称	产品长度（cm）	产品宽度（cm）	产品高度（cm）	实际重量(kg)	状态

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "序号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 30)
		set_head_cell_style(cell)
		cell.Value = "中文名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "是否开票"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "收取税点(1表示1%)"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "退税点(1表示1%)"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++
	}

	for index, product := range product_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 序号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product.Sku

		// 中文名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product.CnName

		// 	是否开票
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if product.NeedInvoice == 1 {
			cell.Value = "是"
		} else {
			cell.Value = "否"
		}

		// 收取税点
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", product.TaxPoint)

		// 退税点
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", product.DrawbackPoint)

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = find_product_status_name(product.StatusId)
	}

	return xlsx_file
}
