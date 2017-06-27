package models

import (
	"matrix/core"
)

type IssueLogInfo struct {
	Id          int64              `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	IssueId     int64              `xorm:"biging notnull index 'issue_id'" json:"issue_id" smith:"verbose_name=缺陷id"`
	Content     string             `xorm:"text notnull 'content'" json:"content" smith:"verbose_name=操作内容"`
	UserId      int64              `xorm:"biging notnull index 'user_id'" json:"user_id" smith:"verbose_name=录入人用户id"`
	UserCode    string             `xorm:"nvarchar(255) notnull index 'user_code'" json:"user_code" smith:"verbose_name=录入人工号"`
	CreateTime core.Time          `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time          `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e IssueLogInfo) TableName() string {
	return TablePrefix + "issue_log"
}

func (e IssueLogInfo) ModelDesc() string {
	return "verbose_name=缺陷日志,entity_json=issue_log,route=issue/log"
}
