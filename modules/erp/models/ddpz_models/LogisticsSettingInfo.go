package ddpz_models

import (
	"matrix/core"
)

type LogisticsSettingInfo struct {
	Id                 int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	AddressSpecialWord string          `xorm:"nvarchar(300) notnull 'address_special_word'" json:"address_special_word" smith:"verbose_name=地址特殊字段,max=300"` // 如果有多个，用逗号隔开
	MaxDeficitPrice    float64         `xorm:"decimal(18,6) notnull 'max_deficit_price'" json:"max_deficit_price" smith:"verbose_name=最大亏损价格"`               //澳币？？
	MaxDeficitRate     float64         `xorm:"decimal(18,6) notnull 'max_deficit_rate'" json:"max_deficit_rate" smith:"verbose_name=最大亏损率"`
	EnvelopeMaxLength  float64         `xorm:"decimal(18,6) notnull 'envelope_max_length'" json:"envelope_max_length" smith:"verbose_name=信封最大长度"`
	EnvelopeMaxWidth   float64         `xorm:"decimal(18,6) notnull 'envelope_max_width'" json:"envelope_max_width" smith:"verbose_name=信封最大宽度"`
	EnvelopeMaxHeight  float64         `xorm:"decimal(18,6) notnull 'envelope_max_height'" json:"envelope_max_height" smith:"verbose_name=信封最大高度"`
	EnvelopeMaxWeight  float64         `xorm:"decimal(18,6) notnull 'envelope_max_weight'" json:"envelope_max_weight" smith:"verbose_name=信封最大重量"`
	MinExpressFee      float64         `xorm:"decimal(18,6) notnull 'min_express_fee'" json:"min_express_fee" smith:"verbose_name=最少快递费"` //澳币
	Ppqzldz            int             `xorm:"int notnull 'ppqzldz'" json:"ppqzldz"`                                                      // 匹配前重量打折  是否匹配物流前对重量打折  1 是 0 否
	Ppqzlzl            float64         `xorm:"decimal(18,6) notnull 'ppqzlzl'" json:"ppqzlzl"`                                            // 匹配前重量折率 匹配物流前对重量打多少折
	Ppqzbdz            int             `xorm:"int notnull 'ppqzbdz'" json:"ppqzbdz"`                                                      // 匹配前中边打折  是否匹配物流前对中边打折   1 是 0 否
	Ppqzbzl            float64         `xorm:"decimal(18,6) notnull 'ppqzbzl'" json:"ppqzbzl"`                                            // 匹配前中边折率 匹配物流前对中边打多少折
	Pphzldz            int             `xorm:"int notnull 'pphzldz'" json:"pphzldz"`                                                      // 匹配后重量打折  是否匹配物流后对重量打折    1 是 0 否
	Pphzlzl            float64         `xorm:"decimal(18,6) notnull 'pphzlzl'" json:"pphzlzl"`                                            // 匹配后重量折率 匹配物流后对重量打多少折
	Pphzbdz            int             `xorm:"int notnull 'pphzbdz'" json:"pphzbdz"`                                                      // 匹配后中边打折  是否匹配物流后对中边打折    1 是 0 否
	Pphzbzl            float64         `xorm:"decimal(18,6) notnull 'pphzbzl'" json:"pphzbzl"`                                            // 匹配后中边折率 匹配物流后对中边打多少折
	Ptfl               float64         `xorm:"decimal(18,6) notnull 'ptfl'" json:"ptfl"`                                                  // 平台费率
	Jyfl               float64         `xorm:"decimal(18,6) notnull 'jyfl'" json:"jyfl"`                                                  // 交易费率
	FeeMaxDiff         float64         `xorm:"decimal(18,6) notnull 'fee_max_diff'" json:"fee_max_diff"`                                  // 运费允许偏差
	Qyppts             int64           `xorm:"bigint notnull 'qyppts'" json:"qyppts"`                                                     // 启用匹配调试， >0 是， 其他 否
	CreateTime         core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime         core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e LogisticsSettingInfo) TableName() string {
	return "erp_logistics_setting"
}

func (e LogisticsSettingInfo) ModelDesc() string {
	return "verbose_name=物流配置,entity_json=logistics_setting,route=logistics/setting"
}
