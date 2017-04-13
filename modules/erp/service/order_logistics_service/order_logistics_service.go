package order_logistics_service

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/core/collections"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/logistics_setting_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/pick_fee_service"
)

type OrderLogistics struct {
	au_stock_list                     []models.StockInfo                      // 所有澳洲仓库列表
	aud_exchange_rate                 float64                                 // 澳币汇率
	logistics_setting                 models.LogisticsSettingInfo             // 物流设置
	pick_fee_setting                  models.PickFeeInfo                      // 拣货费用设置
	address_special_word_list_setting []string                                // 收货地址特殊字段
	logistics_mode_list               []models.LogisticsModeInfo              // 物流方式列表
	logistics_fee_rate_map            map[int64][]models.LogisticsFeeRateInfo // 物流费率列表的映射，映射的 key 是， 物流方式的 id
}

func (p_ol *OrderLogistics) init_data(db_session *xorm.Session) {
	// 加载所有仓库
	var au_stock_list = make([]models.StockInfo, 0)
	err := db_session.Where(fmt.Sprintf("stock_type_id = %d", models.STOCK_TYPE_AU)).Find(&au_stock_list)
	core.HandleError(err)

	p_ol.au_stock_list = au_stock_list

	// 获取澳币的汇率，如果定单里面的汇率没有记录下来，就用这个汇率
	var aud_exchange_rate_setting = exchange_rate_service.Get_aud_exchange_rate(db_session)
	if aud_exchange_rate_setting == 0 {
		panic(errors.New("澳币汇率未设置！"))
	}
	p_ol.aud_exchange_rate = aud_exchange_rate_setting

	// 获取物流参数设置
	var logistics_setting = logistics_setting_service.Get_logistics_setting(db_session)
	p_ol.logistics_setting = logistics_setting

	// 获取拣货人工费用
	p_ol.pick_fee_setting = pick_fee_service.Get_pick_fee_setting(db_session)

	// 获得地址特殊字段的列表
	var address_special_word_list_setting = strings.Split(logistics_setting.AddressSpecialWord, ",")
	for index, address_special_word := range address_special_word_list_setting {
		address_special_word_list_setting[index] = strings.TrimSpace(address_special_word)
	}
	p_ol.address_special_word_list_setting = address_special_word_list_setting

	// 加载 物流方式 和 物流费率

	// 加载 物流方式
	var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
	err = db_session.Find(&logistics_mode_list)
	core.HandleError(err)
	p_ol.logistics_mode_list = logistics_mode_list

	// 加载 物流费率
	var logistics_fee_rate_all_list = make([]models.LogisticsFeeRateInfo, 0)
	err = db_session.Find(&logistics_fee_rate_all_list)
	core.HandleError(err)

	var logistics_fee_rate_map = make(map[int64][]models.LogisticsFeeRateInfo)
	for _, logistics_fee_rate := range logistics_fee_rate_all_list {
		logistics_fee_rate_list, found := logistics_fee_rate_map[logistics_fee_rate.ModeId]
		if found {
			logistics_fee_rate_list = append(logistics_fee_rate_list, logistics_fee_rate)
			logistics_fee_rate_map[logistics_fee_rate.ModeId] = logistics_fee_rate_list
		} else {
			logistics_fee_rate_map[logistics_fee_rate.ModeId] = logistics_fee_rate_list
		}
	}
	p_ol.logistics_fee_rate_map = logistics_fee_rate_map
}

