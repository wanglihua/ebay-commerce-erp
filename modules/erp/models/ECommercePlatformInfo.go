package models

//ECommercePlatformInfo

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

const (
	PLATFORM_ID_EBAY = 100
)

type ECommercePlatformInfo struct {
	Id         int64           `xorm:"bigint notnull pk 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e ECommercePlatformInfo) TableName() string {
	return TablePrefix + "ecommerce_platform"
}

func (e ECommercePlatformInfo) ModelDesc() string {
	return "verbose_name=电商平台,entity_json=platform,route=ecommerce/platform"
}

func (e ECommercePlatformInfo) InitData(db_session *xorm.Session) {
	platform := new(ECommercePlatformInfo)
	platform.Id = PLATFORM_ID_EBAY
	platform.Name = "ebay"
	platform.Desc = core.NewNullString("ebay", true)
	db_session.Insert(platform)
}
