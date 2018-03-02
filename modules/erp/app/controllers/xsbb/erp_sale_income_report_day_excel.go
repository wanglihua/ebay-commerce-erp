package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xsbb_models"
	"fmt"
	"time"
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpSaleIncomeReportDayExcel struct {
	*revel.Controller
}

func (c ErpSaleIncomeReportDayExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSBB_SALE_INCOME_REPORT_DAY, c.Session)

	db_session := c.DbSession

	var begin_time_query core.NullLongTime
	c.Params.Bind(&begin_time_query, "begin_time_query")

	var end_time_query core.NullLongTime
	c.Params.Bind(&end_time_query, "end_time_query")

	var platform_id_query int64
	c.Params.Bind(&platform_id_query, "platform_id_query")

	var shop_id_query int64
	c.Params.Bind(&shop_id_query, "shop_id_query")

	var user_type_query int64
	c.Params.Bind(&user_type_query, "user_type_query")

	var user_id_query int64
	c.Params.Bind(&user_id_query, "user_id_query")

	var filter = "1=1"

	if begin_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.paid_time >= '%s'", begin_time_query.Time.Format("2006-01-02")))
	}

	if end_time_query.Valid == true {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.paid_time < '%s'", end_time_query.Time.AddDate(0, 0, 1).Format("2006-01-02")))
	}

	if platform_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.platform_id = %d", platform_id_query))
	}

	if shop_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("r.shop_id = %d", shop_id_query))
	}

	if user_type_query != 0 && user_id_query != 0 {
		if user_type_query == user_query_type_sale {
			filter = core.SqlAnd(filter, fmt.Sprintf("r.sale_user_id = %d", user_id_query))
		}

		if user_type_query == user_query_type_service {
			filter = core.SqlAnd(filter, fmt.Sprintf("r.service_user_id = %d", user_id_query))
		}

		if user_type_query == user_query_type_dev {
			filter = core.SqlAnd(filter, fmt.Sprintf("r.dev_user_id = %d", user_id_query))
		}
	}

	/*

	SELECT
		date_format(r.paid_time, '%Y-%m-%d')            AS paid_time,
		sum(r.paid_amount)                              AS paid_amount,
		sum(r.express_fee)                              AS express_fee,
		sum(r.income)                                   AS income,
		sum(r.cost_amount)                              AS cost_amount,
		sum(r.estimate_logistics_fee)                   AS estimate_logistics_fee,
		sum(r.actual_logistics_fee)                     AS actual_logistics_fee,
		sum(r.jyfy_amount)                              AS jyfy_amount,
		sum(r.ptfy_amount)                              AS ptfy_amount,
		sum(r.pick_fee)                                 AS pick_fee,
		sum(r.gst_fee)                                  AS gst_fee,
		sum(r.outcome)                                  AS outcome,
		sum(r.income - r.outcome)                       AS profit,
		sum(r.income - r.outcome) / sum(r.income) * 100 AS profit_rate,
	    count(DISTINCT r.order_id)                      AS order_qty,
    	sum(r.product_qty)                              AS product_qty
	FROM erp_sale_income_report r
	WHERE r.platform_id = 0
		  AND r.shop_id = 0
		  AND r.paid_time >= '2006-01-01'
		  AND r.paid_time < '2017-12-30'
		  AND r.product_id = 0
		  AND r.sale_user_id = 0
		  AND r.service_user_id = 0
		  AND r.dev_user_id = 0
	GROUP BY date_format(r.paid_time, '%Y-%m-%d')
	ORDER BY paid_time DESC;

	 */

	var sql_list = `
SELECT
    date_format(r.paid_time, '%%Y-%%m-%%d')         AS paid_time,
    sum(r.paid_amount)                              AS paid_amount,
    sum(r.express_fee)                              AS express_fee,
    sum(r.income)                                   AS income,
    sum(r.cost_amount)                              AS cost_amount,
    sum(r.estimate_logistics_fee)                   AS estimate_logistics_fee,
    sum(r.actual_logistics_fee)                     AS actual_logistics_fee,
    sum(r.jyfy_amount)                              AS jyfy_amount,
    sum(r.ptfy_amount)                              AS ptfy_amount,
    sum(r.pick_fee)                                 AS pick_fee,
    sum(r.gst_fee)                                  AS gst_fee,
    sum(r.outcome)                                  AS outcome,
    sum(r.income - r.outcome)                       AS profit,
    sum(r.income - r.outcome) / sum(r.income) * 100 AS profit_rate,
    count(DISTINCT r.order_id)                      AS order_qty,
    sum(r.product_qty)                              AS product_qty
FROM erp_sale_income_report r
WHERE %s
GROUP BY date_format(r.paid_time, '%%Y-%%m-%%d')
ORDER BY %s;
`

	var order = "paid_time DESC"

	sql_list = fmt.Sprintf(sql_list, filter, order)

	var sale_amount_report_list = make([]xsbb_models.SaleAmountReportInfo, 0)
	core.HandleError(db_session.Sql(sql_list).Find(&sale_amount_report_list))

	var sql_total = `
	SELECT
		''                                              AS paid_time,
		sum(r.paid_amount)                              AS paid_amount,
		sum(r.express_fee)                              AS express_fee,
		sum(r.income)                                   AS income,
		sum(r.cost_amount)                              AS cost_amount,
		sum(r.estimate_logistics_fee)                   AS estimate_logistics_fee,
		sum(r.actual_logistics_fee)                     AS actual_logistics_fee,
		sum(r.jyfy_amount)                              AS jyfy_amount,
		sum(r.ptfy_amount)                              AS ptfy_amount,
		sum(r.pick_fee)                                 AS pick_fee,
		sum(r.gst_fee)                                  AS gst_fee,
		sum(r.outcome)                                  AS outcome,
		sum(r.income - r.outcome)                       AS profit,
		sum(r.income - r.outcome) / sum(r.income) * 100 AS profit_rate,
	    count(DISTINCT r.order_id)                      AS order_qty,
    	sum(r.product_qty)                              AS product_qty
	FROM erp_sale_income_report r
	WHERE %s;
`
	sql_total = fmt.Sprintf(sql_total, filter)

	var total xsbb_models.SaleAmountReportInfo
	_, err := db_session.Sql(sql_total).Get(&total)
	core.HandleError(err)

	var xlsx_file = c.export_sale_amount_report_to_excel(db_session, sale_amount_report_list, total)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "日收入报表.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpSaleIncomeReportDayExcel) export_sale_amount_report_to_excel(db_session *xorm.Session, sale_amount_report_day_list []xsbb_models.SaleAmountReportInfo, total xsbb_models.SaleAmountReportInfo) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品日销售额")
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
		cell.Value = "收入小计（AUD）"
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
		cell.Value = "支出小计（AUD）"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "订单数量"
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
		cell.Value = report.PaidTime

		// 交易金额
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.PaidAmount)

		// 运费收入
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.ExpressFee)

		// 收入小计
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.Income)

		// 商品澳仓成本
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.CostAmount)

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

		// 支出小计
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.Outcome)

		// 订单数量
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprint(report.OrderQty)

		// 商品数量
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprint(report.ProductQty)

		// 利润
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", report.Profit)

		// 利润率
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f%%", report.ProfitRate)
	}

	{
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 日期
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = "合计"

		// 交易金额
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.PaidAmount)

		// 运费收入
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.ExpressFee)

		// 收入小计
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.Income)

		// 商品澳仓成本
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.CostAmount)

		// 预估运费
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.EstimateLogisticsFee)

		// 实际运费
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.ActualLogisticsFee)

		// 交易费用
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.JyfyAmount)

		// 平台费用
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.PtfyAmount)

		// 拣货费
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.PickFee)

		// GST费用
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.GstFee)

		// 支出小计
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.Outcome)

		// 订单数量
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprint(total.OrderQty)

		// 商品数量
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprint(total.ProductQty)

		// 利润
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", total.Profit)

		// 利润率
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f%%", total.ProfitRate)
	}

	return xlsx_file
}
