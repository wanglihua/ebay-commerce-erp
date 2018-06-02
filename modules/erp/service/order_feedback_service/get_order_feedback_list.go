package order_feedback_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
	"time"
)

func Get_order_feedback_list(order_id int64, db_session *xorm.Session) []OrderFeedbackViewItem {
	var sql = `
SELECT
    f.*,
    oi.*
FROM erp_order_feedback f
    INNER JOIN erp_order_item oi ON f.order_item_id = oi.id
WHERE f.order_id = %d;
`
	sql = fmt.Sprintf(sql, order_id)

	var order_feedback_list = make([]OrderFeedbackViewItem, 0)
	core.HandleError(db_session.Sql(sql).Find(&order_feedback_list))

	for index, order_feedback := range order_feedback_list {
		order_feedback_list[index].OrderFeedback.CommentTypeText = feedback_type_to_type_str(order_feedback.OrderFeedback.CommentType)

		order_feedback_list[index].OrderFeedback.CommentTimeText = time.Time(order_feedback.OrderFeedback.CommentTime).Format("2006-01-02 15:04:05")
	}

	return order_feedback_list
}
