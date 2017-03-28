package models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

type LogisticsCompanyInfo struct {
	Id         int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name       string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	Desc       core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	CreateTime core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version    int             `xorm:"version notnull 'version'" json:"version"`
}

func (e LogisticsCompanyInfo) TableName() string {
	return TablePrefix + "logistics_company"
}

func (e LogisticsCompanyInfo) ModelDesc() string {
	return "verbose_name=物流公司,entity_json=logistics_company,route=logistics/company"
}

func (e LogisticsCompanyInfo) InitData(db_session *xorm.Session) {

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "EP"
		logistics_company.Desc = core.MakeNullString("EP", true)
		db_session.Insert(&logistics_company)
	}

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "HUNTER"
		logistics_company.Desc = core.MakeNullString("HUNTER", true)
		db_session.Insert(&logistics_company)
	}

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "TOLL"
		logistics_company.Desc = core.MakeNullString("TOLL", true)
		db_session.Insert(&logistics_company)
	}

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "FASTWAY"
		logistics_company.Desc = core.MakeNullString("FASTWAY", true)
		db_session.Insert(&logistics_company)
	}
}
