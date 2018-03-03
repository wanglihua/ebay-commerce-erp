package ddpz_models

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

type LogisticsCompanyInfo struct {
	Id                 int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	Name               string          `xorm:"nvarchar(200) notnull index unique 'name'" json:"name" smith:"verbose_name=名称,min=2,max=100"`
	EbayName           string          `xorm:"nvarchar(200) null 'ebay_name'" json:"ebay_name"`
	AddressSpecialWord int             `xorm:"int notnull 'address_special_word'" json:"address_special_word" smith:"verbose_name=地址特殊字段"` // 地址特殊字段 0 否， 1 是
	Desc               core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=描述,max=500"`
	StatusId           int64           `xorm:"bigint notnull index 'status_id'" json:"status_id" smith:"verbose_name=状态"`
	IconFile           string          `xorm:"nvarchar(200) null 'icon_file'" json:"icon_file" smith:"verbose_name=图标"`
	Ord                int64           `xorm:"bigint null 'ord'" json:"ord" smith:"verbose_name=排序"`
	CreateTime         core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime         core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e LogisticsCompanyInfo) TableName() string {
	return "erp_logistics_company"
}

func (e LogisticsCompanyInfo) ModelDesc() string {
	return "verbose_name=物流公司,entity_json=logistics_company,route=logistics/company"
}

func (e LogisticsCompanyInfo) InitData(db_session *xorm.Session) {
	/*
	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "AU POST"
		logistics_company.AddressSpecialWord = 1
		logistics_company.Desc = core.MakeNullString("Australia Post", true)
		logistics_company.StatusId = LOGISTICS_COMPANY_ENABLED
		db_session.Insert(&logistics_company)
	}

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "EP"
		logistics_company.AddressSpecialWord = 0
		logistics_company.Desc = core.MakeNullString("EP", true)
		db_session.Insert(&logistics_company)
	}

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "HUNTER"
		logistics_company.AddressSpecialWord = 0
		logistics_company.Desc = core.MakeNullString("HUNTER", true)
		db_session.Insert(&logistics_company)
	}

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "TOLL"
		logistics_company.AddressSpecialWord = 0
		logistics_company.Desc = core.MakeNullString("TOLL", true)
		db_session.Insert(&logistics_company)
	}

	{
		var logistics_company LogisticsCompanyInfo
		logistics_company.Name = "FASTWAY"
		logistics_company.AddressSpecialWord = 0
		logistics_company.Desc = core.MakeNullString("FASTWAY", true)
		db_session.Insert(&logistics_company)
	}
	*/
}
