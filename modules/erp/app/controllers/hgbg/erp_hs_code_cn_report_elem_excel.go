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

type ErpHsCodeCnReportElemExcel struct {
	*revel.Controller
}

func (c ErpHsCodeCnReportElemExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM, c.Session)

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

	type product_hs_code_t struct {
		sku        string
		cn_hs_code string
		au_hs_code string
	}

	product_hs_code_list := make([]product_hs_code_t, 0)

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

			var product_hs_code product_hs_code_t

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
			product_hs_code.sku = cell_val
			col_index++

			// cn_hs_code
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.cn_hs_code = cell_val
			col_index++

			// au_hs_code
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_hs_code.au_hs_code = cell_val
			col_index++

			product_hs_code_list = append(product_hs_code_list, product_hs_code)
		}
	}

	if len(product_hs_code_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	for _, product_hs_code := range product_hs_code_list {
		var product models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_hs_code.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_hs_code.sku
		product.CnHsCode = product_hs_code.cn_hs_code
		product.AuHsCode = product_hs_code.au_hs_code

		if has == true {
			update_col_list := []string{"cn_hs_code", "au_hs_code"}
			_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
			core.HandleError(err)
		} else {
			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_hs_code.sku))
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功",
	})
}

func (c ErpHsCodeCnReportElemExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM, c.Session)

	db_session := c.DbSession

	product_list := make([]models.ProductInfo, 0)

	err := db_session.Find(&product_list)
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

func (c ErpHsCodeCnReportElemExcel) export_product_list_to_excel(db_session *xorm.Session, product_list []models.ProductInfo) *xlsx.File {

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
