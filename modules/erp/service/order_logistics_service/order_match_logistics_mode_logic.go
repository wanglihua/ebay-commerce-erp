package order_logistics_service

import (
	"math"
	"sort"
	"strings"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/modules/erp/service/pick_fee_service"
	"fmt"
	"matrix/core"
	"time"
	"runtime/debug"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/ddlc_models"
	xorm_core "github.com/go-xorm/core"
	"matrix/modules/erp/models"
)

var db_driver string
var db_url string

func (p_ol *OrderLogisticsService) match_logistics_mode_logic(order ddlc_models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) {

	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("匹配订单物流理论值出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	if db_session == nil {
		// get db driver and db url

		if db_driver == "" {
			db_driver, _ = revel.Config.String("db.driver")
		}

		if db_url == "" {
			db_url, _ = revel.Config.String("db.url")
		}

		db_engine, err := xorm.NewEngine(db_driver, db_url)
		core.HandleError(err)

		db_engine.TZLocation = time.Local
		db_engine.SetMaxIdleConns(3)
		db_engine.SetMaxOpenConns(100)

		if revel.DevMode == true {
			db_engine.ShowSQL()
			db_engine.ShowExecTime()
			xorm_core.LogQuerySql = true
		}

		db_session = db_engine.NewSession()
		db_session = db_session.NoAutoCondition(true)

		defer func() {
			db_session.Close()
		}()
	}

	if order.StatusId == models.ORDER_STATUS_YQX {
		return
	}

	// 下面的三个 Hold 也应该加上理论发货仓库的记录，方便进行理论仓发货的统计

	// pick up
	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PICK_UP {
		var order_logistics_logic ddlc_models.OrderLogisticsLogicInfo
		order_logistics_logic.OrderId = order.Id
		order_logistics_logic.StockId = order.StockId

		_, err := db_session.Insert(order_logistics_logic)
		core.HandleError(err)

		return
	}

	// 信封的忽略
	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_ENVELOPE || p_ol.is_envelope(&order, p_order_logistics_data) == true {
		save_special_hold_order_logistics_logic(p_ol, order.Id, p_order_logistics_data, db_session)
		return
	}

	// 国际件的忽略
	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL || p_ol.is_international_order(&order, p_order_logistics_data) == true {
		save_special_hold_order_logistics_logic(p_ol, order.Id, p_order_logistics_data, db_session)
		return
	}

	// 快递类的忽略
	if order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_EXPRESS || p_ol.is_express(&order, p_order_logistics_data) == true {
		save_special_hold_order_logistics_logic(p_ol, order.Id, p_order_logistics_data, db_session)
		return
	}

	// 邮政编码信息不全 17 收货人信息不全的，比如 没有邮编等等
	if p_ol.is_address_info_post_code_miss(p_order_logistics_data) == true {
		revel.TRACE.Println("is_address_info_post_code_miss")
		return
	}

	// 产品信息不全（长宽高重量不全或者有误），这个检查应该放在任何使用长宽高重量数据的检查之前。
	if p_ol.is_product_info_miss(&order, p_order_logistics_data) == true {
		revel.TRACE.Println("is_product_info_miss")
		return
	}

	// 在这之前要加载好 物流方式表里面的信息，物流费率里面的信息，防止对数据库的频繁访问
	var logistics_setting = p_ol.logistics_setting

	type stock_has_product_and_cost_amount_t struct {
		has_product bool
		cost_amount float64
	}

	// 找出所有有商品数量的仓库
	var au_stock_list = p_ol.au_stock_list
	var stock_id_has_products_map = make(map[int64]stock_has_product_and_cost_amount_t)
	for _, stock := range au_stock_list {
		var stock_has_product_and_cost_amount stock_has_product_and_cost_amount_t
		var cost_amount = get_product_cost_amount(p_ol, stock.Id, p_order_logistics_data)
		stock_has_product_and_cost_amount.has_product = true
		stock_has_product_and_cost_amount.cost_amount = cost_amount

		stock_id_has_products_map[stock.Id] = stock_has_product_and_cost_amount
	}

	var logistics_mode_list = p_ol.logistics_mode_list
	var logistics_fee_rate_map = p_ol.logistics_fee_rate_map
	var logistics_fee_rate2_map = p_ol.logistics_fee_rate2_map

	var final_result match_logistics_mode_result
	var match_result_list = make([]match_logistics_mode_result, 0)

	final_result.order_id = order.Id
	final_result.order_code = order.Code
	final_result.total_fee =  math.MaxFloat64 // 初始化为一个很大的不可能的值，澳币

	// 匹配的结果中，有库存记录的那些记录列表
	for _, logistics_mode := range logistics_mode_list {
		var mode_result match_logistics_mode_result
		mode_result.order_id = order.Id
		mode_result.order_code = order.Code
		mode_result.address_id = p_order_logistics_data.order_shipping_address.Id
		mode_result.address = p_order_logistics_data.order_shipping_address
		mode_result.company_id = logistics_mode.CompanyId

		for _, company := range p_ol.logistics_company_list {
			if company.Id == logistics_mode.CompanyId {
				mode_result.company = company
				break
			}
		}

		mode_result.stock_id = logistics_mode.StockId
		mode_result.has_stock = true

		mode_result.stock.Priority = -100 // 默认给个很低的优先级，这样万一在列表里面没找到 stock 也不会出错

		for _, stock := range au_stock_list {
			if stock.Id == logistics_mode.StockId {
				mode_result.stock = stock
				break
			}
		}

		mode_result.mode_id = logistics_mode.Id
		mode_result.mode = logistics_mode

		logistics_fee_rate_list, found_rate1 := logistics_fee_rate_map[logistics_mode.Id]
		logistics_fee_rate2_list, found_rate2 := logistics_fee_rate2_map[logistics_mode.Id]
		if found_rate1 == false && found_rate2 == false {
			continue
		}

		var stock_id = logistics_mode.StockId
		mode_result.cost_amount = 0
		stock_has_product_and_cost_amount, _ := stock_id_has_products_map[stock_id];
		/*
		if found == false {
			continue // 是不是应该是 continue
		}
		*/

		// 记录库存成本价，如果前面订单 found 为 false 这个成本价有可能是 0
		mode_result.cost_amount = stock_has_product_and_cost_amount.cost_amount

		// 获取所有商品的整体的 长，宽，高，重量 等的数据
		mode_result.length, mode_result.width, mode_result.height, mode_result.actual_weight = get_order_products_whole_logistics_data(p_order_logistics_data)

		mode_result.pick_fee = pick_fee_service.Get_product_pick_fee(mode_result.actual_weight, p_ol.pick_fee_setting)

		// 匹配前打折
		{
			// 匹配前重量打折
			if logistics_setting.Ppqzldz == 1 {
				// 实际重量打折
				mode_result.actual_weight = mode_result.actual_weight * logistics_setting.Ppqzlzl
			}

			// 匹配前中边打折
			if logistics_setting.Ppqzbdz == 1 {
				var side_length_slice = sort.Float64Slice([]float64{mode_result.length, mode_result.width, mode_result.height})
				side_length_slice.Sort()
				var mid_length = side_length_slice[1] // 0, 1, 2 中间那条边

				if math.Abs(mid_length-mode_result.length) < core.FLOAT_MIN_DIFF {
					mode_result.length = mode_result.length * logistics_setting.Ppqzbzl
				} else if math.Abs(mid_length-mode_result.width) < core.FLOAT_MIN_DIFF {
					mode_result.width = mode_result.width * logistics_setting.Ppqzbzl
				} else if math.Abs(mid_length-mode_result.height) < core.FLOAT_MIN_DIFF {
					mode_result.height = mode_result.height * logistics_setting.Ppqzbzl
				}
			} else {
			}
		}

		// 体积重是利用体积重公式计算得到的结果
		// 只有这个参数大于零的时候才有意义
		mode_result.dim_weight = 0
		if logistics_mode.Tjzgs > 0 {
			mode_result.dim_weight = mode_result.length * mode_result.width * mode_result.height / logistics_mode.Tjzgs
		}

		mode_result.weight = mode_result.actual_weight // 用户存储最终的计费重量， 初始为实际重量。如果不考虑体积重的话，就用实际重量来算重量的

		// 体积重与实际重是否有比值要求
		if logistics_mode.Zlsfybzyq == 1 { //如果有比值要求
			// #35 18kg的 比值是 0.5， 实际重量限制是18kg，体积重限制是9kg，在这个限制范围内，他们的体积重和实际重量的比值他们是不关心的。
			var actual_weight_limit = logistics_mode.Zlbncg // 实际重量的限重

			var dim_weight_limit = logistics_mode.Zlbncg * logistics_mode.Tjzsjzbz // 通过体积重和实际重量的比值计算 体积重的限重

			if mode_result.dim_weight > dim_weight_limit || mode_result.actual_weight > actual_weight_limit {
				continue
			}
		} else {
			//如果 体积重与实际重 “有” 比值要求， 是否考虑体积重，体积重与实际重取向都是没有意义的，
			//所以只在 体积重与实际重 “没有” 比值要求的时候，才考虑这两个量

			// 如果考虑体积重，就把计算好的体积重，用 体积重系数 打折
			if logistics_mode.Sfkltjz == 1 { // 是否考虑体积重
				mode_result.sfkltjz = 1

				if logistics_mode.Tjzxs > 0 { // 如果 体积重系数 有设置的值
					mode_result.dim_weight = mode_result.dim_weight / logistics_mode.Tjzxs // 这里是除以体积重系数
				} else {
				}

				// 体积重与实际重取向， 如果考虑体积重，这个系数一般都会设置的？
				if logistics_mode.Tjzysjzqx == 1 { // 取大
					mode_result.weight = math.Max(mode_result.dim_weight, mode_result.actual_weight)
				} else if logistics_mode.Tjzysjzqx == 2 { // 取小
					mode_result.weight = math.Min(mode_result.dim_weight, mode_result.actual_weight)
				}

			} else {
				mode_result.sfkltjz = -1
			}
		}

		// 接下来计算都是针对 final_weight_in_mode 来做，不管前面有没有考虑体积重, 后来就是 mode_result.weight

		var extra_weight_fee_in_mode float64 = 0 // 超重的费用， 澳币
		// 重量是否有限制
		if logistics_mode.Zlsfyxz == 1 { //如果有限制
			var extra_weight = mode_result.weight - logistics_mode.Zlbncg
			if extra_weight > 0 { // 重量不能超过，如果重量有限制，物流方式里面一般都有超重限制
				// 查看是否允许超重
				if logistics_mode.Sfyxcz == 0 { // 不允许超重
					continue // 不允许超重，就排除当前物流方式s
				}
				// 允许超重 计算超重费用
				var extra_weight_step float64 = 25 // 每 25 kg一个超重阶梯
				//超重加收运费
				var cz_step_count = math.Ceil(extra_weight / extra_weight_step) // 这个除数不可能为零
				extra_weight_fee_in_mode = cz_step_count * logistics_mode.Czjsyf
			}
			// 如更没有超过限重，就不需要继续考虑
		}

		var extra_length_fee_in_mode float64 = 0 // 超长的费用， 澳币
		// 最长边是否有限制

		// 长度是否有限制
		if logistics_mode.Zcbsfyxz == 1 { // 如果有限制
			var max_length = math.Max(math.Max(mode_result.length, mode_result.width), mode_result.height)
			var extra_length = max_length - logistics_mode.Zcbxzcd
			if extra_length > 0 { // 长度不能超过，如果最长边有限制，物流方式里面一般都有最长边超长限制
				// 是否允许超长
				if logistics_mode.Sfyxcc == 0 { // 如果不允许超长
					continue
				}

				// 允许超长，计算超长费用
				var extra_length_step float64 = logistics_mode.Ccsfjt
				if extra_length_step > 0 {
					var cc_step_count = math.Ceil(extra_length / extra_length_step)
					extra_length_fee_in_mode = cc_step_count * logistics_mode.Ccjsyf
				} else {
					revel.ERROR.Printf("id 为 %d 名称为 %s 的物流方式的超长收费阶梯设置有误！\n", logistics_mode.Id, logistics_mode.Name)
				}
			}
			// 如果没有超过限长，就不需要继续考虑
		}

		mode_result.fee_rate_id = 0
		mode_result.fee_rate = ddpz_models.LogisticsFeeRateInfo{}
		var total_fee_in_fee_rate float64 = 10000000 //非超重部分的费用，初始为一个很大的不可能的数字， 澳币
		if found_rate1 {
			for _, logistics_fee_rate := range logistics_fee_rate_list {

				// 邮编要对上，忽略开头的 0
				var postcode_in_address = strings.TrimPrefix(strings.TrimSpace(p_order_logistics_data.order_shipping_address.PostalCode), "0")
				var postcode_in_fee_rate = strings.TrimPrefix(strings.TrimSpace(logistics_fee_rate.Postcode), "0")
				if postcode_in_address == postcode_in_fee_rate {
					var fee float64 = logistics_fee_rate.Qbj //起步价
					if mode_result.weight-logistics_fee_rate.Qbjzl > 0 { //如果超过了起步重量，就算超过部分的费率；不超过，就不用算了
						var extra_weight_in_fee_rate = mode_result.weight - logistics_fee_rate.Qbjzl
						if logistics_fee_rate.Xzjt > 0 {
							var xz_fee = math.Ceil(extra_weight_in_fee_rate/logistics_fee_rate.Xzjt) * logistics_fee_rate.Xzjtsf
							fee = fee + xz_fee
						} else {
							revel.ERROR.Println("续重阶梯设置有误，不能计算超重费用！", logistics_fee_rate.Xzjt)
						}
					}
					// 计算最低费用的时候，还不考虑偏远地区的费用
					fee = math.Max(fee, logistics_fee_rate.Zdsf) //有最低收费

					// 偏远地区
					if logistics_fee_rate.Sfpydq == 1 {
						fee = fee + logistics_fee_rate.Pydqjsyf
					}

					if fee < total_fee_in_fee_rate { // 取得更优的费率方案
						total_fee_in_fee_rate = fee
						mode_result.fee_rate_id = logistics_fee_rate.Id
						mode_result.fee_rate = logistics_fee_rate
					}
				}
			} // for _, logistics_fee_rate := range logistics_fee_rate_list {
		}

		if found_rate2 {
			for _, logistics_fee_rate2 := range logistics_fee_rate2_list {

				// 邮编要对上，忽略开头的 0
				var postcode_in_address = strings.TrimPrefix(strings.TrimSpace(p_order_logistics_data.order_shipping_address.PostalCode), "0")
				var postcode_in_fee_rate = strings.TrimPrefix(strings.TrimSpace(logistics_fee_rate2.Postcode), "0")
				if postcode_in_address == postcode_in_fee_rate {
					var fee float64 = 0

					if 0 < mode_result.weight && mode_result.weight <= 0.5 { // 0-500g
						fee = logistics_fee_rate2.Fee01
					} else if 0.5 < mode_result.weight && mode_result.weight <= 1 { // 500g-1kg
						fee = logistics_fee_rate2.Fee02
					} else if 1 < mode_result.weight && mode_result.weight <= 2 { // 1kg-2kg
						fee = logistics_fee_rate2.Fee03
					} else if 2 < mode_result.weight && mode_result.weight <= 3 { // 2kg-3kg
						fee = logistics_fee_rate2.Fee04
					} else if 3 < mode_result.weight && mode_result.weight <= 4 { // 3kg-4kg
						fee = logistics_fee_rate2.Fee05
					} else if 4 < mode_result.weight && mode_result.weight <= 5 { // 4kg-5kg
						fee = logistics_fee_rate2.Fee06
					} else if 5 < mode_result.weight && mode_result.weight <= 7 { // 5kg-7kg
						fee = logistics_fee_rate2.Fee07
					} else if 7 < mode_result.weight && mode_result.weight <= 10 { // 7kg-10kg
						fee = logistics_fee_rate2.Fee08
					} else if 10 < mode_result.weight && mode_result.weight <= 15 { // 10kg-15kg
						fee = logistics_fee_rate2.Fee09
					} else if 15 < mode_result.weight && mode_result.weight <= 22 { // 15kg-22kg
						fee = logistics_fee_rate2.Fee10
					} else if 22 < mode_result.weight && mode_result.weight <= 30 { // 22kg-30kg
						fee = logistics_fee_rate2.Fee11
					} else if 30 < mode_result.weight && mode_result.weight <= 50 { // 30kg-50kg
						fee = logistics_fee_rate2.Fee12
					}

					if fee == 0 {
						continue
					}

					// 偏远地区
					if logistics_fee_rate2.Sfpydq == 1 {
						fee = fee + logistics_fee_rate2.Pydqjsyf
					}

					if fee < total_fee_in_fee_rate { // 取得更优的费率方案
						total_fee_in_fee_rate = fee
						mode_result.fee_rate2_id = logistics_fee_rate2.Id
						mode_result.fee_rate2 = logistics_fee_rate2
					}
				}
			} // for _, logistics_fee_rate2 := range logistics_fee_rate2_list {
		}

		if found_rate1 == true && mode_result.fee_rate_id == 0 { //没有费率匹配成功, 比如说，还没有开通该地区的物流
			continue // 排除该物流方式
		}

		if found_rate2 == true && mode_result.fee_rate2_id == 0 { //没有费率匹配成功, 比如说，还没有开通该地区的物流
			continue // 排除该物流方式
		}

		// 超重 和 超长 费用中， 取一个， 取大的一个
		var extra_fee_in_mode = math.Max(extra_weight_fee_in_mode, extra_length_fee_in_mode)
		mode_result.total_fee = total_fee_in_fee_rate + extra_fee_in_mode
		mode_result.total_amount = mode_result.total_fee + mode_result.cost_amount

		// 匹配后打折
		{
			//这里的计算结果，不对物流的挑选产生影响，但会影响输出
			//现在的做法是，最后输出的匹配结果的数据，是受前后打折两者影响的。

			// 匹配后重量打折
			if logistics_setting.Pphzldz == 1 {
				// 是不是指定实际重量打折
				mode_result.actual_weight = mode_result.actual_weight * logistics_setting.Pphzlzl
			}

			// 匹配后中边打折
			if logistics_setting.Pphzbdz == 1 {
				var side_length_slice = sort.Float64Slice([]float64{mode_result.length, mode_result.width, mode_result.height})
				side_length_slice.Sort()
				var mid_length = side_length_slice[1] // 0, 1, 2 中间那条边

				if (math.Abs(mid_length - mode_result.length)) < core.FLOAT_MIN_DIFF {
					mode_result.length = mode_result.length * logistics_setting.Pphzbzl
				} else if (math.Abs(mid_length - mode_result.width)) < core.FLOAT_MIN_DIFF {
					mode_result.width = mode_result.width * logistics_setting.Pphzbzl
				} else if (math.Abs(mid_length - mode_result.height)) < core.FLOAT_MIN_DIFF {
					mode_result.height = mode_result.height * logistics_setting.Pphzbzl
				}

				// 匹配后中边打折时，同时对体积重打折。
				mode_result.dim_weight = mode_result.dim_weight * logistics_setting.Pphzbzl
			}
		}

		// var total_fee_diff = lowest_total_fee - mode_result.total_fee // 就是当前模式比之前的省了多少钱

		/*
		这里不考虑仓库优先级 20180621
		if total_fee_diff > p_ol.logistics_setting.FeeMaxDiff {
			final_result = mode_result
		} else if math.Abs(total_fee_diff) <= p_ol.logistics_setting.FeeMaxDiff {
			//取优先级高的仓库

			if mode_result.stock.Priority > final_result.stock.Priority { // 当前的仓库的优先级要高
				final_result = mode_result
			}
		}
		*/

		if final_result.total_fee > mode_result.total_fee { // 当前的方案费用更低
			final_result = mode_result
		}

		match_result_list = append(match_result_list, mode_result) // 所有匹配结果都放入列表

	} // for _, logistics_mode := range logistics_mode_list {

	var final_result_second match_logistics_mode_result

	if final_result.fee_rate_id != 0 || final_result.fee_rate2_id != 0 { // 最佳物流方式里面有值
		for _, match_result := range match_result_list {
			if final_result.stock_id != match_result.stock_id { // 和最佳物流不是同一个仓库的
				if (final_result_second.fee_rate_id == 0 && final_result_second.fee_rate2_id == 0 ) ||
					(final_result_second.total_fee > match_result.total_fee && ( match_result.fee_rate_id != 0 || match_result.fee_rate2_id != 0)) { // 如果次佳仓还没有值，或者次佳仓的费用比当前匹配值要高
					final_result_second = match_result
				}
			}
		}
	}

	// 保存最后的匹配结果
	save_order_logistics_logic(p_ol, final_result, final_result_second, db_session)
}

