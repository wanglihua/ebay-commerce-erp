package models

import (
	"matrix/core"
)

type StockBillEntryInfo struct {
	Id         int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e StockBillEntryInfo) TableName() string {
	return TablePrefix + "stock_bill_entry"
}

func (e StockBillEntryInfo) ModelDesc() string {
	return "verbose_name=出入库单明细,entity_json=bill_entry,route=stock/bill/entry"
}
