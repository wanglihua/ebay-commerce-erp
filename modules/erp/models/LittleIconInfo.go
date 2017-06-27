package models

import (
	"matrix/core"
)

type LittleIconInfo struct {
	Id         int64      `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	IconFile   string     `xorm:"nvarchar(300) notnull unique 'icon_file'" json:"icon_file"` //图标文件
	CreateTime core.Time  `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time  `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e LittleIconInfo) TableName() string {
	return TablePrefix + "little_icon"
}

func (e LittleIconInfo) ModelDesc() string {
	return "verbose_name=小图标,entity_json=icon,route=little/icon"
}
