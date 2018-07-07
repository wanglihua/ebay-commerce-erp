package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/core"
	"fmt"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/erp/service/excel_service"
	"strings"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

const logistics_report_type_day int64 = 1   // 按天
const logistics_report_type_month int64 = 2 // 按月

type ErpLogisticsReport struct {
	*revel.Controller
}

func (c ErpLogisticsReport) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	c.RenderArgs["logistics_report_type_day"] = logistics_report_type_day
	c.RenderArgs["logistics_report_type_month"] = logistics_report_type_month

	return c.RenderTemplate("erp/ddbb/logistics_report/logistics_report_index.html")
}

func (c ErpLogisticsReport) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	var db_session = c.DbSession

	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)

	var logistics_company_list = make([]ddpz_models.LogisticsCompanyInfo, 0)
	err := db_session.OrderBy("ord asc").Find(&logistics_company_list)
	core.HandleError(err)

	var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err = db_session.OrderBy("name asc").Find(&logistics_mode_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"stock_id":               0,
		"stock_list":             stock_list,
		"logistics_company_id":   0,
		"logistics_company_list": logistics_company_list,
		"logistics_mode_id":      0,
		"logistics_mode_list":    logistics_mode_list,
	})
}

func (c ErpLogisticsReport) GetLogisticsModeList() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	var db_session = c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var logistics_company_id int64
	c.Params.Bind(&logistics_company_id, "logistics_company_id")

	var filter = "1=1"

	if stock_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("stock_id = %d", stock_id))
	}

	if logistics_company_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("company_id = %d", logistics_company_id))
	}

	var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.Where(filter).OrderBy("name asc").Find(&logistics_mode_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"logistics_mode_list": logistics_mode_list,
	})
}

func (c ErpLogisticsReport) apply_query_filter(filter string, is_export bool) string {
	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var logistics_company_id int64
	c.Params.Bind(&logistics_company_id, "logistics_company_id")

	var postalcode_query string
	c.Params.Bind(&postalcode_query, "postalcode_query")
	postalcode_query = strings.TrimSpace(postalcode_query)

	var begin_time_query core.NullLongTime
	var end_time_query core.NullLongTime
	c.Params.Bind(&begin_time_query, "begin_time_query")
	c.Params.Bind(&end_time_query, "end_time_query")

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	if (is_export == true && len(id_list) == 0) || is_export == false {

		if stock_id != 0 {
			// 这里是查询上的对仓库的限制，如果是拣货人员的话，下面还有是不是这个仓库的拣货人员的限制
			filter = core.SqlAnd(filter, fmt.Sprintf("stock_id = %d", stock_id))
		}

		if logistics_company_id != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("company_id = %d", logistics_company_id))
		}

		if postalcode_query != "" {
			// todo 分析多个邮编的情况
			filter = core.SqlAnd(filter, fmt.Sprintf("postcode = '%s'", postalcode_query))
		}

		if begin_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("day_date >= '%s'", begin_time_query.Time.Format("2006-01-02")))
		}

		if end_time_query.Valid == true {
			filter = core.SqlAnd(filter, fmt.Sprintf("day_date <= '%s'", end_time_query.Time.Format("2006-01-02")))
		}
	}

	if len(id_list) != 0 && is_export == true {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("id IN (%s)", strings.Join(id_str_list, ",")))
	}

	return filter
}

type logistics_report_view_item_t struct {
	GroupTime         string  `xorm:"'group_time'" json:"group_time"`
	WeightSum         float64 `xorm:"'weight_sum'" json:"weight_sum"`
	VolumeSum         float64 `xorm:"'volume_sum'" json:"volume_sum"`
	PrdCountSum       int64   `xorm:"'prd_count_sum'" json:"prd_count_sum"`
	OrderCount        int64   `xorm:"'order_count'" json:"order_count"`
	TotalFeeSum       float64 `xorm:"'total_fee_sum'" json:"total_fee_sum"`
	ActualTotalFeeSum float64 `xorm:"'actual_total_fee_sum'" json:"actual_total_fee_sum"`
	ProfitSum         float64 `xorm:"'profit_sum'" json:"profit_sum"`
}

