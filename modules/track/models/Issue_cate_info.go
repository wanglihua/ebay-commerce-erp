package models

import (
	"matrix/core"
)

type IssueCateInfo struct {
	Id          int64              `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string             `xorm:"nvarchar(255) notnull index 'cate_name'" json:"cate_name" smith:"verbose_name=名称,min=2,max=100"`
	Desc        string             `xorm:"nvarchar(4000) notnull 'cate_desc'" json:"cate_desc" smith:"verbose_name=描述,min=2,max=4000"`
	Create_time core.Time          `xorm:"created notnull 'create_time'" json:"create_time"`
	Update_time core.Time          `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int                `xorm:"version notnull 'version'" json:"version"`
}

func (e IssueCateInfo) TableName() string {
	return TablePrefix + "issue_cate"
}

func (e IssueCateInfo) ModelDesc() string {
	return "verbose_name=缺陷类别,entity_json=issue_type,route=issue/type"
}
