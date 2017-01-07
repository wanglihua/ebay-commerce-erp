package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"bytes"
	"time"
)

type ErpProductExport struct {
	*revel.Controller
}

func (c ErpProductExport) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_export/product_export.html")
}

func (c ErpProductExport) ListData() revel.Result {
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

	product_list := make([]models.ProductInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: count,
	})
}

//todo: 做个QueryForm来接收查询参数? ListData() 和 Export() 里面

func (c ErpProductExport) Export() revel.Result {
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

	file_name := "商品导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpProductExport) export_product_list_to_excel(db_session *xorm.Session, product_list []models.ProductInfo) *xlsx.File {

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
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

		font := * xlsx.NewFont(9, "微软雅黑")
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

	cell := row.AddCell()
	sheet.SetColWidth(col_index, col_index, 10)
	set_head_cell_style(cell)
	cell.Value = "商品名称"
	col_index ++

	for _, product := range product_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		cell := row.AddCell() //商品名
		set_text_cell_style(cell, "left")
		cell.Value = product.Name
	}

	return xlsx_file
}
