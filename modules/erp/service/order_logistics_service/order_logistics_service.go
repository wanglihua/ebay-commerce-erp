package order_logistics_service

import (
	"errors"
	"fmt"
	"strings"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/envelope_fee_rate_service"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/logistics_fee_rate_service"
	"matrix/modules/erp/service/logistics_mode_service"
	"matrix/modules/erp/service/logistics_setting_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/pick_fee_service"
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/service/order_history_service"
	"matrix/modules/erp/service/order_item_service"
	"matrix/modules/erp/models/xtpz_models"
)

type OrderLogisticsService struct {
	au_stock_list                     []models.StockInfo                      // 所有澳洲仓库列表
	aud_exchange_rate                 float64                                 // 澳币汇率
	logistics_setting                 models.LogisticsSettingInfo             // 物流设置
	pick_fee_setting                  models.PickFeeInfo                      // 拣货费用设置
	address_special_word_list_setting []string                                // 收货地址特殊字段
	logistics_company_list            []models.LogisticsCompanyInfo           // 物流公司列表
	logistics_mode_list               []models.LogisticsModeInfo              // 物流方式列表
	logistics_fee_rate_map            map[int64][]models.LogisticsFeeRateInfo // 物流费率列表的映射，映射的 key 是， 物流方式的 id
	envelope_fee_rate_list            []models.EnvelopeFeeRateInfo            // 信封物流费率， 按费用从小到大排序
	shop_list                         []xtpz_models.ShopAccountInfo           // 店铺列表，系统中的店铺列表
	special_product_list              []models.SpecialProductInfo             // 特殊商品SKU列表
	pending_keeper                    *pending_inventory_temp_keeper          // 设置过程全局变量，每次开始匹配计算的时候一定要记得初始化
	mel_stock_id                      int64                                   // 墨尔本仓库 id
	LoginUser                         core.LoginUser
}

func (p_ol *OrderLogisticsService) init_data(db_session *xorm.Session) {
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

	// 加载物流公司
	var logistics_company_list = make([]models.LogisticsCompanyInfo, 0)
	err = db_session.Find(&logistics_company_list)
	core.HandleError(err)
	p_ol.logistics_company_list = logistics_company_list

	// 加载 物流方式 和 物流费率

	// 加载 物流方式
	// 只加载 物流公司处于启用状态的物流方式
	p_ol.logistics_mode_list = logistics_mode_service.Get_logistics_mode_enabled(db_session)

	// 加载 物流费率
	// 只加载 物流公司处于启用状态的复溜费率
	var logistics_fee_rate_all_list = logistics_fee_rate_service.Get_logistics_fee_rate_enabled(db_session)

	// 考虑到数据量的问题，就不采用中间变量了 var logistics_fee_rate_map = make(map[int64][]models.LogisticsFeeRateInfo)
	p_ol.logistics_fee_rate_map = make(map[int64][]models.LogisticsFeeRateInfo)
	for _, logistics_fee_rate := range logistics_fee_rate_all_list {
		logistics_fee_rate_list, found := p_ol.logistics_fee_rate_map[logistics_fee_rate.ModeId]
		if found {
			logistics_fee_rate_list = append(logistics_fee_rate_list, logistics_fee_rate)
			p_ol.logistics_fee_rate_map[logistics_fee_rate.ModeId] = logistics_fee_rate_list
		} else {
			p_ol.logistics_fee_rate_map[logistics_fee_rate.ModeId] = logistics_fee_rate_list
		}
	}
	// 不采用中间变量了 p_ol.logistics_fee_rate_map = logistics_fee_rate_map

	p_ol.envelope_fee_rate_list = envelope_fee_rate_service.Get_envelope_fee_rate_list(db_session)

	var shop_list = make([]xtpz_models.ShopAccountInfo, 0)
	err = db_session.Find(&shop_list)
	core.HandleError(err)
	p_ol.shop_list = shop_list

	var special_product_list = make([]models.SpecialProductInfo, 0)
	err = db_session.Find(&special_product_list)
	core.HandleError(err)
	p_ol.special_product_list = special_product_list

	p_ol.mel_stock_id = stock_service.Get_mel_stock_id(db_session)
}

