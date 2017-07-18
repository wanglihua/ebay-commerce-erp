package models

import (
	"matrix/core"
)

type SpecialProductInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId  int64           `xorm:"bigint notnull index 'product_id'" json:"product_id"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e SpecialProductInfo) TableName() string {
	return TablePrefix + "special_product"
}

func (e SpecialProductInfo) ModelDesc() string {
	return "verbose_name=特殊SKU,entity_json=special_product,route=special/product"
}
