package xsbb_models

import "matrix/core"

type SaleReportTimeDayInfo struct {
	Id        int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ReportDay core.Time `xorm:"datetime notnull index 'report_day'" json:"report_day"`
}

func (e SaleReportTimeDayInfo) TableName() string {
	return "erp_sale_report_time_day"
}

func (e SaleReportTimeDayInfo) ModelDesc() string {
	return "verbose_name=销售报表时间天,entity_json=time_day,route=sale/report/time/day"
}
