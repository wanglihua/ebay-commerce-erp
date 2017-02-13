package models

import (
	"matrix/core"
)

type ProductCateInfo struct {
	Id          int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name        string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	FullName    string          `xorm:"nvarchar(400) notnull index 'full_name'" json:"full_name" smith:"verbose_name=全称,min=2,max=100"`
	EnglishName string          `xorm:"nvarchar(200) notnull index 'english_name'" json:"english_name" smith:"verbose_name=英文名,min=2,max=100"`
	Desc        core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	ParentId    core.NullInt    `xorm:"bigint null index 'parent_id'" json:"parent_id" smith:"verbose_name=父目录"`
	StatusId    int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	CreateTime  core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int             `xorm:"version notnull 'version'" json:"version"`
}

func (e ProductCateInfo) TableName() string {
	return TablePrefix + "product_cate"
}

func (e ProductCateInfo) ModelDesc() string {
	return "verbose_name=商品目录,entity_json=product_cate,route=product/cate"
}
