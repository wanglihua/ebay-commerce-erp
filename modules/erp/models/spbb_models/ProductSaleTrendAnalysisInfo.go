package spbb_models

import (
	"matrix/core"
)

type ProductSaleTrendAnalysisInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	TimeRangeId int64     `xorm:"bigint notnull index 'time_range_id'" json:"time_range_id" smith:"verbose_name=时间间隔Id"`
	TypeId      int64     `xorm:"bigint notnull index 'type_id'" json:"type_id" smith:"verbose_name=时间间隔类型"`
	BeginTime   core.Time `xorm:"datetime notnull index 'begin_time'" json:"begin_time"`
	EndTime     core.Time `xorm:"datetime notnull 'end_time'" json:"end_time"`
	ProductId   int64     `xorm:"bigint notnull index 'product_id'" json:"product_id"`
	Qty         int64     `xorm:"bigint notnull 'qty'" json:"qty"`
}

func (e ProductSaleTrendAnalysisInfo) TableName() string {
	return "erp_product_sale_trend_analysis"
}

func (e ProductSaleTrendAnalysisInfo) ModelDesc() string {
	return "verbose_name=商品销售趋势分析,entity_json=analysis,route=product/sale/trend/analysis"
}
