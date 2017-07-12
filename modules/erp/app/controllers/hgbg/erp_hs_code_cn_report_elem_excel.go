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

	type hs_code_import_t struct {
		hs_code string
		cn_name string
		en_name string
		elem01  string
		elem02  string
		elem03  string
		elem04  string
		elem05  string
		elem06  string
		elem07  string
		elem08  string
		elem09  string
		elem10  string
		elem11  string
		elem12  string
		elem13  string
		elem14  string
	}

	hs_code_import_list := make([]hs_code_import_t, 0)

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

			var hs_code_import hs_code_import_t

			col_index := 0

			// 中国海关编码
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			hs_code_import.hs_code = cell_val
			col_index++

			// 申报中文品名
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.cn_name = cell_val
			col_index++

			// 申报英文品名
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.en_name = cell_val
			col_index++

			// 要素1
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem01 = cell_val
			col_index++

			// 要素2
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem02 = cell_val
			col_index++

			// 要素3
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem03 = cell_val
			col_index++

			// 要素4
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem04 = cell_val
			col_index++

			// 要素5
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem05 = cell_val
			col_index++

			// 要素6
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem06 = cell_val
			col_index++

			// 要素7
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem07 = cell_val
			col_index++

			// 要素8
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem08 = cell_val
			col_index++

			// 要素9
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem09 = cell_val
			col_index++

			// 要素10
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem10 = cell_val
			col_index++

			// 要素11
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem11 = cell_val
			col_index++

			// 要素12
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem12 = cell_val
			col_index++

			// 要素13
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem13 = cell_val
			col_index++

			// 要素14
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_import.elem14 = cell_val
			col_index++

			hs_code_import_list = append(hs_code_import_list, hs_code_import)
		}
	}

	if len(hs_code_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入海关编码列表为空！",
		})
	}

	//todo: 验证导入的数据

	for _, hs_code_import := range hs_code_import_list {
		var hs_code models.HsCodeCnReportElemInfo
		has, err := db_session.Where(fmt.Sprintf("hs_code = '%s'", hs_code_import.hs_code)).Get(&hs_code)
		core.HandleError(err)

		hs_code.HsCode = hs_code_import.hs_code
		hs_code.CnName = hs_code_import.cn_name
		hs_code.EnName = hs_code_import.en_name
		hs_code.Elem01 = hs_code_import.elem01
		hs_code.Elem02 = hs_code_import.elem02
		hs_code.Elem03 = hs_code_import.elem03
		hs_code.Elem04 = hs_code_import.elem04
		hs_code.Elem05 = hs_code_import.elem05
		hs_code.Elem06 = hs_code_import.elem06
		hs_code.Elem07 = hs_code_import.elem07
		hs_code.Elem08 = hs_code_import.elem08
		hs_code.Elem09 = hs_code_import.elem09
		hs_code.Elem10 = hs_code_import.elem10
		hs_code.Elem11 = hs_code_import.elem11
		hs_code.Elem12 = hs_code_import.elem12
		hs_code.Elem13 = hs_code_import.elem13
		hs_code.Elem14 = hs_code_import.elem14

		if has == true {
			update_col_list := []string{"hs_code", "cn_name", "en_name", "elem01", "elem02", "elem03", "elem04", "elem05",
			                            "elem06", "elem07", "elem08", "elem09", "elem10", "elem11", "elem12", "elem13", "elem14"}
			_, err = db_session.Id(hs_code.Id).Cols(update_col_list...).Update(&hs_code)
			core.HandleError(err)
		} else {
			_, err = db_session.Insert(&hs_code)
			core.HandleError(err)
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

	product_list := make([]models.HsCodeCnReportElemInfo, 0)

	err := db_session.Find(&product_list)
	core.HandleError(err)

	xlsx_file := c.export_hs_code_list_to_excel(db_session, product_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "中国海关编码导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpHsCodeCnReportElemExcel) export_hs_code_list_to_excel(db_session *xorm.Session, hs_code_list []models.HsCodeCnReportElemInfo) *xlsx.File {

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

	// 中国海关编码 	中文品名	英文品名	要素1	要素2	要素3	要素4	要素5	要素6	要素7	要素8	要素9	要素10	要素11	要素12	要素13	要素14

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("中国海关编码")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("中文品名")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("英文品名")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素1")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素2")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素3")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素4")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素5")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素6")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素7")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素8")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素9")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素10")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素11")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素12")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素13")
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.SetString("要素14")
		col_index++
	}

	for _, product := range hs_code_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 中国海关编码 	中文品名	英文品名	要素1	要素2	要素3	要素4	要素5	要素6	要素7	要素8	要素9	要素10	要素11	要素12	要素13	要素14

		// 中国海关编码
		cell := row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.HsCode)

		// 中文品名
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.CnName)

		// 英文品名
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.EnName)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem01)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem02)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem03)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem04)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem05)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem06)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem07)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem08)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem09)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem10)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem11)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem12)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem13)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(product.Elem14)
	}

	return xlsx_file
}
