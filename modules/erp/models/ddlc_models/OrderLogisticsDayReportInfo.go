package ddlc_models

import (
	"matrix/core"
)

type OrderLogisticsDayReportInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	DayDate    core.Date `xorm:"date notnull index 'day_date'" json:"day_date"`
	StockId    int64     `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	CompanyId  int64     `xorm:"bigint notnull index 'company_id'" json:"company_id" smith:"verbose_name=物流公司"`
	ModeId     int64     `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流模式"`
	PrdCount   int64     `xorm:"bigint notnull 'prd_count'" json:"prd_count" smith:"verbose_name=产品数量"`
	Volume     float64   `xorm:"decimal(18,3) notnull 'volume'" json:"volume" smith:"verbose_name=体积"`
	Weight     float64   `xorm:"decimal(18,3) notnull 'weight'" json:"weight" smith:"verbose_name=重量"`
	Amount     float64   `xorm:"decimal(18,2) notnull 'amount'" json:"amount" smith:"verbose_name=金额"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderLogisticsDayReportInfo) TableName() string {
	return "erp_order_logistics_day_report"
}

func (e OrderLogisticsDayReportInfo) ModelDesc() string {
	return "verbose_name=日物流报表,entity_json=order_logistics_day_report,route=order/logistics/day/report"
}
