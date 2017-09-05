package spbb_models

import (
	"matrix/core"
)

type ProductSaleAnalysisSettingInfo struct {
	Id              int64      `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name            string     `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Level1Days      int64      `xorm:"bigint notnull 'level1_days'" json:"level1_days" smith:"verbose_name=1档统计天数,min=2,max=300"`
	Level2Days      int64      `xorm:"bigint notnull 'level2_days'" json:"level2_days" smith:"verbose_name=2档统计天数,min=2,max=300"`
	Level3Days      int64      `xorm:"bigint notnull 'level3_days'" json:"level3_days" smith:"verbose_name=3档统计天数,min=2,max=300"`
	Level4Days      int64      `xorm:"bigint notnull 'level4_days'" json:"level4_days" smith:"verbose_name=4档统计天数,min=2,max=300"`
	Level5Days      int64      `xorm:"bigint notnull 'level5_days'" json:"level5_days" smith:"verbose_name=5档统计天数,min=2,max=300"`
	Level1Weight    float64    `xorm:"decimal(18,6) notnull 'level1_weight'" json:"level1_weight" smith:"verbose_name=1档统计权重,min=0,max=1"`
	Level2Weight    float64    `xorm:"decimal(18,6) notnull 'level2_weight'" json:"level2_weight" smith:"verbose_name=2档统计权重,min=0,max=1"`
	Level3Weight    float64    `xorm:"decimal(18,6) notnull 'level3_weight'" json:"level3_weight" smith:"verbose_name=3档统计权重,min=0,max=1"`
	Level4Weight    float64    `xorm:"decimal(18,6) notnull 'level4_weight'" json:"level4_weight" smith:"verbose_name=4档统计权重,min=0,max=1"`
	Level5Weight    float64    `xorm:"decimal(18,6) notnull 'level5_weight'" json:"level5_weight" smith:"verbose_name=5档统计权重,min=0,max=1"`
	DayAnalysisDays int64      `xorm:"bigint notnull 'day_analysis_days'" json:"day_analysis_days" smith:"verbose_name=日统计天数,min=2,max=1000"`
	CreateTime      core.Time  `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime      core.Time  `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductSaleAnalysisSettingInfo) TableName() string {
	return "erp_product_sale_analysis_setting"
}

func (e ProductSaleAnalysisSettingInfo) ModelDesc() string {
	return "verbose_name=商品销售统计设置,entity_json=setting,route=product/sale/analysis/setting"
}
