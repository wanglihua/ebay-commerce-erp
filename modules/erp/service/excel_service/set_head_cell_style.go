package excel_service

import "github.com/tealeg/xlsx"

func Set_head_cell_style(cell *xlsx.Cell) {
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
