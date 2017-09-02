package xsgl_models

import (
	"matrix/core"
)

type ProductSaleSettingInfo struct {
	Id               int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	PostageCalWeight float64         `xorm:"decimal(18,6) notnull 'postage_cal_weight'" json:"postage_cal_weight" smith:"verbose_name=客户端运费权重"` // 客户端运费计算权重
	CreateTime       core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime       core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductSaleSettingInfo) TableName() string {
	return "erp_product_sale_setting"
}

func (e ProductSaleSettingInfo) ModelDesc() string {
	return "verbose_name=商品销售设置,entity_json=product_sale_setting,route=product/sale/setting"
}

/*

改成按IP计算

 */
