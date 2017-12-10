package cggl_models

import "matrix/core"

/*
	商品入库长宽高重量
 */

type ProductCheckinWeightInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BillId     int       `xorm:"bigint notnull 'bill_id'" json:"bill_id" smith:"verbose_name=采购单"`
	EntryId    int       `xorm:"bigint notnull 'entry_id'" json:"entry_id" smith:"verbose_name=采购单明细"`
	ProductId  int64     `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	Length     float64   `xorm:"float notnull 'length' default(0)" json:"length" smith:"verbose_name=长度,min=0"` //cm
	Width      float64   `xorm:"float notnull 'width' default(0)" json:"width" smith:"verbose_name=宽度,min=0"`   //cm
	Height     float64   `xorm:"float notnull 'height' default(0)" json:"height" smith:"verbose_name=高度,min=0"` //cm
	Weight     float64   `xorm:"float notnull 'weight' default(0)" json:"weight" smith:"verbose_name=重量,min=0"` //kg
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductCheckinWeightInfo) TableName() string {
	return "erp_product_checkin_weight"
}

func (e ProductCheckinWeightInfo) ModelDesc() string {
	return "verbose_name=商品入库长宽高重量,entity_json=product_checkin_weight,route=product/checkin/weight"
}