// 计算订单的物流方式 和 Hold 类型，调用这个函数的时候应该开启数据库事务，
func (p_ol *OrderLogistics) Calculate(db_session *xorm.Session) {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("匹配物流方式，出现错误： %s\n", err)
		}
	}()

	p_ol.init_data(db_session)

	// 拆分和合并订单，在这之前就已经完成了？ 不是在这里做，就可以了，做完进入 hold

	// 获取 已支付 状态的订单, 匹配完成后，订单的状态就是 Hold 或者 待配货
	var paid_order_list = order_service.Get_paid_order_list(db_session)

	for _, order := range paid_order_list {
		var p_order_logistics_data = new(order_logistics_data_t)
		//　一起加载需要的数据的方式，对性能的影响应该不会很大的
		// p_order_logistics_data.load_data(&order, p_ol, db_session)

		// 加载订单数据 汇率 和 订单收货地址
		p_order_logistics_data.load_exchange_rate(&order, p_ol)
		p_order_logistics_data.load_order_shipping_address(order.Id, db_session)

		// 17 收货人信息不全的，比如没有手机号，没有邮编等等，放在这里前面一点的地方检查比较好
		if p_ol.hold_check_address_info_miss(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 2 有买家留言
		if p_ol.hold_check_buyer_checkout_message(&order, db_session) {
			continue
		}

		// 3 有平台留言 -- 导入的订单里面还没有留言
		// 4 PICK UP   -- 一般都是手工创建的订单

		// 加载订单数据 订单产品信息 和 订单产品的库存信息
		// 加载接下来的判断需要的数据，如果在循环开始时加载，就会有可能的性能浪费
		p_order_logistics_data.load_order_product_info_list(order.Id, db_session)
		p_order_logistics_data.load_products_stock_qty_list(order.Id, db_session)

		// 5. 系统显示所有仓都缺货（中国仓除外）
		if p_ol.hold_check_oss_all_stock(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 6. 系统显示最佳仓缺货（次佳、海上有货）
		if p_ol.hold_check_oss_best_stock(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 7 实际仓库缺货（是在捡货过程中，虽然系统显示有货，但实际仓库没有货了。就需要人工手动将这些订单搬至这里的模块，等待跟客户协商，是退款还是等待集装箱）

		// 8, 9, 10 临时的解决方案
		if p_ol.hold_check_multi_product(&order, p_order_logistics_data, db_session) {
			continue
		}

		match_success, order_logistics_id := p_ol.match_logistics_mode(&order, p_order_logistics_data, db_session)

		// 11 无法匹配物流（长宽高重量没有物流满足）
		if p_ol.hold_check_no_match_logistics_mode(&order, match_success, db_session) {
			continue
		}

		// 1 地址有特殊字段的， 因为要查看匹配的物流方式，所以就放在这里了
		if p_ol.hold_check_address_special_word(&order, p_order_logistics_data, order_logistics_id, db_session) {
			continue
		}

		// 12 产品信息不全（长宽高重量不全或者有误），这个检查应该放在任何使用长宽高重量数据的检查之前
		if p_ol.hold_check_product_info_miss(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 13 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损 hold
		if p_ol.hold_check_great_deficit(&order, p_order_logistics_data, order_logistics_id, db_session) {
			continue
		}

		// 14 判断是否是信件
		if p_ol.hold_check_envelope(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 15 判断是否是快递
		if p_ol.hold_check_express(&order, db_session) {
			continue
		}

		// 16 判断是否是国际件，如果不是澳洲的就是国际件，一般就是新西兰的
		if p_ol.hold_check_international(&order, p_order_logistics_data, db_session) {
			continue
		}
	}

	// 仓库的商品库存状态里面是不是也可以加一个待配货数量，实际库存没有减掉，但是已经这个商品被预定配货了，实际配货成功后，就减掉实际配货 ?? -- 通过在库存表里面插入一个特殊库位的库存记录的方法，这个特殊库位的Id为 -1

}

// 1 地址有特殊字段的
func (p_ol *OrderLogistics) hold_check_address_special_word(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, order_logistics_id int64, db_session *xorm.Session) (is_hold_return bool) {
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
		// 检查匹配的物流方式对应的物流公司是不是 需要考虑 地址特殊字段的
		var address_special_word struct {
			AddressSpecialWord int64 `xorm:"'address_special_word'"`
		}

		/*
	SELECT
		lc.address_special_word
	FROM erp_order_logistics ol
		INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
		INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
	WHERE ol.id = 1
		 */

		sql := `SELECT
		lc.address_special_word
	FROM erp_order_logistics ol
		INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
		INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
	WHERE ol.id = %d`
		sql = fmt.Sprintf(sql, order_logistics_id)

		_, err := db_session.Sql(sql).Get(&address_special_word)
		core.HandleError(err)

		if address_special_word.AddressSpecialWord == 0 { // 地址有特殊字段的。比如 PO BOX, LOCKER。但最佳物流   “ 非“   AU POST的，需要在基础设置里面有特殊字段设置功能
			p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_SPECIAL_ADDRESS, db_session)
			return true
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

// 3 有平台留言 -- 导入的订单里面还没有留言

// 4 PICK UP   -- 一般都是手工创建的订单

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

	if p_ol.hold_check_oss(p_order, p_order_logistics_data, fun_is_not_cn_stock, db_session) {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_OOS_ALL_STOCK, db_session)
		return true
	}

	return false
}

// 6. 系统显示最佳仓缺货（次佳、海上有货）
func (p_ol *OrderLogistics) hold_check_oss_best_stock(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
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
	if p_ol.hold_check_oss(p_order, p_order_logistics_data, fun_is_au_stock, db_session) == true && p_ol.hold_check_oss(p_order, p_order_logistics_data, fun_is_zt_stock, db_session) == false {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_OOS_BEST_STOCK, db_session)
		return true
	}

	return false
}

//返回 true, 表明考虑的是符合要求的
type fun_stock_type_filter_t func(stock_type_id int64) bool

// 订单商品库存检查，用户给5和6的调用
func (p_ol *OrderLogistics) hold_check_oss(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, fun_stock_type_filter fun_stock_type_filter_t, db_session *xorm.Session) (is_oss_return bool) {
	//todo: 暂时理解成是，所有商品，所有仓库(中国仓除外)都缺货。这里的缺货也理解成是，没有一个仓库能凑齐所有商品的所有数量
	var products_stock_qty_list = p_order_logistics_data.products_stock_qty_list

	// 找出所有包含这些商品库存的仓库id列表
	var get_product_stock_id_list = func() []int64 {
		var stock_id_set = collections.NewIntSet()
		for _, product_stock_qty := range products_stock_qty_list {
			stock_id_set.Add(product_stock_qty.StockId)
		}
		return stock_id_set.List()
	}
	var product_stock_id_list = get_product_stock_id_list()

	var all_product_oss = true //下面的逻辑中，如果发现仓库有货，就设置为 false

	// 对每一个仓库循环，看是不是这个仓库里面, 这些商品都有货
	for _, stock_id := range product_stock_id_list {
		var stock_has_all_product = check_stock_has_order_products(stock_id, p_order_logistics_data)

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

// 7 实际仓库缺货（是在捡货过程中，虽然系统显示有货，但实际仓库没有货了。就需要人工手动将这些订单搬至这里的模块，等待跟客户协商，是退款还是等待集装箱）

// 8, 9, 10 临时的解决方案
func (p_ol *OrderLogistics) hold_check_multi_product(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	if len(p_order_logistics_data.order_product_info_list) > 1 {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_MULTI_PRODUCT, db_session)
		return true
	}

	return false
}

// 11 无法匹配物流（长宽高重量没有物流满足）
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
		// 长宽高重量数据不全
		if order_product_info.Length == 0 || order_product_info.Width == 0 || order_product_info.Height == 0 || order_product_info.Weight == 0 {
			p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_PRODUCT_INFO_MISS, db_session)
			return true
		}
	}

	return false
}

