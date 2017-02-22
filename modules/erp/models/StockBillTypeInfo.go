package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	// STOCK_BILL_TYPE_MANUAL_CHECKIN 手工入库单
	STOCK_BILL_TYPE_MANUAL_CHECKIN = 100

	// STOCK_BILL_TYPE_PURCHASE_CHECKIN 采购入库单
	STOCK_BILL_TYPE_PURCHASE_CHECKIN = 150

	// STOCK_BILL_TYPE_MANUAL_CHECKOUT 手工出库单
	STOCK_BILL_TYPE_MANUAL_CHECKOUT = 200

	// STOCK_BILL_TYPE_TRANSFER 调拨单
	STOCK_BILL_TYPE_TRANSFER = 300

	// STOCK_BILL_TYPE_TRANSFER_ZT 调拨单 在途
	STOCK_BILL_TYPE_TRANSFER_ZT = 400
)

type StockBillTypeInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,min=2,max=400"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e StockBillTypeInfo) TableName() string {
	return TablePrefix + "stock_bill_type"
}

func (e StockBillTypeInfo) ModelDesc() string {
	return "verbose_name=出入库单类型,entity_json=bill_type,route=stock/bill/type"
}

func (e StockBillTypeInfo) InitData(db_session *xorm.Session) {
	bill_type := new(StockBillTypeInfo)
	bill_type.Id = STOCK_BILL_TYPE_MANUAL_CHECKIN
	bill_type.Name = "手工入库单"
	bill_type.Desc = core.NewNullString("手工入库单", true)
	db_session.Insert(bill_type)

	bill_type = new(StockBillTypeInfo)
	bill_type.Id = STOCK_BILL_TYPE_MANUAL_CHECKOUT
	bill_type.Name = "手工出库单"
	bill_type.Desc = core.NewNullString("手工出库单", true)
	db_session.Insert(bill_type)

	bill_type = new(StockBillTypeInfo)
	bill_type.Id = STOCK_BILL_TYPE_TRANSFER
	bill_type.Name = "调拨单"
	bill_type.Desc = core.NewNullString("调拨单", true)
	db_session.Insert(bill_type)

	bill_type = new(StockBillTypeInfo)
	bill_type.Id = STOCK_BILL_TYPE_TRANSFER_ZT
	bill_type.Name = "调拨单在途"
	bill_type.Desc = core.NewNullString("调拨单 在途", true)
	db_session.Insert(bill_type)
}
