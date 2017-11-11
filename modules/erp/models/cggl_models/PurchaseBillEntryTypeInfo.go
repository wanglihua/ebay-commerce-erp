package cggl_models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	PURCHASE_BILL_TYPE_NORMAL  = 100 // 正常采购单明细
	PURCHASE_BILL_TYPE_IN_TURN = 200 // 分批入库采购单明细
)

type PurchaseBillEntryTypeInfo struct {
	Id          int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=类型名称,min=2,max=100"`
	Description core.NullString `xorm:"nvarchar(500) 'description'" json:"description" smith:"verbose_name=描述,min=2,max=400"`
}

func (e PurchaseBillEntryTypeInfo) TableName() string {
	return "erp_purchase_bill_entry_type"
}

func (e PurchaseBillEntryTypeInfo) ModelDesc() string {
	return "verbose_name=采购单明细类型,entity_json=purchase_bill_entry_type,route=purchase/bill/entry/type"
}

func (e PurchaseBillEntryTypeInfo) InitData(db_session *xorm.Session) {
	purchase_bille_entry_type := new(PurchaseBillEntryTypeInfo)
	purchase_bille_entry_type.Id = PURCHASE_BILL_TYPE_NORMAL
	purchase_bille_entry_type.Name = "正常"
	purchase_bille_entry_type.Description = core.MakeNullString("正常", true)
	db_session.Insert(purchase_bille_entry_type)

	purchase_bille_entry_type = new(PurchaseBillEntryTypeInfo)
	purchase_bille_entry_type.Id = PURCHASE_BILL_TYPE_IN_TURN
	purchase_bille_entry_type.Name = "分批采购"
	purchase_bille_entry_type.Description = core.MakeNullString("分批采购", true)
	db_session.Insert(purchase_bille_entry_type)
}
