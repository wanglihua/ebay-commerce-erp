package sppz_models

import (
	"matrix/core"
)

type ProductStoreItemInfo struct {
	Id            int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId     int64     `xorm:"bigint notnull index 'product_id'" json:"product_id"`
	ProductBindId int64     `xorm:"bigint index 'product_bind_id'" json:"product_bind_id"`
	ShopId        int64     `xorm:"bigint notnull index 'shop_id'" json:"shop_id"`
	ItemId        string    `xorm:"nvarchar(200) notnull index 'item_id'" json:"item_id"` // ebay系统的 ItemId
	Title         string    `xorm:"nvarchar(300) notnull default('')" json:"title"`
	ViewItemUrl   string    `xorm:"nvarchar(500) notnull 'view_item_url'" json:"view_item_url"` // ebay系统的 查看商品的 url
	CreateTime    core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductStoreItemInfo) TableName() string {
	return "erp_product_store_item"
}

func (e ProductStoreItemInfo) ModelDesc() string {
	return "verbose_name=商品店铺项目,entity_json=product_store_item,route=product/store/item"
}
