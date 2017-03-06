package ebay

import (
	"encoding/xml"
)

type GetOrdersResponse struct {
	Xmlns                    string                 `xml:"xmlns,attr"`
	Timestamp                string                 `xml:"Timestamp"`
	Ack                      string                 `xml:"Ack"`
	Build                    string                 `xml:"Build"`
	Errors                   Errors                 `xml:"Errors"`
	HasMoreOrders            bool                   `xml:"HasMoreOrders"`
	CorrelationID            string                 `xml:"CorrelationID"`
	HardExpirationWarning    string                 `xml:"HardExpirationWarning"`
	Pagination               PaginationResult       `xml:"PaginationResult"`
	OrdersArray              OrderArray             `xml:"OrderArray"`
	OrdersPerPage            int                    `xml:"OrdersPerPage"`
	PageNumber               int                    `xml:"PageNumber"`
	ReturnedOrderCountActual int                    `xml:"ReturnedOrderCountActual"`
}

type Amount struct {
	Currency string   `xml:"currencyID,attr"`
	Value    float64   `xml:",chardata"`
}

type BuyerPackageEnclosures struct {
}

type PaymentInstructionCode struct {
}

type TaxIdentifierAttributeCode struct {
}

type ValueTypeCode struct {
}

type OrderIDArray struct {
	OrderID     string
	BuyerUserID string
}

type OrderArray struct {
	XMLName xml.Name `xml:"OrderArray"`
	Orders  []Order  `xml:"Order"`
}

type Order struct {
	XMLName                             xml.Name `xml:"Order"`
	OrderID                             string
	OrderStatus                         string
	BuyerCheckoutMessage                string
	AdjustmentAmount                    Amount `xml:"AdjustmentAmount"`
	AmountPaid                          Amount `xml:"AmountPaid"`
	AmountSaved                         Amount `xml:"AmountSaved"`
	BuyerPackageEnclosures              BuyerPackageEnclosure `xml:"BuyerPackageEnclosure"`
	BuyerTaxIdentifiers                 BuyerTaxIdentifier `xml:"BuyerTaxIdentifier"`
	BuyerUserID                         string
	CancelDetails                       CancelDetail `xml:"CancelDetail"`
	CancelReason                        string
	CancelReasonDetails                 string
	CancelStatus                        string
	CheckoutStatuss                     CheckoutStatus `xml:"CheckoutStatus"`
	ContainseBayPlusTransaction         bool
	CreatedTime                         string
	ModifiedTime                        string
	CreatingUserRole                    string
	EIASToken                           string
	ExtendedOrderID                     string
	ExternalTransactions                ExternalTransaction `xml:"ExternalTransaction"`
	IntegratedMerchantCreditCardEnabled string
	IsMultiLegShipping                  string
	LogisticsPlanType                   string
	MonetaryDetailss                    MonetaryDetails `xml:"MonetaryDetails"`
	MultiLegShippingDetailss            MultiLegShippingDetails `xml:"MultiLegShippingDetails"`
	PaidTime                            string
	PaymentHoldDetailss                 PaymentHoldDetails `xml:"PaymentHoldDetails"`
	PaymentHoldStatus                   string
	PaymentMethods                      string
	PickupDetailss                      PickupDetails `xml:"PickupDetails"`
	PickupMethodSelected
	RefundArrays                        RefundArray `xml:"RefundArray"`
	SellerEIASToken                     string
	SellerEmail                         string
	SellerUserID                        string
	ShippedTime                         string
	ShippingAddress                     OrderShippingAddress `xml:"ShippingAddress"`
	ShippingConvenienceCharge           string
	ShippingDetailss                    ShippingDetails `xml:"ShippingDetails"`
	ShippingServiceSelected             OrderShippingServiceSelected
	Subtotal                            Amount
	Total                               Amount
	TransactionArray                    TransactionArray `xml:"TransactionArray"`
}

/*
<ShippingAddress>
	<Name>Som Chanthasenh</Name>
	<Street1>23 Graziers Way</Street1>
	<Street2/>
	<CityName>Carnes Hill</CityName>
	<StateOrProvince>New South Wales</StateOrProvince>
	<Country>AU</Country>
	<CountryName>Australia</CountryName>
	<Phone>04 19292611</Phone>
	<PostalCode>2171</PostalCode>
	<AddressID>1255435139019</AddressID>
	<AddressOwner>eBay</AddressOwner>
	<ExternalAddressID/>
</ShippingAddress>
 */

