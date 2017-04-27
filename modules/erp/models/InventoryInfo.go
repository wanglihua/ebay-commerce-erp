package models

import (
	"matrix/core"
)

type InventoryInfo struct {
	Id           int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId      int64     `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	StorageLocId int64     `xorm:"bigint notnull index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=库位"`
	ProductId    int64     `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Qty          int64     `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	CostPrice    float64   `xorm:"float notnull 'cost_price'" json:"cost_price" smith:"verbose_name=成本价,min=0"` // 到仓库的成本价格，如果各个批次的商品成本价不同，会人工做评价处理
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
	这里可以一个商品在一个仓库里面有多个仓位，有时有库存需要预留，比如说订单进入待配货状态时需要预留商品库存，真的配货时再减去真正的库存

	预留库存，可以采用在这个表里面插入一个特殊库位的库存的方式， storage_loc_id == -1 的方式，订单匹配物流方式算总库存的时候需要把这部分给减去

 */
