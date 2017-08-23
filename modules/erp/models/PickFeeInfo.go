package models

import (
	"matrix/core"
)

// 拣货人工费用的设置
type PickFeeInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Fee1       float64         `xorm:"decimal(18,6) notnull 'fee1'" json:"fee1" smith:"verbose_name=5Kg以下"`
	Fee2       float64         `xorm:"decimal(18,6) notnull 'fee2'" json:"fee2" smith:"verbose_name=5.1KG-10KG"`
	Fee3       float64         `xorm:"decimal(18,6) notnull 'fee3'" json:"fee3" smith:"verbose_name=10.1KG-25KG"`
	Fee4       float64         `xorm:"decimal(18,6) notnull 'fee4'" json:"fee4" smith:"verbose_name=25.1KG-50KG"`
	Fee5       float64         `xorm:"decimal(18,6) notnull 'fee5'" json:"fee5" smith:"verbose_name=50.1KG以上"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PickFeeInfo) TableName() string {
	return "erp_pick_fee"
}

func (e PickFeeInfo) ModelDesc() string {
	return "verbose_name=拣货费用,entity_json=pick_fee,route=pick/fee"
}

/*
拣货费是人工成本，这个可以在配置里面设置下。目前我们的设置是按照订单重量划分的 ： 5Kg以下：2AUD ； 5.1KG-10KG 是3澳币；10.1KG-25KG 是4澳币；25.1KG-50KG 是5澳币；50.1KG以上 是10澳币
 */
