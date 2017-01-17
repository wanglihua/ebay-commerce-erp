package models

import (
	"matrix/core"
)

type ProductInfo struct {
	Id          int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Sku         string          `xorm:"nvarchar(200) notnull index unique 'sku'" json:"sku" smith:"verbose_name=SKU,min=2,max=100"`
	Name        string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	EnglishName string          `xorm:"nvarchar(200) notnull index 'english_name'" json:"english_name" smith:"verbose_name=英文名称,min=2,max=100"`
	CateId      core.NullInt    `xorm:"bigint index 'cate_id'" json:"cate_id" smith:"verbose_name=商品目录"`
	CostPrice   float64         `xorm:"float notnull 'cost_price'" json:"cost_price" smith:"verbose_name=成本价,min=0"` //其实就是采购价
	Weight      float64         `xorm:"float notnull 'weight'" json:"weight" smith:"verbose_name=重量,min=0"`
	SupplierId  core.NullInt    `xorm:"bigint index 'supplier_id'" json:"supplier_id" smith:"verbose_name=供应商"`
	Desc        core.NullString `xorm:"nvarchar(500) 'bzxx'" json:"bzxx" smith:"verbose_name=备注,max=500"`
	StatusId    int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int             `xorm:"version notnull 'version'" json:"version"`
}

func (e ProductInfo) TableName() string {
	return TablePrefix + "product"
}

func (e ProductInfo) ModelDesc() string {
	return "verbose_name=商品,entity_json=product,route=product"
}
