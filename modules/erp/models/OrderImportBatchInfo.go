package models

import (
	"matrix/core"
)

type OrderImportBatchInfo struct {
	Id         int64               `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BatchCode  int64               `xorm:"bigint notnull index 'batch_code'" json:"batch_code" smith:"verbose_name=批次号"`
	OrderId    int64               `xorm:"bigint notnull index 'order_id'" json:"order_id" smith:"verbose_name=订单Id"`
	CreateTime core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                 `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderImportBatchInfo) TableName() string {
	return TablePrefix + "order_import_batch"
}

func (e OrderImportBatchInfo) ModelDesc() string {
	return "verbose_name=订单导入批次,entity_json=order_import_batch,route=order/import/batch"
}
