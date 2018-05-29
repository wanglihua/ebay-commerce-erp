package product_logistics_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/kcgl_models"
	"strings"
	"matrix/modules/erp/service/pick_fee_service"
	"math"
	"github.com/revel/revel"
	"matrix/modules/erp/service/product_sale_setting_service"
	"errors"
	"fmt"
	"matrix/modules/erp/models/ddpz_models/logistics_fee_type"
	"matrix/modules/erp/service/logistics_mode_service"
)

func Match_product_postage(product_id int64, postcode string, db_session *xorm.Session) (err error, fee float64) {

	postcode = strings.TrimSpace(postcode)

	/*
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("计算商品地区运费出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()
	*/

	var product kcgl_models.ProductInfo
	has, err := db_session.Where("id = ?", product_id).Get(&product)
	if err != nil {
		return err, 0
	}

	if has == false {
		return errors.New("product not exist!"), 0
	}

	//var logistics_setting = logistics_setting_service.Get_logistics_setting(db_session)

	logistics_mode_list, err := logistics_mode_service.Get_logistics_mode_enabled(db_session)
	if err != nil {
		return err, 0
	}

	var logistics_fee_rate_list = make([]ddpz_models.LogisticsFeeRateInfo, 0)
	err = db_session.Where(fmt.Sprintf("postcode = '%s' OR postcode = '%s'", postcode, strings.TrimPrefix(postcode, "0"))).Find(&logistics_fee_rate_list)
	if err != nil {
		return err, 0
	}

	var logistics_fee_rate2_list = make([]ddpz_models.LogisticsFeeRate2Info, 0)
	err = db_session.Where(fmt.Sprintf("postcode = '%s' OR postcode = '%s'", postcode, strings.TrimPrefix(postcode, "0"))).Find(&logistics_fee_rate2_list)
	if err != nil {
		return err, 0
	}

	// postcode = strings.TrimPrefix(strings.TrimSpace(postcode), "0")
	postcode = strings.TrimPrefix(postcode, "0")

	var pick_fee_setting = pick_fee_service.Get_pick_fee_setting(db_session)

	var final_result match_result_t
	final_result.total_fee = math.MaxFloat64 // 初始化为一个很大的不可能的值，澳币

	var final_pick_fee = pick_fee_service.Get_product_pick_fee(product.Weight, pick_fee_setting) // 拣货费用的是不打折的重量

	for _, logistics_mode := range logistics_mode_list {

		var mode_result match_result_t

		mode_result.length, mode_result.width, mode_result.height, mode_result.actual_weight = product.Length, product.Width, product.Height, product.Weight

		/*
		这里不需要打折
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
		*/

		// 体积重是利用体积重公式计算得到的结果
		// 只有这个参数大于零的时候才有意义
		var product_dim_weight float64 = 0
		if logistics_mode.Tjzgs > 0 {
			product_dim_weight = mode_result.length * mode_result.width * mode_result.height / logistics_mode.Tjzgs
		}

		mode_result.weight = mode_result.actual_weight // 用户存储最终的计费重量， 初始为实际重量。如果不考虑体积重的话，就用实际重量来算重量的

		// 体积重与实际重是否有比值要求
		if logistics_mode.Zlsfybzyq == 1 { //如果有比值要求
			// #35 18kg的 比值是 0.5， 实际重量限制是18kg，体积重限制是9kg，在这个限制范围内，他们的体积重和实际重量的比值他们是不关心的。
			var actual_weight_limit = logistics_mode.Zlbncg // 实际重量的限重

			var dim_weight_limit = logistics_mode.Zlbncg * logistics_mode.Tjzsjzbz // 通过体积重和实际重量的比值计算 体积重的限重

			if product_dim_weight > dim_weight_limit || mode_result.actual_weight > actual_weight_limit {
				continue
			}
		} else {
			//如果 体积重与实际重 “有” 比值要求， 是否考虑体积重，体积重与实际重取向都是没有意义的，
			//所以只在 体积重与实际重 “没有” 比值要求的时候，才考虑这两个量

			// 如果考虑体积重，就把计算好的体积重，用 体积重系数 打折
			if logistics_mode.Sfkltjz == 1 { // 是否考虑体积重
				//mode_result.sfkltjz = 1

				if logistics_mode.Tjzxs > 0 { // 如果 体积重系数 有设置的值
					product_dim_weight = product_dim_weight / logistics_mode.Tjzxs // 这里是除以体积重系数
				} else {
				}

				// 体积重与实际重取向， 如果考虑体积重，这个系数一般都会设置的？
				if logistics_mode.Tjzysjzqx == 1 { // 取大
					mode_result.weight = math.Max(product_dim_weight, mode_result.actual_weight)
				} else if logistics_mode.Tjzysjzqx == 2 { // 取小
					mode_result.weight = math.Min(product_dim_weight, mode_result.actual_weight)
				}

			} else {
				//mode_result.sfkltjz = -1
			}
		}

		// 接下来计算都是针对 final_weight_in_mode 来做，不管前面有没有考虑体积重, 后来就是 product_weight

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

		var total_fee_in_fee_rate float64 = math.MaxFloat64 //非超重部分的费用，初始为一个很大的不可能的数字， 澳币

		if logistics_mode.FeeType == logistics_fee_type.FEE_TYPE1 {
			for _, logistics_fee_rate := range logistics_fee_rate_list {

				if logistics_fee_rate.ModeId != logistics_mode.Id {
					continue
				}

				// 邮编要对上，忽略开头的 0
				var postcode_in_fee_rate = strings.TrimPrefix(strings.TrimSpace(logistics_fee_rate.Postcode), "0")

				// postcode 最前面可能的 0，已经在前面的代码里面去掉了，不能去第二次
				if postcode == postcode_in_fee_rate {
					var fee float64 = logistics_fee_rate.Qbj //起步价
					if mode_result.weight-logistics_fee_rate.Qbjzl > 0 { //如果超过了起步重量，就算超过部分的费率；不超过，就不用算了
						var extra_weight_in_fee_rate = mode_result.weight - logistics_fee_rate.Qbjzl
						if logistics_fee_rate.Xzjt > 0 {
							var xz_fee = math.Ceil(extra_weight_in_fee_rate/logistics_fee_rate.Xzjt) * logistics_fee_rate.Xzjtsf
							fee = fee + xz_fee
						} else {
							revel.ERROR.Println(fmt.Sprintf("费率 %d 续重阶梯设置有误，不能计算超重费用！", logistics_fee_rate.Id), logistics_fee_rate.Xzjt)
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
		} else if logistics_mode.FeeType == logistics_fee_type.FEE_TYPE2 {
			for _, logistics_fee_rate2 := range logistics_fee_rate2_list {

				var postcode_in_fee_rate = strings.TrimPrefix(strings.TrimSpace(logistics_fee_rate2.Postcode), "0")
				if postcode == postcode_in_fee_rate {

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
					} else {
					}

					if fee < total_fee_in_fee_rate { // 取得更优的费率方案
						total_fee_in_fee_rate = fee
						mode_result.fee_rate2_id = logistics_fee_rate2.Id
						mode_result.fee_rate2 = logistics_fee_rate2
					}
				}
			} // for _, logistics_fee_rate2 := range logistics_fee_rate2_list {
		}

		if (logistics_mode.FeeType == logistics_fee_type.FEE_TYPE1 && mode_result.fee_rate_id == 0) ||
			(logistics_mode.FeeType == logistics_fee_type.FEE_TYPE2 && mode_result.fee_rate2_id == 0) { //没有费率匹配成功, 比如说，还没有开通该地区的物流
			continue // 排除该物流方式
		}

		// 超重 和 超长 费用中， 取一个， 取大的一个
		var extra_fee_in_mode = math.Max(extra_weight_fee_in_mode, extra_length_fee_in_mode)
		mode_result.total_fee = total_fee_in_fee_rate + extra_fee_in_mode

		/*
		这里不需要打折
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
		*/

		if mode_result.total_fee < final_result.total_fee {
			final_result = mode_result
		}

	} // for _, logistics_mode := range logistics_mode_list {

	if final_result.fee_rate_id == 0 {
		return nil, -1 // 如果返回的费用小于零就表示，没有匹配的物流方式和费率，不能送达
	}

	var postage_cal_weight = product.PostageCalWeight

	if postage_cal_weight == 0 {
		var product_sale_setting = product_sale_setting_service.Get_product_sale_setting(db_session)
		postage_cal_weight = product_sale_setting.PostageCalWeight
	}

	var total_amount = (final_result.total_fee + final_pick_fee) * postage_cal_weight

	return nil, total_amount
}

type match_result_t struct {
	fee_rate_id   int64
	fee_rate      ddpz_models.LogisticsFeeRateInfo
	fee_rate2_id  int64
	fee_rate2     ddpz_models.LogisticsFeeRate2Info
	length        float64
	width         float64
	height        float64
	actual_weight float64
	dim_weight    float64
	weight        float64 //计费的重量，实际重量或者体积重
	total_fee     float64
}
