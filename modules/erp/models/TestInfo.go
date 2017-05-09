package models

import (
	"time"

	"matrix/core"
)

type TestInfo struct {
	Id         int64               `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	PaidTime1  core.NullLongTime   `xorm:"datetime 'paid_time1'" json:"paid_time1"`
	PaidTime2  time.Time           `xorm:"datetime 'paid_time2'" json:"paid_time2"`
	PaidTime3  core.Time           `xorm:"datetime 'paid_time3'" json:"paid_time3"`
	CreateTime core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                 `xorm:"version notnull 'version'" json:"version"`
}

func (e TestInfo) TableName() string {
	return TablePrefix + "test"
}
