package models

import "matrix/core"

type ContainerTransferEntryInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	TransferId int64           `xorm:"bigint notnull index 'transfer_id'" json:"transfer_id" smith:"verbose_name=集装箱调拨Id"`
	ProductId  int64           `xorm:"bigint notnull 'transfer_id'" json:"transfer_id" smith:"verbose_name=集装箱调拨Id"`
	ProductSku string          `xorm:"nvarchar(200) notnull 'product_sku'" json:"product_sku" smith:"verbose_name=SKU"`
	Qty        int64           `xorm:"bigint notnull 'qty'" json:"qty" smith:"verbose_name=数量"`
	BoxQty     int64           `xorm:"bigint notnull 'box_qty'" json:"box_qty" smith:"verbose_name=箱数"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ContainerTransferEntryInfo) TableName() string {
	return TablePrefix + "container_transfer_entry"
}

func (e ContainerTransferEntryInfo) ModelDesc() string {
	return "verbose_name=集装箱调拨明细,entity_json=transfer_entry,route=container/transfer/entry"
}
