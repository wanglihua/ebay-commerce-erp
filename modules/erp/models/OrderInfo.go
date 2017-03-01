package models

import (
	"matrix/core"
)

type OrderInfo struct {
	Id           int64               `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code         string              `xorm:"nvarchar(200) notnull index unique 'code'" json:"code" smith:"verbose_name=编码"`
	Desc         core.NullString     `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
	TypeId       int64               `xorm:"bigint notnull index 'type_id'" json:"type_id" smith:"verbose_name=类型"`
	EBayStatus   string              `xorm:"nvarchar(200) notnull 'ebay_status'" json:"ebay_status" smith:"verbose_name=ebay状态"`
	StatusId     int64               `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreatedTime  core.NullString     `xorm:"nvarchar(100) 'created_time'" json:"created_time" smith:"verbose_name=创建时间"`
	ModifiedTime core.NullString     `xorm:"nvarchar(100) 'modified_time'" json:"modified_time" smith:"verbose_name=修改时间"`
	PaidTime     core.NullString     `xorm:"nvarchar(100) 'paid_time'" json:"paid_time" smith:"verbose_name=支付时间"`
	ShippedTime  core.NullString     `xorm:"nvarchar(100) 'shipped_time'" json:"shipped_time" smith:"verbose_name=创建时间"`
	CreateTime   core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version      int                 `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderInfo) TableName() string {
	return TablePrefix + "order"
}

func (e OrderInfo) ModelDesc() string {
	return "verbose_name=订单,entity_json=order,route=order"
}
