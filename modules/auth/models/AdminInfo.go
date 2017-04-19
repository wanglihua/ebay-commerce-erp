package models

import (
	"matrix/core"
)

type AdminInfo struct {
	Id      int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	UserId  int64     `xorm:"bigint notnull unique index 'user_id'" json:"user_id"`
	AddTime core.Time `xorm:"created notnull 'add_time'" json:"add_time"`
	IsSuper int64     `xorm:"bigint null 'is_super'" json:"is_super"`
}

func (e AdminInfo) TableName() string {
	return TablePrefix + "admin"
}
