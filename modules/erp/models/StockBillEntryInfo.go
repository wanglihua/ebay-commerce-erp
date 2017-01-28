package models

import (
	"matrix/core"
)

type StockBillEntryInfo struct {
	Id           int64          `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BillId       int64          `xorm:"bigint notnull index 'bill_id'" json:"bill_id" smith:"verbose_name=批次"`
	StorageLocId int64          `xorm:"bigint notnull index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=仓位"`
	ProductId    int64          `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	PreQty       int64          `xorm:"bigint notnull 'pre_qty'" json:"pre_qty" smith:"verbose_name=当前库存"`
	Qty          int64          `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	Price        core.NullFloat `xorm:"decimal(18,4) null 'price'" json:"price" smith:"verbose_name=单价"`
	Amount       core.NullFloat `xorm:"decimal(18,4) null 'amount'" json:"amount" smith:"verbose_name=金额"`
	CreateTime   core.Time      `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time      `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e StockBillEntryInfo) TableName() string {
	return TablePrefix + "stock_bill_entry"
}

func (e StockBillEntryInfo) ModelDesc() string {
	return "verbose_name=出入库单明细,entity_json=bill_entry,route=stock/bill/entry"
}
