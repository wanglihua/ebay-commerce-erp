package order_combine_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
	"strings"
)

func update_combine_order_amount(combine_order_id int64, order_id_str_list []string, db_session *xorm.Session) {

	// 修改 合并后的 order 各种款项
	/*
UPDATE erp_order combine_order
    INNER JOIN (
                   SELECT
                       sum(order_all.adjustment_amount) AS adjustment_amount_sum,
                       sum(order_all.amount_paid)       AS amount_paid_sum,
                       sum(order_all.amount_saved)      AS amount_saved_sum,
                       sum(order_all.total_amount)      AS total_amount_sum,
                       sum(order_all.subtotal_amount)   AS subtotal_amount_sum,
                       sum(order_all.express_fee)       AS express_fee_sum
                   FROM erp_order order_all
                   WHERE order_all.id IN (0, -1, -2, -3)
               ) s
        ON 1 = 1
SET combine_order.adjustment_amount     = s.adjustment_amount_sum,
    combine_order.adjustment_amount_rmb = s.adjustment_amount_sum * combine_order.exchange_rate,
    combine_order.amount_paid           = s.amount_paid_sum,
    combine_order.amount_paid_rmb       = s.amount_paid_sum * combine_order.exchange_rate,
    combine_order.amount_saved          = s.amount_saved_sum,
    combine_order.amount_saved_rmb      = s.amount_saved_sum * combine_order.exchange_rate,
    combine_order.total_amount          = s.total_amount_sum,
    combine_order.total_amount_rmb      = s.total_amount_sum * combine_order.exchange_rate,
    combine_order.subtotal_amount       = s.subtotal_amount_sum,
    combine_order.subtotal_amount_rmb   = s.subtotal_amount_sum * combine_order.exchange_rate,
    combine_order.express_fee           = s.express_fee_sum
WHERE id = 0
	 */

	// 假设各订单的币种是一样的
	var sql_sum_amount = `
UPDATE erp_order combine_order
    INNER JOIN (
                   SELECT
                       sum(order_all.adjustment_amount) AS adjustment_amount_sum,
                       sum(order_all.amount_paid)       AS amount_paid_sum,
                       sum(order_all.amount_saved)      AS amount_saved_sum,
                       sum(order_all.total_amount)      AS total_amount_sum,
                       sum(order_all.subtotal_amount)   AS subtotal_amount_sum,
                       sum(order_all.express_fee)       AS express_fee_sum
                   FROM erp_order order_all
                   WHERE order_all.id IN (%s)
               ) s
        ON 1 = 1
SET combine_order.adjustment_amount     = s.adjustment_amount_sum,
    combine_order.adjustment_amount_rmb = s.adjustment_amount_sum * combine_order.exchange_rate,
    combine_order.amount_paid           = s.amount_paid_sum,
    combine_order.amount_paid_rmb       = s.amount_paid_sum * combine_order.exchange_rate,
    combine_order.amount_saved          = s.amount_saved_sum,
    combine_order.amount_saved_rmb      = s.amount_saved_sum * combine_order.exchange_rate,
    combine_order.total_amount          = s.total_amount_sum,
    combine_order.total_amount_rmb      = s.total_amount_sum * combine_order.exchange_rate,
    combine_order.subtotal_amount       = s.subtotal_amount_sum,
    combine_order.subtotal_amount_rmb   = s.subtotal_amount_sum * combine_order.exchange_rate,
    combine_order.express_fee           = s.express_fee_sum
WHERE id = %d`
	sql_sum_amount = fmt.Sprintf(sql_sum_amount, strings.Join(order_id_str_list, ","), combine_order_id)
	_, err := db_session.Exec(sql_sum_amount)
	core.HandleError(err)
}
