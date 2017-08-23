package ddlc_models

import (
	"matrix/core"
)

type OrderShippingAddressInfo struct {
	Id                int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId           int64           `xorm:"bigint notnull index 'order_id'" json:"order_id"` // 订单
	Name              string          `xorm:"nvarchar(200) notnull 'name'" json:"name"`
	Street1           string          `xorm:"nvarchar(200) notnull 'street1'" json:"street1"`
	Street2           string          `xorm:"nvarchar(200) notnull 'street2'" json:"street2"`
	CityName          string          `xorm:"nvarchar(200) notnull 'city_name'" json:"city_name"`
	Province          string          `xorm:"nvarchar(200) notnull 'province'" json:"province"`
	Country           string          `xorm:"nvarchar(200) notnull 'country'" json:"country"`
	CountryName       string          `xorm:"nvarchar(200) notnull 'country_name'" json:"country_name"`
	Phone             string          `xorm:"nvarchar(200) notnull 'phone'" json:"phone"`
	PostalCode        string          `xorm:"nvarchar(200) notnull 'postal_code'" json:"postal_code"`
	AddressId         string          `xorm:"nvarchar(200) notnull 'address_id'" json:"address_id"`
	AddressOwner      string          `xorm:"nvarchar(200) notnull 'address_owner'" json:"address_owner"`
	ExternalAddressId string          `xorm:"nvarchar(200) notnull 'external_address_id'" json:"external_address_id"`
	CreateTime        core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime        core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderShippingAddressInfo) TableName() string {
	return "erp_order_shipping_address"
}

func (e OrderShippingAddressInfo) ModelDesc() string {
	return "verbose_name=订单送货地址,entity_json=shipping_address,route=order/shipping/address"
}
