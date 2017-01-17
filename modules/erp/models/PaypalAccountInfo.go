package models

import (
	"matrix/core"
)

type PaypalAccountInfo struct {
	Id            int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	PaypalAccount string          `xorm:"nvarchar(200) notnull index 'paypal_account'" json:"paypal_account" smith:"verbose_name=Paypal账号,min=2,max=100"`
	ApiAccount    string          `xorm:"nvarchar(200) notnull 'api_account'" json:"api_account" smith:"verbose_name=API账号,min=2,max=100"`
	ApiToken      string          `xorm:"nvarchar(200) notnull 'api_token'" json:"api_token" smith:"verbose_name=API令牌,min=2,max=100"`
	Balance       core.NullFloat  `xorm:"float null 'balance'" json:"balance" smith:"verbose_name=余额,min=0"`
	Desc          core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	StatusId      int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime    core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version       int             `xorm:"version notnull 'version'" json:"version"`
}

func (e PaypalAccountInfo) TableName() string {
	return TablePrefix + "paypal_account"
}

func (e PaypalAccountInfo) ModelDesc() string {
	return "verbose_name=Paypal账号,entity_json=paypal_account,route=paypal/account"
}
