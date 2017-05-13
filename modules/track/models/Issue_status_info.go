package models

import (
	"matrix/core"
)

type Issue_status_info struct {
	Id          int64              `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string             `xorm:"nvarchar(255) notnull index 'status_name'" json:"status_name" smith:"verbose_name=名称,min=2,max=100"`
	Desc        string             `xorm:"nvarchar(4000) notnull 'status_desc'" json:"status_desc" smith:"verbose_name=描述,min=2,max=4000"`
	Create_time core.Time          `xorm:"created notnull 'create_time'" json:"create_time"`
	Update_time core.Time          `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int                `xorm:"version notnull 'version'" json:"version"`
}

func (e Issue_status_info) TableName() string {
	return TablePrefix + "issue_status"
}

func (e Issue_status_info) ModelDesc() string {
	return "verbose_name=缺陷状态,entity_json=issue_status,route=issue/status"
}
