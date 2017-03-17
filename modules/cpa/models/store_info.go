package models

import (
	"matrix/core"
)

type Store_info struct {
	Id          int64              `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name        string             `xorm:"nvarchar(255) notnull unique index 'name'" json:"name"`
	SellerId    string             `xorm:"nvarchar(200) notnull 'seller_id'" json:"seller_id"`
	Desc        core.NullString    `xorm:"nvarchar(500) null 'desc'" json:"desc"`
	Create_time core.Time          `xorm:"created notnull 'create_time'" json:"create_time"`
	Update_time core.Time          `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version     int                `xorm:"version notnull 'version'" json:"version"`
}

func (e Store_info) TableName() string {
	return TablePrefix + "store"
}
