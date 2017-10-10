package kcgl_models

import (
	"matrix/core"
)

type ProductMelTempSkuInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId   int64     `xorm:"bigint notnull index 'product_id'" json:"product_id"`
	Description string    `xorm:"nvarchar(500) null 'description'" json:"description"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e ProductMelTempSkuInfo) TableName() string {
	return "erp_product_mel_temp_sku"
}

func (e ProductMelTempSkuInfo) ModelDesc() string {
	return "verbose_name=MEL临时SKU,entity_json=product_mel_temp_sku,route=product/mel/temp/sku"
}
