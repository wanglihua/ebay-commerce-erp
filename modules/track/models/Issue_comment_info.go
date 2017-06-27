package models

import (
	"matrix/core"
)

type IssueCommentInfo struct {
	Id         int64             `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	IssueId    int64             `xorm:"biging notnull index 'issue_id'" json:"issue_id" smith:"verbose_name=缺陷id"`
	Body       string            `xorm:"text notnull 'body'" json:"body" smith:"verbose_name=正文,min=2,max=4000"`
	UserId     int64             `xorm:"biging notnull index 'user_id'" json:"user_id" smith:"verbose_name=提报人用户id"`
	UserCode   string            `xorm:"nvarchar(255) notnull index 'user_code'" json:"user_code" smith:"verbose_name=提报人工号"`
	CreateTime core.Time         `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time         `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e IssueCommentInfo) TableName() string {
	return TablePrefix + "issue_comment"
}

func (e IssueCommentInfo) ModelDesc() string {
	return "verbose_name=缺陷评论,entity_json=issue_comment,route=issue/comment"
}
