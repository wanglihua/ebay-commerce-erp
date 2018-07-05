package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"time"
)

type ErpRefundReport struct {
	*revel.Controller
}

const REPORT_TYPE_DAY int64 = 1

const REPORT_TYPE_MONTH int64 = 2

func (c ErpRefundReport) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	// var db_session = c.DbSession

	c.RenderArgs["report_type_day"] = REPORT_TYPE_DAY
	c.RenderArgs["report_type_month"] = REPORT_TYPE_MONTH

	return c.RenderTemplate("erp/ppgl/refund_report/refund_report_index.html")
}

func (c ErpRefundReport) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	var report_type int64
	c.Params.Bind(&report_type, "report_type")

	var sales_user_id int64
	c.Params.Bind(&sales_user_id, "sales_user_id")

	var shop_id int64
	c.Params.Bind(&shop_id, "shop_id")

	var time_begin core.NullLongTime
	c.Params.Bind(&time_begin, "time_begin")

	var time_end core.NullLongTime
	c.Params.Bind(&time_end, "time_end")

	filter, _, offset, limit := core.GetGridRequestParam(c.Request)

	var db_session = c.DbSession

	var update_apply_date_month_sql = `UPDATE erp_order_refund
SET apply_date  = date(apply_time),
    apply_month = date_add(date(apply_time), INTERVAL (1 - day(apply_time)) DAY)
WHERE apply_date IS NULL OR apply_month IS NULL;`

	_, err := db_session.Exec(update_apply_date_month_sql)
	core.HandleError(err)

	type refund_report_item_t struct {
		MonthDate       string  `xorm:"'month_date'" json:"month_date"`
		RefundAmountSum float64 `xorm:"'refund_amount_sum'" json:"refund_amount_sum"`
	}

	/*

SELECT apply_date AS month_date, sum(refund_amount) AS refund_amount_sum
FROM erp_order_refund
WHERE shop_id = 1 AND sales_user_id = 1
GROUP BY apply_date
ORDER BY month_date DESC
LIMIT 0, 10;

SELECT apply_month AS month_date, sum(refund_amount) AS refund_amount_sum
FROM erp_order_refund
WHERE shop_id = 1 AND sales_user_id = 1
GROUP BY apply_month
ORDER BY month_date DESC
LIMIT 0, 10;

	 */

	var sb = core.NewSqlBuilder()
	if report_type == REPORT_TYPE_MONTH {
		sb.Select(`apply_month AS month_date, sum(refund_amount) AS refund_amount_sum`)
	} else {
		sb.Select(`apply_date AS month_date, sum(refund_amount) AS refund_amount_sum`)
	}

	sb.Count("count(*) as count")
	sb.From(`erp_order_refund`)

	sb.Where(filter)

	if report_type == REPORT_TYPE_MONTH {
		sb.GroupBy("apply_month")
	} else {
		sb.GroupBy("apply_date")
	}

	sb.OrderBy("month_date DESC")
	sb.Limit(offset, limit)

	var refund_report_item_list = make([]refund_report_item_t, 0)
	err = db_session.Sql(sb.SqlList()).Find(refund_report_item_list)
	core.HandleError(err)

	var count core.Count
	_, err = db_session.Sql(sb.SqlCount()).Get(&count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  refund_report_item_list,
		Total: count.Count,
	})
}

func (c ErpRefundReport) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	// var db_session = c.DbSession

	/*

	xlsx_file := c.get_export_excel_content(db_session, order_refund_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("订单退款导出_%s.xlsx", time.Now().Format("2016-01-02 15:04"))

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
	*/

	return c.RenderBinary(nil, "", revel.Attachment, time.Now())
}
