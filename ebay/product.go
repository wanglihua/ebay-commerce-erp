package ebay

//GetSellerList
/*
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
 */

type GetSellerListResponse struct {
	Xmlns        string                 `xml:"xmlns,attr"`
	Timestamp    string                 `xml:"Timestamp"`
	Ack          string                 `xml:"Ack"`
	Build        string                 `xml:"Build"`
	Errors       Errors                 `xml:"Errors"`
	HasMoreItems bool                   `xml:"HasMoreItems"`
	Pagination   PaginationResult       `xml:"PaginationResult"`
}

//<HasMoreItems>false</HasMoreItems>
/*
<ItemsPerPage>200</ItemsPerPage>
<PageNumber>1</PageNumber>
<ReturnedItemCountActual>181</ReturnedItemCountActual>
 */
