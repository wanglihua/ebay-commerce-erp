package models

import (
	"matrix/core"
)

// 这个客户不需要
type LogisticsNumberBatchInfo struct {
	Id         int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ModeId     int64                 `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流方式"`
	Prefix     string                `xorm:"nvarchar(20) notnull 'prefix'" json:"prefix" smith:"verbose_name=单号前缀"`
	Suffix     string                `xorm:"nvarchar(20) notnull 'suffix'" json:"suffix" smith:"verbose_name=单号后缀"`
	BeginNum   int64                 `xorm:"bigint notnull 'begin_num'" json:"begin_num" smith:"verbose_name=开始号段,min=0"`
	EndNum     int64                 `xorm:"bigint notnull 'end_num'" json:"end_num" smith:"verbose_name=结束号段,min=0"`
	CurNum     int64                 `xorm:"bigint notnull 'cur_num'" json:"cur_num" smith:"verbose_name=已使用号段,min=0"`
	Remind     int64                 `xorm:"bigint notnull 'remind'" json:"remind" smith:"verbose_name=剩余号段,min=0"`
	StatusId   int64                 `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsNumberBatchInfo) TableName() string {
	return TablePrefix + "logistics_number_batch"
}

func (e LogisticsNumberBatchInfo) ModelDesc() string {
	return "verbose_name=订单运单号,entity_json=number_batch,route=logistics/number/batch"
}
