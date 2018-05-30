package kcgl_models

import (
	"matrix/core"
)

type ProductLogisticsInfo struct {
	Id            int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId     int64     `xorm:"bigint notnull index 'product_id'" json:"product_id"`
	ProductCode   string    `xorm:"nvarchar(100) notnull index 'product_code'" json:"product_code"`
	Qty           int64     `xorm:"bigint notnull 'qty'" json:"qty"`                        // 数量
	PostalCode    string    `xorm:"nvarchar(100) notnull 'postal_code'" json:"postal_code"` // 录入邮编
	StockId       int64     `xorm:"bigint notnull index 'stock_id'" json:"stock_id"`
	StockName     string    `xorm:"nvarchar(100) notnull 'stock_name'" json:"stock_name"`
	CompanyId     int64     `xorm:"bigint notnull index 'company_id'" json:"company_id"`
	CompanyName   string    `xorm:"nvarchar(100) notnull 'company_name'" json:"company_name"`
	ModeId        int64     `xorm:"bigint notnull index 'mode_id'" json:"mode_id"`
	ModeName      string    `xorm:"nvarchar(100) notnull 'mode_name'" json:"mode_name"`
	FeePostalCode string    `xorm:"nvarchar(100) notnull 'fee_postal_code'" json:"fee_postal_code"` // 费用邮编
	Fee           float64   `xorm:"decimal(18,6) notnull 'fee'" json:"fee"`                         // 物流费用
	Cost          float64   `xorm:"decimal(18,6) notnull 'cost'" json:"cost"`                       // 成本费用
	UserId        int64     `xorm:"bigint notnull index 'user_id'" json:"user_id"`
	LoginName     string    `xorm:"nvarchar(100) notnull index 'login_name'" json:"login_name"`
	CreateTime    core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time `xorm:"created notnull 'update_time'" json:"update_time"`
}

func (e ProductLogisticsInfo) TableName() string {
	return "erp_product_logistics"
}

func (e ProductLogisticsInfo) ModelDesc() string {
	return "verbose_name=商品物流费用,entity_json=product_logistics,route=product/logistics"
}
