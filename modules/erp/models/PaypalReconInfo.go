package models

import (
	"matrix/core"
)

// PaypalReconciliationInfo
type PaypalReconInfo struct {
	Id         int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code       string                `xorm:"nvarchar(100) notnull index 'code'" json:"code" smith:"verbose_name=单号,min=2,max=100"`
	Fee        float64               `xorm:"decimal(18,6) 'fee'" json:"fee" smith:"verbose_name=交易费"`
	Amount     float64               `xorm:"decimal(18,6) 'amount'" json:"amount" smith:"verbose_name=金额"`
	ResultId   core.NullInt          `xorm:"bigint null 'result_id'" json:"result_id" smith:"verbose_name=匹配结果"`
	Time       core.Time             `xorm:"datetime notnull 'time'" json:"time" smith:"verbose_name=时间"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e PaypalReconInfo) TableName() string {
	return TablePrefix + "paypal_recon"
}

func (e PaypalReconInfo) ModelDesc() string {
	return "verbose_name=Paypal对账,entity_json=paypal_recon,route=paypal/recon"
}
