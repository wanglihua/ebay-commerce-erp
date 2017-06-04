package order_logistics_service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/core/collections"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
)

// 1 地址有特殊字段的
func (p_ol *OrderLogistics) hold_check_address_special_word(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, order_logistics models.OrderLogisticsInfo, db_session *xorm.Session) (is_hold_return bool) {
	// 检查地址中是否有特殊字段
	var has_address_special_word = false
	for _, address_special_word := range p_ol.address_special_word_list_setting {
		// 地址包含特殊字段
		address_special_word_lower_case := strings.ToLower(address_special_word)
		address_detail_lower_case := strings.ToLower(fmt.Sprintf("%s %s", p_order_logistics_data.order_shipping_address.Street1, p_order_logistics_data.order_shipping_address.Street2))
		if strings.Index(address_detail_lower_case, address_special_word_lower_case) != -1 {
			has_address_special_word = true
			break
		}
	}
	if has_address_special_word == true {
		var logistics_company_list = p_ol.logistics_company_list
		for _, logistics_company := range logistics_company_list {
			if order_logistics.CompanyId == logistics_company.Id && logistics_company.AddressSpecialWord == 0 {
				p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_SPECIAL_ADDRESS, db_session)
				return true
			}
		}
	}

	// 不符合 地址有特殊字段的，但最佳物流 “ 非“ AU POST的
	return false
}

// 2 有买家留言
func (p_ol *OrderLogistics) hold_check_buyer_checkout_message(p_order *models.OrderInfo, db_session *xorm.Session) (is_hold_return bool) {
	if strings.TrimSpace(p_order.BuyerCheckoutMessage) != "" {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_BUYER_CHECKOUT_MESSAGE, db_session)
		return true
	}

	return false
}

// todo: 3 有平台留言 -- 导入的订单里面还没有留言

// todo: 4 PICK UP  -- 一般都是手工创建的订单，做一个 PICK UP的订单类型？

// 5 系统显示所有仓都缺货（中国仓除外）
func (p_ol *OrderLogistics) hold_check_oss_all_stock(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	// 只检查不是 中国仓 的仓库
	var fun_is_not_cn_stock = func(stock_type_id int64) bool {
		if stock_type_id != models.STOCK_TYPE_CN {
			return true
		} else {
			return false
		}
	}

	if p_ol.hold_check_oss_by_stock_type(p_order, p_order_logistics_data, fun_is_not_cn_stock, db_session) {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_OOS_ALL_STOCK, db_session)
		return true
	}

	return false
}

// 6. 系统显示最佳仓缺货（次佳、海上有货）
// 这个在需求上有更改，不用这个了
// 后来改成在途才有货

func (p_ol *OrderLogistics) hold_check_oss_au_stock_but_zt_stock_has(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	// 检查 澳洲 的仓库
	// todo: 这里还有点问题，应该是检查对应的在途仓库？,先不考虑这个问题了

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
	if p_ol.hold_check_oss_by_stock_type(p_order, p_order_logistics_data, fun_is_au_stock, db_session) == true && p_ol.hold_check_oss_by_stock_type(p_order, p_order_logistics_data, fun_is_zt_stock, db_session) == false {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_OOS_AU_STOCK_ZT_HAS, db_session)
		return true
	}

	return false
}

//返回 true, 表明考虑的是符合要求的
type fun_stock_type_filter_t func(stock_type_id int64) bool

