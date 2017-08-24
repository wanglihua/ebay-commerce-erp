package sppz_models

import (
	"matrix/core"
)

type SupplierPayMethodInfo struct {
	Id          int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name        string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Description string          `xorm:"nvarchar(500) 'description'" json:"description" smith:"verbose_name=描述,max=500"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e SupplierPayMethodInfo) TableName() string {
	return "erp_supplier_pay_method"
}

func (e SupplierPayMethodInfo) ModelDesc() string {
	return "verbose_name=供应商付款方式,entity_json=pay_method,route=supplier/pay/method"
}
