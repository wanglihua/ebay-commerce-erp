package models

import (
	"matrix/core"
)

type OrderLogInfo struct {
	Id         int64            `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId    int64            `xorm:"bigint notnull index 'order_id'" json:"order_id" smith:"verbose_name=订单Id"`
	OrderCode  string           `xorm:"nvarchar(100) notnull index 'order_code'" json:"order_code" smith:"verbose_name=订单编号"`
	TransCode  string           `xorm:"nvarchar(100) notnull index 'trans_code'" json:"trans_code" smith:"verbose_name=交易单号"`
	UserId     int64            `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户Id"`
	LoginName  string           `xorm:"nvarchar(100) notnull index 'login_name'" json:"login_name" smith:"verbose_name=操作工号"`
	Operation  string           `xorm:"nvarchar(800) notnull 'operation'" json:"operation" smith:"verbose_name=操作描述"`
	CreateTime core.Time        `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e OrderLogInfo) TableName() string {
	return TablePrefix + "order_log"
}

func (e OrderLogInfo) ModelDesc() string {
	return "verbose_name=订单修改日志,entity_json=order_log,route=order/log"
}
