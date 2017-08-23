package models

import (
	"matrix/core"
)

// 根据电商平台对物流方式的定义，添加对应字段
type PlatformLogisticsModeInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	PlatformId int64           `xorm:"bigint notnull index 'platform_id'" json:"platform_id" smith:"verbose_name=电商平台"`
	Name       string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PlatformLogisticsModeInfo) TableName() string {
	return "erp_platform_logistics_mode"
}

func (e PlatformLogisticsModeInfo) ModelDesc() string {
	return "verbose_name=电商平台物流方式,entity_json=platform_logistics_mode,route=platform/logistics/mode"
}
