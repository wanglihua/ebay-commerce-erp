package models

import (
	"matrix/core"
)

type InventoryInfo struct {
	Id           int64          `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId      int64          `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	StorageLocId int64          `xorm:"bigint notnull index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=库位"`
	ProductId    int64          `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Qty          int64          `xorm:"bigint notnull index 'qty'" json:"qty" smith:"verbose_name=数量"`
	CreateTime   core.Time      `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time      `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version      int            `xorm:"version notnull 'version'" json:"version"`
}

func (e InventoryInfo) TableName() string {
	return TablePrefix + "inventory"
}

func (e InventoryInfo) ModelDesc() string {
	return "verbose_name=商品库存,entity_json=inventory,route=inventory"
}
