package models

import (
	"matrix/core"
)

type PassUserResInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	UserId     int64           `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户"`
	ResId      int64           `xorm:"bigint notnull index 'res_id'" json:"res_id" smith:"verbose_name=权限"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e PassUserResInfo) TableName() string {
	return TablePrefix + "pass_user_res"
}

func (e PassUserResInfo) ModelDesc() string {
	return "verbose_name=用户权限,entity_json=pass_user_res,route=pass/user/res"
}
