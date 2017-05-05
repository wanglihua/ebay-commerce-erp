package models

type RegionInfo struct {
	Id       int64       `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Code     string      `xorm:"nvarchar(40) notnull 'code'" json:"code"`
	Province string      `xorm:"nvarchar(40) notnull 'province'" json:"province"`
	City     string      `xorm:"nvarchar(40) notnull 'city'" json:"city"`
	District string      `xorm:"nvarchar(40) notnull 'district'" json:"district"`
	Parent   string      `xorm:"nvarchar(40) notnull 'parent'" json:"parent"`
}

func (e RegionInfo) TableName() string {
	return TablePrefix + "region"
}
