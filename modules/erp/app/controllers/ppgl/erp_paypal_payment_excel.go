package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/erp/models/ppgl_models"
	"bytes"
	"fmt"
	"time"
)

type ErpPaypalPaymentExcel struct {
	*revel.Controller
}

func (c ErpPaypalPaymentExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_PAYPAL_PAYMENT, c.Session)

	var db_session = c.DbSession

	var account_id_query int64
	c.Params.Bind(&account_id_query, "account_id_query")

	var query_begin_time string
	c.Params.Bind(&query_begin_time, "begin_time_query")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)
	revel.TRACE.Println("query_begin_time", query_begin_time)

	var query_end_time string
	c.Params.Bind(&query_end_time, "end_time_query")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if account_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("account_id = %d", account_id_query))
	}

	if query_begin_date_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("pay_time >= '%04d-%02d-%02d'", query_begin_date.Year(), query_begin_date.Month(), query_begin_date.Day()))
	}

	if query_end_date_err == nil {
		query_end_date = query_end_date.Add(time.Hour * 24)
		filter = core.SqlAnd(filter, fmt.Sprintf("pay_time < '%04d-%02d-%02d'", query_end_date.Year(), query_end_date.Month(), query_end_date.Day()))
	}

	query := db_session.
		Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	paypal_payment_list := make([]ppgl_models.PaypalPaymentInfo, 0)
	err := data_query.Limit(limit, offset).Find(&paypal_payment_list)
	core.HandleError(err)

	xlsx_file := c.get_export_excel_content(db_session, paypal_payment_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("PayPal付款导出_%s.xlsx", time.Now().Format("2016-01-02 15:04"))

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpPaypalPaymentExcel) get_export_excel_content(db_session *xorm.Session, payment_list []ppgl_models.PaypalPaymentInfo) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("智能补货")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "Paypal Account"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "Transaction ID"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "Type"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "Time"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "Gross Amount"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "Fee Amount"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "Net Amount"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "Status"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "Payer"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "Payer Display Name"
		col_index++
	}

	for _, payment := range payment_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = payment.PaypalAccount

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = payment.TransactionId

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = payment.Type

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = payment.PayTime.String()

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloatWithFormat(payment.GrossAmount, "0.00")

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloatWithFormat(payment.FeeAmount, "0.00")

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloatWithFormat(payment.NetAmount, "0.00")

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = payment.Status

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = payment.Payer

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = payment.PayerDisplayName
	}

	return xlsx_file
}
