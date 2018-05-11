package ddlc_models

import (
	"matrix/core"
)

type OrderRefundInfo struct {
	Id            int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId       int64     `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	TransactionId string    `xorm:"nvarchar(100) notnull index 'transaction_id'" json:"transaction_id"`
	Amount        float64   `xorm:"decimal(18,6) notnull 'amount'" json:"amount"`
	Description   string    `xorm:"nvarchar(500) notnull 'description'" json:"description"`
	CreateTime    core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e OrderRefundInfo) TableName() string {
	return "erp_order_refund"
}

func (e OrderRefundInfo) ModelDesc() string {
	return "verbose_name=订单退款信息,entity_json=order_refund,route=order/refund"
}
