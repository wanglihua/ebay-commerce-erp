package models

import (
	"matrix/core"
)

type ProductImageInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId   int64     `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	ProductCode string    `xorm:"nvarchar(200) notnull index 'product_code'" json:"product_code" smith:"verbose_name=商品编码,min=2,max=100"`
	TypeId      int64     `xorm:"bigint notnull index 'type_id'" json:"type_id" smith:"verbose_name=图片属性"`
	ShopId      int64     `xorm:"bigint notnull index 'shop_id'" json:"shop_id" smith:"verbose_name=店铺"`
	Url         string    `xorm:"nvarchar(400) notnull index 'url'" json:"url" smith:"verbose_name=图片链接,min=2,max=300"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int       `xorm:"version notnull 'version'" json:"version"`
}

func (e ProductImageInfo) TableName() string {
	return TablePrefix + "product_image"
}

func (e ProductImageInfo) ModelDesc() string {
	return "verbose_name=商品图片,entity_json=product_image,route=product/image"
}
