package models

import (
	"matrix/core"
)

type LogisticsModeInfo struct {
	Id         int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name       string                `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString       `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsModeInfo) TableName() string {
	return TablePrefix + "logistics_mode"
}

func (e LogisticsModeInfo) ModelDesc() string {
	return "verbose_name=物流方式,entity_json=logistics_mode,route=logistics/mode"
}
