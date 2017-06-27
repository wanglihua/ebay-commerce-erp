package models

import (
	"matrix/core"
)

type OrderTagInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	BgColor    string          `xorm:"nvarchar(50) notnull 'bg_color'" json:"bg_color" smith:"verbose_name=背景色,max=500"`
	Desc       string          `xorm:"nvarchar(500) notnull 'desc'" json:"desc" smith:"verbose_name=备注,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderTagInfo) TableName() string {
	return TablePrefix + "order_tag"
}

func (e OrderTagInfo) ModelDesc() string {
	return "verbose_name=订单标签,entity_json=order_tag,route=order/tag"
}
