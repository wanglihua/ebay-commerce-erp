package models

import (
	"matrix/core"
)

type EnvelopeFeeRateInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull 'name'" json:"name" smith:"verbose_name=信封名称"`
	Length     float64         `xorm:"decimal(18,6) notnull 'length'" json:"length" smith:"verbose_name=信封长度（cm）"`
	Width      float64         `xorm:"decimal(18,6) notnull 'width'" json:"width" smith:"verbose_name=信封宽度（cm）"`
	Height     float64         `xorm:"decimal(18,6) notnull 'height'" json:"height" smith:"verbose_name=信封高度（cm）"`
	Weight     float64         `xorm:"decimal(18,6) notnull 'weight'" json:"weight" smith:"verbose_name=信封重量（kg）"`
	Fee        float64         `xorm:"decimal(18,6) notnull 'fee'" json:"fee" smith:"verbose_name=费用（澳币）"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e EnvelopeFeeRateInfo) TableName() string {
	return "erp_envelope_fee_rate"
}

func (e EnvelopeFeeRateInfo) ModelDesc() string {
	return "verbose_name=信封费率,entity_json=envelope_fee_rate,route=envelope/fee/rate"
}
