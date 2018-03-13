package order_combine_service

import (
	"strconv"
	"strings"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func combine_order(order_ids string, user core.LoginUser, db_session *xorm.Session) {
	var order_id_list = make([]int64, 0)
	var order_id_str_list = strings.Split(order_ids, ",")
	for _, order_id_str := range order_id_str_list {
		order_id, _ := strconv.ParseInt(order_id_str, 10, 64)
		order_id_list = append(order_id_list, order_id)
	}

	var first_order_id = order_id_list[0]

	var combine_order_id = insert_combine_order_from_first_order(first_order_id, order_id_str_list, db_session)

	insert_combine_order_address_from_first_order(combine_order_id, first_order_id, db_session)

	update_combine_order_amount(combine_order_id, order_id_str_list, db_session)

	insert_combine_order_items(combine_order_id, order_id_str_list, db_session)

	cancel_combined_orders(order_id_str_list, db_session)

	insert_order_history(combine_order_id, order_id_list, user, db_session)
}
