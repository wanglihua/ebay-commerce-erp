package cggl_models

import (
	"matrix/core"
)

type PurchaseBillLogInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BillId     int64     `xorm:"bigint notnull index 'bill_id'" json:"bill_id" smith:"verbose_name=采购单id"`
	BillCode   string    `xorm:"nvarchar(200) notnull index 'bill_code'" json:"bill_code" smith:"verbose_name=采购批次号"`
	UserId     int64     `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户Id"`
	LoginName  string    `xorm:"nvarchar(100) notnull index 'login_name'" json:"login_name" smith:"verbose_name=用户工号"`
	Operation  string    `xorm:"nvarchar(200) notnull 'operation'" json:"operation"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PurchaseBillLogInfo) TableName() string {
	return "erp_purchase_bill_log"
}

func (e PurchaseBillLogInfo) ModelDesc() string {
	return "verbose_name=采购单操作日志,entity_json=purchase_bill_log,route=purchase/bill/log"
}
