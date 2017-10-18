package cggl_models

import (
	"matrix/core"
)

/*
	智能补货
	Intelligent replenishment
 */

type IntelReplInfo struct {
	Id           int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId    int64     `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	StockId      int64     `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=仓库"`
	Qty          int64     `xorm:"bigint notnull  'qty'" json:"qty" smith:"verbose_name=数量"`
	Level1Profit float64   `xorm:"decimal(18,6) notnull 'level1_profit'" json:"level1_profit" smith:"verbose_name=1档毛利率 "`
	Level2Profit float64   `xorm:"decimal(18,6) notnull 'level2_profit'" json:"level2_profit" smith:"verbose_name=2档毛利率 "`
	Level3Profit float64   `xorm:"decimal(18,6) notnull 'level3_profit'" json:"level3_profit" smith:"verbose_name=3档毛利率 "`
	StatusId     int       `xorm:"bigint notnull 'status_id'" json:"status_id" smith:"verbose_name=状态"` // 新生成，已采购等
	CreateTime   core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime   core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e IntelReplInfo) TableName() string {
	return "erp_intel_repl"
}

func (e IntelReplInfo) ModelDesc() string {
	return "verbose_name=智能补货,entity_json=intel_repl,route=intel/repl"
}