// 13 //todo: 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损hold
func (p_ol *OrderLogistics) hold_check_great_deficit(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, order_logistics_id int64, db_session *xorm.Session) (is_hold_return bool) {
	//todo: to complete
	if true {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_GREAT_DEFICIT, db_session)
		return true
	}

	return false
}

// 14 信封类。长宽高在35cm X 25cm X 2cm, 重量在0.5Kg以内的订单为信封。若是多数量或者多产品，以前我们设计的思路是，最大产品的长边为最长边，最大产品的中边为中边，所有产品的最小边之和为最小边，判断是否在35cm X 25cm X 2cm以内，重量之和是否在0.5Kg以内
func (p_ol *OrderLogistics) hold_check_envelope(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	// 这里假定信封的长宽高限制： 长 > 宽 > 高

	var envelope_max_length = p_ol.logistics_setting.EnvelopeMaxLength
	var envelope_max_width = p_ol.logistics_setting.EnvelopeMaxWidth
	var envelope_max_height = p_ol.logistics_setting.EnvelopeMaxHeight
	var envelope_max_weight = p_ol.logistics_setting.EnvelopeMaxWeight

	// 如果没有设置信封的最大长度、宽度、高度、重量
	if envelope_max_length == 0 || envelope_max_width == 0 || envelope_max_height == 0 || envelope_max_weight == 0 {
		return false
	}

	// get_order_products_whole_logistics_data的输出中 长 > 宽 > 高
	var whole_length, whole_width, whole_height, _, whole_weight float64 = get_order_products_whole_logistics_data(p_order_logistics_data, 0)

	if whole_length <= envelope_max_length && whole_width <= envelope_max_width && whole_height <= envelope_max_height && whole_weight <= envelope_max_weight {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_ENVELOPE, db_session)
		return true
	}

	return false
}

