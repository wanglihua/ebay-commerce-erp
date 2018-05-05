package xtpz_models

import (
	"matrix/core"
)

type PaypalAccountInfo struct {
	Id            int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	PaypalAccount string          `xorm:"nvarchar(200) notnull index unique 'paypal_account'" json:"paypal_account"`
	ApiUserName   string          `xorm:"nvarchar(200) notnull 'api_username'" json:"api_username"`
	ApiPassword   string          `xorm:"nvarchar(200) notnull 'api_password'" json:"api_password"`
	Signature     string          `xorm:"nvarchar(200) notnull 'signature'" json:"signature"`
	Desc          core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc"`
	StatusId      int64           `xorm:"bigint notnull index 'status_id'" json:"status_id"`
	CreateTime    core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PaypalAccountInfo) TableName() string {
	return "erp_paypal_account"
}

func (e PaypalAccountInfo) ModelDesc() string {
	return "verbose_name=Paypal账号,entity_json=paypal_account,route=paypal/account"
}

/*

ALTER TABLE erp_paypal_account
   CHANGE client_id api_username VARCHAR(200) NOT NULL;

ALTER TABLE erp_paypal_account
   CHANGE secret api_password VARCHAR(200) NOT NULL;

ALTER TABLE erp_paypal_account
   CHANGE balance signature VARCHAR(200) NOT NULL;

 */
