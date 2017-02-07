package models

import (
	"matrix/core"
)

type OrderItemInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId    int64           `xorm:"bigint notnull index 'order_id'" json:"order_id" smith:"verbose_name=订单"`
	ProductId  int64           `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=产品"`
	Sku        string          `xorm:"nvarchar(200) notnull index unique 'sku'" json:"sku" smith:"verbose_name=SKU"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderItemInfo) TableName() string {
	return TablePrefix + "order_item"
}

func (e OrderItemInfo) ModelDesc() string {
	return "verbose_name=订单明细,entity_json=order_item,route=order_item"
}
