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

// GetItem

type GetItemResponse struct {
	Xmlns     string                 `xml:"xmlns,attr"`
	Timestamp string                 `xml:"Timestamp"`
	Ack       string                 `xml:"Ack"`
	Build     string                 `xml:"Build"`
	Errors    Errors                 `xml:"Errors"`
	Item      Item                   `xml:"Item"`
}

type Item struct {
	XMLName        xml.Name           `xml:"Item"`
	ItemId         string             `xml:"ItemID"`
	Title          string             `xml:"Title"`
	Sku            string             `xml:"SKU"`
	ListingDetails ItemListingDetails `xml:"ListingDetails"`
	PictureDetails ItemPictureDetails `xml:"PictureDetails"`
}

type ItemListingDetails struct {
	XMLName     xml.Name       `xml:"ListingDetails"`
	ViewItemUrl string         `xml:"ViewItemURL"`
}

type ItemPictureDetails struct {
	XMLName      xml.Name     `xml:"PictureDetails"`
	GalleryType  string       `xml:"GalleryType"`
	GalleryUrl   string       `xml:"GalleryURL"`
	PhotoDisplay string       `xml:"PhotoDisplay"`
	PictureUrls  []string     `xml:"PictureURL"`
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<GetItemResponse xmlns="urn:ebay:apis:eBLBaseComponents">
    <Timestamp>2017-03-06T22:44:37.850Z</Timestamp>
    <Ack>Success</Ack>
    <Version>983</Version>
    <Build>E983_INTL_API_18074235_R1</Build>
    <Item>
        <ItemID>282307939082</ItemID>
        <Title>2 Person Beach Instant Pop Up Tent Sun Shade Shelter UV Protection Protable Camp</Title>
        <SKU>EAZ0161LB8AU</SKU>
	    <ListingDetails>
	      <ViewItemURL> anyURI </ViewItemURL>
	    </ListingDetails>
        <PictureDetails>
            <GalleryType>Plus</GalleryType>
            <GalleryURL>http://i.ebayimg.com/00/s/MTAwMFgxMDAw/z/i7kAAOSwtfhYonls/$_12.JPG?set_id=880000500F</GalleryURL>
            <PhotoDisplay>PicturePack</PhotoDisplay>
            <PictureURL>http://i.ebayimg.com/00/s/MTAwMFgxMDAw/z/i7kAAOSwtfhYonls/$_12.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODEwWDgxMA==/z/QR4AAOSwjDZYZfA8/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/zdkAAOSw9GhYZfA-/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/RdgAAOSwUKxYZfA-/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/17kAAOSwt5hYZfA-/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/188AAOSwt5hYZfA~/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/yj0AAOSwImRYZfA-/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/N7IAAOSw2xRYZfA9/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/YOUAAOSw5cNYZfA~/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/itkAAOSwnHZYZfA9/$_1.JPG?set_id=880000500F</PictureURL>
            <PictureURL>http://i.ebayimg.com/00/s/ODAwWDgwMA==/z/RecAAOSwUKxYZfA~/$_1.JPG?set_id=880000500F</PictureURL>
        </PictureDetails>
    </Item>
</GetItemResponse>
 */
