package models

import (
	"matrix/core"
)

type OrderLogisticsInfo struct {
	Id         int64        `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId    int64        `xorm:"bigint notnull index 'order_id'" json:"order_id" smith:"verbose_name=订单"`
	AddressId  int64        `xorm:"bigint notnull index 'address_id'" json:"address_id" smith:"verbose_name=物流地址"`      // OrderShippingAddressInfo 的主键
	ModeId     int64        `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流方式"`            // LogisticsModeInfo  的主键，这个表里面包含了 发货仓库Id 等的信息
	FeeRateId  int64        `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流费率"`            // LogisticsFeeRateInfo  的主键
	TotalFee   float64      `xorm:"decimal(18,6) notnull index 'total_fee'" json:"total_fee" smith:"verbose_name=物流费用"` // 由上面的 物流方式和物流费率，计算出来的物流费用
	CreateTime core.Time    `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time    `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int          `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderLogisticsInfo) TableName() string {
	return TablePrefix + "order_logistics"
}

func (e OrderLogisticsInfo) ModelDesc() string {
	return "verbose_name=订单物流,entity_json=order_logistics,route=order/logistics"
}
