package order_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

// 获取的是 AUD 的成本总额
func GetOrderCostAmountSum(order_id int64, db_session *xorm.Session) float64 {
	var cost_amount_sum struct {
		CostAmount float64 `xorm:"'cost_amount'"`
	}

	var sql = `SELECT sum(oi.cost_price * oi.quantity) AS cost_amount
FROM erp_order_item oi
WHERE oi.order_id = %d;`
	sql = fmt.Sprintf(sql, order_id)

	_, err := db_session.Sql(sql).Get(&cost_amount_sum)
	core.HandleError(err)

	return cost_amount_sum.CostAmount
}
