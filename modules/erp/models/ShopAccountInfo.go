package models

import (
	"matrix/core"
)

type ShopAccountInfo struct {
	Id          int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Account     string                `xorm:"nvarchar(200) notnull index 'account'" json:"account" smith:"verbose_name=账号,min=2,max=100"`
	Platform    string                `xorm:"nvarchar(200) notnull 'platform'" json:"platform" smith:"verbose_name=平台,min=2,max=100"`
	Desc        core.NullString       `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	Manager     core.NullString       `xorm:"nvarchar(30) 'manager'" json:"manager" smith:"verbose_name=负责人,max=30"`
	TokenResult core.NullString       `xorm:"nvarchar(200) 'token_result'" json:"token_result" smith:"verbose_name=Token验证结果"`
	StatusId    int64                 `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime  core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e ShopAccountInfo) TableName() string {
	return TablePrefix + "shop_account"
}

func (e ShopAccountInfo) ModelDesc() string {
	return "verbose_name=店铺账号,entity_json=shop_account,route=shop/account"
}
