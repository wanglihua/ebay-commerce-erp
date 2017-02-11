package models

import (
	"matrix/core"
)

type PurchaseBillEntryInfo struct {
	Id           int64          `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BillId       int64          `xorm:"bigint notnull index unique(unq_bill_loc_prd) 'bill_id'" json:"bill_id" smith:"verbose_name=批次"`
	StorageLocId int64          `xorm:"bigint notnull index unique(unq_bill_loc_prd) 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=仓位"`
	ProductId    int64          `xorm:"bigint notnull index unique(unq_bill_loc_prd) 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	PreQty       int64          `xorm:"bigint notnull 'pre_qty'" json:"pre_qty" smith:"verbose_name=当前库存"`
	Qty          int64          `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	Price        core.NullFloat `xorm:"decimal(18,6) null 'price'" json:"price" smith:"verbose_name=单价"`
	Amount       core.NullFloat `xorm:"decimal(18,6) null 'amount'" json:"amount" smith:"verbose_name=金额"`
	CreateTime   core.Time      `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time      `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseBillEntryInfo) TableName() string {
	return TablePrefix + "purchase_bill_entry"
}

func (e PurchaseBillEntryInfo) ModelDesc() string {
	return "verbose_name=采购单,entity_json=purchase_bill_entry,route=purchase/bill/entry"
}

/*
已创建，已下单，采购中，已入库
 */
