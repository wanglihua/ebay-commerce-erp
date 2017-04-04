package order_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/service/logistics_setting_service"
)

// 计算订单的物流方式 和 Hold 类型
func Calculate_order_logistics(db_session *xorm.Session) {
	// 获取澳币的汇率



	logistics_setting := logistics_setting_service.Get_logistics_setting(db_session)
	_ = logistics_setting

	// 拆分和合并订单，在这之前就已经完成了？

	// 检查 订单状态 是不是 正确，状态应该是 已支付 ？？

	// 获取 已支付 状态的订单, 匹配完成后，订单的状态就是 Hold 或者 待配货
	var paid_order_list = Get_paid_order_list(db_session)

	for _, order := range paid_order_list {
		_ = order
	}

	// 仓库的商品库存状态里面是不是也可以加一个待配货数量，实际库存没有减掉，但是已经这个商品被预定配货了，实际配货成功后，就减掉实际配货 ?? -- 通过在库存表里面插入一个特殊库位的库存记录的方法，这个特殊库位的Id为 -1

}
