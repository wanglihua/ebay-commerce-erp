package models

import (
	"matrix/core"
)

type LoginLogInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	UserId     int64           `xorm:"biging notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户id"`
	LoginName  string          `xorm:"nvarchar(255) notnull index 'login_name'" json:"login_name" smith:"verbose_name=工号"`
	LoginTime  core.Time       `xorm:"datetime notnull 'login_time'" json:"login_time" smith:"verbose_name=登录时间"`
	LogoutTime core.Time       `xorm:"datetime notnull 'logout_time'" json:"logout_time" smith:"verbose_name=登出时间"`
	OnlineTime float64         `xorm:"decimal(18,6) notnull 'online_time'" json:"online_time" smith:"verbose_name=在线时间"`
	IpAddress  string          `xorm:"nvarchar(20) notnull 'ip_address'" json:"ip_address" smith:"verbose_name=IP"`
	IpLoc      string          `xorm:"nvarchar(300) notnull 'ip_loc'" json:"ip_loc" smith:"verbose_name=区域"`
	SessionKey string          `xorm:"nvarchar(300) notnull default('') 'session_key'" json:"session_key" smith:"verbose_name=登录Session"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e LoginLogInfo) TableName() string {
	return "erp_login_log"
}

func (e LoginLogInfo) ModelDesc() string {
	return "verbose_name=登录日志,entity_json=login_log,route=login/log"
}
