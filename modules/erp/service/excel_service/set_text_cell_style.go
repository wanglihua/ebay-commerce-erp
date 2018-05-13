package excel_service

import "github.com/tealeg/xlsx"

//统一设置内容单元格格式的函数
func Set_text_cell_style(cell *xlsx.Cell, horizontal_align string) {
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
