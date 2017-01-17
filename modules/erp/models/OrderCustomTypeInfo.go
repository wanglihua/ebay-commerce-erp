package models

import (
	"matrix/core"
)

type OrderCustomTypeInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code       string          `xorm:"nvarchar(100) notnull index 'code'" json:"code" smith:"verbose_name=编号"`
	Name       string          `xorm:"nvarchar(200) notnull 'name'" json:"name" smith:"verbose_name=名称"`
	Color      string          `xorm:"nvarchar(100) notnull 'color'" json:"color" smith:"verbose_name=颜色"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	StatusId   int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderCustomTypeInfo) TableName() string {
	return TablePrefix + "order_custom_type"
}

func (e OrderCustomTypeInfo) ModelDesc() string {
	return "verbose_name=订单自定义类型,entity_json=order_type,route=order/custom/type"
}
