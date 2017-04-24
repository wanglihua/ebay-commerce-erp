package models

import (
	"matrix/core"
)

type ShopAccountLogInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ShopId     int64     `xorm:"bigint notnull index 'shop_id'" json:"shop_id" smith:"verbose_name=店铺Id"`
	ShopCode   string    `xorm:"nvarchar(100) notnull index 'shop_code'" json:"shop_code" smith:"verbose_name=店铺编号"`
	UserId     int64     `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户Id"`
	LoginName  string    `xorm:"nvarchar(100) notnull index 'login_name'" json:"login_name" smith:"verbose_name=用户工号"`
	Operation  string    `xorm:"nvarchar(800) notnull 'operation'" json:"operation" smith:"verbose_name=操作"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e ShopAccountLogInfo) TableName() string {
	return TablePrefix + "shop_account_log"
}

func (e ShopAccountLogInfo) ModelDesc() string {
	return "verbose_name=店铺账号修改日志,entity_json=shop_account_log,route=shop/account/log"
}
