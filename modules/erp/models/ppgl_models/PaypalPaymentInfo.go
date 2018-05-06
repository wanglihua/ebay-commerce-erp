package ppgl_models

import (
	"matrix/core"
)

type PaypalPaymentInfo struct {
	Id               int64         `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	TransactionId    string        `xorm:"nvarchar(200) notnull index unique 'transaction_id'" json:"transaction_id"`
	Type             string        `xorm:"nvarchar(100) notnull 'type'" json:"type"`
	Payer            string        `xorm:"nvarchar(200) notnull 'payer'" json:"payer"`
	PayerDisplayName string        `xorm:"nvarchar(200) notnull 'payer_display_name'" json:"payer_display_name"`
	Currency         string        `xorm:"nvarchar(10) notnull 'currency'" json:"currency"`
	GrossAmount      float64       `xorm:"decimal(18,6) notnull 'gross_amount'" json:"gross_amount"`
	FeeAmount        float64       `xorm:"decimal(18,6) notnull 'fee_amount'" json:"fee_amount"`
	NetAmount        float64       `xorm:"decimal(18,6) notnull 'net_amount'" json:"net_amount"`
	Status           string        `xorm:"nvarchar(100) notnull index 'status'" json:"status"`
	PayTime          core.LongTime `xorm:"datetime notnull 'pay_time'" json:"pay_time"`
	CreateTime       core.Time     `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime       core.Time     `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PaypalPaymentInfo) TableName() string {
	return "erp_paypal_payment"
}

func (e PaypalPaymentInfo) ModelDesc() string {
	return "verbose_name=Paypal付款,entity_json=paypal_payment,route=paypal/payment"
}
