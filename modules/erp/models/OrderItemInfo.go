package models

import (
	"matrix/core"
)

type OrderItemInfo struct {
	Id                       int64           `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId                  int64           `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	OrderLineItemId          string          `xorm:"nvarchar(200) notnull 'order_line_item_id'" json:"order_line_item_id"`
	Quantity                 int             `xorm:"int notnull 'quantity'" json:"quantity"`
	TransactionId            string          `xorm:"nvarchar(200) notnull 'transaction_id'" json:"transaction_id"`
	TransactionPriceCurrency string          `xorm:"nvarchar(50) notnull 'transaction_price_currency'" json:"transaction_price_currency"`
	TransactionPrice         float64         `xorm:"decimal(18,6) notnull 'transaction_price'" json:"transaction_price"`
	ItemId                   string          `xorm:"nvarchar(50) notnull 'item_id'" json:"item_id"`
	ItemProductId            int64           `xorm:"bigint notnull index 'item_product_id'" json:"item_product_id"`
	ItemSku                  string          `xorm:"nvarchar(100) notnull index 'item_sku'" json:"item_sku"`
	ItemTitle                string          `xorm:"nvarchar(200) notnull 'item_title'" json:"item_title"`
	ItemConditionId          string          `xorm:"nvarchar(50) notnull 'item_condition_id'" json:"item_condition_id"`
	ItemConditionDisplayName string          `xorm:"nvarchar(100) notnull 'item_condition_display_name'" json:"item_condition_display_name"`
	CreatedDate              core.LongTime   `xorm:"datetime notnull 'created_date'" json:"created_date"`
	Desc                     core.NullString `xorm:"nvarchar(500) 'desc'" json:"desc" smith:"verbose_name=备注"`
	CreateTime               core.Time       `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime               core.Time       `xorm:"updated notnull 'update_time'" json:"update_time"`
	Version                  int             `xorm:"version notnull 'version'" json:"version"`
}

func (e OrderItemInfo) TableName() string {
	return TablePrefix + "order_item"
}

func (e OrderItemInfo) ModelDesc() string {
	return "verbose_name=订单明细,entity_json=order_item,route=order_item"
}

/*
	Buyer             TransactionBuyer `xml:"Buyer"`
	ShippingDetails   ShippingDetails
	CreatedDate       string
	Item              TransactionItem `xml:"Item"`
	QuantityPurchased int
	Status            TransactionStatus `xml:"Status"`
	TransactionID     string
	TransactionPrice  Amount
	//ShippingServiceSelected ShippingServiceSelected
	FinalValueFee       Amount
	TransactionSiteID   string
	Platform            string
	Taxes               TransactionTaxes `xml:"Taxes"`
	ActualShippingCost  Amount
	ActualHandlingCost  Amount
	OrderLineItemID     string
	ExtendedOrderID     string
	EBayPlusTransaction bool `xml:"eBayPlusTransaction"`

<Item>
	<ApplicationData>ic6d9i846065473221313283232b</ApplicationData>
	<ItemID>282360182333</ItemID>
	<Site>Australia</Site>
	<Title>3/4/5/6/7/Tier Clear Acrylic Round Cupcake Cake Stand Birthday Wedding Party AU[4 Tier]</Title>
	<SKU>FCA042</SKU>
	<ConditionID>1000</ConditionID>
	<ConditionDisplayName>Brand New</ConditionDisplayName>
</Item>

 */
