package models

import (
	"matrix/core"
)

type HsCodeCnReportElemInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	HsCode     string    `xorm:"nvarchar(100) notnull index" json:"hs_code" smith:"verbose_name=海关编码,min=1,max=10"`
	CnName     string    `xorm:"nvarchar(200) notnull 'cn_name'" json:"cn_name" smith:"verbose_name=中文品名"`
	EnName     string    `xorm:"nvarchar(200) notnull 'en_name'" json:"en_name" smith:"verbose_name=英文品名"`
	Elem01     string    `xorm:"nvarchar(500) 'elem01'" json:"elem01" smith:"verbose_name=要素1"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e HsCodeCnReportElemInfo) TableName() string {
	return TablePrefix + "hs_code_cn"
}

func (e HsCodeCnReportElemInfo) ModelDesc() string {
	return "verbose_name=中国海关编码,entity_json=hs_code_cn,route=hs/code/cn"
}
