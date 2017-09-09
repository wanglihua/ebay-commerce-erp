package cggl_models

import (
	"matrix/core"
)

type PurchaseBillEntryInfo struct {
	Id           int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BillId       int64           `xorm:"bigint notnull index unique(unq_purchase_bill_prd) 'bill_id'" json:"bill_id" smith:"verbose_name=批次"`
	ProductId    int64           `xorm:"bigint notnull index unique(unq_purchase_bill_prd) 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Qty          int64           `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	Price        float64         `xorm:"decimal(18,6) notnull 'price'" json:"price" smith:"verbose_name=单价"`
	Amount       float64         `xorm:"decimal(18,6) notnull 'amount'" json:"amount" smith:"verbose_name=金额"`
	PurchaseDate core.NullDate   `xorm:"datetime null 'purchase_date'" json:"purchase_date" smith:"verbose_name=采购日期"`
	CheckinDate  core.NullDate   `xorm:"datetime null 'checkin_date'" json:"checkin_date" smith:"verbose_name=入库日期"`
	Desc         core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,min=2,max=400"`
	StatusId     int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime   core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseBillEntryInfo) TableName() string {
	return "erp_purchase_bill_entry"
}

func (e PurchaseBillEntryInfo) ModelDesc() string {
	return "verbose_name=采购单,entity_json=purchase_bill_entry,route=purchase/bill/entry"
}

/*
已创建，已下单，采购中，已入库
 */