// 订单商品库存检查，用于给 5 和 6 的调用
func (p_ol *OrderLogistics) hold_check_oss_by_stock_type(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, fun_stock_type_filter fun_stock_type_filter_t, db_session *xorm.Session) (is_oss_return bool) {
	var products_stock_qty_list = p_order_logistics_data.products_stock_qty_list

	// 找出所有包含这些商品库存的仓库id列表
	var get_product_stock_id_list = func() []int64 {
		var stock_id_set = collections.NewIntSet()
		for _, product_stock_qty := range products_stock_qty_list {
			if fun_stock_type_filter(product_stock_qty.StockTypeId) { // 只取符合要求的仓库类型。比如，匹配订单的时候，只取澳洲的仓库什么的
				stock_id_set.Add(product_stock_qty.StockId)
			}
		}
		return stock_id_set.List()
	}

	var product_stock_id_list = get_product_stock_id_list()

	var all_product_oss = true //下面的逻辑中，如果发现仓库有货，就设置为 false

	// 对每一个仓库循环，看是不是这个仓库里面, 这些商品都有货
	for _, stock_id := range product_stock_id_list {
		var stock_has_all_product, _ = check_stock_has_order_products(p_ol, stock_id, p_order_logistics_data)

		// 找到仓库全部商品都有货了，就退出；如果这个仓库里面不是全部商品都有货，就继续找
		if stock_has_all_product {
			all_product_oss = false
			break
		}
	}

	return all_product_oss
}

func (p_ol *OrderLogistics) hold_check_has_some_product(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	var products_stock_qty_list = p_order_logistics_data.products_stock_qty_list

	// 找出所有包含这些商品库存的仓库id列表
	var get_product_stock_id_list = func() []int64 {
		var stock_id_set = collections.NewIntSet()
		for _, product_stock_qty := range products_stock_qty_list {
			if product_stock_qty.StockTypeId == models.STOCK_TYPE_AU { // 只取符合要求的仓库类型。比如，匹配订单的时候，只取澳洲的仓库什么的
				stock_id_set.Add(product_stock_qty.StockId)
			}
		}
		return stock_id_set.List()
	}

	var product_stock_id_list = get_product_stock_id_list()

	var has_some = false
	// 对每一个仓库循环，看是不是这个仓库里面, 这些商品都有货
	for _, stock_id := range product_stock_id_list {
		var stock_has_some_product = check_stock_has_some_order_products(p_ol, stock_id, p_order_logistics_data)

		// 找到仓库全部商品 部分有货，就退出
		if stock_has_some_product {
			has_some = true
			break
		}
	}

	if has_some == true {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_OOS_SOME_PRODUCT, db_session)
		return true
	}

	return false
}

// todo: 7 实际仓库缺货（是在捡货过程中，虽然系统显示有货，但实际仓库没有货了。就需要人工手动将这些订单搬至这里的模块，等待跟客户协商，是退款还是等待集装箱）, 这个 hold 由拣货过程中产生

// todo: 8, 9, 10 临时的解决方案
func (p_ol *OrderLogistics) hold_check_multi_product(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	if len(p_order_logistics_data.order_product_info_list) > 1 {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_MULTI_PRODUCT, db_session)
		return true
	}

	return false
}

// 11 无法匹配物流（长宽高重量没有物流满足）， 对这个 函数 的调用放在匹配物流结束之后，这样才知道是不是匹配成功了
func (p_ol *OrderLogistics) hold_check_no_match_logistics_mode(p_order *models.OrderInfo, match_success bool, db_session *xorm.Session) (is_hold_return bool) {
	if match_success == false {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH, db_session)
		return true
	}

	return false
}

// 12 产品信息不全（长宽高重量不全或者有误）
func (p_ol *OrderLogistics) hold_check_product_info_miss(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	var order_product_info_list = p_order_logistics_data.order_product_info_list
	for _, order_product_info := range order_product_info_list {
		// 长 宽 高 重量 数据不全
		if order_product_info.Length == 0 || order_product_info.Width == 0 || order_product_info.Height == 0 || order_product_info.Weight == 0 {
			p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_PRODUCT_INFO_MISS, db_session)
			return true
		}
	}

	return false
}

