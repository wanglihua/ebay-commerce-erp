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
		// 下载订单的时候就把当时的汇率记录了下来了
		var exchange_rate = order.ExchangeRate
		if exchange_rate == 0.0 {
			exchange_rate = aud_exchange_rate
		}

		// 1 地址有特殊字段的

		// 16 判断是否是国际件



		// 15 判断是否是快递



		var order_product_info_list = get_order_product_info_list(order.Id, db_session)

		// 14 判断是否是信件

		for _, order_product_info := range order_product_info_list {
			product_inventory_service.Get_product_stock_qty_list(order_product_info.OrderId, db_session)
		}

		// 检查 订单状态 是不是 正确，状态应该是 已支付 ？？
	}

	// 仓库的商品库存状态里面是不是也可以加一个待配货数量，实际库存没有减掉，但是已经这个商品被预定配货了，实际配货成功后，就减掉实际配货 ?? -- 通过在库存表里面插入一个特殊库位的库存记录的方法，这个特殊库位的Id为 -1

}


/*

1. 地址有特殊字段的。比如 PO BOX, LOCKER。但最佳物流非AU POST的，需要在基础设置里面有特殊字段设置功能
2. 有买家留言 -- 导入的订单里面还没有留言
3. 有平台留言 -- 导入的订单里面还没有留言
4. PICK UP   -- 一般都是手工创建的订单
5. 系统显示所有仓都缺货（中国仓除外）
6. 系统显示最佳仓缺货（次佳、海上有货）
7. 实际仓库缺货
8. 多数量或多产品，但部分缺货（是否拆分订单）
9. 多数量或多产品，虽不缺货，但拆分订单发货，运费更便宜
10. 多数量或多产品，虽不缺货，但拆分订单发货，运费虽然便宜，但仓库拣货更方便
11. 无法匹配物流（长宽高重量没有物流满足）
12. 产品信息不全（长宽高重量不全或者有误）
13. 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损hold
14. 信封类。长宽高在35cm X 25cm X 2cm, 重量在0.5Kg以内的订单为信封。若是多数量或者多产品，以前我们设计的思路是，最大产品的长边为最长边，最大产品的中边为中边，所有产品的最小边之和为最小边，判断是否在35cm X 25cm X 2cm以内，重量之和是否在0.5Kg以内
15. 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包裹
16. 国际件。对于非澳大利亚的订单，直接到国际件模块

 */