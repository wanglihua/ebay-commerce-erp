package models

import (
	"matrix/core"
)

// 这个客户不需要
type LogisticsNumberLogInfo struct {
	Id         int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	NumId      int64                 `xorm:"bigint notnull index 'num_id'" json:"num_id" smith:"verbose_name=运单号Id"`
	ModeName   string                `xorm:"nvarchar(200) notnull 'mode_name'" json:"mode_name" smith:"verbose_name=货运方式"`
	NumRange   string                `xorm:"nvarchar(200) notnull 'num_range'" json:"num_range" smith:"verbose_name=号段"`
	Desc       core.NullString       `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	UserId     int64                 `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户Id"`
	LoginName  string                `xorm:"nvarchar(100) notnull index 'login_name'" json:"login_name" smith:"verbose_name=用户工号"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsNumberLogInfo) TableName() string {
	return TablePrefix + "logistics_number_log"
}

func (e LogisticsNumberLogInfo) ModelDesc() string {
	return "verbose_name=运单号操作日志,entity_json=number_log,route=logistics/number/log"
}
