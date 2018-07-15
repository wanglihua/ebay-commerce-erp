package service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"fmt"
)

func Get_order_amount(shop_id int64, paypal_account_id int64, month_date string, report_type int64, db_session *xorm.Session) float64 {
	/*
SELECT sum(o.amount_paid) AS amount_sum
FROM erp_order o
WHERE o.shop_id = 1
      AND o.id IN (
    SELECT pp.order_id
    FROM erp_paypal_payment pp
    WHERE pp.account_id = 1)
      AND o.paid_time >= '2018-01-01 00:00:00'
      # AND o.paid_time < date_add('2018-01-01 00:00:00', INTERVAL 1 DAY);
      AND o.paid_time < date_add('2018-01-01 00:00:00', INTERVAL 1 MONTH);
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("sum(o.amount_paid) AS amount_sum")
	sb.From("erp_order o")

	var filter = "1=1"

	if shop_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.shop_id = %d", shop_id))
	}

	if paypal_account_id != 0 {
		var paypal_account_id_filter = `
o.id IN (
    SELECT pp.order_id
    FROM erp_paypal_payment pp
    WHERE pp.account_id = %d)
`
		paypal_account_id_filter = fmt.Sprintf(paypal_account_id_filter, paypal_account_id)
	}

	var time_filter = ""
	var time_begin = ""
	if report_type == REPORT_TYPE_DAY {
		time_begin = month_date + " 00:00:00"
		time_filter = `o.paid_time >= '%s' AND o.paid_time < date_add('%s', INTERVAL 1 DAY)`
	} else if report_type == REPORT_TYPE_MONTH {
		time_begin = month_date + "-01 00:00:00"
		time_filter = `o.paid_time >= '%s' AND o.paid_time < date_add('%s', INTERVAL 1 MONTH)`
	} else {
		panic("指定的报表类型不正确，日报表或者月报表！")
	}

	time_filter = fmt.Sprintf(time_filter, time_begin, time_begin)

	var order_paid_amount_sum struct {
		AmountSum float64 `xorm:"'amount_sum'"`
	}

	_, err := db_session.Sql(sb.SqlGet()).Get(&order_paid_amount_sum)
	core.HandleError(err)

	return order_paid_amount_sum.AmountSum
}