// 15 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包裹
func (p_ol *OrderLogistics) hold_check_express(p_order *models.OrderInfo, db_session *xorm.Session) (is_hold_return bool) {
	// 如果没有设置最小运费
	if p_ol.logistics_setting.MinExpressFee == 0 {
		return false
	}

	if p_order.ExpressFee >= p_ol.logistics_setting.MinExpressFee { //假定都是澳币
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_EXPRESS, db_session)
		return true
	}

	return false
}

// 16 判断是否是国际件，如果不是澳洲的就是国际件，一般就是新西兰的
func (p_ol *OrderLogistics) hold_check_international(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	if strings.ToLower(p_order_logistics_data.order_shipping_address.Country) == strings.ToLower("AU") {
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_TYPE_INTERNATIONAL, db_session)
		return true
	}

	return false
}

// 17. 收货人信息不全的，比如没有手机号，没有邮编等等
func (p_ol *OrderLogistics) hold_check_address_info_miss(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (is_hold_return bool) {
	var shipping_address = p_order_logistics_data.order_shipping_address
	if shipping_address.Id == 0 || // 数据库中地址信息不存在
		strings.TrimSpace(shipping_address.Phone) == "" || // 地址信息中缺少电话号码
		strings.TrimSpace(shipping_address.PostalCode) == "" || // 地址信息中缺少邮政编码
		strings.TrimSpace(shipping_address.Name) == "" || // 地址信息中缺少收件人
		(strings.TrimSpace(shipping_address.Street1) == "" && strings.TrimSpace(shipping_address.Street2) == "" ) { // 地址信息中缺少详细地址信息
		p_ol.update_order_to_hold_type(p_order, models.ORDER_HOLD_ADDRESS_INFO_MISS, db_session)
		return true
	}

	return false
}

func (p_ol *OrderLogistics) update_order_to_hold_type(p_order *models.OrderInfo, hold_type_id int64, db_session *xorm.Session) {
	p_order.HoldTypeId = core.MakeNullInt(hold_type_id, true)
	_, err := db_session.Id(p_order.Id).Cols("hold_type_id").Update(p_order)
	core.HandleError(err)
}

func (p_ol *OrderLogistics) match_logistics_mode(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (success_return bool, order_logistics_id_return int64) {

	// 在这里面只考虑最佳仓库的库存情况，次佳仓库有货的情况，别的地方已经考虑并 hold

	// 在这之前要加载好 物流方式表里面的信息，物流费率里面的信息，防止对数据库的频繁访问

	var logistics_setting = p_ol.logistics_setting

	// 找出所有有商品数量的仓库
	var au_stock_list = p_ol.au_stock_list
	var stock_id_has_products_map = make(map[int64]bool)
	for _, stock := range au_stock_list {
		if check_stock_has_order_products(stock.Id, p_order_logistics_data) { // 仓库没有全部的商品数量的库存
			stock_id_has_products_map[stock.Id] = true
		}
	}

	var logistics_mode_list = p_ol.logistics_mode_list
	var logistics_fee_rate_map = p_ol.logistics_fee_rate_map

	var final_selected_company_id int64 = 0
	var final_selected_stock_id int64 = 0
	var final_selected_mode_id int64 = 0
	var final_selected_fee_rate_id int64 = 0
	var final_total_fee float64 = 100000; // 最终的物流费用，初始化为一个很大的不可能的值，澳币
	var final_length float64 = 0
	var final_width float64 = 0
	var final_height float64 = 0
	var final_weight float64 = 0
	var final_actual_weight float64 = 0 //就是商品的实际重量合计，用 final 只是用于和下面的循环里面的 actual_weight 区别开来罢了
	for _, logistics_mode := range logistics_mode_list {
		logistics_fee_rate_list, found := logistics_fee_rate_map[logistics_mode.Id]
		if found == false {
			continue
		}

		var stock_id = logistics_mode.StockId
		if _, found := stock_id_has_products_map[stock_id]; found == false { // 如果物流方式对应的仓库，没有全部的商品数量的库存，就继续下一个物流方式
			continue
		}

		// 获取所有商品的整体的 长，宽，高，重量 等的数据，体积重是利用体积重公式计算得到的结果
		var length, width, height, dim_weight, actual_weight float64 = get_order_products_whole_logistics_data(p_order_logistics_data, logistics_mode.Tjzgs)

		// 匹配前打折
		{
			// 匹配前重量打折
			if logistics_setting.Ppqzldz == 1 {
				// 是不是指定实际重量打折
				actual_weight = actual_weight * logistics_setting.Ppqzlzl
			}

			// 匹配前中边打折
			if logistics_setting.Ppqzbdz == 1 {
				var side_length_slice = sort.Float64Slice([]float64{length, width, height})
				side_length_slice.Sort()
				var mid_length = side_length_slice[1] // 0, 1, 2 中间那条边

				const MIN_FLOAT_DIFF = 0.0000000001
				if (mid_length - length) < MIN_FLOAT_DIFF {
					length = length * logistics_setting.Ppqzbzl
				}
				if (mid_length - width) < MIN_FLOAT_DIFF {
					width = width * logistics_setting.Ppqzbzl
				}
				if (mid_length - height) < MIN_FLOAT_DIFF {
					height = height * logistics_setting.Ppqzbzl
				}
			}
		}

		var final_weight_in_mode = actual_weight // 用户存储最终的计费重量， 初始为实际重量。如果不考虑体积重的话，就用实际重量来算重量的

		// 如果考虑体积重，就把计算好的体积重，用 体积重系数 打折
		if logistics_mode.Sfkltjz == 1 { // 是否考虑体积重
			if logistics_mode.Tjzxs > 0 { // 如果 体积重系数 有设置的值
				dim_weight = dim_weight / logistics_mode.Tjzxs // 这里是除以体积重系数
			}

			// 体积重与实际重取向， 如果考虑体积重，这个系数一般都会设置的？
			if logistics_mode.Tjzysjzqx == 1 { // 取大
				final_weight_in_mode = math.Max(dim_weight, actual_weight)
			} else if logistics_mode.Tjzysjzqx == 2 { // 取小
				final_weight_in_mode = math.Min(dim_weight, actual_weight)
			}

			// 体积重与实际重是否有比值要求，这个比值是要大于的，就是说这个东西不能太沉
			if logistics_mode.Zlsfybzyq == 1 { //如果有比值要求
				if dim_weight/actual_weight < logistics_mode.Tjzsjzbz { //如果小于了比值就是不符合要求，排除这个物流方式，就是这个包裹不能太沉
					continue
				}
			}
		}

		// 接下来计算都是针对 final_weight_in_mode 来做，不管前面有没有考虑体积重

		var extra_weight_fee_in_mode float64 = 0 // 超重的费用， 澳币
		// 重量是否有限制
		if logistics_mode.Zlsfyxz == 1 { //如果有限制
			var extra_weight = final_weight_in_mode - logistics_mode.Zlbncg
			if extra_weight > 0 { // 重量不能超过，如果重量有限制，物流方式里面一般都有超重限制
				// 查看是否允许超重
				if logistics_mode.Sfyxcz == 0 { // 不允许超重
					continue // 不允许超重，就排除当前物流方式s
				}
				// 允许超重 计算超重费用
				var extra_weight_step float64 = 25 // 每 25 kg一个超重阶梯
				extra_weight_fee_in_mode = math.Ceil(extra_weight/extra_weight_step) * logistics_mode.Czjsyf
			}
			// 如更没有超过限重，就不需要继续考虑
		}

		var extra_length_fee_in_mode float64 = 0 // 超长的费用， 澳币
		// 最长边是否有限制
		var max_length = math.Max(math.Max(length, width), height)
		// 长度是否有限制
		if logistics_mode.Zcbsfyxz == 1 { // 如果有限制
			var extra_length = max_length - logistics_mode.Zcbxzcd
			if extra_length > 0 { // 长度不能超过，如果最长边有限制，物流方式里面一般都有最长边超长限制
				// 是否允许超长
				if logistics_mode.Sfyxcc == 0 { // 如果不允许超长
					continue
				}
				// 允许超长，计算超长费用
				var extra_length_step float64 = logistics_mode.Ccsfjt
				extra_length_fee_in_mode = math.Ceil(extra_length/extra_length_step) * logistics_mode.Ccjsyf
			}
			// 如果没有超过限长，就不需要继续考虑
		}

		var selected_fee_rate_id_in_mode int64 = 0
		var total_fee_in_fee_rate float64 = 10000000 //非超重部分的费用，初始为一个很大的不可能的数字， 澳币
		for _, logistics_fee_rate := range logistics_fee_rate_list {
			// 邮编要对上
			if strings.TrimSpace(p_order_logistics_data.order_shipping_address.PostalCode) == strings.TrimSpace(logistics_fee_rate.Postcode) {
				var fee float64 = logistics_fee_rate.Qbj //起步价
				if final_weight_in_mode-logistics_fee_rate.Qbjzl > 0 { //如果超过了起步重量，就算超过部分的费率；不超过，就不用算了
					fee = math.Ceil((final_weight_in_mode-logistics_fee_rate.Qbjzl)/logistics_fee_rate.Xzjt) * logistics_fee_rate.Xzjtsf
				}
				// 计算最低费用的时候，还不考虑偏远地区的费用
				fee = math.Max(fee, logistics_fee_rate.Zdsf) //有最低收费

				// 偏远地区
				if logistics_fee_rate.Sfpydq == 1 {
					fee = fee + logistics_fee_rate.Pydqjsyf
				}

				if fee < total_fee_in_fee_rate { // 取得更优的费率方案
					total_fee_in_fee_rate = fee
					selected_fee_rate_id_in_mode = logistics_fee_rate.Id
				}
			}
		} // for _, logistics_fee_rate := range logistics_fee_rate_list {

		if selected_fee_rate_id_in_mode == 0 { //没有费率匹配, 比如说，还没有开通该地区的物流
			continue // 排除该物流方式
		}

		// 超重 和 超长 费用中， 取一个， 取大的一个
		var extra_fee_in_mode = math.Max(extra_weight_fee_in_mode, extra_length_fee_in_mode)
		var total_fee_in_mode float64 = total_fee_in_fee_rate + extra_fee_in_mode

		if total_fee_in_mode < final_total_fee { // 取得更优的模式方案
			final_selected_fee_rate_id = selected_fee_rate_id_in_mode
			final_selected_company_id = logistics_mode.CompanyId
			final_selected_stock_id = logistics_mode.StockId
			final_selected_mode_id = logistics_mode.Id
			final_total_fee = total_fee_in_mode

			final_length = length
			final_width = width
			final_height = height
			final_weight = final_weight_in_mode
			final_actual_weight = actual_weight
		}
	} // for _, logistics_mode := range logistics_mode_list {

	if final_selected_mode_id == 0 || final_selected_fee_rate_id == 0 { // 匹配失败
		return false, 0
	}

	// 现在的做法是，最后输出的匹配结果的数据，是受前后打折两者影响的。

	// 匹配后打折
	{
		// 匹配后重量打折
		if logistics_setting.Pphzldz == 1 {
			// 是不是指定实际重量打折
			final_weight = final_weight * logistics_setting.Pphzlzl
		}

		// 匹配后中边打折
		if logistics_setting.Pphzbdz == 1 {
			var side_length_slice = sort.Float64Slice([]float64{final_length, final_width, final_height})
			side_length_slice.Sort()
			var mid_length = side_length_slice[1] // 0, 1, 2 中间那条边

			const MIN_FLOAT_DIFF = 0.0000000001
			if (mid_length - final_length) < MIN_FLOAT_DIFF {
				final_length = final_length * logistics_setting.Pphzbzl
			}
			if (mid_length - final_width) < MIN_FLOAT_DIFF {
				final_width = final_width * logistics_setting.Pphzbzl
			}
			if (mid_length - final_height) < MIN_FLOAT_DIFF {
				final_height = final_height * logistics_setting.Pphzbzl
			}
		}
	}

	var address_id = p_order_logistics_data.order_shipping_address.Id

	// 匹配成功，记录匹配成功数据
	var order_logistics models.OrderLogisticsInfo
	has, err := db_session.Where(fmt.Sprintf("order_id = %d", p_order.Id)).Get(&order_logistics)
	core.HandleError(err)

	order_logistics.OrderId = p_order.Id
	order_logistics.AddressId = address_id
	order_logistics.CompanyId = final_selected_company_id
	order_logistics.StockId = final_selected_stock_id
	order_logistics.ModeId = final_selected_mode_id
	order_logistics.FeeRateId = final_selected_fee_rate_id
	order_logistics.Length = final_length
	order_logistics.Width = final_width
	order_logistics.Height = final_height
	order_logistics.Weight = final_weight
	order_logistics.TotalFee = final_total_fee
	order_logistics.PickFee = pick_fee_service.Get_product_pick_fee(final_actual_weight, p_ol.pick_fee_setting)

	if has == true {
		_, err := db_session.Id(order_logistics.Id).Update(&order_logistics)
		core.HandleError(err)
	} else {
		_, err := db_session.Insert(&order_logistics)
		core.HandleError(err)
	}

	// 多个商品时，考虑拆包的方式

	// 多个商品时，能不能分仓库发货的？

	return true, order_logistics.Id
}

// 检查仓库是否有订单中的所有商品数量的库存
func check_stock_has_order_products(stock_id int64, p_order_logistics_data *order_logistics_data_t) bool {
	var products_stock_qty_list = p_order_logistics_data.products_stock_qty_list
	var order_product_info_list = p_order_logistics_data.order_product_info_list

	var stock_has_all_product = true

	// 对订单中的每个商品信息循环，看是不是这个商品在这个仓库里面有货且有足够的数量
	for _, order_product_info := range order_product_info_list {
		var product_has_stock = false
		for _, product_stock_qty := range products_stock_qty_list {
			if product_stock_qty.StockId == stock_id && product_stock_qty.ProductId == order_product_info.ProductId {
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

	return stock_has_all_product
}

// 输出的结果中，保证 length > width > height
func get_order_products_whole_logistics_data(p_order_logistics_data *order_logistics_data_t, dim_weight_ratio float64) (length_return float64, width_return float64, height float64, dim_weight_return float64, weight_return float64) {
	var whole_length, whole_width, whole_height, whole_dim_weight, whole_weight float64 = 0, 0, 0, 0, 0
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
		whole_height = whole_height + logic_height
		whole_weight = whole_weight + order_product_info.Weight
	}

	// 只有这个参数大于零的时候才有意义
	if dim_weight_ratio > 0 {
		whole_dim_weight = whole_length * whole_width * whole_height / dim_weight_ratio
	}

	return whole_length, whole_width, whole_height, whole_dim_weight, whole_weight
}

//todo: 同一个商品的一半数量可以从这个仓库发货，另一半数量可以从哪个仓库发货的怎么处理
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
(完成)15. 快递类。需要在基础设置里面有 收到快递费金额的设置功能，比如收到的运费高于什么金额，视为快递。否则就算收了快递费也视为普通包裹
(完成)16. 国际件。对于非澳大利亚的订单，直接到国际件模块
(完成)17. 收货人信息不全的，比如没有手机号，没有邮编等等

 */

//todo: 地址特殊字段的问题，与最后匹配的最佳物流方式有关
//todo: 体积重 实际重 比值的问题，可以继续关注，虽然感觉是理解的是对的，就是要确认
