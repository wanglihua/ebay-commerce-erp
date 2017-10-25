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
	"matrix/modules/erp/models/hgbg_models"
	"strings"
	"strconv"
)

type ErpHsCodeAuDutyExcel struct {
	*revel.Controller
}

func (c ErpHsCodeAuDutyExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_AU_DUTY, c.Session)

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

	type hs_code_au_import_t struct {
		HsCode       string `json:"hs_code"`
		CnName       string `json:"cn_name"`
		EnName       string `json:"en_name"`
		DutyRate     string `json:"duty_rate"`
		ErrorMessage string `json:"error_message"`

		duty_rate float64 `json:"-"`
	}

	hs_code_au_import_list := make([]hs_code_au_import_t, 0)

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

			var hs_code_au hs_code_au_import_t

			col_index := 0

			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val = cell.Value
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			hs_code_au.HsCode = cell_val
			col_index++

			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val = cell.Value
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_au.CnName = cell_val
			col_index++

			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val = cell.Value
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_au.EnName = cell_val
			col_index++

			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val = cell.Value
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			hs_code_au.DutyRate = cell_val // 5% 会自动变为 0.05
			col_index++

			hs_code_au_import_list = append(hs_code_au_import_list, hs_code_au)
		}
	}

	if len(hs_code_au_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入海关编码列表为空！",
		})
	}

	// 验证导入的数据
	var hs_code_au_import_error_list = make([]hs_code_au_import_t, 0)
	for index, hs_code_au_import := range hs_code_au_import_list {
		if strings.TrimSpace(hs_code_au_import.DutyRate) == "" {
			hs_code_au_import.duty_rate = 0
		} else {
			duty_rate, err := strconv.ParseFloat(hs_code_au_import.DutyRate, 64)
			revel.TRACE.Println(`strings.TrimSpace(strings.Trim(hs_code_au_import.DutyRate, "%%"))`, strings.TrimSpace(strings.Trim(hs_code_au_import.DutyRate, "%")))
			revel.TRACE.Println("duty_rate", duty_rate)
			if err != nil {
				hs_code_au_import.ErrorMessage = "关税税率输入有误！"
				hs_code_au_import_error_list = append(hs_code_au_import_error_list, hs_code_au_import)
			}
			hs_code_au_import_list[index].duty_rate = duty_rate
		}
	}

	if len(hs_code_au_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "导入出现错误，详细请查看页面下部列表！",
			"error_data_list": hs_code_au_import_error_list,
		})
	}

	for _, hs_code_au_import := range hs_code_au_import_list {
		var hs_code_au hgbg_models.HsCodeAuDutyInfo
		has, err := db_session.Where("en_name = ?", hs_code_au_import.EnName).Get(&hs_code_au)
		core.HandleError(err)

		hs_code_au.HsCode = hs_code_au_import.HsCode
		hs_code_au.CnName = hs_code_au_import.CnName
		hs_code_au.EnName = hs_code_au_import.EnName
		hs_code_au.DutyRate = hs_code_au_import.duty_rate

		if has == true {
			update_col_list := []string{"hs_code", "cn_name", "en_name", "duty_rate"}
			_, err = db_session.Id(hs_code_au.Id).Cols(update_col_list...).Update(&hs_code_au)
			core.HandleError(err)
		} else {
			_, err = db_session.Insert(&hs_code_au)
			core.HandleError(err)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpHsCodeAuDutyExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_AU_DUTY, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "(1=1)"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("id IN (%s)", strings.Join(id_str_list, ",")))
	}

	hs_code_au_list := make([]hgbg_models.HsCodeAuDutyInfo, 0)

	err := db_session.Where(filter).Find(&hs_code_au_list)
	core.HandleError(err)

	xlsx_file := c.export_hs_code_list_to_excel(db_session, hs_code_au_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "澳洲海关编码导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpHsCodeAuDutyExcel) export_hs_code_list_to_excel(db_session *xorm.Session, hs_code_au_list []hgbg_models.HsCodeAuDutyInfo) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("澳洲海关编码列表")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	//澳洲海关编码	申报中文品名	申报英文品名	关税税率

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "澳洲海关编码"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "申报中文品名"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "申报英文品名"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "关税税率"
		col_index++
	}

	for _, hs_code := range hs_code_au_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 澳洲海关编码
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetString(hs_code.HsCode)

		// 申报中文品名
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetString(hs_code.CnName)

		// 申报英文品名
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetString(hs_code.EnName)

		// 关税税率
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloat(hs_code.DutyRate)
	}

	return xlsx_file
}
