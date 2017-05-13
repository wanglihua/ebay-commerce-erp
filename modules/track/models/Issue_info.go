package models

import (
	"matrix/core"
)

type Issue_info struct {
	Id             int64             `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CateId         int64             `xorm:"biging notnull index 'cate_id'" json:"cate_id" smith:"verbose_name=类别id"`
	TypeId         int64             `xorm:"biging notnull index 'type_id'" json:"type_id" smith:"verbose_name=类型id"`
	Title          string            `xorm:"nvarchar(255) notnull index 'title'" json:"title" smith:"verbose_name=标题,min=2,max=100"`
	Body           string            `xorm:"nvarchar(4000) notnull 'body'" json:"body" smith:"verbose_name=正文,min=2,max=4000"`
	SubmitUserId   int64             `xorm:"biging notnull index 'submit_user_id'" json:"submit_user_id" smith:"verbose_name=提报人用户id"`
	SubmitUserCode string            `xorm:"nvarchar(255) notnull index 'submit_user_code'" json:"submit_user_code" smith:"verbose_name=提报人工号"`
	EnterUserId    int64             `xorm:"biging notnull index 'enter_user_id'" json:"enter_user_id" smith:"verbose_name=录入人用户id"`
	EnterUserCode  string            `xorm:"nvarchar(255) notnull index 'enter_user_code'" json:"enter_user_code" smith:"verbose_name=录入人工号"`
	StatusId       int64             `xorm:"biging notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态id"`
	Create_time    core.Time         `xorm:"created notnull 'create_time'" json:"create_time"`
	Update_time    core.Time         `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version        int               `xorm:"version notnull 'version'" json:"version"`
}

func (e Issue_info) TableName() string {
	return TablePrefix + "issue"
}

func (e Issue_info) ModelDesc() string {
	return "verbose_name=软件缺陷,entity_json=issue,route=issue"
}
