package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"github.com/tealeg/xlsx"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
)

type ErpContainerTransferReportElemExcel struct {
	*revel.Controller
}

func (c ErpContainerTransferReportElemExcel) Export() revel.Result {
	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var containerTransfer models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&containerTransfer)
	core.HandleError(err)

	xlsx_file := c.export_report_elem_list_to_excel(db_session, transfer_id)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("中国海关申报元素导出_%s.xlsx", containerTransfer.BatchCode) //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpContainerTransferReportElemExcel) export_report_elem_list_to_excel(db_session *xorm.Session, transfer_id int64) *xlsx.File {
	/*
SELECT
    DISTINCT
    hcc.*
FROM erp_container_transfer_entry cte
    INNER JOIN erp_hs_code_cn hcc ON hcc.id = cte.cn_hs_id
WHERE cte.transfer_id = 1
ORDER BY hcc.hs_code;
	 */

	var sql = `
SELECT
    DISTINCT
    hcc.*
FROM erp_container_transfer_entry cte
    INNER JOIN erp_hs_code_cn hcc ON hcc.id = cte.cn_hs_id
WHERE cte.transfer_id = %d
ORDER BY hcc.hs_code;`
	sql = fmt.Sprintf(sql, transfer_id)

	var hs_code_cn_list = make([]models.HsCodeCnReportElemInfo, 0)
	core.HandleError(db_session.Sql(sql).Find(&hs_code_cn_list))

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := *xlsx.NewFont(11, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "0092D050", "00000000") //92D050
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

	//统一设置内容单元格格式的函数
	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := *xlsx.NewFont(11, "微软雅黑")
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
	sheet, err := xlsx_file.AddSheet("集装箱调拨明细")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	for _, hs_code_cn := range hs_code_cn_list {
		// 新建行
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		cell := row.AddCell()
		set_head_cell_style(cell, "center")
		cell.Value = "商品编码"

		cell = row.AddCell()
		set_head_cell_style(cell, "left")
		cell.SetString(hs_code_cn.HsCode)

		// 新建行
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.SetString("品名")

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.SetString(hs_code_cn.CnName)

		if hs_code_cn.Elem01 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素1")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem01)
		}

		if hs_code_cn.Elem02 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素2")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem02)
		}

		if hs_code_cn.Elem03 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素3")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem03)
		}

		if hs_code_cn.Elem04 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素4")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem04)
		}

		if hs_code_cn.Elem05 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素5")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem05)
		}

		if hs_code_cn.Elem06 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素6")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem06)
		}

		if hs_code_cn.Elem06 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素6")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem06)
		}

		if hs_code_cn.Elem07 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素7")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem07)
		}

		if hs_code_cn.Elem08 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素8")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem08)
		}

		if hs_code_cn.Elem09 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素9")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem09)
		}

		if hs_code_cn.Elem10 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素10")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem10)
		}

		if hs_code_cn.Elem11 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素11")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem11)
		}

		if hs_code_cn.Elem12 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素12")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem12)
		}

		if hs_code_cn.Elem13 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素13")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem13)
		}

		if hs_code_cn.Elem14 != "" {
			row = sheet.AddRow()
			row.SetHeightCM(common_row_height)

			cell = row.AddCell()
			set_text_cell_style(cell, "center")
			cell.SetString("申报元素14")

			cell = row.AddCell()
			set_text_cell_style(cell, "left")
			cell.SetString(hs_code_cn.Elem14)
		}
	}

	return xlsx_file
}