// 13 // 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损hold
func (p_ol *OrderLogistics) hold_check_great_deficit(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, order_logistics models.OrderLogisticsInfo, db_session *xorm.Session) (is_hold_return bool) {
	var total_amount = p_order.TotalAmount                                      //收款金额，澳币
	var express_fee = p_order.ExpressFee                                        // 运费收入，澳币
	var order_gst float64 = order_service.Get_order_gst_fee(p_order.AmountPaid) // GST消费税，澳币
	var platform_fee = p_order.TotalAmount * p_ol.logistics_setting.Ptfl        // 平台交易费， 澳币
	var payment_fee = p_order.TotalAmount * p_ol.logistics_setting.Jyfl         // 交易费率， 澳币
	var logistics_fee = order_logistics.TotalFee                                // 物流费用，不包括拣货费用，澳币
	var pick_fee = order_logistics.PickFee                                      // 拣货费用，澳币
	var cost_amount = order_logistics.CostAmount                                // 到仓总成本，澳币

	var profit = (total_amount + express_fee) - (order_gst + platform_fee + payment_fee + logistics_fee + pick_fee + cost_amount)

	// 利润小于零，并且亏损价格大于最大价格或者亏损率大于最大亏损率
	if profit < 0 && ((math.Abs(profit) > p_ol.logistics_setting.MaxDeficitPrice) || (math.Abs(profit)/profit > p_ol.logistics_setting.MaxDeficitRate)) {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_GREAT_DEFICIT, db_session)
		return true
	}

	return false
}

// 14 信封类。长宽高在35cm X 25cm X 2cm, 重量在0.5Kg以内的订单为信封。若是多数量或者多产品，以前我们设计的思路是，最大产品的长边为最长边，最大产品的中边为中边，所有产品的最小边之和为最小边，判断是否在35cm X 25cm X 2cm以内，重量之和是否在0.5Kg以内
func (p_ol *OrderLogistics) hold_check_envelope(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	// 这里假定信封的长宽高限制： 长 > 宽 > 高，不能这么假定

	revel.TRACE.Println("hold_check_envelope begin")

	var envelope_max_length = p_ol.logistics_setting.EnvelopeMaxLength
	var envelope_max_width = p_ol.logistics_setting.EnvelopeMaxWidth
	var envelope_max_height = p_ol.logistics_setting.EnvelopeMaxHeight
	var envelope_max_weight = p_ol.logistics_setting.EnvelopeMaxWeight

	// 如果没有设置信封的最大长度、宽度、高度、重量
	if envelope_max_length == 0 || envelope_max_width == 0 || envelope_max_height == 0 || envelope_max_weight == 0 {
		revel.ERROR.Printf("系统的对信封的长宽高重量限制的设定不完整！") // 在匹配开始时就提示错误
		return false
	}

	// 进行长宽高的排序，结果中 长 > 宽 > 高，这样在信封长宽高限制的设定中，就不用一定要保证 长 > 宽 > 高
	var envelope_side_length_slice = sort.Float64Slice([]float64{envelope_max_length, envelope_max_width, envelope_max_height})
	envelope_side_length_slice.Sort()
	envelope_max_length = envelope_side_length_slice[2] // 最长
	envelope_max_width = envelope_side_length_slice[1]
	envelope_max_height = envelope_side_length_slice[0] // 最短

	// get_order_products_whole_logistics_data的输出中 长 > 宽 > 高
	var whole_length, whole_width, whole_height, whole_weight float64 = get_order_products_whole_logistics_data(p_order_logistics_data)

	revel.TRACE.Println("whole_length, whole_width, whole_height, whole_weight", whole_length, whole_width, whole_height, whole_weight)
	revel.TRACE.Println("envelope_max_length, envelope_max_width, envelope_max_height, envelope_max_weight", envelope_max_length, envelope_max_width, envelope_max_height, envelope_max_weight)

	if whole_length <= envelope_max_length && whole_width <= envelope_max_width && whole_height <= envelope_max_height && whole_weight <= envelope_max_weight {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_ENVELOPE, db_session)
		// 匹配具体信封费率
		var final_fee_rate_id int64 = 0
		var final_fee float64 = math.MaxFloat64
		var envelope_fee_rate_list = p_ol.envelope_fee_rate_list // 已经按费用从小到大排序好了
		for _, envelope_fee_rate := range envelope_fee_rate_list {
			var fee_rate_length = envelope_fee_rate.Length
			var fee_rate_width = envelope_fee_rate.Width
			var fee_rate_height = envelope_fee_rate.Height
			var fee_rate_weight = envelope_fee_rate.Weight

			var envelope_side_length_slice = sort.Float64Slice([]float64{fee_rate_length, fee_rate_width, fee_rate_height})
			envelope_side_length_slice.Sort()
			fee_rate_length = envelope_side_length_slice[2] // 最长
			fee_rate_width = envelope_side_length_slice[1]
			fee_rate_height = envelope_side_length_slice[0] // 最短

			if whole_length < fee_rate_length && whole_width < fee_rate_width && whole_height < fee_rate_height && whole_weight < fee_rate_weight {
				final_fee = envelope_fee_rate.Fee
				final_fee_rate_id = envelope_fee_rate.Id
				break
			}
		}

		if final_fee_rate_id != 0 {
			var order_logistics models.OrderLogisticsInfo
			order_logistics.Length = whole_length
			order_logistics.Width = whole_width
			order_logistics.Height = whole_height
			order_logistics.Weight = whole_weight
			order_logistics.EnvelopeFeeId = final_fee_rate_id
			order_logistics.TotalFee = final_fee

			_, err := db_session.Insert(&order_logistics)
			core.HandleError(err)
		}

		return true
	}

	return false
}

