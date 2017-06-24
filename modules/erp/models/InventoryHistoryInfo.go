package models

import (
	"matrix/core"
)

type InventoryHistoryInfo struct {
	Id           int64         `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId      int64         `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	StorageLocId int64         `xorm:"bigint notnull index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=库位"`
	ProductId    int64         `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Qty          int64         `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	Description  string        `xorm:"nvarchar(500) notnull 'description'" json:"description" smith:"verbose_name=描述,min=2,max=400"`
	CreateTime   core.Time     `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e InventoryHistoryInfo) TableName() string {
	return TablePrefix + "inventory_history"
}

func (e InventoryHistoryInfo) ModelDesc() string {
	return "verbose_name=库存历史,entity_json=inventory_history,route=inventory/history"
}