func (p_ol *OrderLogisticsService) CalculateListPaidOrders(user core.LoginUser, db_session *xorm.Session) error {

	// 需要先合并
	err := Combine_orders(fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_YZF), user, db_session)
	core.HandleError(err)

	// 获取 已支付 状态的订单, 匹配完成后，订单的状态就是 Hold 或者 待配货
	var order_list = order_service.Get_paid_order_list(db_session)

	if len(order_list) == 0 {
		return errors.New("订单列表为空！")
	}

	//先判断订单状态
	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_YZF {
			return errors.New("订单状态为已支付状态的，才能匹配物流！")
		}
	}

	var order_item_zero_count = order_item_service.Get_order_item_zero_count_in_order_list(order_list, db_session)

	if order_item_zero_count > 0 {
		return errors.New("订单明细不能为空！")
	}

	return p_ol.CalculateList(order_list, db_session)
}

func (p_ol *OrderLogisticsService) CalculateListImportBatch(batch_code int64, user core.LoginUser, db_session *xorm.Session) error {

	// 先合并
	err := Combine_orders(fmt.Sprintf("o.batch_code = %d", batch_code), user, db_session)
	core.HandleError(err)

	//获得该批次号的订单列表
	/*
SELECT o.* FROM erp_order o WHERE o.batch_code = 0
	 */
	var sql_import_batch_order = `SELECT o.* FROM erp_order o WHERE o.batch_code = %d AND o.status_id <> %d ORDER BY o.paid_time ASC, o.id ASC`
	sql_import_batch_order = fmt.Sprintf(sql_import_batch_order, batch_code, models.ORDER_STATUS_YQX)

	var order_list = make([]models.OrderInfo, 0)
	err = db_session.Sql(sql_import_batch_order).Find(&order_list)
	core.HandleError(err)

	if len(order_list) == 0 {
		return errors.New("订单列表为空！")
	}

	//先判断订单状态
	for _, order := range order_list {
		if order.StatusId != models.ORDER_STATUS_YZF {
			return errors.New("订单状态为已支付状态的，才能匹配物流！")
		}
	}

	return p_ol.CalculateList(order_list, db_session)
}

/*
type Interface interface {

    Len() int    // Len 为集合内元素的总数

    Less(i, j int) bool　//如果index为i的元素小于index为j的元素，则返回true，否则返回false

    Swap(i, j int)  // Swap 交换索引为 i 和 j 的元素
}
 */
// 按 id 从小到大 进行排序
type order_slice_t []models.OrderInfo

func (s order_slice_t) Len() int {
	return len(s)
}

func (s order_slice_t) Less(i int, j int) bool {
	return s[i].Id < s[j].Id
}

func (s order_slice_t) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (p_ol *OrderLogisticsService) CalculateList(order_list []models.OrderInfo, db_session *xorm.Session) error {
	var match_option Match_option_t
	match_option.Prefer_stock_id = 0
	match_option.Prefer_mode_id = 0
	match_option.Unlock_express = false

	return p_ol.CalculateListWithOption(order_list, match_option, db_session)
}

type Match_option_t struct {
	Prefer_stock_id int64
	Prefer_mode_id  int64
	Unlock_express  bool
}

