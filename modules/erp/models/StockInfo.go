package models

import (
	"matrix/core"
)

const (
	// STOCK_TYPE_CN 国内仓库
	STOCK_TYPE_CN int64 = 100

	// STOCK_TYPE_ZAITU 在途仓库
	STOCK_TYPE_ZT int64 = 200

	// STOCK_TYPE_AU 澳洲仓库
	STOCK_TYPE_AU int64 = 300
)

type StockInfo struct {
	Id            int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name          string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	ShortName     core.NullString `xorm:"nvarchar(200) 'short_name'" json:"short_name" smith:"verbose_name=简称,min=2,max=100"`
	Contact       core.NullString `xorm:"nvarchar(200) 'contact'" json:"contact" smith:"verbose_name=联系人,min=2,max=100"`
	Address       core.NullString `xorm:"nvarchar(200) 'address'" json:"address" smith:"verbose_name=地址,min=2,max=200"`
	Phone         core.NullString `xorm:"nvarchar(200) 'phone'" json:"phone" smith:"verbose_name=电话,min=2,max=100"`
	StockTypeId   int64           `xorm:"bigint notnull 'stock_type_id'" json:"stock_type_id" smith:"verbose_name=仓库类型"`  // 国内仓库 在途仓库 澳洲仓库
	TargetStockId core.NullInt    `xorm:"bigint null 'target_stock_id'" json:"target_stock_id" smith:"verbose_name=目的仓库"` //在途虚拟仓库的目的仓库
	StatusId      int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime    core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version       int             `xorm:"version notnull 'version'" json:"version"`
}

func (e StockInfo) TableName() string {
	return TablePrefix + "stock"
}

func (e StockInfo) ModelDesc() string {
	return "verbose_name=仓库,entity_json=stock,route=stock"
}

type StockTypeInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func GetStockTypeList() []StockTypeInfo {
	stock_type_list := make([]StockTypeInfo, 0)

	var stock_type_cn StockTypeInfo
	stock_type_cn.Id = STOCK_TYPE_CN
	stock_type_cn.Name = "国内仓库"

	var stock_type_zaitu StockTypeInfo
	stock_type_zaitu.Id = STOCK_TYPE_ZT
	stock_type_zaitu.Name = "在途仓库"

	var stock_type_au StockTypeInfo
	stock_type_au.Id = STOCK_TYPE_AU
	stock_type_au.Name = "澳洲仓库"

	return append(stock_type_list, stock_type_cn, stock_type_zaitu, stock_type_au)
}
