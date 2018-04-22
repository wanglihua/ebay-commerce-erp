package order_feedback_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func Get_order_feedback_list(order_id int64, db_session *xorm.Session) []OrderFeedbackViewItem {
	var sql = `
SELECT
    of.*,
    oi.*
FROM erp_order_feedback of
    INNER JOIN erp_order_item oi ON of.order_item_id = oi.id
WHERE of.order_id = %d;
`
	sql = fmt.Sprintf(sql, order_id)

	var order_feedback_list = make([]OrderFeedbackViewItem, 0)
	core.HandleError(db_session.Sql(sql).Find(&order_feedback_list))

	return order_feedback_list
}