// 15 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包裹
func (p_ol *OrderLogistics) hold_check_express(p_order *models.OrderInfo, db_session *xorm.Session) (is_hold_return bool) {
	// 系统没有设置快递的最小运费
	if p_ol.logistics_setting.MinExpressFee == 0 {
		revel.ERROR.Println("系统没有设置快递的最小运费！")
		return false
	}

	// 这里假定都是澳币
	if p_order.ExpressFee >= p_ol.logistics_setting.MinExpressFee {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_EXPRESS, db_session)
		return true
	}

	return false
}

// 16 判断是否是国际件，如果不是澳洲的就是国际件，一般就是新西兰的
func (p_ol *OrderLogistics) hold_check_international(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	var country = p_order_logistics_data.order_shipping_address.Country
	country = strings.TrimSpace(country)
	var country_name = p_order_logistics_data.order_shipping_address.CountryName
	country_name = strings.TrimSpace(country_name)

	if (country != "" && strings.ToLower(country) != strings.ToLower("AU")) ||
		(country_name != "" && strings.ToLower(country_name) != strings.ToLower("Australia")) { // != "" 是因为有时导入的时候没提供国家信息的，就默认是国内件了
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_INTERNATIONAL, db_session)
		return true
	}

	return false
}

// 17. 收货人信息不全的，比如没有手机号，没有邮编等等
func (p_ol *OrderLogistics) hold_check_address_info_miss(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	var shipping_address = p_order_logistics_data.order_shipping_address
	if shipping_address.Id == 0 || // 数据库中地址信息不存在
	// strings.TrimSpace(shipping_address.Phone) == "" || // 地址信息中缺少电话号码，后来说电话号码缺失，不算收货人信息不全
		strings.TrimSpace(shipping_address.PostalCode) == "" || // 地址信息中缺少邮政编码
		strings.TrimSpace(shipping_address.Name) == "" || // 地址信息中缺少收件人
		(strings.TrimSpace(shipping_address.Street1) == "" && strings.TrimSpace(shipping_address.Street2) == "" || // 地址信息中缺少详细地址信息
			utf8.RuneCountInString(strings.TrimSpace(shipping_address.Street1)+strings.TrimSpace(shipping_address.Street2)) < 7  ) { // 详细地址信息少于 ７ 个的时候进入 hold， 由客服介入
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_ADDRESS_INFO_MISS, db_session)
		return true
	}

	return false
}

