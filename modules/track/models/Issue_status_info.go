package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	ISSUE_STATUS_YTJ = 100 // 已提交
	ISSUE_STATUS_YCL = 200 // 已处理
)

type IssueStatusInfo struct {
	Id          int64              `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string             `xorm:"nvarchar(255) notnull index 'status_name'" json:"status_name" smith:"verbose_name=名称,min=2,max=100"`
	Desc        string             `xorm:"nvarchar(4000) notnull 'status_desc'" json:"status_desc" smith:"verbose_name=描述,min=2,max=4000"`
	Create_time core.Time          `xorm:"created notnull 'create_time'" json:"create_time"`
	Update_time core.Time          `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int                `xorm:"version notnull 'version'" json:"version"`
}

func (e IssueStatusInfo) TableName() string {
	return TablePrefix + "issue_status"
}

func (e IssueStatusInfo) ModelDesc() string {
	return "verbose_name=缺陷状态,entity_json=issue_status,route=issue/status"
}

func (e IssueStatusInfo) InitData(db_session *xorm.Session) {
	var p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_YTJ
	p_issue_status.Name = "已提交"
	p_issue_status.Desc = "已提交"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_YCL
	p_issue_status.Name = "已处理"
	p_issue_status.Desc = "已处理"
	db_session.Insert(p_issue_status)
}
