package order_service

import (
	"errors"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/logistics_setting_service"
	"matrix/modules/erp/service/product_inventory_service"
)

// 计算订单的物流方式 和 Hold 类型
func Calculate_order_logistics(db_session *xorm.Session) {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("匹配物流方式，出现错误： %s\n", err)
		}
	}()

	// 获取澳币的汇率，如果定单里面的汇率没有记录下来，就用这个汇率
	var aud_exchange_rate = exchange_rate_service.Get_aud_exchange_rate(db_session)
	if aud_exchange_rate == 0 {
		panic(errors.New("澳币汇率未设置！"))
	}

	logistics_setting := logistics_setting_service.Get_logistics_setting(db_session)
	_ = logistics_setting

	// 拆分和合并订单，在这之前就已经完成了？ 不是在这里做，就可以了，做完进入 hold

	// 获取 已支付 状态的订单, 匹配完成后，订单的状态就是 Hold 或者 待配货
	var paid_order_list = get_paid_order_list(db_session)

	for _, order := range paid_order_list {

		// 下载订单的时候就可把当时的汇率记录了下来
		var exchange_rate = order.ExchangeRate
		_ = exchange_rate

		var order_product_info_list = get_order_product_info_list(order.Id, db_session)

		for _, order_product_info := range order_product_info_list {
			product_inventory_service.Get_product_stock_qty_list(order_product_info.OrderId, db_session)
		}

		// 检查 订单状态 是不是 正确，状态应该是 已支付 ？？
	}

	// 仓库的商品库存状态里面是不是也可以加一个待配货数量，实际库存没有减掉，但是已经这个商品被预定配货了，实际配货成功后，就减掉实际配货 ?? -- 通过在库存表里面插入一个特殊库位的库存记录的方法，这个特殊库位的Id为 -1

}