type OrderShippingAddress struct {
	Name              string
	Street1           string
	Street2           string
	CityName          string
	StateOrProvince   string
	Country           string
	CountryName       string
	Phone             string
	PostalCode        string
	AddressID         string
	AddressOwner      string
	ExternalAddressID string
}

/*
<ShippingServiceSelected>
	<ShippingService>AU_Regular</ShippingService>
	<ShippingServiceCost currencyID="AUD">0.0</ShippingServiceCost>
</ShippingServiceSelected>
 */

type OrderShippingServiceSelected struct {
	ShippingService     string
	ShippingServiceCost Amount
}

type TransactionArray struct {
	XMLName      xml.Name      `xml:"TransactionArray"`
	Transactions []Transaction `xml:"Transaction"`
}

type Transaction struct {
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
}

/*
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
type TransactionItem struct {
	ApplicationData      string
	ItemID               string
	Site                 string
	Title                string
	SKU                  string
	ConditionID          string
	ConditionDisplayName string
}

/*
<Buyer>
	<Email>carolboys@rocketmail.com</Email>
	<UserFirstName>CAROL</UserFirstName>
	<UserLastName>BOYS</UserLastName>
</Buyer>
 */
type TransactionBuyer struct {
	XMLName       xml.Name `xml:"Buyer"`
	Email         string
	UserFirstName string
	UserLastName  string
}

/*
<Status>
	<PaymentHoldStatus>None</PaymentHoldStatus>
	<InquiryStatus>NotApplicable</InquiryStatus>
	<ReturnStatus>NotApplicable</ReturnStatus>
</Status>
 */
type TransactionStatus struct {
	XMLName           xml.Name `xml:"Status"`
	PaymentHoldStatus string
	InquiryStatus     string
	ReturnStatus      string
}

/*
<Taxes>
	<TotalTaxAmount currencyID="AUD">0.0</TotalTaxAmount>
	<TaxDetails>
		<Imposition>SalesTax</Imposition>
		<TaxDescription>SalesTax</TaxDescription>
		<TaxAmount currencyID="AUD">0.0</TaxAmount>
		<TaxOnSubtotalAmount currencyID="AUD">0.0</TaxOnSubtotalAmount>
		<TaxOnShippingAmount currencyID="AUD">0.0</TaxOnShippingAmount>
		<TaxOnHandlingAmount currencyID="AUD">0.0</TaxOnHandlingAmount>
	</TaxDetails>
	<TaxDetails>
		<Imposition>WasteRecyclingFee</Imposition>
		<TaxDescription>ElectronicWasteRecyclingFee</TaxDescription>
		<TaxAmount currencyID="AUD">0.0</TaxAmount>
	</TaxDetails>
</Taxes>
 */
type TransactionTaxes struct {
	XMLName        xml.Name `xml:"Taxes"`
	TotalTaxAmount Amount
	TaxDetails     []TransactionTaxDetails
}

/*
	<TaxDetails>
		<Imposition>SalesTax</Imposition>
		<TaxDescription>SalesTax</TaxDescription>
		<TaxAmount currencyID="AUD">0.0</TaxAmount>
		<TaxOnSubtotalAmount currencyID="AUD">0.0</TaxOnSubtotalAmount>
		<TaxOnShippingAmount currencyID="AUD">0.0</TaxOnShippingAmount>
		<TaxOnHandlingAmount currencyID="AUD">0.0</TaxOnHandlingAmount>
	</TaxDetails>
 */
type TransactionTaxDetails struct {
	XMLName             xml.Name `xml:"TaxDetails"`
	Imposition          string
	TaxDescription      string
	TaxAmount           Amount
	TxOnSubtotalAmount  Amount
	TaxOnShippingAmount Amount
	TaxOnHandlingAmount Amount
}

type PickupMethodSelected struct {
	XMLName               xml.Name `xml:"PickupMethodSelected"`
	MerchantPickupCode    string
	PickupFulfillmentTime string
	PickupLocationUUID    string
	PickupMethod          string
	PickupStatus          string
	PickupStoreID         string
}

type ShippingDetails struct {
	CalculatedShippingRate
	CODCost         string
	InsuranceFee    string
	InsuranceOption string
	InsuranceWanted string

	InternationalShippingServiceOption
	//SellingManagerSalesRecordNumber
	ShipmentTrackingDetails
	ShippingServiceOptions
	TaxTable
}

type TaxJurisdiction struct {
	JurisdictionID        string
	SalesTaxPercent       string
	ShippingIncludedInTax string
}

