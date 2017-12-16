package xsbb_models

import "matrix/core"

type SaleAmountReportMonthInfo struct {
	Id            int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ReportMonthId int64     `xorm:"bigint notnull index 'report_month_id' default(0)" json:"report_month_id"`
	ReportMonth   core.Time `xorm:"datetime notnull index 'report_month' default(0)" json:"report_month"`
	ProductId     int64     `xorm:"bigint notnull index 'product_id' default(0)" json:"product_id"`
	Amount        float64   `xorm:"decimal(18,2) notnull 'amount' default(0)" json:"amount"`
}

func (e SaleAmountReportMonthInfo) TableName() string {
	return "erp_sale_amount_report_month"
}

func (e SaleAmountReportMonthInfo) ModelDesc() string {
	return "verbose_name=月销售额报表,entity_json=sale_amount_report_month,route=sale/amount/report/month"
}
