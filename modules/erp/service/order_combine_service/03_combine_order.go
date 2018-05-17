package order_combine_service

import (
	"strconv"
	"strings"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models"
)

func combine_order(order_ids string, user core.LoginUser, db_session *xorm.Session) int64 {
	var order_id_list = make([]int64, 0)
	var order_id_str_list = strings.Split(order_ids, ",")
	for _, order_id_str := range order_id_str_list {
		order_id, _ := strconv.ParseInt(order_id_str, 10, 64)
		order_id_list = append(order_id_list, order_id)
	}

	for _, order_id := range order_id_list {
		var order ddlc_models.OrderInfo

		// 合并前再次检查订单是否符合合并要求，可以防止重复合并的的问题
		has, err := db_session.
			Where("(id = ?) AND (type_id <> ? AND type_id <> ? AND type_id <> ?) AND (status_id = ?) AND status_id <> ?", order_id, ddpz_models.ORDER_TYPE_BHDD, ddpz_models.ORDER_TYPE_HBDD, ddpz_models.ORDER_TYPE_CFDD, models.ORDER_STATUS_YZF, models.ORDER_STATUS_YQX).
			Get(&order)
		core.HandleError(err)

		if has == false {
			return 0
		}
	}

	var first_order_id = order_id_list[0]

	var combine_order_id = insert_combine_order_from_first_order(first_order_id, order_id_str_list, db_session)

	insert_combine_order_address_from_first_order(combine_order_id, first_order_id, db_session)

	update_combine_order_amount(combine_order_id, order_id_str_list, db_session)

	insert_combine_order_items(combine_order_id, order_id_str_list, db_session)

	cancel_combined_orders(order_id_str_list, db_session)

	insert_order_history(combine_order_id, order_id_list, user, db_session)

	return combine_order_id
}
