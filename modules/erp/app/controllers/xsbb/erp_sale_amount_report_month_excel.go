package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
	"fmt"
	"time"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
)

type ErpSaleAmountReportMonthExcel struct {
	*revel.Controller
}

func (c ErpSaleAmountReportMonthExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPZL_PRODUCT_FEEDBACK_QUERY, c.Session)

	var db_session = c.DbSession

	var begin_time_query core.NullLongTime
	c.Params.Bind(&begin_time_query, "begin_time_query")

	var end_time_query core.NullLongTime
	c.Params.Bind(&end_time_query, "end_time_query")

	filter := "1=1"

	if begin_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("report_month >= '%s'", begin_time_query.Time.Format("2006-01-02")))
	}

	if end_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("report_month < '%s'", end_time_query.Time.AddDate(0, 0, 1).Format("2006-01-02")))
	}

	var query = db_session.Where(filter)

	var sale_amount_report_month_list = make([]xsbb_models.SaleAmountReportMonthInfo, 0)
	err := query.OrderBy("report_month ASC").Find(&sale_amount_report_month_list)
	core.HandleError(err)

	var xlsx_file = c.export_sale_amount_report_to_excel(db_session, sale_amount_report_month_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "月销售额报表.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpSaleAmountReportMonthExcel) export_sale_amount_report_to_excel(db_session *xorm.Session, sale_amount_report_day_list []xsbb_models.SaleAmountReportMonthInfo) *xlsx.File {

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

	head_row_height := 1.6
	common_row_height := 0.6

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	//添加工作表
	sheet, err := xlsx_file.AddSheet("商品质量反馈")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "日期"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "交易金额（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "运费收入（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "商品澳仓成本（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "预估运费（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "实际运费（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "交易费用（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "平台费用（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "拣货费（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "GST费用（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "商品数量（件）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "利润（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "利润率"
		col_index++
	}

	for _, report := range sale_amount_report_day_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 日期
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf(time.Time(report.ReportMonth).Format("2006-01"))

		// 交易金额
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.PaidAmount)

		// 运费收入
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.ExpressFee)

		// 商品澳仓成本
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.CostPrice)

		// 预估运费
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.EstimateLogisticsFee)

		// 实际运费
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.ActualLogisticsFee)

		// 交易费用
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.JyfyAmount)

		// 平台费用
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.PtfyAmount)

		// 拣货费
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.PickFee)

		// GST费用
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.GstFee)

		// 商品数量
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%d", report.ProductQty)

		// 利润
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.0f", report.Profit)

		// 利润率
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f%%", report.ProfitRate)
	}

	return xlsx_file
}
