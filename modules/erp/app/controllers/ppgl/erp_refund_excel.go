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
	"matrix/modules/erp/service/excel_service"
)

type ErpRefundExcel struct {
	*revel.Controller
}

func (c ErpRefundExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	var db_session = c.DbSession

	var query_begin_time string
	c.Params.Bind(&query_begin_time, "begin_time_query")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)

	var query_end_time string
	c.Params.Bind(&query_end_time, "end_time_query")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)

	var filter = "1=1"

	if query_begin_date_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("apply_time >= '%04d-%02d-%02d'", query_begin_date.Year(), query_begin_date.Month(), query_begin_date.Day()))
	}

	if query_end_date_err == nil {
		query_end_date = query_end_date.Add(time.Hour * 24)
		filter = core.SqlAnd(filter, fmt.Sprintf("apply_time < '%04d-%02d-%02d'", query_end_date.Year(), query_end_date.Month(), query_end_date.Day()))
	}

	query := db_session.
		Where(filter)

	data_query := *query

	data_query = *data_query.Asc("id")

	order_refund_list := make([]ppgl_models.OrderRefundInfo, 0)
	err := data_query.Find(&order_refund_list)
	core.HandleError(err)

	xlsx_file := c.get_export_excel_content(db_session, order_refund_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("订单退款导出_%s.xlsx", time.Now().Format("2006-01-02 15:04"))

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpRefundExcel) get_export_excel_content(db_session *xorm.Session, order_refund_list []ppgl_models.OrderRefundInfo) *xlsx.File {

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		excel_service.Set_head_cell_style(cell)
	}

	//统一设置内容单元格格式的函数
	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		excel_service.Set_text_cell_style(cell, horizontal_align)
	}

	head_row_height := 0.8
	common_row_height := 0.6

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	//添加工作表
	sheet, err := xlsx_file.AddSheet("订单退款")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "店铺"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单编号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "付款流水号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "退款流水号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "订单金额"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 14)
		set_head_cell_style(cell)
		cell.Value = "退款金额"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "退款方式"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "退款类型"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "退款原因"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "申请人"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 28)
		set_head_cell_style(cell)
		cell.Value = "申请时间"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "备注"
		col_index++
	}

	for _, refund := range order_refund_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.ShopName

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.OrderCode

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.PayTransId

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.TransactionId

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloatWithFormat(refund.OrderAmount, "0.00")

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.SetFloatWithFormat(refund.RefundAmount, "0.00")

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.MethodName

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.TypeName

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.Reason

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.StatusName

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.ApplyUserName

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.ApplyTime.String()

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = refund.Description
	}

	return xlsx_file
}
