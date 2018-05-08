package ddlc_models

import (
	"matrix/core"
)

type OrderExternalTransactionInfo struct {
	Id              int64         `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId         int64         `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	TransactionId   string        `xorm:"nvarchar(100) notnull index 'transaction_id'" json:"transaction_id"`
	TransactionTime core.LongTime `xorm:"datetime notnull index 'transaction_time'" json:"transaction_time"`
	PaymentId       int64         `xorm:"bigint notnull index 'payment_id'" json:"payment_id"`
	CreateTime      core.Time     `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime      core.Time     `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderExternalTransactionInfo) TableName() string {
	return "erp_order_external_transaction"
}

func (e OrderExternalTransactionInfo) ModelDesc() string {
	return "verbose_name=order external transaction,entity_json=order_external_transaction,route=order/external/transaction"
}
