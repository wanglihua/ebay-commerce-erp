package xsbb_models

import "matrix/core"

type SaleAmountReportMonthInfo struct {
	Id                   int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ReportMonthId        int64     `xorm:"bigint notnull index 'report_month_id' default(0)" json:"report_month_id"`
	ReportMonth          core.Time `xorm:"datetime notnull index 'report_month' default(0)" json:"report_month"`
	PaidAmount           float64   `xorm:"decimal(18,2) notnull 'paid_amount' default(0)" json:"paid_amount"`
	ExpressFee           float64   `xorm:"decimal(18,2) notnull 'express_fee' default(0)" json:"express_fee"`
	CostPrice            float64   `xorm:"decimal(18,2) notnull 'cost_price' default(0)" json:"cost_price"` // RMB ，其他都是 AUD
	EstimateLogisticsFee float64   `xorm:"decimal(18,2) notnull 'estimate_logistics_fee' default(0)" json:"estimate_logistics_fee"`
	ActualLogisticsFee   float64   `xorm:"decimal(18,2) notnull 'actual_logistics_fee' default(0)" json:"actual_logistics_fee"`
	JyfyAmount           float64   `xorm:"decimal(18,2) notnull 'jyfy_amount' default(0)" json:"jyfy_amount"`
	PtfyAmount           float64   `xorm:"decimal(18,2) notnull 'ptfy_amount' default(0)" json:"ptfy_amount"`
	PickFee              float64   `xorm:"decimal(18,2) notnull 'pick_fee' default(0)" json:"pick_fee"`
	GstFee               float64   `xorm:"decimal(18,2) notnull 'gst_fee' default(0)" json:"gst_fee"`
	Profit               float64   `xorm:"decimal(18,2) notnull 'profit' default(0)" json:"profit"`
	ProductQty           float64   `xorm:"decimal(18,2) notnull 'product_qty' default(0)" json:"product_qty"`
}

func (e SaleAmountReportMonthInfo) TableName() string {
	return "erp_sale_amount_report_month"
}

func (e SaleAmountReportMonthInfo) ModelDesc() string {
	return "verbose_name=月销售额报表,entity_json=sale_amount_report_month,route=sale/amount/report/month"
}
