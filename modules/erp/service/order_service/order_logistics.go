package order_service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/core/collections"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/service/logistics_setting_service"
)

type OrderLogistics struct {
	aud_exchange_rate_setting         float64
	logistics_setting                 models.LogisticsSettingInfo
	address_special_word_list_setting []string
}

func (ol *OrderLogistics) init_data(db_session *xorm.Session) {
	// 获取澳币的汇率，如果定单里面的汇率没有记录下来，就用这个汇率
	var aud_exchange_rate_setting = exchange_rate_service.Get_aud_exchange_rate(db_session)
	if aud_exchange_rate_setting == 0 {
		panic(errors.New("澳币汇率未设置！"))
	}
	ol.aud_exchange_rate_setting = aud_exchange_rate_setting

	// 获取物流参数设置
	var logistics_setting = logistics_setting_service.Get_logistics_setting(db_session)
	ol.logistics_setting = logistics_setting

	// 获得地址特殊字段的列表
	var address_special_word_list_setting = strings.Split(logistics_setting.AddressSpecialWord, ",")
	for index, address_special_word := range address_special_word_list_setting {
		address_special_word_list_setting[index] = strings.TrimSpace(address_special_word)
	}
	ol.address_special_word_list_setting = address_special_word_list_setting
}

type order_logistics_data_t struct {
	exchange_rate           float64
	order_shipping_address  models.OrderShippingAddressInfo
	order_product_info_list []Order_product_info
}

// 计算订单的物流方式 和 Hold 类型
func (ol *OrderLogistics) Calculate(db_session *xorm.Session) {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("匹配物流方式，出现错误： %s\n", err)
		}
	}()

	ol.init_data(db_session)

	// 拆分和合并订单，在这之前就已经完成了？ 不是在这里做，就可以了，做完进入 hold

	// 获取 已支付 状态的订单, 匹配完成后，订单的状态就是 Hold 或者 待配货
	var paid_order_list = get_paid_order_list(db_session)

	for _, order := range paid_order_list {
		var order_logistics_data order_logistics_data_t

		// 下载订单的时候就把当时的汇率记录了下来了
		var exchange_rate = order.ExchangeRate
		if exchange_rate == 0.0 {
			exchange_rate = ol.aud_exchange_rate_setting
		}
		order_logistics_data.exchange_rate = exchange_rate

		// 订单收货地址
		var order_shipping_address = get_order_shipping_address(order.Id, db_session)
		order_logistics_data.order_shipping_address = order_shipping_address

		// 1 地址有特殊字段的
		if ol.hold_check_address_special_word(&order, &order_logistics_data, db_session) {
			continue
		}

		// 2 有买家留言
		if ol.hold_check_buyer_checkout_message(&order, db_session) {
			continue
		}

		// 16 判断是否是国际件，如果不是澳洲的就是国际件，一般就是新西兰的
		if ol.hold_check_international(&order, &order_logistics_data, db_session) {
			continue
		}

		// 15 判断是否是快递
		// order 里面要加个字段，运费收入

		order_logistics_data.order_product_info_list = get_order_product_info_list(order.Id, db_session)

		// 5. 系统显示所有仓都缺货（中国仓除外）
		if ol.hold_check_oss_all_stock(&order, &order_logistics_data, db_session) {
			continue
		}

		// 6. 系统显示最佳仓缺货（次佳、海上有货）
		if ol.hold_check_oss_best_stock(&order, &order_logistics_data, db_session) {
			continue
		}

		// 多个商品的情况

		// 14 判断是否是信件

		/*
		for _, order_product_info := range order_product_info_list {
			inventory_service.Get_product_stock_qty_list(order_product_info.OrderId, db_session)
		}
		*/

		// 检查 订单状态 是不是 正确，状态应该是 已支付 ？？
	}

	// 仓库的商品库存状态里面是不是也可以加一个待配货数量，实际库存没有减掉，但是已经这个商品被预定配货了，实际配货成功后，就减掉实际配货 ?? -- 通过在库存表里面插入一个特殊库位的库存记录的方法，这个特殊库位的Id为 -1

}

