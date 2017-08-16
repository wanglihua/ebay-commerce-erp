package order_item_service

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"fmt"
	"strings"
	"matrix/core"
)

// 获取订单列表中订单明细项是空的订单数目
func Get_order_item_zero_count_in_order_list(order_list []models.OrderInfo, db_session *xorm.Session) int64 {

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
         FROM erp_order_item oi
         WHERE oi.order_id IN (88415, 88416, 88417)
         GROUP BY oi.order_id
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
         FROM erp_order_item oi
         WHERE oi.order_id IN (%s)
         GROUP BY oi.order_id
         HAVING count(oi.id) = 0
     ) order_item_counts;`

	sql = fmt.Sprintf(sql, strings.Join(order_id_str_list, ","))

	_, err := db_session.Sql(sql).Get(&order_item_zero_count)
	core.HandleError(err)

	return order_item_zero_count.Count
}
