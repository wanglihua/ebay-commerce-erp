package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	ISSUE_STATUS_YTJ = 100  // 已提交
	ISSUE_STATUS_CLZ = 200  // 处理中
	ISSUE_STATUS_YCL = 300  // 已处理
	ISSUE_STATUS_YJJ = 400  // 已解决
	ISSUE_STATUS_WJJ = 500  // 未解决
	ISSUE_STATUS_YCX = 600  // 又出现
	ISSUE_STATUS_KFZ = 700  // 开发中
	ISSUE_STATUS_YWC = 800  // 已完成
	ISSUE_STATUS_YSX = 900  // 已实现
	ISSUE_STATUS_WSX = 1000 // 未实现
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
	p_issue_status.Desc = "Bug、需求、变更 已提交"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_CLZ
	p_issue_status.Name = "处理中"
	p_issue_status.Desc = "Bug 处理中"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_YCL
	p_issue_status.Name = "已处理"
	p_issue_status.Desc = "Bug 已处理"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_YJJ
	p_issue_status.Name = "已解决"
	p_issue_status.Desc = "Bug 已解决"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_WJJ
	p_issue_status.Name = "未解决"
	p_issue_status.Desc = "Bug 未解决"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_YCX
	p_issue_status.Name = "又出现"
	p_issue_status.Desc = "Bug 又出现"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_KFZ
	p_issue_status.Name = "开发中"
	p_issue_status.Desc = "需求、变更 开发中"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_YWC
	p_issue_status.Name = "已完成"
	p_issue_status.Desc = "需求、变更 已完成"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_YSX
	p_issue_status.Name = "已实现"
	p_issue_status.Desc = "需求、变更 已实现"
	db_session.Insert(p_issue_status)

	p_issue_status = new(IssueStatusInfo)
	p_issue_status.Id = ISSUE_STATUS_WSX
	p_issue_status.Name = "未实现"
	p_issue_status.Desc = "需求、变更 未实现"
	db_session.Insert(p_issue_status)
}
