package models

import (
	"matrix/core"
)

type EBayAPIConfigInfo struct {
	Id         int64       `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ShopId     int64       `xorm:"int64 notnull 'shop_id'" json:"shop_id" smith:"verbose_name=店铺"`
	AppId      string      `xorm:"nvarchar(100)  notnull 'app_id'" json:"app_id" smith:"verbose_name=AppID"`
	CertId     string      `xorm:"nvarchar(100)  notnull 'cert_id'" json:"cert_id" smith:"verbose_name=CertID"`
	DevId      string      `xorm:"nvarchar(100)  notnull 'dev_id'" json:"dev_id" smith:"verbose_name=DevID"`
	Token      string      `xorm:"nvarchar(2000)  notnull 'token'" json:"token" smith:"verbose_name=Token"`
	CreateTime core.Time   `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time   `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e EBayAPIConfigInfo) TableName() string {
	return "erp_ebay_api_config"
}

func (e EBayAPIConfigInfo) ModelDesc() string {
	return "verbose_name=ebay API 配置,entity_json=ebay_api_config,route=ebay/api/config"
}

func (e EBayAPIConfigInfo) GetEndPoint() string {
	return "https://api.ebay.com/ws/api.dll"
	//return "https://api.sandbox.ebay.com/ws/api.dll"
}

func (e EBayAPIConfigInfo) GetSiteId() string {
	return "3" //澳洲的 Site Id
}

func (e EBayAPIConfigInfo) GetApiVersion() string {
	//return "967" // API 版本
	return "995" // API 版本
}

/*

	AppID:    "***REMOVED***",
	CertID:   "***REMOVED***",
	DevID:    "***REMOVED***",
	RuName:   "",
	Token:    "***REMOVED***",
	EndPoint: "https://api.ebay.com/ws/api.dll",
	SiteId:   "3",   //澳洲的 Site Id
	Version:  "967", // API 版本

 */
