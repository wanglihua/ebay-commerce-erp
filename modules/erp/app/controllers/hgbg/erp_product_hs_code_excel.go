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
)

type ErpProductHsCodeExcel struct {
	*revel.Controller
}

func (c ErpProductHsCodeExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_HS_CODE, c.Session)

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

	type product_hs_code_import_t struct {
		Sku          string  `json:"sku"`
		CnDeclName   string  `json:"cn_decl_name"`
		AuDeclName   string  `json:"au_decl_name"`
		CnHsCode     string  `json:"cn_hs_code"`
		AuHsCode     string  `json:"au_hs_code"`
		ErrorMessage string  `json:"error_message"`

		product models.ProductInfo `json:"-"`
	}

	product_hs_code_import_list := make([]product_hs_code_import_t, 0)

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
			//row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_hs_code product_hs_code_import_t

			col_index := 0

			// SKU
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			product_hs_code.Sku = cell_val
			col_index++

			// cn_decl_name
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.CnDeclName = cell_val
			col_index++

			// au_decl_name
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.AuDeclName = cell_val
			col_index++

			// cn_hs_code
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.CnHsCode = cell_val
			col_index++

			// au_hs_code
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.AuHsCode = cell_val
			col_index++

			product_hs_code_import_list = append(product_hs_code_import_list, product_hs_code)
		}
	}

	if len(product_hs_code_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	//验证导入的数据
	var product_hs_code_import_error_list = make([]product_hs_code_import_t, 0)

	for index, product_hs_code_import := range product_hs_code_import_list {
		var product models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_hs_code_import.Sku)).Get(&product)
		core.HandleError(err)

		if has == false {
			product_hs_code_import.ErrorMessage = fmt.Sprintf("商品SKU %s 数据库中不存在！", product_hs_code_import.Sku)
			product_hs_code_import_error_list = append(product_hs_code_import_error_list, product_hs_code_import)
		} else {
			product_hs_code_import_list[index].product = product
		}
	}

	if len(product_hs_code_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "导入出现错误，详细请查看页面下部列表！",
			"error_data_list": product_hs_code_import_error_list,
		})
	}

	for _, product_hs_code_import := range product_hs_code_import_list {
		var product = product_hs_code_import.product

		product.Sku = product_hs_code_import.Sku
		product.CnHsCode = product_hs_code_import.CnHsCode
		product.AuHsCode = product_hs_code_import.AuHsCode

		update_col_list := []string{"cn_hs_code", "au_hs_code"}
		_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpProductHsCodeExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_HS_CODE, c.Session)

	var filter_type int64
	c.Params.Bind(&filter_type, "filter_type")

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	db_session := c.DbSession

	product_list := make([]models.ProductInfo, 0)

	var filter = "1=1"
	if filter_type == 1 { //有零值数据的记录
		filter_has_zero_data := core.SqlOr("cn_hs_code = ''", "au_hs_code = ''")
		filter = core.SqlAnd(filter, filter_has_zero_data)
	}

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("id IN (%s)", strings.Join(id_str_list, ",")))
	}

	err := db_session.Where(filter).Find(&product_list)
	core.HandleError(err)

	xlsx_file := c.export_product_list_to_excel(db_session, product_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品海关编码导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpProductHsCodeExcel) export_product_list_to_excel(db_session *xorm.Session, product_list []models.ProductInfo) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品海关编码列表")
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

	//序号 SKU 中国海关编码 澳洲海关编码

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "序号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "中国申报品名"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "澳洲申报品名"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "中国海关编码"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "澳洲海关编码"
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
		cell.Value = product.Sku

		// 中国申报品名
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product.CnDeclName

		// 澳洲申报品名
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product.AuDeclName

		// 中国海关编码
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = product.CnHsCode

		// 澳洲海关编码
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = product.AuHsCode

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = find_product_status_name(product.StatusId)
	}

	return xlsx_file
}
