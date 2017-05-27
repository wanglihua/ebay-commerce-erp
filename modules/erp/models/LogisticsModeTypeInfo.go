package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	LOGISTICS_MODE_TYPE_PACKAGE = 100 // 包裹
	LOGISTICS_MODE_TYPE_EXPRESS = 200 // 快递
)

type LogisticsModeTypeInfo struct {
	Id          int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string          `xorm:"nvarchar(200) notnull index unique(c_s_name) 'name'" json:"name" smith:"verbose_name=类型名称,min=2,max=100"`
	Description string          `xorm:"nvarchar(500) null 'description'" json:"description" smith:"verbose_name=描述,min=2,max=400"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int             `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsModeTypeInfo) TableName() string {
	return TablePrefix + "logistics_mode_type"
}

func (e LogisticsModeTypeInfo) ModelDesc() string {
	return "verbose_name=物流方式类型,entity_json=logistics_mode_type,route=logistics/mode/type"
}

func (e LogisticsModeTypeInfo) InitData(db_session *xorm.Session) {
	logistics_mode_type := new(LogisticsModeTypeInfo)
	logistics_mode_type.Id = LOGISTICS_MODE_TYPE_PACKAGE
	logistics_mode_type.Name = "包裹"
	logistics_mode_type.Description = "包裹"
	db_session.Insert(logistics_mode_type)

	logistics_mode_type = new(LogisticsModeTypeInfo)
	logistics_mode_type.Id = LOGISTICS_MODE_TYPE_EXPRESS
	logistics_mode_type.Name = "快递"
	logistics_mode_type.Description = "快递"
	db_session.Insert(logistics_mode_type)
}
