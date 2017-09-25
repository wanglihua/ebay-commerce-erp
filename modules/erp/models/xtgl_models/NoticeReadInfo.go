package xtgl_models

import (
	"matrix/core"
)

type NoticeReadInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	NoticeId   int64     `xorm:"biging notnull index 'notice_id'" json:"notice_id" smith:"verbose_name=通知id"`
	UserId     int64     `xorm:"biging notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户id"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e NoticeReadInfo) TableName() string {
	return "erp_notice_read"
}

func (e NoticeReadInfo) ModelDesc() string {
	return "verbose_name=通知阅读记录,entity_json=notice_read,route=notice_read"
}
