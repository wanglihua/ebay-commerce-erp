package order_logistics_service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/service/pick_fee_service"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/ddlc_models"
	"time"
	"matrix/modules/erp/models/ddpz_models/logistics_fee_type"
)

/*
func (p_ol *OrderLogistics) match_logistics_mode(p_order *ddlc_models.OrderInfo, p_order_logistics_data *order_logistics_data_t, perfer_stock_id int64, db_session *xorm.Session) (success_return bool, order_logistics_return ddlc_models.OrderLogisticsInfo) {
	return p_ol.match_logistics_mode_with_prefer_stock(p_order, p_order_logistics_data, 0, db_session)
}
*/

type match_return_t struct {
	success          bool                           // 就是看有没有匹配的物流方式
	order_logistics  ddlc_models.OrderLogisticsInfo // 匹配的订单物流信息
	is_envelope_hold bool                           // 是否信封HOLD
	has_stock        bool                           // 匹配时仓库是否有库存
}

func (p_ol *OrderLogisticsService) match_logistics_mode_with_option(p_order *ddlc_models.OrderInfo, p_order_logistics_data *order_logistics_data_t, match_option Match_option_t, is_previous_hold bool, db_session *xorm.Session) match_return_t {

	// 在这里面只考虑最佳仓库的库存情况，次佳仓库有货的情况，别的地方已经考虑并 hold

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
		has_product, cost_amount := check_stock_has_order_products_and_get_cost_amount(p_ol, p_order, stock.Id, p_order_logistics_data, true) // 最后一个参数取 true，考虑订单中被 temp pending的数量
		stock_has_product_and_cost_amount.has_product = has_product                                                                           // 仓库有没有全部的商品数量的库存
		stock_has_product_and_cost_amount.cost_amount = cost_amount

		stock_id_has_products_map[stock.Id] = stock_has_product_and_cost_amount
	}

	var logistics_mode_list = p_ol.logistics_mode_list
	var logistics_fee_rate_map = p_ol.logistics_fee_rate_map
	var logistics_fee_rate2_map = p_ol.logistics_fee_rate2_map

	var final_result match_logistics_mode_result
	final_result.order_id = p_order.Id
	final_result.order_code = p_order.Code

	var lowest_total_fee = math.MaxFloat64 // 初始化为一个很大的不可能的值，澳币，

	// 清除订单匹配的调试信息
	if p_ol.logistics_setting.Qyppts > 0 {
		var delete_sql = `DELETE FROM erp_order_logistics_match_debug WHERE order_id = %d;`
		delete_sql = fmt.Sprintf(delete_sql, p_order.Id)
		_, err := db_session.Exec(delete_sql)
		core.HandleError(err)
	}

	// 匹配的结果中，有库存记录的那些记录列表
	var match_result_has_stock_list = make([]match_logistics_mode_result, 0)
	var is_envelope_hold = false
	var is_express = p_ol.is_express(p_order, p_order_logistics_data) // 根据付款的快递费判断是不是快递
	for _, logistics_mode := range logistics_mode_list {

		if match_option.Unlock_express == false && is_express == true { // 根据付款的快递费判断是快递，但不是解锁快递的话，就只匹配快递的物流方式
			if logistics_mode.TypeId != ddpz_models.LOGISTICS_MODE_TYPE_EXPRESS {
				continue
			}
		}

		// 如果指定了仓库，并且和指定的仓库不同，就忽略
		if match_option.Prefer_stock_id != 0 && match_option.Prefer_stock_id != logistics_mode.StockId {
			continue
		}

		// 如果指定了物流方式，并且和指定的物流方式不同，就忽略
		if match_option.Prefer_mode_id != 0 && match_option.Prefer_mode_id != logistics_mode.Id {
			continue
		}

		// 如果指定了 物流方式，并且指定的物流方式是 信封， 就设置为信封Hold，并且匹配信封费率
		if match_option.Prefer_mode_id != 0 && match_option.Prefer_mode_id == logistics_mode.Id && logistics_mode.TypeId == ddpz_models.LOGISTICS_MODE_TYPE_ENVELOPE {
			p_ol.update_order_to_hold_type(p_order, ddpz_models.ORDER_HOLD_TYPE_ENVELOPE, db_session)
			var match_envelope_success = p_ol.match_envelope_fee_rate(p_order, p_order_logistics_data, db_session)
			if match_envelope_success == false {
				// 强制匹配为信封，但是费率没有匹配成功，就是长宽高重量没有符合信封的要求，就把费用最大的那个信封费率设置为最后匹配成功的费率
				var envelope_fee_rate_list = p_ol.envelope_fee_rate_list
				var max_fee_envelop_fee_rate = envelope_fee_rate_list[len(envelope_fee_rate_list)-1]
				save_special_hold_order_logistics(p_ol, p_order, p_order_logistics_data, max_fee_envelop_fee_rate.Id, max_fee_envelop_fee_rate.Fee, db_session)
			}
			is_envelope_hold = true
			break
		}

		// 如果是解锁快递，并且当前的物流方式是快递的，就排除掉
		if match_option.Unlock_express == true && logistics_mode.TypeId == ddpz_models.LOGISTICS_MODE_TYPE_EXPRESS {
			continue
		}

		var mt match_trace
		mt.set_is_trace(p_ol.logistics_setting.Qyppts)

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

		mt.Println("订单编号", p_order.Code)
		mt.Println("匹配单个物流方式开始")
		mt.Println("物流方式数据库id", logistics_mode.Id, "物流方式名称：", logistics_mode.Name)

		if logistics_mode.FeeType != logistics_fee_type.FEE_TYPE1 && logistics_mode.FeeType != logistics_fee_type.FEE_TYPE2 {
			revel.ERROR.Println(fmt.Sprintf("物流方式的费率类型设置值不对！ id: %d, name: %s", logistics_mode.Id, logistics_mode.Name))
		}

		logistics_fee_rate_list, found_rate1 := logistics_fee_rate_map[logistics_mode.Id]
		logistics_fee_rate2_list, found_rate2 := logistics_fee_rate2_map[logistics_mode.Id]
		if logistics_mode.FeeType == logistics_fee_type.FEE_TYPE1 {
			if found_rate1 == false {
				mt.Println("没有找到物流方式的对应邮编的费率设置")
				mt.Println("匹配单个物流方式结束")

				if mt.is_trace == true {
					log_logistics_match_debug(mode_result, "无邮编对应费率设置", mt.String(), db_session)
				}

				continue
			}
		}

		if logistics_mode.FeeType == logistics_fee_type.FEE_TYPE2 {
			if found_rate2 == false {
				mt.Println("没有找到物流方式的对应邮编的费率2设置")
				mt.Println("匹配单个物流方式结束")

				if mt.is_trace == true {
					log_logistics_match_debug(mode_result, "无邮编对应费率2设置", mt.String(), db_session)
				}

				continue
			}
		}

		var stock_id = logistics_mode.StockId
		mt.Println("物流方式对应的仓库数据库id", stock_id)
		mode_result.cost_amount = 0
		stock_has_product_and_cost_amount, found_stock := stock_id_has_products_map[stock_id];
		if found_stock == false || stock_has_product_and_cost_amount.has_product == false { // 如果物流方式对应的仓库，没有全部的商品数量的库存，就继续下一个物流方式
			mt.Println("物流方式对应的仓库没有订单对应的库存")
			//mt.Println("匹配单个物流方式结束")
			mt.Println("继续进行物流方式费用的计算")
			mode_result.has_stock = false

			// todo 是否要 continue?
		} else {
			mt.Println("物流方式对应的仓库存在订单对应的库存，澳仓成本价为", stock_has_product_and_cost_amount.cost_amount)
		}

		// 不管有没有库存都记录成本价
		mode_result.cost_amount = stock_has_product_and_cost_amount.cost_amount

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

				if math.Abs(mid_length-mode_result.length) < core.FLOAT_MIN_DIFF {
					mode_result.length = mode_result.length * logistics_setting.Ppqzbzl
					mt.Println("打折后的中边为", mode_result.length)
				} else if math.Abs(mid_length-mode_result.width) < core.FLOAT_MIN_DIFF {
					mode_result.width = mode_result.width * logistics_setting.Ppqzbzl
					mt.Println("打折后的中边为", mode_result.width)
				} else if math.Abs(mid_length-mode_result.height) < core.FLOAT_MIN_DIFF {
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
				mode_result.sfkltjz = 1
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
				mode_result.sfkltjz = -1
				mt.Println("是否考虑体积重", "否")
			}
		}

		var extra_fee_in_mode float64 = 0

		// 以下是旧的不考虑 toll 新方式的原来的匹配方式

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

		// 超重 和 超长 费用中， 取一个， 取大的一个
		mt.Println("物流方式超重费用", extra_weight_fee_in_mode)
		mt.Println("物流方式超长费用", extra_length_fee_in_mode)
		extra_fee_in_mode = math.Max(extra_weight_fee_in_mode, extra_length_fee_in_mode)
		mt.Println("物流方式超重和超长，取大的一个的结果", extra_fee_in_mode)

		/* 以下是新的考虑 toll 的特殊方式，结果先还是按原来的来

		// 不是 TOLL 物流方式， TOLL的 M 费用需要特殊处理
		if strings.Contains(strings.ToUpper(logistics_mode.Name), "TOLL") == false {
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

			// 超重 和 超长 费用中， 取一个， 取大的一个
			mt.Println("物流方式超重费用", extra_weight_fee_in_mode)
			mt.Println("物流方式超长费用", extra_length_fee_in_mode)
			extra_fee_in_mode = math.Max(extra_weight_fee_in_mode, extra_length_fee_in_mode)
			mt.Println("物流方式超重和超长，取大的一个的结果", extra_fee_in_mode)
		} else { // 是 TOLL 物流方式， TOLL的 M 费用需要特殊处理

			// 重量是否有限制
			if logistics_mode.Zlsfyxz == 1 { //如果有限制
				mt.Println("重量是否有限制", "是")
			} else {
				mt.Println("重量是否有限制", "否")
			}

			// 长度是否有限制
			if logistics_mode.Zcbsfyxz == 1 { // 如果有限制
				mt.Println("长度是否有限制", "是")
			} else {
				mt.Println("长度是否有限制", "否")
			}

			var weight_not_allowed = 0
			var length_not_allowed = 0

			for _, order_product_info := range p_order_logistics_data.order_product_info_list {
				mt.Println("当前商品id", order_product_info.ProductId)

				// 1. 考虑超重

				// 单个SKU，商品超重费用，可能会有多个数量
				var product_extra_weight_fee float64 = 0

				var product_actual_weight = order_product_info.Weight

				var product_weight = product_actual_weight // 先初始化为实际重量

				// 体积重问题
				if logistics_mode.Sfkltjz == 1 { // 是否考虑体积重
					var product_dim_weight float64 = 0

					if logistics_mode.Tjzgs > 0 {
						product_dim_weight = order_product_info.Length * order_product_info.Width * order_product_info.Height / logistics_mode.Tjzgs
					}

					if logistics_mode.Tjzxs > 0 { // 如果 体积重系数 有设置的值
						product_dim_weight = product_dim_weight / logistics_mode.Tjzxs // 这里是除以体积重系数
					}

					if logistics_mode.Tjzysjzqx == 1 { // 取大
						product_weight = math.Max(product_actual_weight, product_dim_weight)
					} else if logistics_mode.Tjzysjzqx == 2 { // 取小
						product_weight = math.Min(product_actual_weight, product_dim_weight)
					}
				}

				var extra_weight = product_weight - logistics_mode.Zlbncg
				if extra_weight > 0 {
					// 查看是否允许超重
					if logistics_mode.Sfyxcz == 0 { // 不允许超重
						weight_not_allowed = 1
						continue // 不允许超重，需要排除当前物流方式
					}

					product_extra_weight_fee = logistics_mode.Czjsyf // 只算一个超重的费用 M 费， 不考虑步长
					product_extra_weight_fee = product_extra_weight_fee * float64(order_product_info.Qty)
				}

				// 2. 考虑超长

				var product_extra_length_fee float64 = 0

				// 取得商品的最长边
				var product_max_length = math.Max(math.Max(order_product_info.Length, order_product_info.Width), order_product_info.Height)
				var extra_length = product_max_length - logistics_mode.Zcbxzcd
				if extra_length > 0 { // 长度不能超过，如果最长边有限制，物流方式里面一般都有最长边超长限制
					if logistics_mode.Sfyxcc == 0 { // 如果不允许超长
						length_not_allowed = 1
						continue
					}

					// 允许超长，计算超长费用
					product_extra_length_fee = logistics_mode.Ccjsyf // 只算一个超长的费用 M 费， 不考虑步长
					product_extra_length_fee = product_extra_length_fee * float64(order_product_info.Qty)
				}

				var product_extra_fee = math.Max(product_extra_weight_fee, product_extra_length_fee) // 两个钟选择费用大的那个
				mt.Println("当前商品物流方式综合后的超重超长的费用为：", product_extra_fee)

				extra_fee_in_mode = extra_fee_in_mode + product_extra_fee // 累加
			}

			if weight_not_allowed > 0 || length_not_allowed > 0 {
				mt.Println("当前商品物流方式不允许超长超重，忽略当前物流方式。")
				continue
			}

			mt.Println("当前物流方式下各商品综合后的超重超长的费用为：", extra_fee_in_mode)
		}
		*/

		mode_result.fee_rate_id = 0
		mode_result.fee_rate = ddpz_models.LogisticsFeeRateInfo{}
		var total_fee_in_fee_rate float64 = 10000000 //非超重部分的费用，初始为一个很大的不可能的数字， 澳币
		if found_rate1 {
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
							revel.ERROR.Println("续重阶梯设置有误，不能计算超重费用！", logistics_fee_rate.Xzjt)
							mt.PrintlnIndent(1, "续重阶梯设置有误，不能计算超重费用！", logistics_fee_rate.Xzjt)
						}
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
						mode_result.fee_rate = logistics_fee_rate
					} else {
						mt.PrintlnIndent(1, "之前匹配的费率，比当前费率更优，弃用当前费率计算结果")
					}
					mt.PrintlnIndent(1, "匹配单个费率结束")
					mt.Println()
				}
			} // for _, logistics_fee_rate := range logistics_fee_rate_list {
		}

		mode_result.fee_rate2_id = 0
		mode_result.fee_rate2 = ddpz_models.LogisticsFeeRate2Info{}
		if found_rate2 {
			for _, logistics_fee_rate2 := range logistics_fee_rate2_list {

				// 邮编要对上，忽略开头的 0
				var postcode_in_address = strings.TrimPrefix(strings.TrimSpace(p_order_logistics_data.order_shipping_address.PostalCode), "0")
				var postcode_in_fee_rate = strings.TrimPrefix(strings.TrimSpace(logistics_fee_rate2.Postcode), "0")
				if postcode_in_address == postcode_in_fee_rate {
					mt.Println()
					mt.PrintlnIndent(1, "匹配单个费率开始")

					mt.PrintlnIndent(1, "订单收货地址的邮编", p_order_logistics_data.order_shipping_address.PostalCode)
					mt.PrintlnIndent(1, "当前物流方式下，符合邮编的费率数据库 id", logistics_fee_rate2.Id)

					var fee float64 = 0

					if 0 < mode_result.weight && mode_result.weight <= 0.5 { // 0-500g
						mt.PrintlnIndent(1, "0-500g 费率匹配成功！")
						fee = logistics_fee_rate2.Fee01
					} else if 0.5 < mode_result.weight && mode_result.weight <= 1 { // 500g-1kg
						mt.PrintlnIndent(1, "500g-1kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee02
					} else if 1 < mode_result.weight && mode_result.weight <= 2 { // 1kg-2kg
						mt.PrintlnIndent(1, "1kg-2kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee03
					} else if 2 < mode_result.weight && mode_result.weight <= 3 { // 2kg-3kg
						mt.PrintlnIndent(1, "2kg-3kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee04
					} else if 3 < mode_result.weight && mode_result.weight <= 4 { // 3kg-4kg
						mt.PrintlnIndent(1, "3kg-4kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee05
					} else if 4 < mode_result.weight && mode_result.weight <= 5 { // 4kg-5kg
						mt.PrintlnIndent(1, "4kg-5kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee06
					} else if 5 < mode_result.weight && mode_result.weight <= 7 { // 5kg-7kg
						mt.PrintlnIndent(1, "5kg-7kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee07
					} else if 7 < mode_result.weight && mode_result.weight <= 10 { // 7kg-10kg
						mt.PrintlnIndent(1, "7kg-10kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee08
					} else if 10 < mode_result.weight && mode_result.weight <= 15 { // 10kg-15kg
						mt.PrintlnIndent(1, "10kg-15kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee09
					} else if 15 < mode_result.weight && mode_result.weight <= 22 { // 15kg-22kg
						mt.PrintlnIndent(1, "15kg-22kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee10
					} else if 22 < mode_result.weight && mode_result.weight <= 30 { // 22kg-30kg
						mt.PrintlnIndent(1, "22kg-30kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee11
					} else if 30 < mode_result.weight && mode_result.weight <= 50 { // 30kg-50kg
						mt.PrintlnIndent(1, "30kg-50kg 费率匹配成功！")
						fee = logistics_fee_rate2.Fee12
					}

					if fee == 0 {
						mt.PrintlnIndent(1, "当前重量该费率不支持！")
						continue
					}

					// 偏远地区
					if logistics_fee_rate2.Sfpydq == 1 {
						mt.PrintlnIndent(1, "是否偏远地区", "是")
						mt.PrintlnIndent(1, "偏远地区加收运费", logistics_fee_rate2.Pydqjsyf)
						fee = fee + logistics_fee_rate2.Pydqjsyf
						mt.PrintlnIndent(1, "考虑偏远地区后的费用", fee)
					} else {
						mt.PrintlnIndent(1, "是否偏远地区", "否")
					}

					mt.PrintlnIndent(1, "最终当前费率的费用", fee)
					if fee < total_fee_in_fee_rate { // 取得更优的费率方案
						mt.PrintlnIndent(1, "之前没有匹配的费率，或当前费率比之前匹配的费率更优，暂时选用当前费率计算结果")
						total_fee_in_fee_rate = fee
						mode_result.fee_rate2_id = logistics_fee_rate2.Id
						mode_result.fee_rate2 = logistics_fee_rate2
					} else {
						mt.PrintlnIndent(1, "之前匹配的费率，比当前费率更优，弃用当前费率计算结果")
					}
					mt.PrintlnIndent(1, "匹配单个费率结束")
					mt.Println()
				}
			} // for _, logistics_fee_rate2 := range logistics_fee_rate2_list {
		}

		if found_rate1 {
			if mode_result.fee_rate_id == 0 { //没有费率匹配成功, 比如说，还没有开通该地区的物流
				mt.Println("该物流方式下没有匹配的物流费率，弃用当前物流方式")
				mt.Println("匹配单个物流方式结束")

				if mt.is_trace == true {
					log_logistics_match_debug(mode_result, "没有匹配的物流费率", mt.String(), db_session)
				}

				continue // 排除该物流方式
			}
		}

		if found_rate2 {
			if mode_result.fee_rate2_id == 0 { //没有费率匹配成功, 比如说，还没有开通该地区的物流
				mt.Println("该物流方式下没有匹配的物流费率，弃用当前物流方式")
				mt.Println("匹配单个物流方式结束")

				if mt.is_trace == true {
					log_logistics_match_debug(mode_result, "没有匹配的物流费率", mt.String(), db_session)
				}

				continue // 排除该物流方式
			}
		}

		mode_result.total_fee = total_fee_in_fee_rate + extra_fee_in_mode
		mt.Println("物流方式中物流费率费用加上物流方式超重超长的结果", mode_result.total_fee)
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

		var total_fee_diff = lowest_total_fee - mode_result.total_fee // 就是当前模式比之前的省了多少钱

		if p_ol.is_now_fav_stock_time() {
			if mode_result.has_stock == true {
				// 三个区间 final_total_amount-mode_result.total_amount，后来改用 fee了，不考虑 库存的cost_price，因为没有库存的也参与筛选
				// 负无穷大 到 -p_ol.logistics_setting.FeeMaxDiff， -p_ol.logistics_setting.FeeMaxDiff到+p_ol.logistics_setting.FeeMaxDiff, +p_ol.logistics_setting.FeeMaxDiff 到 正无穷大
				if total_fee_diff > p_ol.logistics_setting.FeeMaxDiff { // 取得更优的模式方案，挑选最优解时，把商品的到仓成本也考虑进去   +p_ol.logistics_setting.FeeMaxDiff 到 正无穷大
					mt.Println("之前没有匹配的物流方式，或当前物流方式比之前匹配的物流方式节省费用超过最大允许偏差，暂时选用当前物流方式计算结果")
					final_result = mode_result
				} else if math.Abs(total_fee_diff) <= p_ol.logistics_setting.FeeMaxDiff { // 运费偏差在最大允许范围内，其实还是算上了澳仓成本的  -p_ol.logistics_setting.FeeMaxDiff到+p_ol.logistics_setting.FeeMaxDiff,
					//取优先级高的仓库

					if mode_result.stock.Priority > final_result.stock.Priority { // 当前的仓库的优先级要高
						mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，且当前仓库是优先级更高，选用当前物流方式计算结果, 记录次佳仓备选")
						final_result = mode_result
					} else { //
						mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，但当前仓库是优先级并不更高，弃用当前物流方式计算结果，记录次佳仓备选。")
					}
				} else { // 在运费偏差范围之外， 负无穷大 到 -p_ol.logistics_setting.FeeMaxDiff
					mt.Println("之前匹配的物流，比当前物流方式更优, 且超过了允许偏差范围，弃用当前物流方式计算结果，记录次佳仓备选")
				}

				match_result_has_stock_list = append(match_result_has_stock_list, mode_result)
			} else {
				//没有库存
				if total_fee_diff > p_ol.logistics_setting.FeeMaxDiff {
					mt.Println("之前没有匹配的物流方式，或当前物流方式比之前匹配的物流方式节省费用超过最大允许偏差，暂时选用当前物流方式计算结果，仓库无货，不记录次佳仓备选。")
					final_result = mode_result
				} else if math.Abs(total_fee_diff) <= p_ol.logistics_setting.FeeMaxDiff {
					//取优先级高的仓库
					if mode_result.stock.Priority > final_result.stock.Priority { // 当前的仓库的优先级要高
						mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，且当前仓库是优先级更高，选用当前物流方式计算结果，仓库无货，不记录次佳仓备选")
						final_result = mode_result
					} else { //
						mt.Println("当前匹配的物流方式费用和之前最优解在允许偏差范围内，但当前仓库是优先级并不更高，弃用当前物流方式计算结果，仓库无货，不记录次佳仓备选。")
					}
				} else { // 在运费偏差范围之外， 负无穷大 到 -p_ol.logistics_setting.FeeMaxDiff
					mt.Println("之前匹配的物流，比当前物流方式更优, 且超过了允许偏差范围，弃用当前物流方式计算结果，仓库无货。")
				}
			}
		} else {
			if mode_result.has_stock == true {
				if total_fee_diff > 0 {
					mt.Println("之前没有匹配的物流方式，或当前物流方式比之前匹配的物流方式节省费用，暂时选用当前物流方式计算结果")
					final_result = mode_result
				} else {
					mt.Println("之前匹配的物流，比当前物流方式更优,弃用当前物流方式计算结果，记录次佳仓备选")
				}

				match_result_has_stock_list = append(match_result_has_stock_list, mode_result)
			} else {
				if total_fee_diff > 0 {
					mt.Println("之前没有匹配的物流方式，或当前物流方式比之前匹配的物流方式节省费用，暂时选用当前物流方式计算结果，仓库无货，不记录次佳仓备选。")
					final_result = mode_result
				} else {
					mt.Println("之前匹配的物流，比当前物流方式更优, 弃用当前物流方式计算结果，仓库无货。")
				}
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
	}                                                                 // for _, logistics_mode := range logistics_mode_list {

	var order_logistics ddlc_models.OrderLogisticsInfo

	// 这个要放在 匹配失败的判断的前面，否则信封也会返回匹配失败
	if is_envelope_hold == true {
		var match_return match_return_t
		match_return.success = true
		match_return.order_logistics = ddlc_models.OrderLogisticsInfo{} // 前面已经保存 order_Logistics， 这里就随便给个值就行
		match_return.is_envelope_hold = true
		match_return.has_stock = true

		return match_return
	}

	if final_result.mode_id == 0 || (final_result.fee_rate_id == 0 && final_result.fee_rate2_id == 0) { // 匹配失败
		var match_return match_return_t
		match_return.success = false
		match_return.order_logistics = order_logistics
		match_return.is_envelope_hold = false
		match_return.has_stock = false

		return match_return
	}

	order_logistics = save_order_logistics(p_ol, p_order, p_order_logistics_data, final_result, match_option.Unlock_express, p_ol.LoginUser, is_previous_hold, db_session)

	if final_result.has_stock == false {
		if len(match_result_has_stock_list) != 0 {
			// 最后匹配的仓库没有货，但是存在匹配的有货的仓库
			p_ol.update_order_to_hold_type(p_order, ddpz_models.ORDER_HOLD_TYPE_OOS_BEST_STOCK, db_session)
		} else {
			p_ol.update_order_to_hold_type(p_order, ddpz_models.ORDER_HOLD_TYPE_OOS_ALL_STOCK, db_session) // 缺货的Hold，后面的 Hold check 会把它进一步更新成 在途才有货等等。
		}

		var match_return match_return_t
		match_return.success = true
		match_return.order_logistics = order_logistics
		match_return.is_envelope_hold = false
		match_return.has_stock = true

		return match_return
	} else {
		var product_info_list = p_order_logistics_data.order_product_info_list
		for _, product_info := range product_info_list {
			p_ol.pending_keeper.add(order_logistics.StockId, product_info.ProductId, product_info.Qty)
		}

		var match_return match_return_t
		match_return.success = true
		match_return.order_logistics = order_logistics
		match_return.is_envelope_hold = false
		match_return.has_stock = false

		return match_return
	}
	// 多个商品时，考虑拆包的方式

	// 多个商品时，能不能分仓库发货的？
}

func (p_ol *OrderLogisticsService) is_now_fav_stock_time() bool {
	var time_now = time.Now()
	var fav_stock_time_range_list = p_ol.fav_stock_time_range_list

	var weekday_num = int64(time_now.Weekday())
	if weekday_num == 0 {
		weekday_num = 7
	}

	var time_total_minutes = int64(time_now.Hour()*60 + time_now.Minute())

	if len(fav_stock_time_range_list) == 0 {
		return false
	}

	for _, time_range := range fav_stock_time_range_list {
		if time_total_minutes >= time_range.BeginHour*60+time_range.BeginMinute && time_total_minutes <= time_range.EndHour*60+time_range.EndMinute {
			return true
		}
	}

	return false
}
