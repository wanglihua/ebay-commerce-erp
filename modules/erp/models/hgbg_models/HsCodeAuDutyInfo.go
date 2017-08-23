package hgbg_models

import (
	"matrix/core"
)

type HsCodeAuDutyInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	HsCode     string    `xorm:"nvarchar(100) notnull index 'hs_code'" json:"hs_code" smith:"verbose_name=海关编码,min=1,max=10"`
	CnName     string    `xorm:"nvarchar(200) notnull 'cn_name'" json:"cn_name" smith:"verbose_name=中文品名"`
	EnName     string    `xorm:"nvarchar(200) notnull index 'en_name'" json:"en_name" smith:"verbose_name=英文品名"`
	DutyRate   float64   `xorm:"decimal(18,6) notnull 'duty_rate'" json:"duty_rate" smith:"verbose_name=关税税率"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e HsCodeAuDutyInfo) TableName() string {
	return "erp_hs_code_au"
}

func (e HsCodeAuDutyInfo) ModelDesc() string {
	return "verbose_name=澳洲海关编码,entity_json=hs_code_au,route=hs/code/au"
}
