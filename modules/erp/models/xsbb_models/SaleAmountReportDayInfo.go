package xsbb_models

import "matrix/core"

type SaleAmountReportDayInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ReportDayId int64     `xorm:"bigint notnull index 'report_day_id' default(0)" json:"report_day_id"`
	ReportDay   core.Time `xorm:"datetime notnull index 'report_day' default(0)" json:"report_day"`
	Amount      float64   `xorm:"decimal(18,2) notnull 'amount' default(0)" json:"amount"`
}

func (e SaleAmountReportDayInfo) TableName() string {
	return "erp_sale_amount_report_day"
}

func (e SaleAmountReportDayInfo) ModelDesc() string {
	return "verbose_name=日销售额报表,entity_json=sale_amount_report_day,route=sale/amount/report/day"
}
