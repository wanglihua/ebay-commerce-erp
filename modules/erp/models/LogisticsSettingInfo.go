package models

import (
	"matrix/core"
)

type LogisticsSettingInfo struct {
	Id                 int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	AddressSpecialWord string          `xorm:"nvarchar(300) notnull 'address_special_word'" json:"address_special_word" smith:"verbose_name=地址特殊字段,max=300"` // 如果有多个，用逗号隔开
	MaxDeficitPrice    float64         `xorm:"decimal(18,6) notnull 'max_deficit_price'" json:"max_deficit_price" smith:"verbose_name=最大亏损价格"`
	MaxDeficitRate     float64         `xorm:"decimal(18,6) notnull 'max_deficit_rate'" json:"max_deficit_rate" smith:"verbose_name=最大亏损率"`
	EnvelopeMaxLength  float64         `xorm:"decimal(18,6) notnull 'envelope_max_length'" json:"envelope_max_length" smith:"verbose_name=信封最大长度"`
	EnvelopeMaxWidth   float64         `xorm:"decimal(18,6) notnull 'envelope_max_width'" json:"envelope_max_width" smith:"verbose_name=信封最大宽度"`
	EnvelopeMaxHeight  float64         `xorm:"decimal(18,6) notnull 'envelope_max_height'" json:"envelope_max_height" smith:"verbose_name=信封最大高度"`
	EnvelopeMaxWeight  float64         `xorm:"decimal(18,6) notnull 'envelope_max_weight'" json:"envelope_max_weight" smith:"verbose_name=信封最大重量"`
	MinExpressFee      float64         `xorm:"decimal(18,6) notnull 'min_express_fee'" json:"min_express_fee" smith:"verbose_name=最少快递费"`
	CreateTime         core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime         core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version            int             `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsSettingInfo) TableName() string {
	return TablePrefix + "logistics_setting"
}

func (e LogisticsSettingInfo) ModelDesc() string {
	return "verbose_name=状态,entity_json=logistics_setting,route=logistics/setting"
}
