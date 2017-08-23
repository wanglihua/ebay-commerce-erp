package sppz_models

import (
	"matrix/core"
)

type SeaContainerCostInfo struct {
	Id         int64         `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Volume     int64         `xorm:"bigint notnull 'volume'" json:"volume"`
	Cost       int64         `xorm:"bigint notnull 'cost'" json:"cost"`
	CreateTime core.Time     `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time     `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e SeaContainerCostInfo) TableName() string {
	return "erp_sea_container_cost"
}
