package models

import (
	"matrix/core"
)

type ProductFeedbackInfo struct {
	Id            int64          `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CateId        int64          `xorm:"bigint notnull index 'cate_id'" json:"cate_id" smith:"verbose_name=反馈类别"`
	ProductId     int64          `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品"`
	SupplierId    int64          `xorm:"bigint notnull index 'supplier_id'" json:"supplier_id" smith:"verbose_name=供应商"`
	Code          string         `xorm:"nvarchar(200) notnull index 'code'" json:"code" smith:"verbose_name=反馈编码"` //数据插入后， 根据生成的ID更新Code
	Description   string         `xorm:"nvarchar(800) 'description'" json:"description" smith:"verbose_name=描述,min=0,max=800"`
	UserId        int64          `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=反馈人id"`
	UserLoginName string         `xorm:"nvarchar(200) notnull 'user_login_name'" json:"user_login_name" smith:"verbose_name=反馈人工号,min=2,max=100"`
	CreateTime    core.Time      `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime    core.Time      `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductFeedbackInfo) TableName() string {
	return TablePrefix + "product_feedback"
}

func (e ProductFeedbackInfo) ModelDesc() string {
	return "verbose_name=商品质量反馈,entity_json=product_feedback,route=product/feedback"
}
