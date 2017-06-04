package order_logistics_service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/pick_fee_service"
)

func (p_ol *OrderLogistics) match_logistics_mode(p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, db_session *xorm.Session) (success_return bool, order_logistics_return models.OrderLogisticsInfo) {

	// 在这里面只考虑最佳仓库的库存情况，次佳仓库有货的情况，别的地方已经考虑并 hold

	// 在这之前要加载好 物流方式表里面的信息，物流费率里面的信息，防止对数据库的频繁访问

	var logistics_setting = p_ol.logistics_setting

	// 找出所有有商品数量的仓库
	var au_stock_list = p_ol.au_stock_list
	var stock_id_has_products_map = make(map[int64]float64)
	for _, stock := range au_stock_list {
		if has, cost_amount := check_stock_has_order_products(p_ol, stock.Id, p_order_logistics_data); has == true { // 仓库有没有全部的商品数量的库存，有
			stock_id_has_products_map[stock.Id] = cost_amount
		}
	}

	var logistics_mode_list = p_ol.logistics_mode_list
	var logistics_fee_rate_map = p_ol.logistics_fee_rate_map

	var final_result match_logistics_mode_result
	final_result.order_id = p_order.Id
	final_result.order_code = p_order.Code
	final_result.address_id = 0
	final_result.company_id = 0
	final_result.stock_id = 0
	//final_result.stock
	final_result.mode_id = 0
	// final_result.mode
	final_result.fee_rate_id = 0
	final_result.total_fee = 0 // 最终的物流费用
	final_result.length = 0
	final_result.width = 0
	final_result.height = 0
	// var final_weight float64 = 0
	final_result.actual_weight = 0 //就是商品的实际重量合计，用 final 只是用于和下面的循环里面的 actual_weight 区别开来罢了
	final_result.dim_weight = 0
	final_result.cost_amount = 0
	final_result.total_amount = 0 // 总的 费用和成本的和 final_total_fee + final_cost_amount,

	var lowest_total_fee = math.MaxFloat64 // 初始化为一个很大的不可能的值，澳币，

	// 清除订单匹配的调试信息
	if p_ol.logistics_setting.Qyppts > 0 {
		var delete_sql = `DELETE FROM erp_order_logistics_match_debug WHERE order_id = %d;`
		delete_sql = fmt.Sprintf(delete_sql, p_order.Id)
		_, err := db_session.Exec(delete_sql)
		core.HandleError(err)
	}

	var match_result_list = make([]match_logistics_mode_result, 0)
	for _, logistics_mode := range logistics_mode_list {
		var mt match_trace
		mt.set_is_trace(p_ol.logistics_setting.Qyppts)

		var mode_result match_logistics_mode_result
		mode_result.order_id = p_order.Id
		mode_result.order_code = p_order.Code
		mode_result.address_id = p_order_logistics_data.order_shipping_address.Id
		mode_result.company_id = logistics_mode.CompanyId
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

		mt.Println("订单编号", p_order.Code)
		mt.Println("匹配单个物流方式开始")
		mt.Println("物流方式数据库id", logistics_mode.Id, "物流方式名称：", logistics_mode.Name)
		logistics_fee_rate_list, found := logistics_fee_rate_map[logistics_mode.Id]
		if found == false {
			mt.Println("没有找到物流方式对应的费率设置")
			mt.Println("匹配单个物流方式结束")

			if mt.is_trace == true {
				log_logistics_match_debug(mode_result, "无对应费率设置", mt.String(), db_session)
			}

			continue
		}

		var stock_id = logistics_mode.StockId
		mt.Println("物流方式对应的仓库数据库id", stock_id)
		mode_result.cost_amount = 0
		if cost_amount, found := stock_id_has_products_map[stock_id]; found == false { // 如果物流方式对应的仓库，没有全部的商品数量的库存，就继续下一个物流方式
			mt.Println("物流方式对应的仓库没有订单对应的库存")
			//mt.Println("匹配单个物流方式结束")
			mt.Println("继续进行物流方式费用的计算")
			mode_result.has_stock = false

			/*

			if mt.is_trace == true {
				log_logistics_match_debug(mode_result, "无对应库存", mt.String(), db_session)
			}
			continue

			*/
		} else { // 如果库存都有的，就记录下商品成本总价
			mt.Println("物流方式对应的仓库存在订单对应的库存，澳仓成本价为", cost_amount)
			mode_result.cost_amount = cost_amount
		}

		// 获取所有商品的整体的 长，宽，高，重量 等的数据
		mode_result.length, mode_result.width, mode_result.height, mode_result.actual_weight = get_order_products_whole_logistics_data(p_order_logistics_data)
		mt.Println("计算所得的订单商品的整体，长，宽，高，实际重量", mode_result.length, mode_result.width, mode_result.height, mode_result.actual_weight)

		mode_result.pick_fee = pick_fee_service.Get_product_pick_fee(mode_result.actual_weight, p_ol.pick_fee_setting)

		// 匹配前打折
		{
			mt.Println("现在考虑匹配前打折")

			// 匹配前重量打折
			if logistics_setting.Ppqzldz == 1 {
				// 实际重量打折
				mt.Println("匹配前实际重量打折", "是")
				mt.Println("匹配前实际重量折率", logistics_setting.Ppqzlzl)
				mt.Println("匹配前打折前的实际重量", mode_result.actual_weight)
				mode_result.actual_weight = mode_result.actual_weight * logistics_setting.Ppqzlzl
				mt.Println("匹配前打折后的实际重量", mode_result.actual_weight)
			} else {
				mt.Println("匹配实际重量打折", "否")
			}

			// 匹配前中边打折
			if logistics_setting.Ppqzbdz == 1 {
				mt.Println("匹配前中边打折", "是")

				mt.Println("打折时的长宽高", mode_result.length, mode_result.width, mode_result.height)

				var side_length_slice = sort.Float64Slice([]float64{mode_result.length, mode_result.width, mode_result.height})
				side_length_slice.Sort()
				var mid_length = side_length_slice[1] // 0, 1, 2 中间那条边

				mt.Println("匹配前中边折率", logistics_setting.Ppqzbzl)
				mt.Println("打折前的中边为", mid_length)

				const MIN_FLOAT_DIFF = 0.0000000001
				if math.Abs(mid_length-mode_result.length) < MIN_FLOAT_DIFF {
					mode_result.length = mode_result.length * logistics_setting.Ppqzbzl
					mt.Println("打折后的中边为", mode_result.length)
				} else if math.Abs(mid_length-mode_result.width) < MIN_FLOAT_DIFF {
					mode_result.width = mode_result.width * logistics_setting.Ppqzbzl
					mt.Println("打折后的中边为", mode_result.width)
				} else if math.Abs(mid_length-mode_result.height) < MIN_FLOAT_DIFF {
					mode_result.height = mode_result.height * logistics_setting.Ppqzbzl
					mt.Println("打折后的中边为", mode_result.height)
				}
			} else {
				mt.Println("匹配前中边打折", "否")
			}
		}

		// 体积重是利用体积重公式计算得到的结果
		// 只有这个参数大于零的时候才有意义
		mode_result.dim_weight = 0
		if logistics_mode.Tjzgs > 0 {
			mt.Println("物流方式的体积重公式参数为", logistics_mode.Tjzgs)
			mode_result.dim_weight = mode_result.length * mode_result.width * mode_result.height / logistics_mode.Tjzgs
			mt.Println("计算体积重，长，宽，高，体积重", mode_result.length, mode_result.width, mode_result.height, mode_result.dim_weight)
		} else {
			mt.Println("物流方式的体积重公式参数为 0， 不计算体积重。")
		}

		mt.Println("计算所得的订单商品的整体，长，宽，高，体积重，实际重量", mode_result.length, mode_result.width, mode_result.height, mode_result.dim_weight, mode_result.actual_weight)

		mode_result.weight = mode_result.actual_weight // 用户存储最终的计费重量， 初始为实际重量。如果不考虑体积重的话，就用实际重量来算重量的

		// 体积重与实际重是否有比值要求
		if logistics_mode.Zlsfybzyq == 1 { //如果有比值要求
			// #35 18kg的 比值是 0.5， 实际重量限制是18kg，体积重限制是9kg，在这个限制范围内，他们的体积重和实际重量的比值他们是不关心的。
			mt.Println("体积重与实际重是否有比值要求 ", "是")

			var actual_weight_limit = logistics_mode.Zlbncg // 实际重量的限重
			mt.Println("实际重量的限重 ", actual_weight_limit)

			var dim_weight_limit = logistics_mode.Zlbncg * logistics_mode.Tjzsjzbz // 通过体积重和实际重量的比值计算 体积重的限重
			mt.Println("体积重的限重 ", dim_weight_limit)

			if mode_result.dim_weight > dim_weight_limit || mode_result.actual_weight > actual_weight_limit {
				if mode_result.actual_weight > actual_weight_limit {
					mt.Println("体积重与实际重比值有要求，但实际重量超过了限重")
				}

				if mode_result.dim_weight > dim_weight_limit {
					mt.Println("体积重与实际重比值有要求，但体积重超过了限重")
				}

				if mt.is_trace == true {
					log_logistics_match_debug(mode_result, "体积重与实际重比值有要求，但超过了限重", mt.String(), db_session)
				}

				continue
			} else {
				mt.Println("体积重与实际重比值有要求，没有超过限重")
			}
		} else {
			//如果 体积重与实际重 “有” 比值要求， 是否考虑体积重，体积重与实际重取向都是没有意义的，
			//所以只在 体积重与实际重 “没有” 比值要求的时候，才考虑这两个量

			mt.Println("体积重与实际重是否有比值要求", "否")

			// 如果考虑体积重，就把计算好的体积重，用 体积重系数 打折
			if logistics_mode.Sfkltjz == 1 { // 是否考虑体积重
				mt.Println("是否考虑体积重", "是")

				if logistics_mode.Tjzxs > 0 { // 如果 体积重系数 有设置的值
					mt.Println("体积重系数", logistics_mode.Tjzxs)
					mt.Println("体积重除以体积重系数之前", mode_result.dim_weight)
					mode_result.dim_weight = mode_result.dim_weight / logistics_mode.Tjzxs // 这里是除以体积重系数
					mt.Println("体积重除以体积重系数之后", mode_result.dim_weight)
				} else {
					mt.Println("体积重系数未设置或设置不符合要求")
				}

				// 体积重与实际重取向， 如果考虑体积重，这个系数一般都会设置的？
				if logistics_mode.Tjzysjzqx == 1 { // 取大
					mt.Println("体积重与实际重取向", "取大")
					mt.Println("这时的体积重和实际重量分别是", mode_result.dim_weight, mode_result.actual_weight)
					mode_result.weight = math.Max(mode_result.dim_weight, mode_result.actual_weight)
					mt.Println("体积重与实际重取向 取大 后的结果是", mode_result.weight)
				} else if logistics_mode.Tjzysjzqx == 2 { // 取小
					mt.Println("体积重与实际重取向", "取小")
					mt.Println("这时的体积重和实际重量分别是", mode_result.dim_weight, mode_result.actual_weight)
					mode_result.weight = math.Min(mode_result.dim_weight, mode_result.actual_weight)
					mt.Println("体积重与实际重取向 取小 后的结果是", mode_result.weight)
				}

			} else {
				mt.Println("是否考虑体积重", "否")
			}
		}

		// 接下来计算都是针对 final_weight_in_mode 来做，不管前面有没有考虑体积重, 后来就是 mode_result.weight

		var extra_weight_fee_in_mode float64 = 0 // 超重的费用， 澳币
		// 重量是否有限制
		if logistics_mode.Zlsfyxz == 1 { //如果有限制
			mt.Println("重量是否有限制", "是")
			mt.Println("当前重量是", mode_result.weight)
			mt.Println("限制重量是", logistics_mode.Zlbncg)
			var extra_weight = mode_result.weight - logistics_mode.Zlbncg
			if extra_weight > 0 { // 重量不能超过，如果重量有限制，物流方式里面一般都有超重限制
				mt.Println("超重了，超重的重量为", extra_weight)
				// 查看是否允许超重
				if logistics_mode.Sfyxcz == 0 { // 不允许超重
					mt.Println("是否允许超重", "否")
					mt.Println("匹配单个物流方式结束")
					if mt.is_trace == true {
						log_logistics_match_debug(mode_result, "超重", mt.String(), db_session)
					}
					continue // 不允许超重，就排除当前物流方式s
				} else {
					mt.Println("是否允许超重", "是")
				}
				// 允许超重 计算超重费用
				var extra_weight_step float64 = 25 // 每 25 kg一个超重阶梯
				mt.Println("超重阶梯", extra_weight_step)
				//超重加收运费
				var cz_step_count = math.Ceil(extra_weight / extra_weight_step) // 这个除数不可能为零
				mt.Println("超重阶梯数量", cz_step_count)
				mt.Println("超重阶梯收费", logistics_mode.Czjsyf)
				extra_weight_fee_in_mode = cz_step_count * logistics_mode.Czjsyf
				mt.Println("超重收费", extra_weight_fee_in_mode)
			} else {
				mt.Println("没有超重")
			}
			// 如更没有超过限重，就不需要继续考虑
		} else {
			mt.Println("重量是否有限制", "否")
		}

		var extra_length_fee_in_mode float64 = 0 // 超长的费用， 澳币
		// 最长边是否有限制

		// 长度是否有限制
		if logistics_mode.Zcbsfyxz == 1 { // 如果有限制
			mt.Println("长度是否有限制", "是")
			var max_length = math.Max(math.Max(mode_result.length, mode_result.width), mode_result.height)
			mt.Println("这时的边长分别是", mode_result.length, mode_result.width, mode_result.height)
			mt.Println("这时的最长边是", max_length)
			mt.Println("最长边限制 长度", logistics_mode.Zcbxzcd)
			var extra_length = max_length - logistics_mode.Zcbxzcd
			if extra_length > 0 { // 长度不能超过，如果最长边有限制，物流方式里面一般都有最长边超长限制
				mt.Println("超长了，超过的长度为", extra_length)
				// 是否允许超长
				if logistics_mode.Sfyxcc == 0 { // 如果不允许超长
					mt.Println("是否允许超长", "否")
					mt.Println("匹配单个物流方式结束")
					if mt.is_trace == true {
						log_logistics_match_debug(mode_result, "超长", mt.String(), db_session)
					}
					continue
				} else {
					mt.Println("是否允许超长", "是")
				}

				// 允许超长，计算超长费用
				var extra_length_step float64 = logistics_mode.Ccsfjt
				if extra_length_step > 0 {
					mt.Println("超长阶梯", logistics_mode.Ccsfjt)
					var cc_step_count = math.Ceil(extra_length / extra_length_step)
					mt.Println("超长阶梯数量", cc_step_count)
					mt.Println("超长阶梯收费", logistics_mode.Ccjsyf)
					extra_length_fee_in_mode = cc_step_count * logistics_mode.Ccjsyf
					mt.Println("超长收费", extra_length_fee_in_mode)
				} else {
					revel.ERROR.Printf("id 为 %d 名称为 %s 的物流方式的超长收费阶梯设置有误！\n", logistics_mode.Id, logistics_mode.Name)
				}
			} else {
				mt.Println("没有超过限长")
			}
			// 如果没有超过限长，就不需要继续考虑
		} else {
			mt.Println("长度是否有限制", "否")
		}

		mode_result.fee_rate_id = 0
		var total_fee_in_fee_rate float64 = 10000000 //非超重部分的费用，初始为一个很大的不可能的数字， 澳币
		for _, logistics_fee_rate := range logistics_fee_rate_list {

			// 邮编要对上，忽略开头的 0
			var postcode_in_address = strings.TrimPrefix(strings.TrimSpace(p_order_logistics_data.order_shipping_address.PostalCode), "0")
			var postcode_in_fee_rate = strings.TrimPrefix(strings.TrimSpace(logistics_fee_rate.Postcode), "0")
			if postcode_in_address == postcode_in_fee_rate {
				mt.Println()
				mt.PrintlnIndent(1, "匹配单个费率开始")

				mt.PrintlnIndent(1, "订单收货地址的邮编", p_order_logistics_data.order_shipping_address.PostalCode)
				mt.PrintlnIndent(1, "当前物流方式下，符合邮编的费率数据库 id", logistics_fee_rate.Id)

				mt.PrintlnIndent(1, "费率起步价", logistics_fee_rate.Qbj)
				mt.PrintlnIndent(1, "费率起步价重量", logistics_fee_rate.Qbjzl)

				var fee float64 = logistics_fee_rate.Qbj //起步价
				if mode_result.weight-logistics_fee_rate.Qbjzl > 0 { //如果超过了起步重量，就算超过部分的费率；不超过，就不用算了
					mt.PrintlnIndent(1, "超过起步价重量了")
					var extra_weight_in_fee_rate = mode_result.weight - logistics_fee_rate.Qbjzl
					mt.PrintlnIndent(1, "超过重量", extra_weight_in_fee_rate)
					mt.PrintlnIndent(1, "续重阶梯", logistics_fee_rate.Xzjt)
					if logistics_fee_rate.Xzjt > 0 {
						var xz_step_count = math.Ceil(extra_weight_in_fee_rate / logistics_fee_rate.Xzjt)
						mt.PrintlnIndent(1, "续重阶梯数量", xz_step_count)
						mt.PrintlnIndent(1, "续重阶梯收费", logistics_fee_rate.Xzjtsf)
						var xz_fee = math.Ceil(extra_weight_in_fee_rate/logistics_fee_rate.Xzjt) * logistics_fee_rate.Xzjtsf
						mt.PrintlnIndent(1, "续重费用", xz_fee)
						fee = fee + xz_fee
						mt.PrintlnIndent(1, "续重后总费用", fee)
					} else {
						mt.PrintlnIndent(1, "续重后总费用", fee)
					}
				} else {
					revel.ERROR.Println("续重阶梯设置有误，不能计算超重费用！", logistics_fee_rate.Xzjt)
					mt.PrintlnIndent(1, "续重阶梯设置有误，不能计算超重费用！", logistics_fee_rate.Xzjt)
				}
				mt.PrintlnIndent(1, "费率最低收费", logistics_fee_rate.Zdsf)
				// 计算最低费用的时候，还不考虑偏远地区的费用
				fee = math.Max(fee, logistics_fee_rate.Zdsf) //有最低收费
				mt.PrintlnIndent(1, "考虑费率最低收费后的费用", fee)

				// 偏远地区
				if logistics_fee_rate.Sfpydq == 1 {
					mt.PrintlnIndent(1, "是否偏远地区", "是")
					mt.PrintlnIndent(1, "偏远地区加收运费", logistics_fee_rate.Pydqjsyf)
					fee = fee + logistics_fee_rate.Pydqjsyf
					mt.PrintlnIndent(1, "考虑偏远地区后的费用", fee)
				} else {
					mt.PrintlnIndent(1, "是否偏远地区", "否")
				}

				mt.PrintlnIndent(1, "最终当前费率的费用", fee)
				if fee < total_fee_in_fee_rate { // 取得更优的费率方案
					mt.PrintlnIndent(1, "之前没有匹配的费率，或当前费率比之前匹配的费率更优，暂时选用当前费率计算结果")
					total_fee_in_fee_rate = fee
					mode_result.fee_rate_id = logistics_fee_rate.Id
				} else {
					mt.PrintlnIndent(1, "之前匹配的费率，比当前费率更优，弃用当前费率计算结果")
				}
				mt.PrintlnIndent(1, "匹配单个费率结束")
				mt.Println()
			}
		} // for _, logistics_fee_rate := range logistics_fee_rate_list {

		if mode_result.fee_rate_id == 0 { //没有费率匹配成功, 比如说，还没有开通该地区的物流
			mt.Println("该物流方式下没有匹配的物流费率，弃用当前物流方式")
			mt.Println("匹配单个物流方式结束")

			if mt.is_trace == true {
				log_logistics_match_debug(mode_result, "没有匹配的物流费率", mt.String(), db_session)
			}

			continue // 排除该物流方式
		}

		// 超重 和 超长 费用中， 取一个， 取大的一个
		mt.Println("物流方式超重费用", extra_weight_fee_in_mode)
		mt.Println("物流方式超长费用", extra_length_fee_in_mode)
		var extra_fee_in_mode = math.Max(extra_weight_fee_in_mode, extra_length_fee_in_mode)
		mt.Println("物流方式超重和超长，取大的一个的结果", extra_fee_in_mode)
		mode_result.total_fee = total_fee_in_fee_rate + extra_fee_in_mode
		mt.Println("物流方式中物流费率费用加上物流方式超重超长大的一个后的结果", mode_result.total_fee)
		mode_result.total_amount = mode_result.total_fee + mode_result.cost_amount
		mt.Println("考虑澳仓成本后的，物流方式总费用", mode_result.total_amount)

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

		if mode_result.has_stock == true {
			// 三个区间 final_total_amount-mode_result.total_amount，后来改用 fee了，不考虑 库存的cost_price，因为没有库存的也参与筛选
			// 负无穷大 到 -p_ol.logistics_setting.FeeMaxDiff， -p_ol.logistics_setting.FeeMaxDiff到+p_ol.logistics_setting.FeeMaxDiff, +p_ol.logistics_setting.FeeMaxDiff 到 正无穷大
			if totol_fee_diff > p_ol.logistics_setting.FeeMaxDiff { // 取得更优的模式方案，挑选最优解时，把商品的到仓成本也考虑进去   +p_ol.logistics_setting.FeeMaxDiff 到 正无穷大
				mt.Println("之前没有匹配的物流方式，或当前物流方式比之前匹配的物流方式节省费用超过最大允许偏差，暂时选用当前物流方式计算结果")
				final_result = mode_result
				match_result_list = append(match_result_list, mode_result)
			} else if math.Abs(totol_fee_diff) <= p_ol.logistics_setting.FeeMaxDiff { // 运费偏差在最大允许范围内，其实还是算上了澳仓成本的  -p_ol.logistics_setting.FeeMaxDiff到+p_ol.logistics_setting.FeeMaxDiff,
				//取优先级高的仓库

				if mode_result.stock.Priority > final_result.stock.Priority { // 当前的仓库的优先级要高
					mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，且当前仓库是优先级更高，选用当前物流方式计算结果。")
					final_result = mode_result
				} else { //
					mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，但当前仓库是优先级并不更高，弃用当前物流方式计算结果，但记录次佳仓备选。")
				}
				match_result_list = append(match_result_list, mode_result)
			} else { // 在运费偏差范围之外， 负无穷大 到 -p_ol.logistics_setting.FeeMaxDiff
				mt.Println("之前匹配的物流，比当前物流方式更优, 且超过了允许偏差范围，弃用当前物流方式计算结果")
			}
		} else {
			//没有库存
			if totol_fee_diff > p_ol.logistics_setting.FeeMaxDiff {
				mt.Println("之前没有匹配的物流方式，或当前物流方式比之前匹配的物流方式节省费用超过最大允许偏差，暂时选用当前物流方式计算结果，仓库无货。")
				final_result = mode_result
				match_result_list = append(match_result_list, mode_result)
			} else if math.Abs(totol_fee_diff) <= p_ol.logistics_setting.FeeMaxDiff {
				//取优先级高的仓库

				if mode_result.stock.Priority > final_result.stock.Priority { // 当前的仓库的优先级要高
					mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，且当前仓库是优先级更高，选用当前物流方式计算结果，仓库无货。")
					final_result = mode_result
					match_result_list = append(match_result_list, mode_result)
				} else { //
					mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，但当前仓库是优先级并不更高，弃用当前物流方式计算结果，仓库无货，不记录次佳仓备选。")
				}
			} else { // 在运费偏差范围之外， 负无穷大 到 -p_ol.logistics_setting.FeeMaxDiff
				mt.Println("之前匹配的物流，比当前物流方式更优, 且超过了允许偏差范围，弃用当前物流方式计算结果，但仓库无货。")
			}
		}

		if mode_result.total_fee < lowest_total_fee {
			lowest_total_fee = mode_result.total_fee
		}

		mt.Println("匹配单个物流方式结束")

		if mt.is_trace == true {
			if mode_result.has_stock == true {
				log_logistics_match_debug(mode_result, "备选", mt.String(), db_session)
			} else {
				log_logistics_match_debug(mode_result, "无对应库存", mt.String(), db_session)
			}
		}
	} // for _, logistics_mode := range logistics_mode_list {

	if final_result.mode_id == 0 || final_result.fee_rate_id == 0 || final_result.has_stock == false { // 匹配失败
		var order_logistics models.OrderLogisticsInfo
		if final_result.has_stock == false {
			//todo: 记录物流信息 和 hold
		}

		return false, order_logistics
	}

	var order_logistics = save_order_logistics(p_order, final_result, true, db_session)

	// 多个商品时，考虑拆包的方式

	// 多个商品时，能不能分仓库发货的？

	return true, order_logistics
}
