package order_combine_service

import (
	"errors"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

// 合并订单
// order_filter: o.batch_code = '20170428120352'
func Combine_orders(order_filter string, user core.LoginUser, db_session *xorm.Session) error {
	// 不合并订单，先提供个假的实现
	return nil

	if len(order_filter) == 0 {
		return errors.New("未提供订单过滤条件！")
	}

	var order_ids_list = get_order_ids_list_to_combine(order_filter, db_session)

	for _, order_ids := range order_ids_list {
		combine_order_list(order_ids, user, db_session)
	}

	return nil
}
