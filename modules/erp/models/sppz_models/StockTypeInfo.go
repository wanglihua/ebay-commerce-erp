package sppz_models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	// STOCK_TYPE_CN 国内仓库
	STOCK_TYPE_CN int64 = 100

	// STOCK_TYPE_ZT 在途仓库
	STOCK_TYPE_ZT int64 = 200

	// STOCK_TYPE_AU 澳洲仓库
	STOCK_TYPE_AU int64 = 300
)

type StockTypeInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name"`
	Desc       string          `xorm:"nvarchar(500) notnull 'desc'" json:"desc"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e StockTypeInfo) TableName() string {
	return "erp_stock_type"
}

func (e StockTypeInfo) ModelDesc() string {
	return "verbose_name=仓库,entity_json=stock_type,route=stock/type"
}

func (e StockTypeInfo) InitData(db_session *xorm.Session) {
	var stock_type_cn StockTypeInfo
	stock_type_cn.Id = STOCK_TYPE_CN
	stock_type_cn.Name = "国内仓库"
	stock_type_cn.Desc = "国内仓库"
	db_session.Insert(&stock_type_cn)

	var stock_type_zaitu StockTypeInfo
	stock_type_zaitu.Id = STOCK_TYPE_ZT
	stock_type_zaitu.Name = "在途仓库"
	stock_type_zaitu.Desc = "在途仓库"
	db_session.Insert(&stock_type_zaitu)

	var stock_type_au StockTypeInfo
	stock_type_au.Id = STOCK_TYPE_AU
	stock_type_au.Name = "澳洲仓库"
	stock_type_au.Desc = "澳洲仓库"
	db_session.Insert(&stock_type_au)
}

func GetStockTypeList(db_session *xorm.Session) []StockTypeInfo {
	var stock_type_list = make([]StockTypeInfo, 0)
	err := db_session.OrderBy("id ASC").Find(&stock_type_list)
	core.HandleError(err)
	return stock_type_list
}
