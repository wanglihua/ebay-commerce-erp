package models

import (
	"matrix/core"
)

type StockBillInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code       string          `xorm:"nvarchar(100) index 'code'" json:"code" smith:"verbose_name=批次号"`
	StockId    int64           `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`        //入库单的入单仓库；出库单的出单仓库；调拨单的调入仓库
	SrcStockId core.NullInt    `xorm:"bigint null index 'src_stock_id'" json:"src_stock_id" smith:"verbose_name=调出仓库"` //调拨单的调出仓库
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注,min=2,max=400"`
	Date       core.Date       `xorm:"datetime notnull 'date'" json:"date" smith:"verbose_name=日期"`
	UserId     int64           `xorm:"bigint notnull 'user_id'" json:"user_id" smith:"verbose_name=操作人"`
	TypeId     int64           `xorm:"bigint notnull 'type_id'" json:"type_id" smith:"verbose_name=类型"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e StockBillInfo) TableName() string {
	return TablePrefix + "stock_bill"
}

func (e StockBillInfo) ModelDesc() string {
	return "verbose_name=出入库单,entity_json=stock_bill,route=stock/bill"
}
