package controllers

import (
	"github.com/revel/revel"
	"encoding/json"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/service/product_sale_setting_service"
	"fmt"
	"matrix/modules/erp/models/kcgl_models"
	"strings"
	"matrix/modules/erp/service/sea_container_cost_service"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/inventory_service"
)

type ErpProductPriceEstimate struct {
	*revel.Controller
}

func (c ErpProductPriceEstimate) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_PRICE_ESTIMATE, c.Session)

	return c.RenderTemplate("erp/xsgl/product_price_estimate/product_price_estimate_index.html")
}

func (c ErpProductPriceEstimate) GetSkuInfo() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_PRICE_ESTIMATE, c.Session)

	var db_session = c.DbSession

	var product_price_estimate ProductPriceEstimate
	err := json.Unmarshal(c.GetRequestBody(), &product_price_estimate)
	core.HandleError(err)

	var product kcgl_models.ProductInfo
	has, err := db_session.Where("sku = ?", product_price_estimate.Sku).Get(&product)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success":                false,
			"message":                "SKU信息不存在！",
			"product_price_estimate": nil,
		})
	}

	product_price_estimate.Length = product.Length
	product_price_estimate.Width = product.Width
	product_price_estimate.Height = product.Height
	product_price_estimate.ActualWeight = product.Weight
	product_price_estimate.CnInventoryPrice = inventory_service.GetProductCnStockCostPrice(product.Id, db_session)

	return c.RenderJson(map[string]interface{}{
		"success":                true,
		"message":                "获取SKU信息成功！",
		"product_price_estimate": product_price_estimate,
	})
}