func save_order_logistics_logic(p_ol *OrderLogisticsService, match_result match_logistics_mode_result, match_result_second match_logistics_mode_result, db_session *xorm.Session) {
	// 匹配成功，记录匹配成功数据
	var order_logistics_logic ddlc_models.OrderLogisticsLogicInfo
	has, err := db_session.Where(fmt.Sprintf("order_id = %d", match_result.order_id)).Get(&order_logistics_logic)
	core.HandleError(err)

	order_logistics_logic.OrderId = match_result.order_id
	order_logistics_logic.AddressId = match_result.address_id
	order_logistics_logic.CompanyId = match_result.company_id
	order_logistics_logic.StockId = match_result.stock_id
	order_logistics_logic.ModeId = match_result.mode_id

	if match_result.fee_rate_id != 0 {
		order_logistics_logic.FeeRateId = match_result.fee_rate_id
	}

	if match_result.fee_rate2_id != 0 {
		order_logistics_logic.FeeRateId = match_result.fee_rate2_id
	}

	order_logistics_logic.Length = match_result.length
	order_logistics_logic.Width = match_result.width
	order_logistics_logic.Height = match_result.height
	order_logistics_logic.DimWeight = match_result.dim_weight
	order_logistics_logic.Sfkltjz = match_result.sfkltjz
	order_logistics_logic.ActualWeight = match_result.actual_weight
	order_logistics_logic.Weight = match_result.weight
	order_logistics_logic.TotalFee = match_result.total_fee
	order_logistics_logic.PickFee = match_result.pick_fee
	order_logistics_logic.CostAmount = match_result.cost_amount
	order_logistics_logic.Ptfl = p_ol.logistics_setting.Ptfl
	order_logistics_logic.Jyfl = p_ol.logistics_setting.Jyfl

	order_logistics_logic.CompanyIdSecond = match_result_second.company_id
	order_logistics_logic.StockIdSecond = match_result_second.stock_id
	order_logistics_logic.ModeIdSecond = match_result_second.mode_id

	if match_result_second.fee_rate_id != 0 {
		order_logistics_logic.FeeRateIdSecond = match_result_second.fee_rate_id
	}

	if match_result_second.fee_rate2_id != 0 {
		order_logistics_logic.FeeRateIdSecond = match_result_second.fee_rate2_id
	}

	if order_logistics_logic.FeeRateIdSecond != 0 {
		order_logistics_logic.TotalFeeSecond = match_result_second.total_fee
		order_logistics_logic.FeeDiff = match_result_second.total_fee - match_result.total_fee
	}

	revel.TRACE.Println("order_logistics_logic.FeeRateIdSecond", order_logistics_logic.FeeRateIdSecond)

	if has == true {
		revel.TRACE.Println("order_logistics_logic.FeeRateIdSecond", order_logistics_logic.FeeRateIdSecond)
		_, err := db_session.Id(order_logistics_logic.Id).Update(&order_logistics_logic)
		core.HandleError(err)
	} else {
		_, err := db_session.Insert(&order_logistics_logic)
		core.HandleError(err)
	}
}

