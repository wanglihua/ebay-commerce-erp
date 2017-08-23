package order_logistics_service

import (
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/ddlc_models"
)

type match_logistics_mode_result struct {
	order_id      int64
	order_code    string
	address_id    int64
	address       ddlc_models.OrderShippingAddressInfo
	company_id    int64
	company       ddpz_models.LogisticsCompanyInfo
	stock_id      int64
	stock         sppz_models.StockInfo
	mode_id       int64
	mode          ddpz_models.LogisticsModeInfo
	fee_rate_id   int64
	fee_rate      ddpz_models.LogisticsFeeRateInfo
	length        float64
	width         float64
	height        float64
	actual_weight float64
	dim_weight    float64
	sfkltjz       int64
	weight        float64 //计算时用的重量体积重或者实际重量
	total_fee     float64
	pick_fee      float64
	cost_amount   float64
	total_amount  float64
	has_stock     bool
}

//排序相关的函数实现
