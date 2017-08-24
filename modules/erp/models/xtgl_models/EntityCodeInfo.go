package xtgl_models

import (
	"matrix/core"
)

const (
	ENTITY_ID_BILL = 100100 // 入库单，出库单，调拨单，盘点单的统一单号

)

type EntityCodeInfo struct {
	Id       int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	EntityId string          `xorm:"bigint notnull index unique(unq_entity_code) 'entity_id'" json:"entity_id" smith:"verbose_name=实体Id"`
	Serial   int64           `xorm:"bigint notnull 'serial'" json:"serial" smith:"verbose_name=序列号"`
	Date     core.NullDate   `xorm:"datetime null index unique(unq_entity_code) 'date'" json:"date" smith:"verbose_name=日期"`
}

func (e EntityCodeInfo) TableName() string {
	return "erp_entity_code"
}

func (e EntityCodeInfo) ModelDesc() string {
	return "verbose_name=实体编码,entity_json=entity_code,route=entity_code"
}
