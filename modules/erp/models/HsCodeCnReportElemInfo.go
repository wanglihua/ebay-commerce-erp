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
	Elem02     string    `xorm:"nvarchar(500) 'elem02'" json:"elem02" smith:"verbose_name=要素2"`
	Elem03     string    `xorm:"nvarchar(500) 'elem03'" json:"elem03" smith:"verbose_name=要素3"`
	Elem04     string    `xorm:"nvarchar(500) 'elem04'" json:"elem04" smith:"verbose_name=要素4"`
	Elem05     string    `xorm:"nvarchar(500) 'elem05'" json:"elem05" smith:"verbose_name=要素5"`
	Elem06     string    `xorm:"nvarchar(500) 'elem06'" json:"elem06" smith:"verbose_name=要素6"`
	Elem07     string    `xorm:"nvarchar(500) 'elem07'" json:"elem07" smith:"verbose_name=要素7"`
	Elem08     string    `xorm:"nvarchar(500) 'elem08'" json:"elem08" smith:"verbose_name=要素8"`
	Elem09     string    `xorm:"nvarchar(500) 'elem09'" json:"elem09" smith:"verbose_name=要素9"`
	Elem10     string    `xorm:"nvarchar(500) 'elem10'" json:"elem10" smith:"verbose_name=要素10"`
	Elem11     string    `xorm:"nvarchar(500) 'elem11'" json:"elem11" smith:"verbose_name=要素11"`
	Elem12     string    `xorm:"nvarchar(500) 'elem12'" json:"elem12" smith:"verbose_name=要素12"`
	Elem13     string    `xorm:"nvarchar(500) 'elem13'" json:"elem13" smith:"verbose_name=要素13"`
	Elem14     string    `xorm:"nvarchar(500) 'elem14'" json:"elem14" smith:"verbose_name=要素14"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e HsCodeCnReportElemInfo) TableName() string {
	return TablePrefix + "hs_code_cn"
}

func (e HsCodeCnReportElemInfo) ModelDesc() string {
	return "verbose_name=中国海关编码,entity_json=hs_code_cn,route=hs/code/cn"
}
