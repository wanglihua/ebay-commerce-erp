package ppgl_models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const PAYMENT_METHOD_PALPAY int64 = 1
const PAYMENT_METHOD_BANK int64 = 2

type OrderRefundMethodInfo struct {
	Id          int64     `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string    `xorm:"nvarchar(200) notnull index 'name'" json:"name"`
	Description string    `xorm:"nvarchar(500) notnull 'description'" json:"description"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderRefundMethodInfo) TableName() string {
	return "erp_order_refund_method"
}

func (e OrderRefundMethodInfo) ModelDesc() string {
	return "verbose_name=订单退款方式,entity_json=order_refund_method,route=order/refund/method"
}

func (e OrderRefundMethodInfo) InitData(db_session *xorm.Session) {

	/*
	退款渠道选择：paypal接口/银行转账
	 */

	var refund_method = new(OrderRefundMethodInfo)
	refund_method.Id = PAYMENT_METHOD_PALPAY
	refund_method.Name = "paypal接口"
	refund_method.Description = "paypal接口"
	db_session.Insert(refund_method)

	refund_method = new(OrderRefundMethodInfo)
	refund_method.Id = PAYMENT_METHOD_BANK
	refund_method.Name = "银行转账"
	refund_method.Description = "银行转账"
	db_session.Insert(refund_method)
}
