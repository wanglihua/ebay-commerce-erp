package models

import (
	"matrix/core"
)

type ProductImportBatchInfo struct {
	Id         int64               `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BatchCode  int64               `xorm:"bigint notnull index 'batch_code'" json:"batch_code" smith:"verbose_name=批次号"`
	ProductId  int64               `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品Id"`
	CreateTime core.Time           `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time           `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int                 `xorm:"version notnull 'version'" json:"version"`
}

func (e ProductImportBatchInfo) TableName() string {
	return TablePrefix + "product_import_batch"
}

func (e ProductImportBatchInfo) ModelDesc() string {
	return "verbose_name=商品导入批次,entity_json=product_import_batch,route=product/import/batch"
}