// 1 地址有特殊字段的
func (ol *OrderLogistics) hold_check_address_special_word(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	var has_address_special_word = false
	for _, address_special_word := range ol.address_special_word_list_setting {
		// 地址包含特殊字段
		address_special_word_lower_case := strings.ToLower(address_special_word)
		address_detail_lower_case := strings.ToLower(fmt.Sprintf("%s %s", p_order_logistics_data.order_shipping_address.Street1, p_order_logistics_data.order_shipping_address.Street2))
		if strings.Index(address_detail_lower_case, address_special_word_lower_case) != -1 {
			ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_SPECIAL_ADDRESS, db_session)
			has_address_special_word = true
			break
		}
	}

	return has_address_special_word
}

// 2 有买家留言
func (ol *OrderLogistics) hold_check_buyer_checkout_message(p_order *models.OrderInfo, db_session *xorm.Session) (is_hold_return bool) {
	if strings.TrimSpace(p_order.BuyerCheckoutMessage) != "" {
		ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_BUYER_CHECKOUT_MESSAGE, db_session)
		return true
	}

	return false
}

// 5 系统显示所有仓都缺货（中国仓除外）
func (ol *OrderLogistics) hold_check_oss_all_stock(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	// 只检查不是 中国仓 的仓库
	var fun_is_not_cn_stock = func(stock_type_id int64) bool {
		if stock_type_id != models.STOCK_TYPE_CN {
			return true
		} else {
			return false
		}
	}

	if ol.hold_check_oss(p_order, p_order_logistics_data, fun_is_not_cn_stock, db_session) {
		ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_OOS_ALL_STOCK, db_session)
		return true
	}

	return false
}

// 6. 系统显示最佳仓缺货（次佳、海上有货）
func (ol *OrderLogistics) hold_check_oss_best_stock(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	// 检查 澳洲 的仓库
	var fun_is_au_stock = func(stock_type_id int64) bool {
		if stock_type_id == models.STOCK_TYPE_AU {
			return true
		} else {
			return false
		}
	}

	var fun_is_zt_stock = func(stock_type_id int64) bool {
		if stock_type_id == models.STOCK_TYPE_ZT {
			return true
		} else {
			return false
		}
	}

	// 澳洲仓库缺货，但在途仓库有货
	if ol.hold_check_oss(p_order, p_order_logistics_data, fun_is_au_stock, db_session) == true && ol.hold_check_oss(p_order, p_order_logistics_data, fun_is_zt_stock, db_session) == false {
		ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_OOS_BEST_STOCK, db_session)
		return true
	}

	return false
}

//返回 true, 表明考虑的是符合要求的
type fun_stock_type_filter_t func(stock_type_id int64) bool

