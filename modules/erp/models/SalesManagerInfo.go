package models

import (
	"matrix/core"
)

type SalesManagerInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ShopId     int64           `xorm:"bigint notnull index unique(shop_product_user) 'shop_id'" json:"shop_id" smith:"verbose_name=店铺"`
	ProductId  int64           `xorm:"bigint notnull index unique(shop_product_user) 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	UserId     int64           `xorm:"bigint notnull index unique(shop_product_user) 'user_id'" json:"user_id" smith:"verbose_name=销售"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e SalesManagerInfo) TableName() string {
	return TablePrefix + "sales_manager"
}

func (e SalesManagerInfo) ModelDesc() string {
	return "verbose_name=销售负责人,entity_json=sales_manager,route=sales/manager"
}
