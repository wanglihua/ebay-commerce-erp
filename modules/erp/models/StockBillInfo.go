package models

import (
	"matrix/core"
)

type StockBillInfo struct {
	Id         int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CreateTime core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e StockBillInfo) TableName() string {
	return TablePrefix + "stock_bill"
}

func (e StockBillInfo) ModelDesc() string {
	return "verbose_name=出入库单,entity_json=stock_bill,route=stock/bill"
}
