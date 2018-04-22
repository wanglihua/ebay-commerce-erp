package order_feedback

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func Get_order_feedback_type_list(order_id int64, db_session *xorm.Session) []int64 {
	var sql = `
SELECT of.comment_type AS comment_type 
FROM erp_order_feedback of
    INNER JOIN erp_order_item oi ON of.order_item_id = oi.id
WHERE of.order_id = %d;
`
	sql = fmt.Sprintf(sql, order_id)

	type order_feedback_type_t struct {
		CommentType int64 `xorm:"'comment_type'"`
	}

	var order_feedback_type_item_list = make([]order_feedback_type_t, 0)
	core.HandleError(db_session.Sql(sql).Find(&order_feedback_type_item_list))

	var order_feedback_type_list = make([]int64, 0)
	for _, order_feedback_type := range order_feedback_type_item_list {
		order_feedback_type_list = append(order_feedback_type_list, order_feedback_type.CommentType)
	}

	return order_feedback_type_list
}
