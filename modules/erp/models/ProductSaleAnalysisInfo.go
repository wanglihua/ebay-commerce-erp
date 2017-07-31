package models

import (
	"matrix/core"
)

const (
	// 统计实际值
	PRODUCT_SALE_ANALYSIS_TYPE_ACTUAL = 100
	// 统计理论值
	PRODUCT_SALE_ANALYSIS_TYPE_LOGIC = 200
	// 只统计SKU，不分仓库统计
	PRODUCT_SALE_ANALYSIS_TYPE_SKU = 300
)

type ProductSaleAnalysisInfo struct {
	Id         int64             `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId    int64             `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"` // 如果不是分仓库统计的，全局置为 0
	ProductId  int64             `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Sku        string            `xorm:"bigint notnull index 'sku'" json:"sku" smith:"verbose_name=SKU"`
	Level1Qty  int64             `xorm:"bigint notnull 'level1_qty'" json:"level1_qty" smith:"verbose_name=1档数量"`
	Level2Qty  int64             `xorm:"bigint notnull 'level2_qty'" json:"level2_qty" smith:"verbose_name=2档数量"`
	Level3Qty  int64             `xorm:"bigint notnull 'level3_qty'" json:"level3_qty" smith:"verbose_name=3档数量"`
	Level4Qty  int64             `xorm:"bigint notnull 'level4_qty'" json:"level4_qty" smith:"verbose_name=4档数量"`
	Level5Qty  int64             `xorm:"bigint notnull 'level5_qty'" json:"level5_qty" smith:"verbose_name=5档数量"`
	Type       int64             `xorm:"bigint notnull 'type'" json:"type" smith:"verbose_name=统计类型"`
	CreateTime core.Time         `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time         `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductSaleAnalysisInfo) TableName() string {
	return TablePrefix + "product_sale_analysis"
}

func (e ProductSaleAnalysisInfo) ModelDesc() string {
	return "verbose_name=商品销售统计,entity_json=analysis,route=product/sale/analysis"
}
