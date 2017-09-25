package ddlc_models

import (
	"matrix/core"
)

type OrderLogisticsLogicInfo struct {
	Id            int64   `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId       int64   `xorm:"bigint notnull index 'order_id'" json:"order_id" smith:"verbose_name=订单"`
	AddressId     int64   `xorm:"bigint notnull index 'address_id'" json:"address_id" smith:"verbose_name=物流地址"`           // OrderShippingAddressInfo 的主键，手工创建可能为空
	CompanyId     int64   `xorm:"bigint notnull index 'company_id'" json:"company_id" smith:"verbose_name=物流公司"`           // 物流公司Id 等的信息，如果手工建的物流信息，这个字段是通过下拉框选择的
	StockId       int64   `xorm:"bigint notnull index 'stock_id'" json:"stock_id" smith:"verbose_name=发货仓库"`               // 发货仓库，仓库的主键
	ModeId        int64   `xorm:"bigint notnull index 'mode_id'" json:"mode_id" smith:"verbose_name=物流方式"`                 // LogisticsModeInfo  的主键，这个表里面也包含了 发货仓库Id， 物流公司Id 等的信息，手工创建可能为空
	FeeRateId     int64   `xorm:"bigint notnull index 'fee_rate_id'" json:"fee_rate_id" smith:"verbose_name=物流费率"`         // LogisticsFeeRateInfo  的主键， 手工创建可能为空
	Length        float64 `xorm:"float notnull 'length'" json:"length" smith:"verbose_name=长度,min=0"`                      // cm
	Width         float64 `xorm:"float notnull 'width'" json:"width" smith:"verbose_name=宽度,min=0"`                        // cm
	Height        float64 `xorm:"float notnull 'height'" json:"height" smith:"verbose_name=高度,min=0"`                      // cm
	Weight        float64 `xorm:"float notnull 'weight'" json:"weight" smith:"verbose_name=重量,min=0"`                      // kg 就是实际重量
	DimWeight     float64 `xorm:"float notnull 'dim_weight'" json:"dim_weight" smith:"verbose_name=体积重,min=0"`             // kg
	Sfkltjz       int64   `xorm:"bigint notnull 'sfkltjz'" json:"sfkltjz" smith:"verbose_name=是否考虑体积重"`                    // -1 否， 1 是
	TotalFee      float64 `xorm:"decimal(18,6) notnull 'total_fee'" json:"total_fee" smith:"verbose_name=物流费用"`            // 由上面的 物流方式和物流费率，计算出来的物流费用，不包括拣货费用
	PickFee       float64 `xorm:"decimal(18,6) notnull 'pick_fee'" json:"pick_fee" smith:"verbose_name=拣货费用"`              // 拣货费用，没有算在total_fee里面 total_fee里面只有物流公司的运费，匹配完成，匹配后打折完之后计算
	CostAmount    float64 `xorm:"decimal(18,6) notnull 'cost_amount'" json:"cost_amount" smith:"verbose_name=商品成本"`        //商品成本，整个订单的商品成本总和
	EnvelopeFeeId int64   `xorm:"bigint notnull index 'envelope_fee_id'" json:"envelope_fee_id" smith:"verbose_name=信封费率"` // 信封费率id，对于信封发货的，该字段才有设置
	TrackingNum   string  `xorm:"nvarchar(100) null 'tracking_num'" json:"tracking_num" smith:"verbose_name=运单号"`
	Ptfl          float64 `xorm:"decimal(18,6) notnull 'ptfl'" json:"ptfl"` // 平台费率
	Jyfl          float64 `xorm:"decimal(18,6) notnull 'jyfl'" json:"jyfl"` // 交易费率

	CompanyIdSecond int64   `xorm:"bigint notnull index default(0) 'company_id_second'" json:"company_id_second" smith:"verbose_name=次佳物流公司"`
	StockIdSecond   int64   `xorm:"bigint notnull index default(0) 'stock_id_second'" json:"stock_id_second" smith:"verbose_name=次佳发货仓库"`
	ModeIdSecond    int64   `xorm:"bigint notnull index default(0) 'mode_id_second'" json:"mode_id_second" smith:"verbose_name=次佳物流方式"`
	FeeRateIdSecond int64   `xorm:"bigint notnull index default(0) 'fee_rate_id_second'" json:"fee_rate_id_second" smith:"verbose_name=次佳物流费率"`
	TotalFeeSecond  float64 `xorm:"decimal(18,6) notnull default(0) 'total_fee_second'" json:"total_fee_second" smith:"verbose_name=次佳物流费用"`
	FeeDiff         float64 `xorm:"decimal(18,6) notnull default(0) 'fee_diff'" json:"fee_diff" smith:"verbose_name=费用偏差"`

	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderLogisticsLogicInfo) TableName() string {
	return "erp_order_logistics_logic"
}

func (e OrderLogisticsLogicInfo) ModelDesc() string {
	return "verbose_name=订单物流理论值,entity_json=order_logistics,route=order/logistics/logic"
}
