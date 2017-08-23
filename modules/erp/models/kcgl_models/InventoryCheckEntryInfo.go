package kcgl_models

import (
	"matrix/core"
)

type InventoryCheckEntryInfo struct {
	Id           int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CheckId      int64     `xorm:"bigint notnull index 'check_id'" json:"check_id" smith:"verbose_name=盘点批次"`
	StorageLocId int64     `xorm:"bigint notnull index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=库位"`
	ProductId    int64     `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	PreQty       int64     `xorm:"bigint notnull 'pre_qty'" json:"pre_qty" smith:"verbose_name=盘点前库存"`
	CheckQty     int64     `xorm:"bigint notnull 'check_qty'" json:"check_qty" smith:"verbose_name=盘点库存"`
	QtyDiff      int64     `xorm:"bigint notnull 'qty_diff'" json:"qty_diff" smith:"verbose_name=库存差值"`
	CreateTime   core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e InventoryCheckEntryInfo) TableName() string {
	return "erp_inventory_check_entry"
}

func (e InventoryCheckEntryInfo) ModelDesc() string {
	return "verbose_name=库存盘点明细,entity_json=check_entry,route=inventory/check/entry"
}
