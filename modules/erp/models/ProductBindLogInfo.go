package models

import (
	"matrix/core"
)

type ProductBindLogInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ProductId   int64     `xorm:"bigint notnull index 'product_id'" json:"product_id" smith:"verbose_name=商品Id"`
	ProductCode string    `xorm:"nvarchar(100) notnull index 'product_code'" json:"product_code" smith:"verbose_name=商品编号"`
	UserId      int64     `xorm:"bigint notnull index 'user_id'" json:"user_id" smith:"verbose_name=用户Id"`
	LoginName   string    `xorm:"nvarchar(100) notnull index 'login_name'" json:"login_name" smith:"verbose_name=用户工号"`
	Operation   string    `xorm:"nvarchar(800) notnull 'operation'" json:"operation" smith:"verbose_name=操作"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
}

func (e ProductBindLogInfo) TableName() string {
	return TablePrefix + "product_bind_log"
}

func (e ProductBindLogInfo) ModelDesc() string {
	return "verbose_name=绑定商品修改日志,entity_json=product_bind_log,route=product/bind/log"
}
