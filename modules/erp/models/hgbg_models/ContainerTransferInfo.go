package hgbg_models

import (
	"matrix/core"
)

type ContainerTransferInfo struct {
	Id                int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	BatchCode         string          `xorm:"nvarchar(200) notnull index 'batch_code'" json:"batch_code" smith:"verbose_name=批次号"`
	ContainerCode     string          `xorm:"nvarchar(200) notnull 'container_code'" json:"container_code" smith:"verbose_name=集装箱号,min=2,max=100"`
	BlNo              string          `xorm:"nvarchar(200) notnull 'bl_no'" json:"bl_no" smith:"verbose_name=提单号,min=2,max=100"`
	SealNo            string          `xorm:"nvarchar(200) notnull 'seal_no'" json:"seal_no" smith:"verbose_name=铅封号,min=2,max=100"`
	TotalVolumn       float64         `xorm:"float64 notnull 'total_volumn'" json:"total_volumn" smith:"verbose_name=总体积"`
	TotalWeight       float64         `xorm:"float64 notnull 'total_weight'" json:"total_weight" smith:"verbose_name=总重量"`
	TotalAmount1      float64         `xorm:"float64 notnull 'total_amount1'" json:"total_amount1" smith:"verbose_name=总金额1"` // 报关税的金额
	TotalAmount2      float64         `xorm:"float64 notnull 'total_amount2'" json:"total_amount2" smith:"verbose_name=总金额2"` // 报采购成本的金额
	PriceUpRate       float64         `xorm:"float64 notnull 'price_up_rate'" json:"price_up_rate" smith:"verbose_name=非征品上调"`
	PriceDownRate     float64         `xorm:"float64 notnull 'price_down_rate'" json:"price_down_rate"  smith:"verbose_name=征税品下调"`
	CostUpRate        float64         `xorm:"float64 notnull 'cost_up_rate'" json:"cost_up_rate"  smith:"verbose_name=成本上调"`
	PackingTime       core.Time       `xorm:"datetime notnull 'packing_time'" json:"packing_time" smith:"verbose_name=装箱时间"`
	DepartureTime     core.Time       `xorm:"datetime notnull 'departure_time'" json:"departure_time" smith:"verbose_name=到港时间"`
	StockTime         core.Time       `xorm:"datetime notnull 'stock_time'" json:"stock_time" smith:"verbose_name=到仓时间"`
	StartPortId       int64           `xorm:"bigint notnull index 'start_port_id'" json:"start_port_id" smith:"verbose_name=起始港口"`
	StartPortName     string          `xorm:"nvarchar(200) notnull default('') 'start_port_name'" json:"start_port_name" smith:"verbose_name=起始港口"`
	DestPortId        int64           `xorm:"bigint notnull index 'dest_port_id'" json:"dest_port_id" smith:"verbose_name=目的港口"`
	DestPortName      string          `xorm:"nvarchar(200) notnull default('') 'dest_port_name'" json:"dest_port_name" smith:"verbose_name=目的港口"`
	StatusId          int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	ImportCompanyCode string          `xorm:"nvarchar(20) notnull default('') 'import_company_code'" json:"import_company_code" smith:"verbose_name=澳洲清关公司代码"`
	ImportCompanyName string          `xorm:"nvarchar(200) notnull default('') 'import_company_name'" json:"import_company_name" smith:"verbose_name=澳洲清关公司"`
	Description       string          `xorm:"nvarchar(800) notnull default('') 'description'" json:"description"`
	ExchangeRate      float64         `xorm:"decimal(18,6) notnull 'exchange_rate'" json:"exchange_rate"` // 订单下载时把当前的汇率设置放到这个里面
	CreateTime        core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime        core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ContainerTransferInfo) TableName() string {
	return "erp_container_transfer"
}

func (e ContainerTransferInfo) ModelDesc() string {
	return "verbose_name=集装箱调拨,entity_json=transfer,route=container/transfer"
}
