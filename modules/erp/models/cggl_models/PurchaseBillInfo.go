package cggl_models

import (
	"matrix/core"
)

type PurchaseBillInfo struct {
	Id            int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code          string          `xorm:"nvarchar(100) index 'code'" json:"code" smith:"verbose_name=批次号"`
	Desc          core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,min=2,max=400"`
	Date          core.Date       `xorm:"datetime notnull 'date'" json:"date" smith:"verbose_name=日期"`
	UserId        int64           `xorm:"bigint notnull 'user_id'" json:"user_id" smith:"verbose_name=操作人"`
	TransExpenses float64         `xorm:"decimal(18,6) notnull 'trans_expenses'" json:"trans_expenses" smith:"verbose_name=运费"` //transportation expenses
	StatusId      int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime    core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseBillInfo) TableName() string {
	return "erp_purchase_bill"
}

func (e PurchaseBillInfo) ModelDesc() string {
	return "verbose_name=采购单,entity_json=purchase_bill,route=purchase/bill"
}
