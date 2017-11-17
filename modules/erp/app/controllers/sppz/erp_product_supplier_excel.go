package controllers

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"strings"
	"matrix/modules/erp/models/kcgl_models"
	"bytes"
	"time"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/sppz_models"
)

type ErpProductSupplierExcel struct {
	*revel.Controller
}

func (c ErpProductSupplierExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_SUPPLIER, c.Session)

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

	type product_supplier_t struct {
		sku           string
		supplier_name string
		product       kcgl_models.ProductInfo
		supplier_id   int64
	}

	product_supplier_list := make([]product_supplier_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	//revel.TRACE.Printf("sheet_count: %d\n", sheet_count)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//revel.TRACE.Printf("max_row: %d\n", max_row)
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_supplier product_supplier_t

			col_index := 1

			// SKU
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			product_supplier.sku = cell_val
			col_index++

			// 供应商
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_supplier.supplier_name = cell_val
			col_index++

			product_supplier_list = append(product_supplier_list, product_supplier)
		}
	}

	if len(product_supplier_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品负责人列表为空！",
		})
	}

	for index, product_supplier := range product_supplier_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_supplier.sku)).Get(&product)
		core.HandleError(err)

		if has == false {
			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_supplier.sku))
		}

		product_supplier_list[index].product = product

		// 供应商
		if product_supplier.supplier_name == "" {
			product_supplier_list[index].supplier_id = 0
		} else {
			var supplier sppz_models.SupplierInfo
			has, err = db_session.Where("name = ?", product_supplier.supplier_name).Get(&supplier)
			core.HandleError(err)

			if has == false {
				panic(fmt.Sprintf("指定的供应商名称数据库中不存在！"))
			} else {
				product_supplier_list[index].supplier_id = supplier.Id
			}
		}
	}

	for _, product_supplier := range product_supplier_list {
		var product = product_supplier.product
		product.SupplierId = core.MakeNullInt(product_supplier.supplier_id, true)

		update_col_list := []string{"supplier_id"}
		_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpProductSupplierExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_SUPPLIER, c.Session)

	db_session := c.DbSession

	/*
SELECT p.*, sup.*, s.*
FROM erp_product p
    LEFT JOIN erp_supplier sup ON sup.id = p.supplier_id
    INNER JOIN erp_status s ON s.id = p.status_id;
	 */

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("p.id IN (%s)", strings.Join(id_str_list, ",")))
	} else {
		var supplier_name_query string
		c.Params.Bind(&supplier_name_query, "supplier_name_query")
		supplier_name_query = strings.TrimSpace(supplier_name_query)
		if supplier_name_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("sup.name LIKE '%%%s%%'", supplier_name_query))
		}
	}

	var sb = core.NewSqlBuilder()
	sb.Select("p.*, sup.*, s.*")
	sb.From(`erp_product p
    LEFT JOIN erp_supplier sup ON sup.id = p.supplier_id
    INNER JOIN erp_status s ON s.id = p.status_id`)

	sb.Where(filter)
	sb.OrderBy("p.sku ASC")

	var product_supplier_list = make([]product_supplier_export_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&product_supplier_list)
	core.HandleError(err)

	xlsx_file := c.export_product_list_to_excel(db_session, product_supplier_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品供应商导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

type product_supplier_export_t struct {
	Product  kcgl_models.ProductInfo  `xorm:"extends"`
	Supplier sppz_models.SupplierInfo `xorm:"extends"`
	Status   models.StatusInfo        `xorm:"extends"`
}

func (c ErpProductSupplierExcel) export_product_list_to_excel(db_session *xorm.Session, product_supplier_list []product_supplier_export_t) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品供应商")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	// SKU 中文名称 供应商名称 状态
	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 40)
		set_head_cell_style(cell)
		cell.Value = "中文名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 40)
		set_head_cell_style(cell)
		cell.Value = "供应商"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++
	}

	for _, product_supplier := range product_supplier_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// SKU
		cell := row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_supplier.Product.Sku

		// 中文名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_supplier.Product.CnName

		// 供应商
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product_supplier.Supplier.Name

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = product_supplier.Status.Name
	}

	return xlsx_file
}
