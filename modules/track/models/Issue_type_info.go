package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	ISSUE_TYPE_REQ = 100 // 需求
	ISSUE_TYPE_CHG = 200 // 变更
	ISSUE_TYPE_BUG = 300 // Bug
)

type IssueTypeInfo struct {
	Id          int64              `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string             `xorm:"nvarchar(255) notnull index 'type_name'" json:"type_name" smith:"verbose_name=名称,min=2,max=100"`
	Desc        string             `xorm:"nvarchar(4000) notnull 'type_desc'" json:"type_desc" smith:"verbose_name=描述,min=2,max=4000"`
	Create_time core.Time          `xorm:"created notnull 'create_time'" json:"create_time"`
	Update_time core.Time          `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int                `xorm:"version notnull 'version'" json:"version"`
}

func (e IssueTypeInfo) TableName() string {
	return TablePrefix + "issue_type"
}

func (e IssueTypeInfo) ModelDesc() string {
	return "verbose_name=缺陷类型,entity_json=issue_type,route=issue/type"
}

func (e IssueTypeInfo) InitData(db_session *xorm.Session) {
	var p_issue_type = new(IssueTypeInfo)
	p_issue_type.Id = ISSUE_TYPE_REQ
	p_issue_type.Name = "需求"
	p_issue_type.Desc = "需求"
	db_session.Insert(p_issue_type)

	p_issue_type = new(IssueTypeInfo)
	p_issue_type.Id = ISSUE_TYPE_CHG
	p_issue_type.Name = "变更"
	p_issue_type.Desc = "变更"
	db_session.Insert(p_issue_type)

	p_issue_type = new(IssueTypeInfo)
	p_issue_type.Id = ISSUE_TYPE_BUG
	p_issue_type.Name = "Bug"
	p_issue_type.Desc = "Bug"
	db_session.Insert(p_issue_type)
}
