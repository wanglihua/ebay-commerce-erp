package xsbb_models

import "matrix/core"

type SaleReportTimeMonthInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ReportMonth core.Time `xorm:"datetime notnull index 'report_month'" json:"report_month"`
}

func (e SaleReportTimeMonthInfo) TableName() string {
	return "erp_sale_report_time_month"
}

func (e SaleReportTimeMonthInfo) ModelDesc() string {
	return "verbose_name=销售报表时间月,entity_json=time_month,route=sale/report/time/month"
}