func (c ErpProductPriceEstimate) Estimate() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_PRICE_ESTIMATE, c.Session)

	var db_session = c.DbSession

	var product_price_estimate ProductPriceEstimate
	err := json.Unmarshal(c.GetRequestBody(), &product_price_estimate)
	core.HandleError(err)

	product_price_estimate.Sku = strings.TrimSpace(product_price_estimate.Sku)
	if product_price_estimate.Sku != "" {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where("sku = ?", product_price_estimate.Sku).Get(&product)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(map[string]interface{}{
				"success":                false,
				"message":                "SKU信息不存在！",
				"product_price_estimate": product_price_estimate,
			})
		}

		/*
			然后对应历史订单100单及以上的则跳出该SKU的实际平均运费，没有满100单的则提取wollongzone系统体积重、实际重 最接近的5个SKU的平均运费/平均重量 得出每公斤费率。
		 */

		var product_order_item_fee_avg struct {
			ItemCount int64   `xorm:"'item_count'"`
			AvgFee    float64 `xorm:"'avg_fee'"`
		}

		/*
	SELECT count(*) AS item_count, sum(ol.actual_fee) / sum(if(ol.actual_fee = 0, 0, oi.quantity)) AS avg_fee
	FROM erp_order_item oi
	   INNER JOIN erp_order_logistics ol ON oi.order_id = ol.order_id
	WHERE oi.quantity = 1
		 AND (SELECT count(oi_inner.id)
			  FROM erp_order_item oi_inner
			  WHERE oi_inner.order_id = oi.order_id) = 1
		 AND oi.item_product_id = 17578;
		 */

		var sql_get_product_order_item_fee_avg = `SELECT count(*) AS item_count, sum(ol.actual_fee) / sum(if(ol.actual_fee = 0, 0, oi.quantity)) AS avg_fee
FROM erp_order_item oi
   INNER JOIN erp_order_logistics ol ON oi.order_id = ol.order_id
WHERE oi.quantity = 1
	 AND (SELECT count(oi_inner.id)
		  FROM erp_order_item oi_inner
		  WHERE oi_inner.order_id = oi.order_id) = 1
	 AND oi.item_product_id = %d;`
		sql_get_product_order_item_fee_avg = fmt.Sprintf(sql_get_product_order_item_fee_avg, product.Id)

		_, err = db_session.Sql(sql_get_product_order_item_fee_avg).Get(&product_order_item_fee_avg)
		core.HandleError(err)

		if product_order_item_fee_avg.ItemCount >= 100 { // 如果有100个符合条件的物流记录的，就显示这100个的平均值
			product_price_estimate.LogisticsFee = fmt.Sprintf("%0.2f", product_order_item_fee_avg.AvgFee)

			return c.RenderJson(map[string]interface{}{
				"success":                true,
				"message":                "获取商品历史运费成功！",
				"product_price_estimate": product_price_estimate,
				"product_list":           nil,
			})
		}
	}

	// check

	// do estimate

	var product_sale_setting = product_sale_setting_service.Get_product_sale_setting(db_session)

	var dim_weight_formula = product_sale_setting.DimWeightFormula

	if dim_weight_formula == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "体积重公式系数不能为零设置有误！",
		})
	}

	if product_price_estimate.Qty == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "数量不能为零！",
		})
	}

	var dim_weight = product_price_estimate.Length * product_price_estimate.Width * product_price_estimate.Height / product_sale_setting.DimWeightFormula
	product_price_estimate.DimWeight = fmt.Sprintf("%0.2f", dim_weight)

	if dim_weight == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "录入的长宽高不能为零！",
		})
	}

	var actual_weight = product_price_estimate.ActualWeight
	if actual_weight == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "录入的实际重量不能为零！",
		})
	}

	/*
SELECT p.*
FROM erp_product p
ORDER BY abs(p.length * p.width * p.height / 40000 - 100) + abs(p.weight - 100) ASC
LIMIT 5;


SELECT sum(ol.actual_fee) / sum(if(ol.actual_fee = 0, 0, oi.quantity)) AS avg_fee
FROM erp_order_item oi INNER JOIN erp_order_logistics ol ON oi.order_id = ol.order_id
WHERE oi.quantity = 1
  AND (SELECT count(oi_inner.id)
	   FROM erp_order_item oi_inner
	   WHERE oi_inner.order_id = oi.order_id) = 1
  AND oi.item_product_id IN (17578, 23332, 17632, 17583, 23539);
 */

	// 获得最接近的 5 个商品的列表
	var sql_get_five_product = `SELECT p.*
FROM erp_product p
ORDER BY abs(p.length * p.width * p.height / 40000 - %f) + abs(p.weight - %f) ASC
LIMIT 5;`

	sql_get_five_product = fmt.Sprintf(sql_get_five_product, dim_weight, actual_weight)
	var product_list = make([]kcgl_models.ProductInfo, 0)
	core.HandleError(db_session.Sql(sql_get_five_product).Find(&product_list))

	var product_actual_weight_sum float64 = 0
	var product_count float64 = 0
	for _, product := range product_list {
		product_actual_weight_sum = product_actual_weight_sum + product.Weight
		product_count += 1
	}
	var product_avg_actual_weight = product_actual_weight_sum / product_count

	var product_id_str_list = make([]string, 0)
	for _, product := range product_list {
		product_id_str_list = append(product_id_str_list, fmt.Sprint(product.Id))
	}

	// 获得接近的 5 个商品的平均运费
	var sql_get_avg_fee = `SELECT sum(ol.actual_fee) / sum(if(ol.actual_fee = 0, 0, oi.quantity)) AS avg_fee
FROM erp_order_item oi INNER JOIN erp_order_logistics ol ON oi.order_id = ol.order_id
WHERE oi.quantity = 1
  AND (SELECT count(oi_inner.id)
	   FROM erp_order_item oi_inner
	   WHERE oi_inner.order_id = oi.order_id) = 1
  AND oi.item_product_id IN (%s);`
	sql_get_avg_fee = fmt.Sprintf(sql_get_avg_fee, strings.Join(product_id_str_list, ","))

	var avg_fee struct {
		AvgFee float64 `xorm:"'avg_fee'"`
	}
	_, err = db_session.Sql(sql_get_avg_fee).Get(&avg_fee)
	core.HandleError(err)

	// 预估

	var logistics_fee = avg_fee.AvgFee / product_avg_actual_weight * product_price_estimate.ActualWeight * product_price_estimate.Qty
	product_price_estimate.LogisticsFee = fmt.Sprintf("%0.2f", logistics_fee)

	var sea_container_cost = sea_container_cost_service.Get_sea_container_cost(db_session)

	// 计算商品体积，计算平摊海运成本
	var avg_sea_shipping_cost = float64(sea_container_cost.Cost) * product_price_estimate.Length * product_price_estimate.Width * product_price_estimate.Height / float64(sea_container_cost.Volume)
	var au_inventory_price float64 = product_price_estimate.CnInventoryPrice + avg_sea_shipping_cost

	var aud_rate = exchange_rate_service.Get_aud_exchange_rate(db_session)

	product_price_estimate.AuInventoryPrice = fmt.Sprintf("%0.2f", au_inventory_price)

	var au_inventory_price_aud = au_inventory_price / aud_rate //转成澳币

	var price float64 = (au_inventory_price_aud + logistics_fee) / (1 - product_price_estimate.ProfitRate)
	product_price_estimate.Price = fmt.Sprintf("%0.2f", price)
	product_price_estimate.Profit = fmt.Sprintf("%0.2f", price*product_price_estimate.ProfitRate)

	// 获取各个商品的平均运费，是分立的平均

	/*
SELECT oi.item_product_id AS product_id, sum(ol.actual_fee) / sum(if(ol.actual_fee = 0, 0, oi.quantity)) AS avg_fee
FROM erp_order_item oi INNER JOIN erp_order_logistics ol ON oi.order_id = ol.order_id
WHERE oi.quantity = 1
      AND (SELECT count(oi_inner.id)
           FROM erp_order_item oi_inner
           WHERE oi_inner.order_id = oi.order_id) = 1
      AND oi.item_product_id IN (%s)
GROUP BY item_product_id;
	 */

	var sql_get_product_avg_fee = `SELECT oi.item_product_id AS product_id, sum(ol.actual_fee) / sum(if(ol.actual_fee = 0, 0, oi.quantity)) AS avg_fee
FROM erp_order_item oi INNER JOIN erp_order_logistics ol ON oi.order_id = ol.order_id
WHERE oi.quantity = 1
      AND (SELECT count(oi_inner.id)
           FROM erp_order_item oi_inner
           WHERE oi_inner.order_id = oi.order_id) = 1
      AND oi.item_product_id IN (%s)
GROUP BY item_product_id;`
	sql_get_product_avg_fee = fmt.Sprintf(sql_get_product_avg_fee, strings.Join(product_id_str_list, ","))

	type producct_avg_fee_item_t struct {
		ProductId int64   `xorm:"'product_id'"`
		AvgFee    float64 `xorm:"'avg_fee'"`
	}

	var product_avg_fee_list = make([]producct_avg_fee_item_t, 0)
	core.HandleError(db_session.Sql(sql_get_product_avg_fee).Find(&product_avg_fee_list))

	var product_avg_fee_map = make(map[int64]float64)
	for _, product_avg_fee := range product_avg_fee_list {
		product_avg_fee_map[product_avg_fee.ProductId] = product_avg_fee.AvgFee
	}

	// 计算 5 个商品的体积重
	// 设置 5 个商品的平均运费
	for index, product := range product_list {
		// 计算商品的体积重
		var dim_weight = product.Length * product.Width * product.Height / product_sale_setting.DimWeightFormula
		product_list[index].DimWeight = fmt.Sprintf("%0.2f", dim_weight)

		// 设置商品的平均运费
		product_list[index].AvgFee = fmt.Sprintf("%0.2f", product_avg_fee_map[product.Id])
	}

	return c.RenderJson(map[string]interface{}{
		"success":                true,
		"message":                "根据同类商品，售价预估成功！",
		"product_price_estimate": product_price_estimate,
		"product_list":           product_list,
	})
}

type ProductPriceEstimate struct {
	Sku              string  `json:"sku"`                // sku
	Length           float64 `json:"length"`             // 长 cm
	Width            float64 `json:"width"`              // 宽 cm
	Height           float64 `json:"height"`             // 高 cm
	DimWeight        string  `json:"dim_weight"`         // 体积重
	ActualWeight     float64 `json:"actual_weight"`      // 实际重量 kg
	CnInventoryPrice float64 `json:"cn_inventory_price"` // 中国仓采购价 RMB?
	AuInventoryPrice string  `json:"au_inventory_price"` // 预估澳仓成本 RMB?
	Qty              float64 `json:"qty"`                // 数量
	LogisticsFee     string  `json:"logistics_fee"`      // 预估物流费（澳币）
	Profit           string  `json:"profit"`             // 毛利 （澳币）
	ProfitRate       float64 `json:"profit_rate"`        // 毛利率
	Price            string  `json:"price"`              // 预估售价（澳币）
}
