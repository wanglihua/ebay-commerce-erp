package controllers

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
	"strings"
)

type ErpOrderExport struct {
	*revel.Controller
}

func (c ErpOrderExport) Export() revel.Result {
	var status string
	c.Params.Bind(&status, "status")

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var hold_type_id int64
	c.Params.Bind(&hold_type_id, "hold_type_id")

	var db_session = c.DbSession

	var status_id int64 = 0
	if status == "dph" {
		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		status_id = models.ORDER_STATUS_PHZ
	} else if status == "hold" {
		status_id = models.ORDER_STATUS_HOLD
	}

	// 先看 URL 里面的状态分类
	status = strings.ToLower(status) //大小写
	if status == "all" {
		status_id = 0
	} else if status == "ycj" {
		status_id = models.ORDER_STATUS_YCJ
	} else if status == "yzf" {
		status_id = models.ORDER_STATUS_YZF
	} else if status == "dph" {
		status_id = models.ORDER_STATUS_DPH
	} else if status == "phz" {
		status_id = models.ORDER_STATUS_PHZ
	} else if status == "yfh" {
		status_id = models.ORDER_STATUS_YFH
	} else if status == "ywc" {
		status_id = models.ORDER_STATUS_YWC
	} else if status == "yqx" {
		status_id = models.ORDER_STATUS_YQX
	}

	var filter = "1=1"

	if status_id != 0 {
		filter = fmt.Sprintf("o.status_id = %d", status_id)
	}

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
	}

	if hold_type_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
	}

	xlsx_file := c.get_export_excel(db_session, filter)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "订单编码订单代码" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

type order_code_num_export_view_t struct {
	Order    models.OrderInfo          `xorm:"extends" json:"o"`
	HoldType models.OrderHoldTypeInfo  `xorm:"extends" json:"ht"`
}

func (c ErpOrderExport) get_export_excel(db_session *xorm.Session, filter string) *xlsx.File {

	/*
SELECT o.*, ht.*
FROM erp_order o
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
	 */

	var sql = `SELECT o.*, ht.*
FROM erp_order o
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
WHERE `

	sql = sql + filter

	var order_list = make([]order_code_num_export_view_t, 0)
	err := db_session.Sql(sql).Find(&order_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("订单编码订单代码")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
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
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyFont = true

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

	row_height := 0.6

	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0

		var cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "订单编码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单代码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "hold原因"
	}

	for _, order := range order_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 订单编码
		var cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order.Order.Code

		// 订单代码
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order.Order.FullTranId

		// hold原因
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = order.HoldType.Name
	}

	return xlsx_file
}
