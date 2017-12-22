package xsbb_models

import "matrix/core"

type SaleAmountReportDayInfo struct {
	Id                   int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ReportDayId          int64     `xorm:"bigint notnull index 'report_day_id' default(0)" json:"report_day_id"`
	ReportDay            core.Time `xorm:"datetime notnull index 'report_day' default(0)" json:"report_day"`
	PaidAmount           float64   `xorm:"decimal(18,2) notnull 'paid_amount' default(0)" json:"paid_amount"`
	ExpressFee           float64   `xorm:"decimal(18,2) notnull 'express_fee' default(0)" json:"express_fee"`
	CostPrice            float64   `xorm:"decimal(18,2) notnull 'cost_price' default(0)" json:"cost_price"` // AUD，其他的也都是 AUD
	EstimateLogisticsFee float64   `xorm:"decimal(18,2) notnull 'estimate_logistics_fee' default(0)" json:"estimate_logistics_fee"`
	ActualLogisticsFee   float64   `xorm:"decimal(18,2) notnull 'actual_logistics_fee' default(0)" json:"actual_logistics_fee"`
	JyfyAmount           float64   `xorm:"decimal(18,2) notnull 'jyfy_amount' default(0)" json:"jyfy_amount"`
	PtfyAmount           float64   `xorm:"decimal(18,2) notnull 'ptfy_amount' default(0)" json:"ptfy_amount"`
	PickFee              float64   `xorm:"decimal(18,2) notnull 'pick_fee' default(0)" json:"pick_fee"`
	GstFee               float64   `xorm:"decimal(18,2) notnull 'gst_fee' default(0)" json:"gst_fee"`
	Profit               float64   `xorm:"decimal(18,2) notnull 'profit' default(0)" json:"profit"`
	ProductQty           float64   `xorm:"decimal(18,2) notnull 'product_qty' default(0)" json:"product_qty"`
}

func (e SaleAmountReportDayInfo) TableName() string {
	return "erp_sale_amount_report_day"
}

func (e SaleAmountReportDayInfo) ModelDesc() string {
	return "verbose_name=日销售额报表,entity_json=sale_amount_report_day,route=sale/amount/report/day"
}

/*
        paid_amount DECIMAL(18, 2); # 交易金额 AUD
        express_fee DECIMAL(18, 2); # 运费收入 AUD
        cost_price DECIMAL(18, 2); # 商品澳仓成本 RMB
        estimate_logistics_fee DECIMAL(18, 2); # 预估运费 AUD
        actual_logistics_fee DECIMAL(18, 2); # 实际运费 AUD
        jyfy_amount DECIMAL(18, 2); # 交易费用 AUD
        ptfy_amount DECIMAL(18, 2); # 平台费用 AUD
        pick_fee DECIMAL(18, 2); # 拣货费 AUD
        gst_fee DECIMAL(18, 2); # GST费用 AUD
        profit DECIMAL(18, 2); # 利润 AUD
        product_qty INT; # 商品数量
 */
