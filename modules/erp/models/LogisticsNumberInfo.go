package models

import (
	"matrix/core"
)

// 这个客户不需要
type LogisticsNumberInfo struct {
	Id         int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Num        string                `xorm:"nvarchar(100) notnull 'num'" json:"num" smith:"verbose_name=物流单号,min=0, max=100"`
	ModeId     int64                 `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流方式"`
	Desc       core.NullString       `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=详情,max=500"`
	StatusId   int64                 `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsNumberInfo) TableName() string {
	return TablePrefix + "logistics_number"
}

func (e LogisticsNumberInfo) ModelDesc() string {
	return "verbose_name=运单号不规则号段,entity_json=logistics_number,route=logistics/number"
}