// 匹配前要判断一下订单的状态，订单不是已支付状态的话，不能匹配
// 计算订单的物流方式 和 Hold 类型，调用这个函数的时候应该  “开启数据库事务“， 调用完成后，订单的状态就是 Hold 或者 待配货
func (p_ol *OrderLogisticsService) CalculateListWithOption(order_list []models.OrderInfo, match_option Match_option_t, db_session *xorm.Session) error {
	/*
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("匹配物流方式，出现错误： %s\n", err)
		}
	}()
	*/

	if len(order_list) == 0 {
		return errors.New("订单列表为空！")
	}

	p_ol.init_data(db_session)

	{
		// 检查信封限制的相关设置，要放在 init_data() 之后

		var envelope_max_length = p_ol.logistics_setting.EnvelopeMaxLength
		var envelope_max_width = p_ol.logistics_setting.EnvelopeMaxWidth
		var envelope_max_height = p_ol.logistics_setting.EnvelopeMaxHeight
		var envelope_max_weight = p_ol.logistics_setting.EnvelopeMaxWeight

		// 如果没有设置信封的最大长度、宽度、高度、重量
		if envelope_max_length == 0 || envelope_max_width == 0 || envelope_max_height == 0 || envelope_max_weight == 0 {
			return errors.New("系统的对信封的长宽高重量限制的设定不完整！")
		}

		// 检查系统没有设置快递的最小运费，要放在 init_data() 之后
		if p_ol.logistics_setting.MinExpressFee == 0 {
			return errors.New("系统没有设置快递的最小运费！")
		}
	}

	// todo: 拆分和合并订单的问题，在这之前就已经完成了？ 不是在这里做，就可以了，做完进入 hold ??

	// todo 合并订单
	// 同一买家，同一付款paypal，同一收货地址，同一店铺 的订单可以合并

	// 可以先判断除收货地址外的其它信息是否一样，如果这些信息都不一样，就不用加载收货地址信息了

	// 加载所有订单的 收货地址，这样在判断收货地址是否一样的时候，就不会多次访问数据库了
	// 只需要加载前面判断出来的，其它信息一样的订单的收货地址就可以了，不用加载全部订单的收货地址信息
	//var order_shipping_address_list = order_shipping_address_service.Get_orders_shipping_address_list(order_list, db_session)

	p_ol.pending_keeper = new_pending_inventory_temp_keeper()

	for _, order := range order_list {
		//检查订单状态
		//订单状态检查 这些要排除的
		if order.StatusId == models.ORDER_STATUS_YCJ ||
			order.StatusId == models.ORDER_STATUS_PHZ ||
			order.StatusId == models.ORDER_STATUS_YQX ||
			order.StatusId == models.ORDER_STATUS_YWC ||
			order.StatusId == models.ORDER_STATUS_YFH {

			continue
		}

		// 这些是要包含的， 正反两方面进行状态检查
		if order.StatusId != models.ORDER_STATUS_YZF && order.StatusId != models.ORDER_STATUS_HOLD && order.StatusId != models.ORDER_STATUS_DPH {
			continue
		}

		var is_previous_hold = order.StatusId == models.ORDER_STATUS_HOLD &&
			order.HoldTypeId.Int64 != models.ORDER_HOLD_TYPE_ENVELOPE &&
			order.HoldTypeId.Int64 != models.ORDER_HOLD_TYPE_EXPRESS &&
			order.HoldTypeId.Int64 != models.ORDER_HOLD_TYPE_INTERNATIONAL &&
			order.HoldTypeId.Int64 != models.ORDER_HOLD_TYPE_PICK_UP

		var pre_hold_type_id = order.HoldTypeId.Int64

		//移除订单原来匹配的物流信息。如果订单接下来被Hold这些信息可能还会被保留
		//remove_order_logistics(order, db_session) //要禁用，否则会出错

		var p_order_logistics_data = new(order_logistics_data_t)
		//　一起加载需要的数据的方式，对性能的影响应该不会很大的
		p_order_logistics_data.load_data(&order, p_ol, db_session)

		// 匹配理论值
		go p_ol.match_logistics_mode_logic(&order, p_order_logistics_data)

		// 老鼠仓的检查
		if p_ol.hold_check_paypal_account_not_match(&order, db_session) {
			continue
		}

		// 3 有平台留言 -- 导入的订单里面还没有留言
		// 4 PICK UP   -- 一般都是手工创建的订单

		// 8, 9, 10 临时的解决方案
		//todo to complete
		/*
		if p_ol.hold_check_multi_product(&order, p_order_logistics_data, db_session) {
			continue
		}
		*/

		// 12 产品信息不全（长宽高重量不全或者有误），这个检查应该放在任何使用长宽高重量数据的检查之前。 放来匹配物流数据之前
		if p_ol.hold_check_product_info_miss(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 如果不指定匹配的物流方式，就进行信封的 Hold；如果指定了，就不进行信封的 Hold
		if match_option.Prefer_mode_id == 0 {
			// 14 判断是否是信件 放在匹配物流之前
			if p_ol.hold_check_envelope(&order, p_order_logistics_data, db_session) {
				continue
			}
		}

		// 16 判断是否是国际件，如果不是澳洲的就是国际件，一般就是新西兰的， 放在匹配物流之前
		if p_ol.hold_check_international(&order, p_order_logistics_data, db_session) {
			//remove_order_logistics(order, db_session) //要禁用，否则会出错
			continue
		}

		// 邮政编码信息不全 17 收货人信息不全的，比如没有手机号，没有邮编等等，放在这里前面一点的地方检查比较好
		if p_ol.hold_check_address_info_post_code_miss(&order, p_order_logistics_data, db_session) {
			continue
		}

		match_success, order_logistics, is_envelope_hold := p_ol.match_logistics_mode_with_option(&order, p_order_logistics_data, match_option, db_session)

		// 如果通过指定的物流的方式在匹配时进行了信封 Hold，就跳过
		if is_envelope_hold == true {
			//移除刚刚匹配的物流方式
			//remove_order_logistics(order, db_session) //要禁用，否则会出错
			continue
		}

		// 如果匹配过程中， 匹配成快递 就不继续 hold 了。但别的 最佳仓无货，缺货等都还可以继续Hold。
		if order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Int64 == models.ORDER_HOLD_TYPE_EXPRESS {
			continue
		}

		// 11 无匹配的物流方式（长宽高重量没有物流满足）
		if p_ol.hold_check_no_match_logistics_mode(&order, match_success, db_session) {
			continue
		}

		if match_success == true {
			// 13 匹配最佳仓库及物流后订单亏损。需要在基础设置里面有设置亏损额或者亏损率的设置功能，超过设置范围了，自动到订单亏损 hold
			if p_ol.hold_check_great_deficit(&order, p_order_logistics_data, order_logistics, db_session) {
				//把亏损的尽量放在前面
				continue
			}

			// 1 地址有特殊字段的， 因为要查看匹配的物流方式，所以就放在这里了，放在匹配物流方式之后
			if p_ol.hold_check_address_special_word(&order, p_order_logistics_data, order_logistics, db_session) {
				continue
			}
		}

		if match_option.Unlock_express == false {
			// 15 判断是否是快递
			if p_ol.hold_check_express(&order, p_order_logistics_data, db_session) {
				//移除刚刚匹配的物流方式
				//remove_order_logistics(order, db_session) //要禁用，否则会出错
				continue
			}
		}

		// 17 收货人信息不全的，比如没有手机号，没有邮编等等，没有邮编前面也检查了的
		if p_ol.hold_check_address_info_miss(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 2 有买家留言
		if p_ol.hold_check_buyer_checkout_message(&order, db_session) {
			continue
		}

		// 5. 系统显示所有仓都缺货（中国仓除外）
		if p_ol.hold_check_oss_all_stock(&order, p_order_logistics_data, db_session) {
			//部分商品有货的检查，放在所有仓库都缺货的后面
			if p_ol.hold_check_has_some_product(&order, p_order_logistics_data, db_session) {
				continue
			}

			continue
		}

		// 6. 系统显示 在途才有货
		if p_ol.hold_check_oss_au_stock_but_zt_stock_has(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 7. 特殊SKU
		if p_ol.hold_check_special_sku(&order, p_order_logistics_data, db_session) {
			continue
		}

		// 如果状态开始是 hold，现在又变为非Hold的，就把状态改成 已处理 Hold
		if is_previous_hold == true && order.StatusId != models.ORDER_STATUS_HOLD {
			order.PreHoldTypeId = pre_hold_type_id
			order.StatusId = models.ORDER_STATUS_PROCESSED_HOLD
			_, err := db_session.Id(order.Id).Cols("pre_hold_type_id", "status_id").Update(&order)
			core.HandleError(err)

			var login_user = p_ol.LoginUser
			order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("匹配物流订单 Hold 转 已处理Hold， 工号：%s", login_user.LoginName), db_session)
		}
	}

	// 仓库的商品库存状态里面是不是也可以加一个待配货数量，实际库存没有减掉，但是已经这个商品被预定配货了，实际配货成功后，就减掉实际配货 ?? -- 通过在库存表里面插入一个特殊库位的库存记录的方法，这个特殊库位的Id为 -1

	// 更新库存待配货数量
	go inventory_service.Calculate_inventory_pending_qty()

	return nil
}

// 用于单条物流的匹配
func (p_ol *OrderLogisticsService) CalculateSingle(order models.OrderInfo, db_session *xorm.Session) error {
	return p_ol.CalculateList([]models.OrderInfo{order}, db_session)
}

func (p_ol *OrderLogisticsService) CalculateSingleWithOption(order models.OrderInfo, match_option Match_option_t, db_session *xorm.Session) error {
	return p_ol.CalculateListWithOption([]models.OrderInfo{order}, match_option, db_session)
}
