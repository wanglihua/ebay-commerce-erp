package ddpz_models

import (
	"matrix/core"
)

const CITY_ID_MEL int64 = 1
const CITY_ID_BNE int64 = 2
const CITY_ID_SYD int64 = 3

type CityPostcodeInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	CityId     int64     `xorm:"bigint notnull index 'city_id'" json:"city_id"`
	CityName   string    `xorm:"nvarchar(30) notnull 'city_name'" json:"city_name"`
	Postcode   string    `xorm:"nvarchar(10) notnull 'postcode'" json:"postcode"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e CityPostcodeInfo) TableName() string {
	return "erp_city_postcode"
}

func (e CityPostcodeInfo) ModelDesc() string {
	return "verbose_name=城市邮编,entity_json=city_postcode,route=city/postcode"
}
