package spbb_models

import (
	"matrix/core"
)

const (
	// 10 天一个间隔
	TIME_RANGE_TYPE_TEN_DAYS int64 = 100
)

type ProductSaleAnalysisTimeRangeInfo struct {
	Id        int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	TypeId    int64     `xorm:"bigint notnull index 'type_id'" json:"type_id" smith:"verbose_name=时间间隔类型"`
	BeginTime core.Time `xorm:"datetime notnull 'begin_time'" json:"begin_time"`
	EndTime   core.Time `xorm:"datetime notnull 'end_time'" json:"end_time"`
}

func (e ProductSaleAnalysisTimeRangeInfo) TableName() string {
	return "erp_product_sale_analysis_time_range"
}

func (e ProductSaleAnalysisTimeRangeInfo) ModelDesc() string {
	return "verbose_name=商品销售统计时间段定义,entity_json=analysis,route=product/sale/analysis/time/range"
}
