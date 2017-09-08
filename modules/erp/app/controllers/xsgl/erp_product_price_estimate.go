package controllers

import (
	"github.com/revel/revel"
	"encoding/json"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpProductPriceEstimate struct {
	*revel.Controller
}

func (c ErpProductPriceEstimate) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_PRICE_ESTIMATE, c.Session)

	return c.RenderTemplate("erp/xsgl/product_price_estimate/product_price_estimate_index.html")
}

func (c ErpProductPriceEstimate) Post() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_PRICE_ESTIMATE, c.Session)

	var product_price_estimate ProductPriceEstimate
	err := json.Unmarshal(c.GetRequestBody(), &product_price_estimate)
	core.HandleError(err)

	// do estimate

	return c.RenderJson(map[string]interface{}{
		"success":                true,
		"message":                "商品售价预估成功！",
		"product_price_estimate": product_price_estimate,
	})
}

type ProductPriceEstimate struct {
	Length           float64 `json:"length"`             // 长 cm
	Width            float64 `json:"width"`              // 宽 cm
	Height           float64 `json:"height"`             // 高 cm
	DimWeight        float64 `json:"dim_weight"`         // 体积重
	ActualWeight     float64 `json:"actual_weight"`      // 实际重量 kg
	CnInventoryPrice float64 `json:"cn_inventory_price"` // 中国仓采购价 RMB?
	AuInventoryPrice float64 `json:"au_inventory_price"` // 预估澳仓成本 RMB?
	Qty              float64 `json:"qty"`                // 数量
	LogisticsFee     float64 `json:"logistics_fee"`      // 预估物流费（澳币）
	Profit           float64 `json:"profit"`             // 毛利 （澳币）
	ProfitRate       float64 `json:"profit_rate"`        // 毛利率
	Price            float64 `json:"price"`              // 预估售价（澳币）
}
