package sppz_models

import (
	"matrix/core"
	"matrix/modules/erp/models"
)

type BrandInfo struct {
	Id          int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name        string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Description string          `xorm:"nvarchar(800) 'description'" json:"description" smith:"verbose_name=描述,min=0,max=800"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e BrandInfo) TableName() string {
	return models.TablePrefix + "brand"
}

func (e BrandInfo) ModelDesc() string {
	return "verbose_name=品牌,entity_json=brand,route=brand"
}
