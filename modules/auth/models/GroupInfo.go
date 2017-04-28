package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

type GroupInfo struct {
	Id         int64     `xorm:"bigint notnull pk 'id'" json:"id"`
	GroupName  string    `xorm:"nvarchar(255) notnull unique index 'group_name'" json:"group_name"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int       `xorm:"version notnull 'version'" json:"version"`
}

func (e GroupInfo) TableName() string {
	return TablePrefix + "group"
}

const (
	GROUP_CUSTOM_SERVICE = 11100
	GROUP_SALES          = 11200
	GROUP_PICK_CREW      = 11300
)

func (e GroupInfo) InitData(db_session *xorm.Session) {
	group := new(GroupInfo)
	group.Id = GROUP_CUSTOM_SERVICE
	group.GroupName = "客服"
	db_session.Insert(group)

	group = new(GroupInfo)
	group.Id = GROUP_SALES
	group.GroupName = "销售"
	db_session.Insert(group)

	group = new(GroupInfo)
	group.Id = GROUP_PICK_CREW
	group.GroupName = "拣货人员"
	db_session.Insert(group)
}
