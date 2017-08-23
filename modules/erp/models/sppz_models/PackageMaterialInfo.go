package sppz_models

import (
	"matrix/core"
)

type PackageMaterialInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	Level      int64           `xorm:"bigint notnull 'level'" json:"level" smith:"verbose_name=级别"`
	Length     float64         `xorm:"float notnull 'length'" json:"length" smith:"verbose_name=长(cm),min=0"`
	Width      float64         `xorm:"float notnull 'width'" json:"width" smith:"verbose_name=宽(cm),min=0"`
	Height     float64         `xorm:"float notnull 'height'" json:"height" smith:"verbose_name=高(cm),min=0"`
	Weight     float64         `xorm:"float notnull 'weight'" json:"weight" smith:"verbose_name=重量(克),min=0"`
	Price      float64         `xorm:"float notnull 'price'" json:"price" smith:"verbose_name=价格(元),min=0"`
	StatusId   int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e PackageMaterialInfo) TableName() string {
	return "erp_package_material"
}

func (e PackageMaterialInfo) ModelDesc() string {
	return "verbose_name=包材,entity_json=package_material,route=package/material"
}