// 18. 对老鼠仓的检查
func (p_ol *OrderLogistics) hold_check_paypal_account_not_match(p_order *models.OrderInfo, db_session *xorm.Session) (is_hold_return bool) {
	var payment_method = p_order.PaymentMethods

	var find_shop_in_shop_list = func(shop_id int64, shop_list []models.ShopAccountInfo) models.ShopAccountInfo {
		for _, shop := range shop_list {
			if shop.Id == shop_id {
				return shop
			}
		}

		return models.ShopAccountInfo{}
	}

	var shop = find_shop_in_shop_list(p_order.ShopId, p_ol.shop_list)

	// 老鼠仓的判断
	// 如果是通过PayPal支付的，而且店铺设置里面也设置了PayPal列表的
	// 由于导入的订单对 payment_method 没有设置值的，所以这里的老鼠仓的判断对导入订单是无效的
	if strings.ToLower(payment_method) == strings.ToLower("PayPal") && shop.PaypalAccount != "" {
		shop_paypal_email_list := strings.Split(shop.PaypalAccount, ",")

		// 之前订单下载和导入时，已经把SellerEmail放入订单中
		seller_email := p_order.SellerEmail //是用SellerEmail来判断老鼠仓吗？

		seller_email_in_list := false
		for _, shop_paypal_email := range shop_paypal_email_list {
			if strings.ToLower(seller_email) == strings.ToLower(shop_paypal_email) {
				seller_email_in_list = true
			}
		}

		// 设置的账户里面和卖家账户匹配的，就设置成为Hold订单
		if seller_email_in_list == false {
			p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH, db_session)
			return true
		}
	}

	return false
}

func (p_ol *OrderLogistics) update_order_to_hold_type(p_order *models.OrderInfo, hold_type_id int64, db_session *xorm.Session) {
	p_order.HoldTypeId = core.MakeNullInt(hold_type_id, true)
	p_order.StatusId = models.ORDER_STATUS_HOLD // 修改订单状态为 hold 状态
	_, err := db_session.Id(p_order.Id).Cols("hold_type_id", "status_id").Update(p_order)
	core.HandleError(err)
}

// 检查仓库是否有订单中的所有商品数量的库存。这里的有货就理解成这个仓库有订单的所有的商品；这里的缺货也理解成是，没有一个仓库能凑齐所有商品的所有数量。同时返回订单的商品总成本
func check_stock_has_order_products(p_ol *OrderLogistics, stock_id int64, p_order_logistics_data *order_logistics_data_t) (has_return bool, cost_amount_return float64) {
	var products_stock_qty_list = p_order_logistics_data.products_stock_qty_list
	var order_product_info_list = p_order_logistics_data.order_product_info_list

	var stock_has_all_product = true
	var cost_amount float64 = 0

	// 对订单中的每个商品信息循环，看是不是这个商品在这个仓库里面有货且有足够的数量
	for _, order_product_info := range order_product_info_list {
		var product_has_stock = false
		for _, product_stock_qty := range products_stock_qty_list {
			if product_stock_qty.StockId == stock_id && product_stock_qty.ProductId == order_product_info.ProductId {
				if product_stock_qty.Qty >= order_product_info.Qty {
					product_has_stock = true
					cost_amount = cost_amount + float64(order_product_info.Qty)*product_stock_qty.CostPrice/p_ol.aud_exchange_rate // 库存成本换算成澳币
					break                                                                                                          //这个 break 加不加都可以的，因为 这个 stock_id， 这个 product_id 的就一条记录
				}
			}
		}

		// 如果商品在这个仓库里面没有货，就判定 这个仓库里面，并不是所有商品都有货
		if product_has_stock == false {
			stock_has_all_product = false
			break // 订单里面的 接下来的 商品也就不用检查了
		}
	}

	if stock_has_all_product == true {
		return true, cost_amount
	} else {
		return false, 0
	}
}

func check_stock_has_some_order_products(p_ol *OrderLogistics, stock_id int64, p_order_logistics_data *order_logistics_data_t) (has_some_return bool) {
	var products_stock_qty_list = p_order_logistics_data.products_stock_qty_list
	var order_product_info_list = p_order_logistics_data.order_product_info_list

	/*
	检查放在所有商品都缺货之后，所以这些就可以省掉了
	var stock_has_all_product = true

	for _, order_product_info := range order_product_info_list {
		var product_has_stock = false
		for _, product_stock_qty := range products_stock_qty_list {
			if product_stock_qty.StockId == stock_id && product_stock_qty.ProductId == order_product_info.ProductId {
				if product_stock_qty.Qty >= order_product_info.Qty {
					product_has_stock = true
					break
				}
			}
		}

		// 如果商品在这个仓库里面没有货，就判定 这个仓库里面，并不是所有商品都有货
		if product_has_stock == false {
			stock_has_all_product = false
			break // 订单里面的 接下来的 商品也就不用检查了
		}
	}
	*/

	var stock_has_some_product = false
	for _, order_product_info := range order_product_info_list {
		for _, product_stock_qty := range products_stock_qty_list {
			if product_stock_qty.StockId == stock_id && product_stock_qty.ProductId == order_product_info.ProductId {
				if product_stock_qty.Qty >= order_product_info.Qty {
					stock_has_some_product = true
					break
				}
			}
		}

		if stock_has_some_product == true {
			break
		}
	}

	//if stock_has_all_product == false && stock_has_some_product == true {
	if stock_has_some_product == true {
		return true
	} else {
		return false
	}
}

