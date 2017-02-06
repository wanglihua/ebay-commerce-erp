package models


import (
	"matrix/core"
)

type OrderItemInfo struct {
	Id          int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Sku         string          `xorm:"nvarchar(200) notnull index unique 'sku'" json:"sku" smith:"verbose_name=SKU,min=2,max=100"`
	Name        string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc        core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,max=500"`
	StatusId    int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int             `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderItemInfo) TableName() string {
	return TablePrefix + "order_item"
}

func (e OrderItemInfo) ModelDesc() string {
	return "verbose_name=订单明细,entity_json=order_item,route=order_item"
}
