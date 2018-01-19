package xtpz_models

import (
	"matrix/core"
)

type EBayAPIConfigInfo struct {
	Id         int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ShopId     int64     `xorm:"int64 notnull 'shop_id'" json:"shop_id" smith:"verbose_name=店铺"`
	AppId      string    `xorm:"nvarchar(100)  notnull 'app_id'" json:"app_id" smith:"verbose_name=AppID"`
	CertId     string    `xorm:"nvarchar(100)  notnull 'cert_id'" json:"cert_id" smith:"verbose_name=CertID"`
	DevId      string    `xorm:"nvarchar(100)  notnull 'dev_id'" json:"dev_id" smith:"verbose_name=DevID"`
	Token      string    `xorm:"nvarchar(2000)  notnull 'token'" json:"token" smith:"verbose_name=Token"`
	CreateTime core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
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
	return "15" //澳洲的 Site Id
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

/*

eBay Site ID to Global ID Mapping
Each eBay site maps to a unique eBay global ID. You must specify a global ID in the X-EBAY-SOA-GLOBAL-ID HTTP header for each API call. The global ID you use must correspond to one of the eBay site IDs in the following table:

eBay Site ID	Site Name	Global ID
0	eBay United States	EBAY-US
2	eBay Canada (English)	EBAY-ENCA
3	eBay UK	EBAY-GB
15	eBay Australia	EBAY-AU
16	eBay Austria	EBAY-AT
23	eBay Belgium (French)	EBAY-FRBE
71	eBay France	EBAY-FR
77	eBay Germany	EBAY-DE
100	eBay Motors	EBAY-MOTOR
101	eBay Italy	EBAY-IT
123	eBay Belgium (Dutch)	EBAY-NLBE
146	eBay Netherlands	EBAY-NL
186	eBay Spain	EBAY-ES
193	eBay Switzerland	EBAY-CH
201	eBay Hong Kong	EBAY-HK
203	eBay India	EBAY-IN
205	eBay Ireland	EBAY-IE
207	eBay Malaysia	EBAY-MY
210	eBay Canada (French)	EBAY-FRCA
211	eBay Philippines	EBAY-PH
212	eBay Poland	EBAY-PL
216	eBay Singapore	EBAY-SG

 */
