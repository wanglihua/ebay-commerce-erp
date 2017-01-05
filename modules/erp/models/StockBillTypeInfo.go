package models

import (
	"matrix/core"
)

const (
	// STOCK_BILL_TYPE_CHECKIN 入库单
	STOCK_BILL_TYPE_CHECKIN = 100

	// STOCK_BILL_TYPE_CHECKOUT 出库单
	STOCK_BILL_TYPE_CHECKOUT = 200

	// STOCK_BILL_TYPE_TRANSFER 调拨单
	STOCK_BILL_TYPE_TRANSFER = 300
)


type StockBillTypeInfo struct {
	Id         int64                 `xorm:"bigint notnull pk 'id'" json:"id"`
	Name       string                `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString       `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,min=2,max=400"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e StockBillTypeInfo) TableName() string {
	return TablePrefix + "stock_bill_type"
}

func (e StockBillTypeInfo) ModelDesc() string {
	return "verbose_name=出入库单类型,entity_json=bill_type,route=stock/bill/type"
}
