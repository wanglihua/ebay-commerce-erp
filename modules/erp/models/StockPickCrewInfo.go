package models

import "matrix/core"

type StockPickCrewInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	StockId    int64           `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	UserId     int64           `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户"`
	IsManager  int64           `xorm:"bigint notnull 'is_manager'" json:"is_manager" smith:"verbose_name=是否管理员"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e StockPickCrewInfo) TableName() string {
	return TablePrefix + "stock_pick_crew"
}

func (e StockPickCrewInfo) ModelDesc() string {
	return "verbose_name=仓库拣货人员,entity_json=stock_pick_crew,route=stock/pick/crew"
}
