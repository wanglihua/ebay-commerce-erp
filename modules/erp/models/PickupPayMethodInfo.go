package models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	PICKUP_PAY_METHOD_ONLINE = 100 // 在线支付
	PICKUP_PAY_METHOD_CASH   = 200 // 现金支付
	PICKUP_PAY_METHOD_CARD   = 300 // 现场刷卡
)

type PickupPayMethodInfo struct {
	Id         int64           `xorm:"bigint notnull 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PickupPayMethodInfo) TableName() string {
	return "erp_pickup_pay_method"
}

func (e PickupPayMethodInfo) ModelDesc() string {
	return "verbose_name=PICKUP支付类型,entity_json=pay_method,route=pickup/pay/method"
}

func (e PickupPayMethodInfo) InitData(db_session *xorm.Session) {
	// PICK UP 的订单 一般就都是 手工创建的
	order_type := new(PickupPayMethodInfo)
	order_type.Id = PICKUP_PAY_METHOD_ONLINE
	order_type.Name = "在线支付"
	db_session.Insert(order_type)

	order_type = new(PickupPayMethodInfo)
	order_type.Id = PICKUP_PAY_METHOD_CASH
	order_type.Name = "现金支付"
	db_session.Insert(order_type)

	order_type = new(PickupPayMethodInfo)
	order_type.Id = PICKUP_PAY_METHOD_CARD
	order_type.Name = "现场刷卡"
	db_session.Insert(order_type)
}
