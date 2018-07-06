package ddbb_models

import (
	"matrix/core"
)

type LogisticsReportInfo struct {
	Id             int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	DayDate        core.Date `xorm:"date notnull index 'day_date'" json:"day_date"`
	OrderId        int64     `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	ShopId         int64     `xorm:"bigint notnull index 'shop_id'" json:"shop_id"`
	StockId        int64     `xorm:"bigint notnull index 'stock_id'" json:"stock_id"`
	CompanyId      int64     `xorm:"bigint notnull index 'company_id'" json:"company_id"`
	ModeId         int64     `xorm:"bigint notnull index 'mode_id'" json:"mode_id"`
	Postcode       string    `xorm:"nvarchar(200) notnull index 'postcode'" json:"postcode"`
	PrdCount       int64     `xorm:"bigint notnull 'prd_count'" json:"prd_count"`
	Volume         float64   `xorm:"decimal(18,3) notnull 'volume'" json:"volume"`
	Weight         float64   `xorm:"decimal(18,3) notnull 'weight'" json:"weight"`
	Amount         float64   `xorm:"decimal(18,2) notnull 'amount'" json:"amount"`
	TotalFee       float64   `xorm:"decimal(18,2) notnull 'total_fee'" json:"total_fee"`
	ActualTotalFee float64   `xorm:"decimal(18,2) notnull 'actual_total_fee'" json:"actual_total_fee"`
	Profit         float64   `xorm:"decimal(18,2) notnull 'profit'" json:"profit"`
	CreateTime     core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime     core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e LogisticsReportInfo) TableName() string {
	return "erp_logistics_report"
}

func (e LogisticsReportInfo) ModelDesc() string {
	return "verbose_name=物流报表,entity_json=logistics_report,route=logistics/report"
}

/*

OrderCode      string    `xorm:"order_code notnull 'order_code'" json:"order_code"`
StockName      string    `xorm:"nvarchar(200) notnull 'stock_name'" json:"stock_name"`
CompanyName    string    `xorm:"nvarchar(200) notnull 'company_name'" json:"company_name"`
ModeName       string    `xorm:"nvarchar(200) notnull 'mode_name'" json:"mode_name"`

 */
