package service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func Get_user_item_order_item_list(buyer_user_id string, item_id string, db_session *xorm.Session) []MessageOrderItem {
	/*
SELECT
    oi.*,
    o.*,
    s.*
FROM erp_order_item oi
    INNER JOIN erp_order o ON oi.order_id = o.id
    INNER JOIN erp_status s ON o.status_id = s.id
WHERE o.buyer_user_id = '123'
      AND oi.item_id = '123'
ORDER BY o.create_time DESC;
	 */

	var sql = `SELECT
    oi.*,
    o.*,
    s.*
FROM erp_order_item oi
    INNER JOIN erp_order o ON oi.order_id = o.id
    INNER JOIN erp_status s ON o.status_id = s.id
WHERE o.buyer_user_id = '%s'
      AND oi.item_id = '%s'
ORDER BY o.create_time DESC;`
	sql = fmt.Sprintf(sql, buyer_user_id, item_id)

	var message_order_item_list = make([]MessageOrderItem, 0)
	err := db_session.Sql(sql).Find(&message_order_item_list)
	core.HandleError(err)

	calculate_order_item_list_extra_field(message_order_item_list)

	return message_order_item_list
}
