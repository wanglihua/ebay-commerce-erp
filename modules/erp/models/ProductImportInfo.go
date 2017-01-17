package models

import (
	"matrix/core"
)

type ProductImportInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Sku         string    `xorm:"nvarchar(200) notnull index 'sku'" json:"sku" smith:"verbose_name=SKU,min=2,max=100"`
	Name        string    `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	EnglishName string    `xorm:"nvarchar(200) notnull index 'english_name'" json:"english_name" smith:"verbose_name=英文名称,min=2,max=100"`
	CostPrice   float64   `xorm:"float notnull 'cost_price'" json:"cost_price" smith:"verbose_name=成本价,min=0"`
	Weight      float64   `xorm:"float notnull 'weight'" json:"weight" smith:"verbose_name=重量,min=0"`
	UserId      int64     `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=导入人"`
	ImportTime  core.Time `xorm:"created notnull 'import_time'" json:"import_time"`
}

func (e ProductImportInfo) TableName() string {
	return TablePrefix + "product_import"
}

func (e ProductImportInfo) ModelDesc() string {
	return "verbose_name=商品导入,entity_json=product_import,route=product/import"
}
