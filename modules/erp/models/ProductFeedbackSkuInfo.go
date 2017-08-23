package models

import (
	"matrix/core"
)

type ProductFeedbackSkuInfo struct {
	Id         int64          `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId  int64          `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	LastCreate core.Time      `xorm:"datetime notnull 'last_create'" json:"last_create"`
	LastUpdate core.Time      `xorm:"datetime notnull 'last_update'" json:"last_update"`
}

func (e ProductFeedbackSkuInfo) TableName() string {
	return "erp_product_feedback_sku"
}

func (e ProductFeedbackSkuInfo) ModelDesc() string {
	return "verbose_name=商品质量反馈SKU"
}