// 获得多个商品的整体长宽高重量的参数，输出的结果中，保证 length > width > height
func get_order_products_whole_logistics_data(p_order_logistics_data *order_logistics_data_t) (length_return float64, width_return float64, height float64, weight_return float64) {
	var whole_length, whole_width, whole_height, whole_weight float64 = 0, 0, 0, 0
	for _, order_product_info := range p_order_logistics_data.order_product_info_list {
		// 有时产品的长宽高，可能不是 长 > 宽 > 高 的，需要调整一下，暂且 称作 逻辑 长宽高
		var side_length_slice = sort.Float64Slice([]float64{order_product_info.Length, order_product_info.Width, order_product_info.Height})
		side_length_slice.Sort() // 从小到大的排序
		// 排好序后，依次为 高 宽 长
		var logic_height = side_length_slice[0]
		var logic_width = side_length_slice[1]
		var logic_length = side_length_slice[2]

		if logic_length > whole_length {
			whole_length = logic_length
		}
		if logic_width > whole_width {
			whole_width = logic_width
		}
		whole_height = whole_height + logic_height*float64(order_product_info.Qty)
		whole_weight = whole_weight + order_product_info.Weight*float64(order_product_info.Qty)
	}

	return whole_length, whole_width, whole_height, whole_weight
}

// todo: 同一个商品的一半数量可以从这个仓库发货，另一半数量可以从哪个仓库发货的怎么处理
/*

(完成)1. 地址有特殊字段的。比如 PO BOX, LOCKER。但最佳物流非AU POST的，需要在基础设置里面有特殊字段设置功能
(完成)2. 有买家留言 -- 导入的订单里面还没有留言
3. 有平台留言 -- 导入的订单里面还没有留言
4. PICK UP   -- 一般都是手工创建的订单
(完成)5. 系统显示所有仓都缺货（中国仓除外）
(完成)6. 系统显示最佳仓缺货（次佳、海上有货）
7. 实际仓库缺货（是在捡货过程中，虽然系统显示有货，但实际仓库没有货了。就需要人工手动将这些订单搬至这里的模块，等待跟客户协商，是退款还是等待集装箱）
(完成，临时完成)8. 多数量或多产品，但部分缺货（是否拆分订单）
(完成，临时完成)9. 多数量或多产品，虽不缺货，但拆分订单发货，运费更便宜
(完成，临时完成)10. 多数量或多产品，虽不缺货，但拆分订单发货，运费虽然便宜，但仓库拣货更方便
11. 无法匹配物流（长宽高重量没有物流满足）
(12.完成，这个检查应该放在 11 之前， 就是任何是有长宽高重量数据之前）12. 产品信息不全（长宽高重量不全或者有误）
13. 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损hold
(完成)14. 信封类。长宽高在35cm X 25cm X 2cm, 重量在0.5Kg以内的订单为信封。若是多数量或者多产品，以前我们设计的思路是，最大产品的长边为最长边，最大产品的中边为中边，所有产品的最小边之和为最小边，判断是否在35cm X 25cm X 2cm以内，重量之和是否在0.5Kg以内
(完成)15. 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包
(完成)16. 国际件。对于非澳大利亚的订单，直接到国际件模块
(完成)17. 收货人信息不全的，比如没有手机号，没有邮编等等

 */
