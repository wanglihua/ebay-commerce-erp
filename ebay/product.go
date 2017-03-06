package ebay

import "encoding/xml"

//GetSellerList

type GetSellerListResponse struct {
	Xmlns                   string                 `xml:"xmlns,attr"`
	Timestamp               string                 `xml:"Timestamp"`
	Ack                     string                 `xml:"Ack"`
	Build                   string                 `xml:"Build"`
	Errors                  Errors                 `xml:"Errors"`
	HasMoreItems            bool                   `xml:"HasMoreItems"`
	Pagination              PaginationResult       `xml:"PaginationResult"`
	ItemsPerPage            int
	PageNumber              int
	ReturnedItemCountActual int
	ItemArray               SellerListItemArray    `xml:"ItemArray"`
}

type SellerListItemArray struct {
	XMLName xml.Name             `xml:"ItemArray"`
	Items   []SellerListItem     `xml:"Item"`
}

type SellerListItem struct {
	ItemId string `xml:"ItemID"`
}
