package xtpz_models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

const (
	CURRENCY_RMB int64 = 1010 // 人民币
	CURRENCY_USD int64 = 1020 // 美元
	CURRENCY_AUD int64 = 1030 // 澳元
)

type CurrencyInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code       string          `xorm:"nvarchar(200) notnull index 'code'" json:"code" smith:"verbose_name=符号"`
	Name       string          `xorm:"nvarchar(200) notnull 'name'" json:"name" smith:"verbose_name=名称"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
}

func (e CurrencyInfo) TableName() string {
	return "erp_currency"
}

func (e CurrencyInfo) ModelDesc() string {
	return "verbose_name=货币,entity_json=currency,route=currency"
}

func (e CurrencyInfo) InitData(db_session *xorm.Session) {
	currency := new(CurrencyInfo)
	currency.Id = CURRENCY_RMB
	currency.Code = "RMB"
	currency.Name = "人民币"
	currency.Desc = core.MakeNullString("人民币", true)
	db_session.Insert(currency)

	currency = new(CurrencyInfo)
	currency.Id = CURRENCY_USD
	currency.Code = "USD"
	currency.Name = "美元"
	currency.Desc = core.MakeNullString("美元", true)
	db_session.Insert(currency)

	currency = new(CurrencyInfo)
	currency.Id = CURRENCY_AUD
	currency.Code = "AUD"
	currency.Name = "澳元"
	currency.Desc = core.MakeNullString("澳元", true)
	db_session.Insert(currency)
}
