package models

import (
	"matrix/core"
)

type OrderHistoryInfo struct {
	Id          int64          `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId     int64          `xorm:"bigint notnull index 'order_id'" json:"order_id" smith:"verbose_name=订单Id"`
	OrderCode   string         `xorm:"nvarchar(200) notnull index 'order_code'" json:"order_code" smith:"verbose_name=订单编码"`
	StatusId    int64          `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	Name        string         `xorm:"nvarchar(200) notnull 'name'" json:"name" smith:"verbose_name=名称"`
	Description string         `xorm:"nvarchar(500) notnull 'description'" json:"description" smith:"verbose_name=描述"`
	CreateTime  core.Time      `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e OrderHistoryInfo) TableName() string {
	return TablePrefix + "order_history"
}

func (e OrderHistoryInfo) ModelDesc() string {
	return "verbose_name=订单历史,entity_json=order_history,route=order/history"
}
