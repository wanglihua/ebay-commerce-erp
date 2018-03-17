package order_combine_service

import (
	"errors"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

// 合并订单
// order_filter: o.batch_code = '20170428120352'
func Combine_orders(order_filter string, user core.LoginUser, db_session *xorm.Session) []int64 {
	// 不合并订单，就可以先返回
	var temp = make([]int64, 0)
	return temp

	if len(order_filter) == 0 {
		panic(errors.New("未提供订单过滤条件！"))
	}

	var order_ids_list = get_order_ids_list_to_combine(order_filter, db_session)

	var combine_order_id_list = make([]int64, 0)
	for _, order_ids := range order_ids_list {
		var combine_order_id = combine_order(order_ids, user, db_session)
		combine_order_id_list = append(combine_order_id_list, combine_order_id)
	}

	return combine_order_id_list
}
