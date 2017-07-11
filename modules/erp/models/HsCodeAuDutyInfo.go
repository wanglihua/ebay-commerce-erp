package models

import (
	"matrix/core"
)

type HsCodeAuDutyInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e HsCodeAuDutyInfo) TableName() string {
	return TablePrefix + "hs_code_au"
}

func (e HsCodeAuDutyInfo) ModelDesc() string {
	return "verbose_name=澳洲海关编码,entity_json=hs_code_au,route=hs/code/au"
}