func save_special_hold_order_logistics_logic(p_ol *OrderLogisticsService, order_id int64, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) {

	revel.TRACE.Println("save_special_hold_order_logistics_logic")

	var order_logistics_logic ddlc_models.OrderLogisticsLogicInfo
	has, err := db_session.Where("order_id = ?", order_id).Get(&order_logistics_logic)
	core.HandleError(err)

	var whole_length, whole_width, whole_height, whole_weight float64 = get_order_products_whole_logistics_data(p_order_logistics_data)

	order_logistics_logic.OrderId = order_id
	order_logistics_logic.AddressId = p_order_logistics_data.order_shipping_address.Id
	order_logistics_logic.CompanyId = 0
	order_logistics_logic.StockId = p_ol.mel_stock_id //设置发货仓库为墨尔本仓库
	order_logistics_logic.ModeId = 0
	order_logistics_logic.FeeRateId = 0
	order_logistics_logic.Length = whole_length
	order_logistics_logic.Width = whole_width
	order_logistics_logic.Height = whole_height
	order_logistics_logic.Weight = whole_weight
	order_logistics_logic.DimWeight = 0
	order_logistics_logic.Sfkltjz = 0
	order_logistics_logic.Sfkltjz = 0
	order_logistics_logic.Ptfl = p_ol.logistics_setting.Ptfl
	order_logistics_logic.Jyfl = p_ol.logistics_setting.Jyfl

	if has == false {
		_, err = db_session.Insert(&order_logistics_logic)
		core.HandleError(err)
	} else {
		_, err = db_session.Id(order_logistics_logic.Id).Update(&order_logistics_logic)
		core.HandleError(err)
	}
}
