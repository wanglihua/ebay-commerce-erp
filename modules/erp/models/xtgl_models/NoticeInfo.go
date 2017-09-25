package xtgl_models

import (
	"matrix/core"
)

type NoticeInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Title      string    `xorm:"nvarchar(300) notnull 'title'" json:"title"`
	ListTop    int64     `xorm:"bigint notnull 'list_top' default(0)" json:"list_top"`
	Important  int64     `xorm:"bigint notnull 'important' default(0)" json:"important"`
	Body       string    `xorm:"nvarchar(10000) notnull 'body'" json:"body"`
	UserId     int64     `xorm:"biging notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户id"`
	LoginName  string    `xorm:"nvarchar(255) notnull index 'login_name'" json:"login_name" smith:"verbose_name=工号"`
	IsRead     bool      `xorm:"-" json:"is_read"` // 是否已阅
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e NoticeInfo) TableName() string {
	return "erp_notice"
}

func (e NoticeInfo) ModelDesc() string {
	return "verbose_name=通知公告,entity_json=notice,route=notice"
}
