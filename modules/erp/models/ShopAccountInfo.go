package models

import (
	"matrix/core"
)

type ShopAccountInfo struct {
	Id            int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code          string          `xorm:"nvarchar(200) notnull index 'code'" json:"code" smith:"verbose_name=编号,min=2,max=100"`
	Name          string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	PlatformId    int64           `xorm:"int64 notnull index 'platform_id'" json:"platform_id" smith:"verbose_name=平台"`
	Desc          core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	PaypalAccount string          `xorm:"nvarchar(500) null 'paypal_account'" json:"paypal_account" smith:"verbose_name=Paypal账号,max=400"` // 防止老鼠仓，多个用逗号隔开
	Manager       core.NullString `xorm:"nvarchar(30) 'manager'" json:"manager" smith:"verbose_name=负责人,max=30"`
	TokenResult   core.NullString `xorm:"nvarchar(200) 'token_result'" json:"token_result" smith:"verbose_name=Token验证结果"`
	StatusId      int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime    core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version       int             `xorm:"version notnull 'version'" json:"version"`
}

func (e ShopAccountInfo) TableName() string {
	return TablePrefix + "shop_account"
}

func (e ShopAccountInfo) ModelDesc() string {
	return "verbose_name=店铺账号,entity_json=shop_account,route=shop/account"
}
