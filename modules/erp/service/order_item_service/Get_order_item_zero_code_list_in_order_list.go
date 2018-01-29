package order_item_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"strings"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
)

// 获取订单列表中订单明细项是空的订单列表
func Get_order_item_zero_code_list_in_order_list(order_list []ddlc_models.OrderInfo, db_session *xorm.Session) []string {

	if len(order_list) == 0 {
		panic("订单列表不能为空！")
	}

	type order_item_zero_order_t struct {
		OrderId        int64     `xorm:"'order_id'"`
		OrderCode      string    `xorm:"'order_code'"`
		OrderCeateTime core.Time `xorm:"'order_create_time'"`
		ItemCount      int64     `xorm:"'item_count'"`
		OrderStatus    string    `xorm:"'order_status'"`
	}

	/*
SELECT
    o.id          AS order_id,
    o.code        AS order_code,
    o.create_time AS order_create_time,
    count(oi.id)  AS item_count,
    s.name        AS order_status
FROM erp_order o
    INNER JOIN erp_status s ON o.status_id = s.id
    LEFT JOIN erp_order_item oi ON o.id = oi.order_id
WHERE o.id IN (301527, 301530)
GROUP BY o.id
HAVING count(oi.id) = 0
	 */

	var order_id_str_list = make([]string, 0)
	for _, order := range order_list {
		order_id_str_list = append(order_id_str_list, fmt.Sprint(order.Id))
	}

	var sql = `SELECT
    o.id          AS order_id,
    o.code        AS order_code,
    o.create_time AS order_create_time,
    count(oi.id)  AS item_count,
    s.name        AS order_status
FROM erp_order o
    INNER JOIN erp_status s ON o.status_id = s.id
    LEFT JOIN erp_order_item oi ON o.id = oi.order_id
WHERE o.id IN (%s)
GROUP BY o.id
HAVING count(oi.id) = 0`

	sql = fmt.Sprintf(sql, strings.Join(order_id_str_list, ","))

	var order_item_zero_order_list = make([]order_item_zero_order_t, 0)
	err := db_session.Sql(sql).Find(&order_item_zero_order_list)
	core.HandleError(err)

	var order_code_list = make([]string, 0)
	for _, order_item_zero_order := range order_item_zero_order_list {
		order_code_list = append(order_code_list, order_item_zero_order.OrderCode)
	}

	return order_code_list
}
