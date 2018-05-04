package xtpz_models

import (
	"matrix/core"
)

type PaypalPaymentInfo struct {
	Id            int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	PaypalAccount string          `xorm:"nvarchar(200) notnull index unique 'paypal_account'" json:"paypal_account"`
	ClientId      string          `xorm:"nvarchar(200) notnull 'client_id'" json:"client_id"`
	Secret        string          `xorm:"nvarchar(200) notnull 'secret'" json:"secret"`
	Desc          core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc"`
	StatusId      int64           `xorm:"bigint notnull index 'status_id'" json:"status_id"`
	CreateTime    core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PaypalPaymentInfo) TableName() string {
	return "erp_paypal_payment"
}

func (e PaypalPaymentInfo) ModelDesc() string {
	return "verbose_name=Paypal付款,entity_json=paypal_payment,route=paypal/payment"
}

/*
		ClientID: "AZgwu4yt5Ba0gyTu1dGBH3txHCJbMuFNvrmQxBaQbfDncDiCs6W_rwJD8Ir-0pZrN-_eq7n9zVd8Y-5f",
		Secret:   "EBzA1wRl5t73OMugOieDj_tI3vihfJmGl47ukQT-cpctooIzDu0K7IPESNC0cKodlLSOXzwI8qXSM0rd",

 */
