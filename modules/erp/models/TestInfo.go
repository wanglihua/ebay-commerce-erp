package models

import (
	"matrix/core"
)

type TestInfo struct {
	Id         int64               `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	PaidTime   core.NullLongTime   `xorm:"datetime 'paid_time'" json:"paid_time" smith:"verbose_name=支付时间"`
	CreateTime core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                 `xorm:"version notnull 'version'" json:"version"`
}

func (e TestInfo) TableName() string {
	return TablePrefix + "test"
}
