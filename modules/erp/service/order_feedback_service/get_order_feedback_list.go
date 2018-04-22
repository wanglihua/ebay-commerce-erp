package order_feedback_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
	"strings"
	"time"
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

	for _, order_feedback := range order_feedback_list {
		//  // 1: negative, 2: neutral, 3: positive
		if order_feedback.OrderFeedback.CommentType == 1 {
			order_feedback.OrderFeedback.CommentTypeText = strings.ToLower("negative")
		} else if order_feedback.OrderFeedback.CommentType == 2 {
			order_feedback.OrderFeedback.CommentTypeText = strings.ToLower("neutral")
		} else if order_feedback.OrderFeedback.CommentType == 3 {
			order_feedback.OrderFeedback.CommentTypeText = strings.ToLower("positive")
		}

		// time.Now().Format("2006-01-02 15:04:05")
		order_feedback.OrderFeedback.CommentTimeText = time.Time(order_feedback.OrderFeedback.CommentTime).Format("2006-01-02 15:04:05")
	}

	return order_feedback_list
}