// 订单商品库存检查，用户给5和6的调用
func (ol *OrderLogistics) hold_check_oss(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, fun_stock_type_filter fun_stock_type_filter_t, db_session *xorm.Session) (is_oss_return bool) {
	//todo: 暂时理解成是，所有商品，所有仓库(中国仓除外)都缺货。这里的缺货也理解成是，没有一个仓库能凑齐所以商品的所有数量

	// 获取订单中的所有商品的id列表
	var get_order_product_id_list = func() []int64 {
		var product_id_list = make([]int64, 0)
		for _, order_product_info := range p_order_logistics_data.order_product_info_list {
			product_id_list = append(product_id_list, order_product_info.ProductId)
		}
		return product_id_list
	}
	var order_product_id_list = get_order_product_id_list()

	var product_stock_qty_list = inventory_service.Get_products_stock_qty_list(order_product_id_list, db_session)

	// 找出所有包含这些商品库存的仓库id列表
	var get_product_stock_id_list = func() []int64 {
		var stock_id_set = collections.NewIntSet()
		for _, product_stock_qty := range product_stock_qty_list {
			stock_id_set.Add(product_stock_qty.StockId)
		}
		return stock_id_set.List()
	}
	var product_stock_id_list = get_product_stock_id_list()

	var all_product_oss = true //下面的逻辑中，如果发现仓库有货，就设置为 false

	// 对每一个仓库循环，看是不是这个仓库里面, 这些商品都有货
	for _, stock_id := range product_stock_id_list {
		var stock_has_all_product = true

		// 对订单中的每个商品信息循环，看是不是这个商品在这个仓库里面有货且有足够的数量
		for _, order_product_info := range p_order_logistics_data.order_product_info_list {
			var product_has_stock = false
			for _, product_stock_qty := range product_stock_qty_list {
				if product_stock_qty.StockId == stock_id && product_stock_qty.ProductId == order_product_info.ProductId && fun_stock_type_filter(product_stock_qty.StockTypeId) {
					if product_stock_qty.Qty >= order_product_info.Qty {
						product_has_stock = true
						break //这个 break 加不加都可以的，因为 这个 stock_id， 这个 product_id 的就一条记录
					}
				}
			}
			// 如果商品在这个仓库里面没有货，就判定 这个仓库里面，并不是所有商品都有货
			if product_has_stock == false {
				stock_has_all_product = false
				break // 订单里面的 接下来的 商品也就不用检查了
			}
		}

		// 找到仓库全部商品都有货了
		if stock_has_all_product {
			all_product_oss = false
			break
		} else {
			// 如果这个仓库里面不是全部商品都有货，就继续找
			continue // 因为是循环的最后，所以这个continue可以加也可以不加
		}
	}

	if all_product_oss {
		return true
	}

	return false
}

// 16 判断是否是国际件，如果不是澳洲的就是国际件，一般就是新西兰的
func (ol *OrderLogistics) hold_check_international(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	if strings.ToLower(p_order_logistics_data.order_shipping_address.Country) == strings.ToLower("AU") {
		ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_INTERNATIONAL, db_session)
		return true
	}

	return false
}

func (ol *OrderLogistics) update_order_to_hold_type(p_order *models.OrderInfo, hold_type_id int64, db_session *xorm.Session) {
	p_order.HoldTypeId = core.MakeNullInt(hold_type_id, true)
	_, err := db_session.Id(p_order.Id).Cols("hold_type_id").Update(p_order)
	core.HandleError(err)
}

//todo: 同一个商品的一半数量可以从这个仓库发货，另一半数量可以从哪个仓库发货的怎么处理
/*

(完成)1. 地址有特殊字段的。比如 PO BOX, LOCKER。但最佳物流非AU POST的，需要在基础设置里面有特殊字段设置功能
(完成)2. 有买家留言 -- 导入的订单里面还没有留言
3. 有平台留言 -- 导入的订单里面还没有留言
4. PICK UP   -- 一般都是手工创建的订单
(完成)5. 系统显示所有仓都缺货（中国仓除外）
6. 系统显示最佳仓缺货（次佳、海上有货）
7. 实际仓库缺货（是在捡货过程中，虽然系统显示有货，但实际仓库没有货了。就需要人工手动将这些订单搬至这里的模块，等待跟客户协商，是退款还是等待集装箱）
8. 多数量或多产品，但部分缺货（是否拆分订单）
9. 多数量或多产品，虽不缺货，但拆分订单发货，运费更便宜
10. 多数量或多产品，虽不缺货，但拆分订单发货，运费虽然便宜，但仓库拣货更方便
11. 无法匹配物流（长宽高重量没有物流满足）
12. 产品信息不全（长宽高重量不全或者有误）
13. 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损hold
14. 信封类。长宽高在35cm X 25cm X 2cm, 重量在0.5Kg以内的订单为信封。若是多数量或者多产品，以前我们设计的思路是，最大产品的长边为最长边，最大产品的中边为中边，所有产品的最小边之和为最小边，判断是否在35cm X 25cm X 2cm以内，重量之和是否在0.5Kg以内
15. 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包裹
(完成)16. 国际件。对于非澳大利亚的订单，直接到国际件模块

 */
