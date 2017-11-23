package controllers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"strings"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
	auth_models "matrix/modules/auth/models"
)

type ErpProductExport struct {
	*revel.Controller
}

func (c ErpProductExport) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_EXPORT, c.Session)

	return c.RenderTemplate("erp/kcgl/product_export/product_export.html")
}

func (c ErpProductExport) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_EXPORT, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_list := make([]kcgl_models.ProductInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(kcgl_models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: count,
	})
}

func (c ErpProductExport) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_EXPORT, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("p.id IN (%s)", strings.Join(id_str_list, ",")))
	} else {
		var batch_code int64
		c.Params.Bind(&batch_code, "batch_code")

		var product_sku_query string
		c.Params.Bind(&product_sku_query, "product_sku_query")

		if strings.TrimSpace(product_sku_query) != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("p.sku LIKE '%%%s%%'", product_sku_query))
		}

		if batch_code != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("p.batch_code = %d", batch_code))
		}
	}

	/*
SELECT p.*, sup.*, su.*, du.*
FROM erp_product p
    LEFT JOIN erp_supplier sup ON p.supplier_id = sup.id
    LEFT JOIN auth_user su ON p.sale_user_id = su.id
    LEFT JOIN auth_user du ON p.dev_user_id = du.id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("p.*, sup.*, su.*, du.*")
	sb.From(`erp_product p
    LEFT JOIN erp_supplier sup ON p.supplier_id = sup.id
    LEFT JOIN auth_user su ON p.sale_user_id = su.id
    LEFT JOIN auth_user du ON p.dev_user_id = du.id`)
	sb.Where(filter)

	var product_export_list = make([]product_export_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&product_export_list)
	core.HandleError(err)

	xlsx_file := c.export_product_list_to_excel(db_session, product_export_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

type product_export_item_t struct {
	Product  kcgl_models.ProductInfo  `xorm:"extends"`
	Supplier sppz_models.SupplierInfo `xorm:"extends"`
	SaleUser auth_models.UserInfo     `xorm:"extends"`
	DevUser  auth_models.UserInfo     `xorm:"extends"`
}

func (c ErpProductExport) export_product_list_to_excel(db_session *xorm.Session, product_list []product_export_item_t) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品列表")
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

	//序号	SKU	中文名称	采购价（人民币元）	产品长度（cm）	产品宽度（cm）	产品高度（cm）	实际重量(kg)	状态

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
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "销售员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "开发员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 35)
		set_head_cell_style(cell)
		cell.Value = "供应商"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "采购价（人民币元）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "产品长度（cm）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "产品宽度（cm）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "产品高度（cm）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "实际重量(kg)"
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
		set_text_cell_style(cell, "center")
		cell.Value = product.Product.Sku

		// 中文名称
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = product.Product.CnName

		// 销售员
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if product.SaleUser.LoginName != "" {
			cell.Value = fmt.Sprintf("%s(%s)", product.SaleUser.NickName, product.SaleUser.LoginName)
		} else {
			cell.Value = ""
		}

		// 开发员
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if product.DevUser.LoginName != "" {
			cell.Value = fmt.Sprintf("%s(%s)", product.DevUser.NickName, product.DevUser.LoginName)
		} else {
			cell.Value = ""
		}

		// 供应商
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = product.Supplier.Name

		// 采购价（人民币元）
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", product.Product.CostPrice)

		// 	产品长度（cm）
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", product.Product.Length)

		// 产品宽度（cm）
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", product.Product.Width)

		// 	产品高度（cm）
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", product.Product.Height)

		// 	实际重量(kg)
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.3f", product.Product.Weight)

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = find_product_status_name(product.Product.StatusId)
	}

	return xlsx_file
}
