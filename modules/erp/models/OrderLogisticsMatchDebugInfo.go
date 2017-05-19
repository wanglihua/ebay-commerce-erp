package models

import (
	"matrix/core"
)

type OrderLogisticsMatchDebugInfo struct {
	Id          int64        `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId     int64        `xorm:"bigint notnull index 'order_id'" json:"order_id" smith:"verbose_name=订单"`
	OrderCode   string       `xorm:"nvarchar(200) notnull index 'order_code'" json:"order_code"  smith:"verbose_name=订单编号"` //方便查询
	AddressId   int64        `xorm:"bigint notnull index 'address_id'" json:"address_id" smith:"verbose_name=物流地址"`         // OrderShippingAddressInfo 的主键
	CompanyId   int64        `xorm:"bigint notnull index 'company_id'" json:"company_id" smith:"verbose_name=物流公司"`         // 物流公司Id
	StockId     int64        `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=发货仓库"`             // 发货仓库，仓库的主键
	ModeId      int64        `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流方式"`               // LogisticsModeInfo  的主键
	FeeRateId   int64        `xorm:"bigint notnull index 'fee_rate_id'" json:"fee_rate_id" smith:"verbose_name=物流费率"`       // 如果物流方式没有符合要求的费率，值就是 0
	TotalFee    float64      `xorm:"decimal(18,6) notnull 'total_fee'" json:"total_fee" smith:"verbose_name=物流费用"`          // 由上面的 物流方式和物流费率，计算出来的物流费用，不包括拣货费用
	TotalCost   float64      `xorm:"decimal(18,6) notnull 'total_cost'" json:"total_cost" smith:"verbose_name=成本费用"`        // 由上面的 物流方式和物流费率，计算出来的物流费用 和 仓库成本的总和
	MatchResult string       `xorm:"nvarchar(200) notnull 'match_result'" json:"match_result" smith:"verbose_name=匹配结果"`
	MatchDebug  string       `xorm:"text notnull 'match_debug'" json:"match_debug" smith:"verbose_name=匹配调试信息"`
	CreateTime  core.Time    `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time    `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int          `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderLogisticsMatchDebugInfo) TableName() string {
	return TablePrefix + "order_logistics_match_debug"
}

func (e OrderLogisticsMatchDebugInfo) ModelDesc() string {
	return "verbose_name=订单物流匹配调试,entity_json=match_debug,route=order/logistics/match/debug"
}
