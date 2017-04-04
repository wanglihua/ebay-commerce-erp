package models

import (
	"matrix/core"
)

type InventoryInfo struct {
	Id           int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId      int64     `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	StorageLocId int64     `xorm:"bigint notnull index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=库位"`
	ProductId    int64     `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	ReservedQty  int64     `xorm:"bigint notnull 'reserved_qty'" json:"reserved_qty" smith:"verbose_name=保留数量"` // 就是预留的数量，已经被待配货的订单预留。如果订单最后配货失败，需要把预留数量改回来
	Qty          int64     `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	Desc         int64     `xorm:"nvarchar(500) notnull 'desc'" json:"desc" smith:"verbose_name=描述"`
	CreateTime   core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version      int       `xorm:"version notnull 'version'" json:"version"`
}

func (e InventoryInfo) TableName() string {
	return TablePrefix + "inventory"
}

func (e InventoryInfo) ModelDesc() string {
	return "verbose_name=商品库存,entity_json=inventory,route=inventory"
}

/*
一般来讲，一个商品在一个仓库只记录一个库存记录，所以这里面的 storage_loc_id 就很尴尬
 */