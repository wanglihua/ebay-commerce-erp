package models

import (
	"matrix/core"
)

type OrderInfo struct {
	Id                       int64               `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code                     string              `xorm:"nvarchar(200) notnull index unique 'code'" json:"code" smith:"verbose_name=编码"`
	Desc                     core.NullString     `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
	TypeId                   int64               `xorm:"bigint notnull index 'type_id'" json:"type_id" smith:"verbose_name=类型"`
	EBayStatus               string              `xorm:"nvarchar(200) notnull 'ebay_status'" json:"ebay_status" smith:"verbose_name=ebay状态"`
	StatusId                 int64               `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	AdjustmentAmountCurrency string              `xorm:"nvarchar(10) notnull 'adjustment_amount_currency'" json:"adjustment_amount_currency" smith:"verbose_name=调整金额货种"`
	AdjustmentAmount         float64             `xorm:"decimal(18,6) notnull 'adjustment_amount'" json:"adjustment_amount" smith:"verbose_name=调整金额"`
	AmountPaidCurrency       string              `xorm:"nvarchar(10) notnull 'amount_paid_currency'" json:"amount_paid_currency" smith:"verbose_name=支付金额货种"`
	AmountPaid               float64             `xorm:"decimal(18,6) notnull 'amount_paid'" json:"amount_paid" smith:"verbose_name=支付金额"`
	AmountSavedCurrency      string              `xorm:"nvarchar(10) notnull 'amount_saved_currency'" json:"amount_saved_currency" smith:"verbose_name=节省金额货种"`
	AmountSaved              float64             `xorm:"decimal(18,6) notnull 'amount_saved'" json:"amount_saved" smith:"verbose_name=节省金额"`
	CreatedTime              core.NullLongTime   `xorm:"datetime 'created_time'" json:"created_time" smith:"verbose_name=创建时间"`
	ModifiedTime             core.NullLongTime   `xorm:"datetime 'modified_time'" json:"modified_time" smith:"verbose_name=修改时间"`
	PaidTime                 core.NullLongTime   `xorm:"datetime 'paid_time'" json:"paid_time" smith:"verbose_name=支付时间"`
	ShippedTime              core.NullLongTime   `xorm:"datetime 'shipped_time'" json:"shipped_time" smith:"verbose_name=创建时间"`
	HoldTypeId               core.NullInt        `xorm:"bigint null index 'hold_type_id'" json:"hold_type_id" smith:"verbose_name=Hold类型"`
	CreateTime               core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime               core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version                  int                 `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderInfo) TableName() string {
	return TablePrefix + "order"
}

func (e OrderInfo) ModelDesc() string {
	return "verbose_name=订单,entity_json=order,route=order"
}
