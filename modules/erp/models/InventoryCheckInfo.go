package models

import (
	"matrix/core"
)

type InventoryCheckInfo struct {
	Id           int64          `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId      int64          `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	StorageLocId int64          `xorm:"bigint notnull index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=库位"`
	ProductId    int64          `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	PreQty       int64          `xorm:"bigint notnull 'pre_qty'" json:"pre_qty" smith:"verbose_name=盘点前库存"`
	CheckQty     int64          `xorm:"bigint notnull 'check_qty'" json:"check_qty" smith:"verbose_name=盘点库存"`
	QtyDiff      int64          `xorm:"bigint notnull 'qty_diff'" json:"qty_diff" smith:"verbose_name=库存差值"`
	StatusId     int64          `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime   core.Time      `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time      `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version      int            `xorm:"version notnull 'version'" json:"version"`
}

func (e InventoryCheckInfo) TableName() string {
	return TablePrefix + "inventory_check"
}

func (e InventoryCheckInfo) ModelDesc() string {
	return "verbose_name=库存盘点,entity_json=inventory_check,route=inventory/check"
}
