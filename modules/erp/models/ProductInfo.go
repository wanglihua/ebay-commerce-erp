package models

import (
	"matrix/core"
)

type ProductInfo struct {
	Id            int64                 `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Sku           string                `xorm:"nvarchar(200) notnull index 'sku'" json:"sku" smith:"verbose_name=SKU,min=2,max=100"`
	Name          string                `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	EnglishName   string                `xorm:"nvarchar(200) notnull index 'english_name'" json:"english_name" smith:"verbose_name=英文名称,min=2,max=100"`
	BrandId       core.NullInt          `xorm:"bigint index 'brand_id'" json:"brand_id" smith:"verbose_name=品牌"`
	CateId        core.NullInt          `xorm:"bigint index 'cate_id'" json:"cate_id" smith:"verbose_name=商品目录"`
	CostPrice     float64               `xorm:"float notnull 'cost_price'" json:"cost_price" smith:"verbose_name=成本价,min=0"`
	Quantity      int64                 `xorm:"bigint notnull 'quantity'" json:"quantity" smith:"verbose_name=库存,min=0"`
	Weight        float64               `xorm:"float notnull 'weight'" json:"weight" smith:"verbose_name=重量,min=0"`
	StockId       core.NullInt          `xorm:"bigint index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	StorageLocId  core.NullInt          `xorm:"bigint index 'storage_loc_id'" json:"storage_loc_id" smith:"verbose_name=仓位"`
	SupplierId    core.NullInt          `xorm:"bigint index 'supplier_id'" json:"supplier_id" smith:"verbose_name=默认供应商"`
	SupplierPrice float64               `xorm:"float notnull 'supplier_price'" json:"supplier_price" smith:"verbose_name=供应商报价,min=0"`
	StatusId      int64                 `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	Cgts          core.NullInt          `xorm:"bigint 'cgts'" json:"cgts" smith:"verbose_name=采购天数"`
	Cgsx          core.NullInt          `xorm:"bigint 'cgsx'" json:"cgsx" smith:"verbose_name=采购上限"`
	XsyId         core.NullInt          `xorm:"bigint index 'xsy_id'" json:"xsy_id" smith:"verbose_name=销售员"`
	Xsxx          core.NullString       `xorm:"nvarchar(500) 'xsxx'" json:"xsxx" smith:"verbose_name=销售信息,max=500"`
	Ywmjl         core.NullString       `xorm:"nvarchar(500) 'ywmjl'" json:"ywmjl" smith:"verbose_name=英文名记录,max=500"`
	Zjbz          core.NullString       `xorm:"nvarchar(500) 'zjbz'" json:"zjbz" smith:"verbose_name=质检标准,max=500"`
	Bzxx          core.NullString       `xorm:"nvarchar(500) 'bzxx'" json:"bzxx" smith:"verbose_name=备注信息,max=500"`
	CreateTime    core.Time             `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time             `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version       int                   `xorm:"version notnull 'version'" json:"version"`
}

func (e ProductInfo) TableName() string {
	return TablePrefix + "product"
}

func (e ProductInfo) ModelDesc() string {
	return "verbose_name=商品,entity_json=product,route=product"
}

/*

StockName     string                `xorm:"nvarchar(200) notnull index 'stock_name'" json:"stock_name" smith:"verbose_name=库存名称,min=2,max=100"`
RetailPrice   float64               `xorm:"float notnull 'retail_price'" json:"retail_price" smith:"verbose_name=零售价,min=0"`
Zxcgl         core.NullInt          `xorm:"bigint 'zxcgl'" json:"zxcgl" smith:"verbose_name=最小采购量"`
Kcjjts        core.NullInt          `xorm:"bigint 'kcjjts'" json:"kcjjts" smith:"verbose_name=库存警戒天数"`
CgyId         core.NullInt          `xorm:"bigint index 'cgy_id'" json:"cgy_id" smith:"verbose_name=采购员"`
PhyId         core.NullInt          `xorm:"bigint index 'phy_id'" json:"phy_id" smith:"verbose_name=配货员"`
ProducingArea core.NullString       `xorm:"nvarchar(200) 'producing_area'" json:"producing_area" smith:"verbose_name=产地,max=100"`
 */