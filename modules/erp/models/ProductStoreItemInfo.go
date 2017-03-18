package models

import (
	"matrix/core"
)

type ProductStoreItemInfo struct {
	Id          int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ItemId      string          `xorm:"nvarchar(50) notnull index 'item_id'" json:"item_id"` // ebay系统的 ItemId
	ViewItemUrl string          `xorm:"nvarchar(500) 'view_item_url'" json:"view_item_url"`  // ebay系统的 查看商品的 url
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int             `xorm:"version notnull 'version'" json:"version"`
}

func (e ProductStoreItemInfo) TableName() string {
	return TablePrefix + "product_store_item"
}

func (e ProductStoreItemInfo) ModelDesc() string {
	return "verbose_name=商品店铺项目,entity_json=product_store_item,route=product/store/item"
}