func (c ErpLogisticsReport) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	db_session := c.DbSession

	filter, _, offset, limit := core.GetGridRequestParam(c.Request)

	var report_type int64
	c.Params.Bind(&report_type, "report_type")

	/*

SELECT
    date_format(day_date, '%Y-%c-%d') AS group_time,
    sum(weight)                       AS weight_sum,
    sum(volume)                       AS volume_sum,
    sum(prd_count)                    AS prd_count_sum,
    count(id)                         AS order_count,
    sum(total_fee)                    AS total_fee_sum,
    sum(actual_total_fee)             AS actual_total_fee_sum,
    sum(profit)                       AS profit_sum
FROM erp_logistics_report
WHERE stock_id = 1
GROUP BY day_date
ORDER BY day_date DESC;

	 */

	var sb = core.NewSqlBuilder()
	if report_type == logistics_report_type_month {
		sb.Select(`
    date_format(month_date, '%Y-%c')  AS group_time,
    sum(weight)                       AS weight_sum,
    sum(volume)                       AS volume_sum,
    sum(prd_count)                    AS prd_count_sum,
    count(id)                         AS order_count,
    sum(total_fee)                    AS total_fee_sum,
    sum(actual_total_fee)             AS actual_total_fee_sum,
    sum(profit)                       AS profit_sum
`)
	} else {
		sb.Select(`
    date_format(day_date, '%Y-%c-%d') AS group_time,
    sum(weight)                       AS weight_sum,
    sum(volume)                       AS volume_sum,
    sum(prd_count)                    AS prd_count_sum,
    count(id)                         AS order_count,
    sum(total_fee)                    AS total_fee_sum,
    sum(actual_total_fee)             AS actual_total_fee_sum,
    sum(profit)                       AS profit_sum
`)
	}

	sb.Count("count(id) AS count")
	sb.From(`erp_logistics_report`)

	filter = c.apply_query_filter(filter, false)
	sb.Where(filter)

	if report_type == logistics_report_type_month {
		sb.GroupBy("month_date")
		sb.OrderBy("month_date DESC")
	} else {
		sb.GroupBy("day_date")
		sb.OrderBy("day_date DESC")
	}

	sb.Limit(offset, limit)

	var logistics_report_list = make([]logistics_report_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_report_list)
	core.HandleError(err)

	var count core.Count
	_, err = db_session.Sql(sb.SqlCount()).Get(&count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_report_list,
		Total: count.Count,
	})
}

func (c ErpLogisticsReport) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDBB_LOGISTICS_REPORT, c.Session)

	var db_session = c.DbSession

	xlsx_file := c.get_export_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 iReader 接口

	file_name := "物流报表" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpLogisticsReport) get_export_excel(db_session *xorm.Session) *xlsx.File {
	var filter = "(1=1)"

	var report_type int64
	c.Params.Bind(&report_type, "report_type")

	/*

SELECT
    date_format(day_date, '%Y-%c-%d') AS group_time,
    sum(weight)                       AS weight_sum,
    sum(volume)                       AS volume_sum,
    sum(prd_count)                    AS prd_count_sum,
    count(id)                         AS order_count,
    sum(total_fee)                    AS total_fee_sum,
    sum(actual_total_fee)             AS actual_total_fee_sum,
    sum(profit)                       AS profit_sum
FROM erp_logistics_report
WHERE stock_id = 1
GROUP BY day_date
ORDER BY day_date DESC;

	 */

	var sb = core.NewSqlBuilder()
	if report_type == logistics_report_type_month {
		sb.Select(`
    date_format(month_date, '%Y-%c')  AS group_time,
    sum(weight)                       AS weight_sum,
    sum(volume)                       AS volume_sum,
    sum(prd_count)                    AS prd_count_sum,
    count(id)                         AS order_count,
    sum(total_fee)                    AS total_fee_sum,
    sum(actual_total_fee)             AS actual_total_fee_sum,
    sum(profit)                       AS profit_sum
`)
	} else {
		sb.Select(`
    date_format(day_date, '%Y-%c-%d') AS group_time,
    sum(weight)                       AS weight_sum,
    sum(volume)                       AS volume_sum,
    sum(prd_count)                    AS prd_count_sum,
    count(id)                         AS order_count,
    sum(total_fee)                    AS total_fee_sum,
    sum(actual_total_fee)             AS actual_total_fee_sum,
    sum(profit)                       AS profit_sum
`)
	}

	sb.Count("count(id) AS count")
	sb.From(`erp_logistics_report`)

	filter = c.apply_query_filter(filter, false)
	sb.Where(filter)

	if report_type == logistics_report_type_month {
		sb.GroupBy("month_date")
		sb.OrderBy("month_date DESC")
	} else {
		sb.GroupBy("day_date")
		sb.OrderBy("day_date DESC")
	}

	var logistics_report_list = make([]logistics_report_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_report_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("物流报表")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		excel_service.Set_head_cell_style(cell)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		excel_service.Set_text_cell_style(cell, horizontal_align)
	}

	row_height := 0.6

	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "序号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "时间"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "重量"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "体积"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "商品数量"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "单数"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "预估运费"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "实际运费"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "总利润"
	}

	for index, logistics_report := range logistics_report_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 序号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = logistics_report.GroupTime

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", logistics_report.WeightSum)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", logistics_report.VolumeSum)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprint(logistics_report.PrdCountSum)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprint(logistics_report.OrderCount)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", logistics_report.TotalFeeSum)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", logistics_report.ActualTotalFeeSum)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", logistics_report.ProfitSum)
	}

	return xlsx_file
}
