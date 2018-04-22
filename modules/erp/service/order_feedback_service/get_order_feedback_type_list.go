package order_feedback_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
	"strings"
)

func Get_order_feedback_type_map(order_id_list []int64, db_session *xorm.Session) map[int64][]string {

	if len(order_id_list) == 0 {
		return make(map[int64][]string)
	}

	var order_id_str_list = make([]string, 0)
	for _, order_id := range order_id_list {
		order_id_str_list = append(order_id_str_list, fmt.Sprint(order_id))
	}

	var sql = `
SELECT
    oi.order_id     AS order_id,
    of.comment_type AS comment_type 
FROM erp_order_feedback of
    INNER JOIN erp_order_item oi ON of.order_item_id = oi.id
WHERE of.order_id IN (%s);
`
	sql = fmt.Sprintf(sql, strings.Join(order_id_str_list, ","))

	type order_feedback_type_t struct {
		OrderId     int64 `xorm:"'order_id'"`
		CommentType int64 `xorm:"'comment_type'"`
	}

	var order_feedback_type_item_list = make([]order_feedback_type_t, 0)
	core.HandleError(db_session.Sql(sql).Find(&order_feedback_type_item_list))

	var order_feedback_type_map = make(map[int64][]string)
	for _, order_feedback_type := range order_feedback_type_item_list {
		order_feedback_type_list, ok := order_feedback_type_map[order_feedback_type.OrderId]
		if ok == false {
			order_feedback_type_list = make([]string, 0)
		}

		var order_feedback_type_str = feedback_type_to_type_str(order_feedback_type.CommentType)
		order_feedback_type_list = append(order_feedback_type_list, order_feedback_type_str)
		order_feedback_type_map[order_feedback_type.OrderId] = order_feedback_type_list
	}

	return order_feedback_type_map
}
