package xsbb_models

// 这个 SaleAmountReportInfo 只用来接收查询结果，不创建数据库表
type SaleAmountReportInfo struct {
	PaidTime             string  `xorm:"'paid_time'"`
	PaidAmount           float64 `xorm:"'paid_amount'"`
	ExpressFee           float64 `xorm:"'express_fee'"`
	Income               float64 `xorm:"'income'"`
	CostAmount           float64 `xorm:"'cost_amount'"`
	EstimateLogisticsFee float64 `xorm:"'estimate_logistics_fee'"`
	ActualLogisticsFee   float64 `xorm:"'actual_logistics_fee'"`
	JyfyAmount           float64 `xorm:"'jyfy_amount'"`
	PtfyAmount           float64 `xorm:"'ptfy_amount'"`
	PickFee              float64 `xorm:"'pick_fee'"`
	GstFee               float64 `xorm:"'gst_fee'"`
	Outcome              float64 `xorm:"'outcome'"`
	Profit               float64 `xorm:"'profit'"`
	ProfitRate           float64 `xorm:"'profit_rate'"`
}

/*

SELECT
    date_format(r.paid_time, '%Y-%m-%d')            AS paid_time,
    sum(r.paid_amount)                              AS paid_amount,
    sum(r.express_fee)                              AS express_fee,
    sum(r.income)                                   AS income,
    sum(r.cost_amount)                              AS cost_amount,
    sum(r.estimate_logistics_fee)                   AS estimate_logistics_fee,
    sum(r.actual_logistics_fee)                     AS actual_logistics_fee,
    sum(r.jyfy_amount)                              AS jyfy_amount,
    sum(r.ptfy_amount)                              AS ptfy_amount,
    sum(r.pick_fee)                                 AS pick_fee,
    sum(r.gst_fee)                                  AS gst_fee,
    sum(r.outcome)                                  AS outcome,
    sum(r.income - r.outcome)                       AS profit,
    sum(r.income - r.outcome) / sum(r.income) * 100 AS profit_rate
FROM erp_sale_amount_report r
WHERE r.platform_id = 0
      AND r.shop_id = 0
      AND r.paid_time >= '2006-01-01'
      AND r.paid_time < '2017-12-30'
      AND r.product_id = 0
      AND r.sale_user_id = 0
      AND r.service_user_id = 0
      AND r.dev_user_id = 0
GROUP BY date_format(r.paid_time, '%Y-%m-%d');

 */