type TaxTable struct {
	TaxJurisdiction
}

type ShippingPackageInfo struct {
	ActualDeliveryTime       string
	ScheduledDeliveryTimeMax string
	ScheduledDeliveryTimeMin string
	ShippingTrackingEvent    string
	StoreID                  string
}

type ShippingServiceOptions struct {
	ExpeditedService              string
	ImportCharge                  string
	LogisticPlanType              string
	ShippingInsuranceCost         string
	ShippingService               string
	ShippingServiceAdditionalCost string
	ShippingServiceCost           string
	ShippingServicePriority       int
}

type ShipmentTrackingDetails struct {
	ShipmentTrackingNumber string
	ShippingCarrierUsed    string
}

type InternationalShippingServiceOption struct {
	//	ImportCharge  string
	//     ShippingInsuranceCost string
	//     ShippingService string
	//      ShippingServiceAdditionalCost  string
	//     ShippingServiceCost string
	//     ShippingServicePriority int
	//    ShipToLocation string
}

type CalculatedShippingRate struct {
	InternationalPackagingHandlingCosts string
	OriginatingPostalCode               string
	PackageDepth                        string
	PackageLength                       string
	PackageWidth                        string
	PackagingHandlingCosts              string
	ShippingIrregular                   string
	ShippingPackage                     string
	WeightMajor                         string
	WeightMinor                         string
}

type PickupOptions struct {
	PickupMethod   string
	PickupPriority string
}

type RefundArray struct{}

type PickupDetails struct {
	PickupOptions
}

type RequiredSellerActionArray struct {
	RequiredSellerAction string
}

type PaymentHoldDetails struct {
	ExpectedReleaseDate        string
	NumOfReqSellerActions      string
	PaymentHoldReason          string
	RequiredSellerActionArrays RequiredSellerActionArray
}

type ShipToAddress struct {
	XMLName           xml.Name `xml:"ShipToAddress"`
	AddressAttribute  string
	AddressID         string
	AddressOwner      string
	CityName          string
	Country           string
	CountryName       string
	ExternalAddressID string
	Name              string
	Phone             string
	PostalCode        string
	//ReferenceID       string
	StateOrProvince string
	Street1         string
	Street2         string
}

type SellerShipmentToLogisticsProvider struct {
	ShippingServiceDetails
	ShipToAddress
	ShippingTimeMax int
	ShippingTimeMin int
}

type ShippingServiceDetails struct {
	ShippingService   string
	TotalShippingCost string
}

type MultiLegShippingDetails struct {
	SellerShipmentToLogisticsProvider
}

type MonetaryDetails struct {
	Payments
	Refunds
}

type Payments struct {
	Payment []Payment
}

type Refunds struct {
	Refund
}

type Refund struct {
	FeeOrCreditAmount string
	//ReferenceID       string
	RefundAmount string
	RefundStatus string
	RefundTime   string
	RefundTo     string
	RefundType   string

	//RefundAmount string
	RefundFromSeller string
	//	RefundID string
	//	RefundStatus string
	//	RefundTime string
	TotalRefundToBuyer string
}

type SalesTax struct {
	SalesTaxAmount        string
	SalesTaxPercent       string
	SalesTaxState         string
	ShippingIncludedInTax string
}

type Payment struct {
	FeeOrCreditAmount  string
	Payee              string
	Payer              string
	PaymentAmount      string
	PaymentReferenceID string
	PaymentStatus      string
	PaymentTime        string
	//ReferenceID        string
}

type ExternalTransaction struct {
	ExternalTransactionID     string
	ExternalTransactionStatus string
	ExternalTransactionTime   string
	FeeOrCreditAmount         string
	PaymentOrRefundAmount     string
}

type BuyerPackageEnclosure struct {
	BuyerPackageEnclosureType string
}

type BuyerTaxIdentifier struct {
	Attribute string
	ID        string
	Type      string
}

type CancelDetail struct {
	CancelCompleteDate  string
	CancelIntiationDate string
	CancelIntiator      string
	CancelReason        string
	CancelReasonDetails string
}

type CheckoutStatus struct {
	eBayPaymentStatus                   string
	IntegratedMerchantCreditCardEnabled string
	LastModifiedTime                    string
	PaymentInstrument                   string
	PaymentMethod                       string
	Status                              string
}

func (o *OrderIDArray) GetOrderId(orderid string) {
	o.OrderID = orderid

}
