package kcgl_models

import (
	"matrix/core"
)

type ProductBindInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code       string    `xorm:"nvarchar(100) notnull index 'code'" json:"code" smith:"verbose_name=商品编号,min=2,max=100"`
	Name       string    `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=200"`
	BindList   string    `xorm:"nvarchar(500) notnull 'bind_list'" json:"bind_list" smith:"verbose_name=捆绑商品列表,min=1,max=500"`
	Desc       string    `xorm:"nvarchar(500) notnull index 'desc'" json:"desc" smith:"verbose_name=描述,min=2,max=200"`
	StatusId   int64     `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e ProductBindInfo) TableName() string {
	return "erp_product_bind"
}

func (e ProductBindInfo) ModelDesc() string {
	return "verbose_name=绑定商品,entity_json=product_bind,route=product/bind"
}
