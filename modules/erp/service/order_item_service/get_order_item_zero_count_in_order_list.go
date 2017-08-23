package order_item_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"strings"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
)

// 获取订单列表中订单明细项是空的订单数目
func Get_order_item_zero_count_in_order_list(order_list []ddlc_models.OrderInfo, db_session *xorm.Session) int64 {

	if len(order_list) == 0 {
		panic("订单列表不能为空！")
	}

	var order_item_zero_count struct {
		Count int64 `xorm:"'order_item_zero_count'"`
	}

	/*
SELECT count(*) AS order_item_zero_count
FROM (
         SELECT count(oi.id) AS order_item_count
         FROM erp_order o LEFT JOIN erp_order_item oi ON o.id = oi.order_id
         WHERE o.id IN (141408, 141407, 141406, 141405)
         GROUP BY o.id
         HAVING count(oi.id) = 0
     ) order_item_counts;
	 */

	var order_id_str_list = make([]string, 0)
	for _, order := range order_list {
		order_id_str_list = append(order_id_str_list, fmt.Sprint(order.Id))
	}

	var sql = `SELECT count(*) AS order_item_zero_count
FROM (
         SELECT count(oi.id) AS order_item_count
         FROM erp_order o LEFT JOIN erp_order_item oi ON o.id = oi.order_id
         WHERE o.id IN (%s)
         GROUP BY o.id
         HAVING count(oi.id) = 0
     ) order_item_counts;`

	sql = fmt.Sprintf(sql, strings.Join(order_id_str_list, ","))

	_, err := db_session.Sql(sql).Get(&order_item_zero_count)
	core.HandleError(err)

	return order_item_zero_count.Count
}
