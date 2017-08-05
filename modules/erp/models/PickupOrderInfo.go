package models

import (
	"matrix/core"
)

// 拣货人工费用的设置
type PickupOrderInfo struct {
	Id          int64        `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId     int64        `xorm:"bigint notnull 'order_id'" json:"order_id" smith:"verbose_name=订单id"`
	CreatorId   int64        `xorm:"bigint notnull 'creator_id'" json:"creator_id" smith:"verbose_name=建单人员id"`
	Creator     string       `xorm:"nvarchar(255) notnull 'creator'" json:"creator" smith:"verbose_name=建单人员工号"`
	PayeeId     int64        `xorm:"bigint notnull 'payee_id'" json:"payee_id" smith:"verbose_name=收款人员id"`
	Payee       string       `xorm:"nvarchar(255) notnull 'payee'" json:"payee" smith:"verbose_name=收款人员工号"`
	PayMethodId int64        `xorm:"bigint notnull 'pay_method_id'" json:"pay_method_id" smith:"verbose_name=收款方式id"`
	CreateTime  core.Time    `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time    `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PickupOrderInfo) TableName() string {
	return TablePrefix + "pickup_order"
}

func (e PickupOrderInfo) ModelDesc() string {
	return "verbose_name=Pickup订单信息,entity_json=pickup_order,route=pickup/order"
}

// 是否已收款信息放在对应的订单信息里面的状态，如果是 YZF 状态就是已收款

//在线支付、现金支付、现场刷卡

const (
	PICKUP_PAY_METHOD_ONLINE = 100 // 在线支付
	PICKUP_PAY_METHOD_CASH   = 200 // 现金支付
	PICKUP_PAY_METHOD_CARD   = 300 // 现场刷卡
)
