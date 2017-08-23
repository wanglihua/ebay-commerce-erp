package models

import (
	"matrix/core"
)

type ProductFeedbackCateInfo struct {
	Id          int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name        string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Description string          `xorm:"nvarchar(800) 'description'" json:"description" smith:"verbose_name=描述,min=0,max=800"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductFeedbackCateInfo) TableName() string {
	return "erp_product_feedback_cate"
}

func (e ProductFeedbackCateInfo) ModelDesc() string {
	return "verbose_name=商品质量反馈类别,entity_json=product_feedback_cate,route=product/feedback/cate"
}
