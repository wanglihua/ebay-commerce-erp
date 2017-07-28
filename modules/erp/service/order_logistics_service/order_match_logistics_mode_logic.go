package order_logistics_service

import (
	"math"
	"sort"
	"strings"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/pick_fee_service"
)

func (p_ol *OrderLogisticsService) match_logistics_mode_logic(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) {

	// 信封的忽略
	if p_ol.is_envelope(p_order, p_order_logistics_data) == true {
		return
	}

	// 国际件的忽略
	if p_ol.is_international_order(p_order, p_order_logistics_data) == true {
		return
	}

	// 快递类的忽略
	if p_ol.is_express(p_order, p_order_logistics_data) == true {
		return
	}

	// 邮政编码信息不全 17 收货人信息不全的，比如 没有邮编等等
	if p_ol.is_address_info_post_code_miss(p_order_logistics_data) == true {
		return
	}

	// 产品信息不全（长宽高重量不全或者有误），这个检查应该放在任何使用长宽高重量数据的检查之前。
	if p_ol.is_product_info_miss(p_order, p_order_logistics_data) == true {
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
		has_product, cost_amount := check_stock_has_order_products_and_get_cost_amount(p_ol, stock.Id, p_order_logistics_data);
		stock_has_product_and_cost_amount.has_product = has_product // 仓库有没有全部的商品数量的库存
		stock_has_product_and_cost_amount.cost_amount = cost_amount

		stock_id_has_products_map[stock.Id] = stock_has_product_and_cost_amount
	}

	var logistics_mode_list = p_ol.logistics_mode_list
	var logistics_fee_rate_map = p_ol.logistics_fee_rate_map

	var final_result match_logistics_mode_result
	final_result.order_id = p_order.Id
	final_result.order_code = p_order.Code

	var lowest_total_fee = math.MaxFloat64 // 初始化为一个很大的不可能的值，澳币，

	// 匹配的结果中，有库存记录的那些记录列表
	for _, logistics_mode := range logistics_mode_list {
		var mode_result match_logistics_mode_result
		mode_result.order_id = p_order.Id
		mode_result.order_code = p_order.Code
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

		mode_result.stock.Priority = 100 // 默认给个很低的优先级，这样万一在列表里面没找到 stock 也不会出错

		for _, stock := range au_stock_list {
			if stock.Id == logistics_mode.StockId {
				mode_result.stock = stock
				break
			}
		}

		mode_result.mode_id = logistics_mode.Id
		mode_result.mode = logistics_mode

		logistics_fee_rate_list, found := logistics_fee_rate_map[logistics_mode.Id]
		if found == false {
			continue
		}

		var stock_id = logistics_mode.StockId
		mode_result.cost_amount = 0
		stock_has_product_and_cost_amount, found := stock_id_has_products_map[stock_id];
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

				const MIN_FLOAT_DIFF = 0.0000000001
				if math.Abs(mid_length-mode_result.length) < MIN_FLOAT_DIFF {
					mode_result.length = mode_result.length * logistics_setting.Ppqzbzl
				} else if math.Abs(mid_length-mode_result.width) < MIN_FLOAT_DIFF {
					mode_result.width = mode_result.width * logistics_setting.Ppqzbzl
				} else if math.Abs(mid_length-mode_result.height) < MIN_FLOAT_DIFF {
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
		mode_result.fee_rate = models.LogisticsFeeRateInfo{}
		var total_fee_in_fee_rate float64 = 10000000 //非超重部分的费用，初始为一个很大的不可能的数字， 澳币
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

		if mode_result.fee_rate_id == 0 { //没有费率匹配成功, 比如说，还没有开通该地区的物流
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

				const MIN_FLOAT_DIFF = 0.0000000001
				if (math.Abs(mid_length - mode_result.length)) < MIN_FLOAT_DIFF {
					mode_result.length = mode_result.length * logistics_setting.Pphzbzl
				} else if (math.Abs(mid_length - mode_result.width)) < MIN_FLOAT_DIFF {
					mode_result.width = mode_result.width * logistics_setting.Pphzbzl
				} else if (math.Abs(mid_length - mode_result.height)) < MIN_FLOAT_DIFF {
					mode_result.height = mode_result.height * logistics_setting.Pphzbzl
				}

				// 匹配后中边打折时，同时对体积重打折。
				mode_result.dim_weight = mode_result.dim_weight * logistics_setting.Pphzbzl
			}
		}

		var totol_fee_diff = lowest_total_fee - mode_result.total_fee // 就是当前模式比之前的省了多少钱

		// 三个区间 final_total_amount-mode_result.total_amount，后来改用 fee了，不考虑 库存的cost_price，因为没有库存的也参与筛选
		// 负无穷大 到 -p_ol.logistics_setting.FeeMaxDiff， -p_ol.logistics_setting.FeeMaxDiff到+p_ol.logistics_setting.FeeMaxDiff, +p_ol.logistics_setting.FeeMaxDiff 到 正无穷大
		if totol_fee_diff > p_ol.logistics_setting.FeeMaxDiff { // 取得更优的模式方案，挑选最优解时，把商品的到仓成本也考虑进去   +p_ol.logistics_setting.FeeMaxDiff 到 正无穷大
			final_result = mode_result
		} else if math.Abs(totol_fee_diff) <= p_ol.logistics_setting.FeeMaxDiff { // 运费偏差在最大允许范围内，其实还是算上了澳仓成本的  -p_ol.logistics_setting.FeeMaxDiff到+p_ol.logistics_setting.FeeMaxDiff,
			//取优先级高的仓库

			if mode_result.stock.Priority > final_result.stock.Priority { // 当前的仓库的优先级要高
				final_result = mode_result
			}
		}

		if mode_result.total_fee < lowest_total_fee {
			lowest_total_fee = mode_result.total_fee
		}

	} // for _, logistics_mode := range logistics_mode_list {

	//todo: 保存最后的匹配结果
}
