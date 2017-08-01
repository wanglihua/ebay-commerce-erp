package models

import (
	"matrix/core"
)

type ProductSaleDayAnalysisInfo struct {
	Id          int64             `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId     int64             `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"` // 如果不是分仓库统计的，全局置为 0
	ProductId   int64             `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Sku         string            `xorm:"nvarchar(200) notnull index 'sku'" json:"sku" smith:"verbose_name=SKU"`
	ProductName string            `xorm:"-" json:"product_name"`
	Day         core.Date         `xorm:"datetime notnull index 'day'" json:"day"`
	Qty         int64             `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=销售数量"`
	Type        int64             `xorm:"bigint notnull 'type'" json:"type" smith:"verbose_name=统计类型"`
	TypeName    string            `xorm:"-" json:"type_name"`
	CreateTime  core.Time         `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time         `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductSaleDayAnalysisInfo) TableName() string {
	return TablePrefix + "product_sale_day_analysis"
}

func (e ProductSaleDayAnalysisInfo) ModelDesc() string {
	return "verbose_name=商品销售日统计,entity_json=analysis,route=product/sale/day/analysis"
}
