package models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

type OrderHoldTypeInfo struct {
	Id         int64               `xorm:"bigint notnull pk 'id'" json:"id"`
	Name       string              `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称"`
	Desc       core.NullString     `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
	CreateTime core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                 `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderHoldTypeInfo) TableName() string {
	return TablePrefix + "order_hold_type"
}

func (e OrderHoldTypeInfo) ModelDesc() string {
	return "verbose_name=订单Hold类型,entity_json=order_hold_type,route=order/hold/type"
}

const (
	ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH int64 = 10000
)

func (e OrderHoldTypeInfo) InitData(db_session *xorm.Session) {
	order_hold_type := new(OrderHoldTypeInfo)
	order_hold_type.Id = ORDER_HOLD_TYPE_SELLER_PAYPAL_MISMATCH
	order_hold_type.Name = "卖家paypal不匹配"
	order_hold_type.Desc = core.NewNullString("店铺配置的paypal账号和订单中卖家的paypal账号不匹配", true)

	db_session.Insert(order_hold_type)
}
