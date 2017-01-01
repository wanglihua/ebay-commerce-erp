package models

import (
	"matrix/core"
)

type LoginLogInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	LoginName  string          `xorm:"nvarchar(255) notnull index 'login_name'" json:"login_name" smith:"verbose_name=工号,max=30"`
	IpAddress  string          `xorm:"nvarchar(20) notnull 'ip_address'" json:"ip_address" smith:"verbose_name=IP,max=30"`
	Desc       core.NullString `xorm:"nvarchar(300) null 'desc'" json:"desc" smith:"verbose_name=描述,max=300"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e LoginLogInfo) TableName() string {
	return TablePrefix + "login_log"
}

func (e LoginLogInfo) ModelDesc() string {
	return "verbose_name=登录日志,entity_json=login_log,route=login/log"
}
