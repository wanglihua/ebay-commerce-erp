package ddlc_models

import (
	"matrix/core"
)


type OrderStockBestLogisticsInfo struct {
	Id         int64      `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId    int64      `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	OrderCode  string     `xorm:"nvarchar(200) notnull index 'order_code'" json:"order_code"`
	StockId    int64      `xorm:"bigint notnull 'stock_id'" json:"stock_id"`
	StockName  string     `xorm:"nvarchar(200) notnull 'stock_name'" json:"stock_name"`
	Fee        float64    `xorm:"decimal(18,6) notnull 'fee'" json:"fee" smith:"verbose_name=物流费用"` // 由物流方式和物流费率，计算出来的物流费用，不包括拣货费用
	CreateTime core.Time  `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time  `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderStockBestLogisticsInfo) TableName() string {
	return "erp_order_stock_best_logistics"
}

func (e OrderStockBestLogisticsInfo) ModelDesc() string {
	return "verbose_name=订单仓库最佳物流,entity_json=best_logistics,route=order/stock/best/logistics"
}
