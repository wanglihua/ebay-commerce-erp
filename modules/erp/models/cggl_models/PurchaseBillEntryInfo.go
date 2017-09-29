package cggl_models

import (
	"matrix/core"
)

type PurchaseBillEntryInfo struct {
	Id                      int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BillId                  int64           `xorm:"bigint notnull index unique(unq_purchase_bill_prd) 'bill_id'" json:"bill_id" smith:"verbose_name=批次"`
	ProductId               int64           `xorm:"bigint notnull index unique(unq_purchase_bill_prd) 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Qty                     int64           `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=总数量"`                             // 所有仓库的采购数量，各个仓库加起来的总和
	MelQty                  int64           `xorm:"bigint notnull default(0) 'mel_qty'" json:"mel_qty" smith:"verbose_name=墨尔本数量"`        // 墨尔本仓库采购数量
	BneQty                  int64           `xorm:"bigint notnull default(0) 'bne_qty'" json:"bne_qty" smith:"verbose_name=布里斯班数量"`       // 墨尔本仓库采购数量
	CheckinQty              int64           `xorm:"bigint notnull 'checkin_qty' default(0)" json:"checkin_qty" smith:"verbose_name=入库数量"` // 最终入库的数量
	Price                   float64         `xorm:"decimal(18,6) notnull 'price'" json:"price" smith:"verbose_name=单价"`
	Amount                  float64         `xorm:"decimal(18,6) notnull 'amount'" json:"amount" smith:"verbose_name=金额"`
	PurchaseTime            core.NullTime   `xorm:"datetime null 'purchase_time'" json:"purchase_time" smith:"verbose_name=采购时间"`
	CheckinTime             core.NullTime   `xorm:"datetime null 'checkin_time'" json:"checkin_time" smith:"verbose_name=入库时间"`
	Desc                    core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,min=2,max=400"`
	StatusId                int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	AbnormalTypeId          int64           `xorm:"bigint notnull index 'abnormal_type_id' default(0)" json:"abnormal_type_id" smith:"verbose_name=异常类型"`
	AbnormalProcessResultId int64           `xorm:"bigint notnull index 'abnormal_process_result_id' default(0)" json:"abnormal_process_result_id" smith:"verbose_name=异常处理结果"`
	CheckinUserId           int64           `xorm:"bigint notnull 'checkin_user_id' default(0)" json:"checkin_user_id" smith:"verbose_name=入库人"`
	CreateTime              core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime              core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseBillEntryInfo) TableName() string {
	return "erp_purchase_bill_entry"
}

func (e PurchaseBillEntryInfo) ModelDesc() string {
	return "verbose_name=采购单,entity_json=purchase_bill_entry,route=purchase/bill/entry"
}

/*
已创建，已审批，采购中，已入库
 */
