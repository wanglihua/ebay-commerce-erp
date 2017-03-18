package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"github.com/franela/goreq"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/ebay"
	"matrix/modules/erp/models"
)

const log_response_xml_to_file = false

func fetch_product_item(config models.EBayAPIConfigInfo, item_id string) ebay.GetItemResponse {

	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetItemRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<ItemID>%s</ItemID>
</GetItemRequest>`

	req_body = fmt.Sprintf(req_body, config.Token, item_id)

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.GetEndPoint(),
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		ShowDebug:   revel.DevMode,
	}

	req.AddHeader("X-EBAY-API-CALL-NAME", "GetItem")
	req.AddHeader("X-EBAY-API-DEV-NAME", config.DevId)
	req.AddHeader("X-EBAY-API-CERT-NAME", config.CertId)
	req.AddHeader("X-EBAY-API-APP-NAME", config.AppId)
	req.AddHeader("X-EBAY-API-VERSION", config.GetApiVersion())
	req.AddHeader("X-EBAY-API-COMPATIBILITY-LEVEL", config.GetApiVersion())
	req.AddHeader("X-EBAY-API-REQUEST-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-RESPONSE-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-SITEID", config.GetSiteId())
	req.AddHeader("Accept", "application/xml,application/xhtml+xml")

	res, err := req.Do()

	var response ebay.GetItemResponse

	if err != nil {
		revel.ERROR.Println("eBay GetItem ERROR - processUrl -> req.Do: " + err.Error())
		return response
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("eBay GetItem ERROR - ioutil.ReadAll : " + err.Error())
		return response
	}

	if revel.DevMode && log_response_xml_to_file {
		core.Force_write_to_file(fmt.Sprintf(`z:\test_product_response\item_%s.xml`, item_id), string(data))
		//core.Force_write_to_file(fmt.Sprintf(`/home/wanglihua/ebay_sync/test_product_response/item_%s.xml`, item_id), string(data))
	}

	err = xml.Unmarshal(data, &response)
	if err != nil {
		revel.ERROR.Println(err)
	}

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetItem ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return response
	}

	return response
}

/*
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
	PictureDetails ItemPictureDetails `xml:"PictureDetails"`
}

type ItemPictureDetails struct {
	XMLName      xml.Name     `xml:"PictureDetails"`
	GalleryType  string       `xml:"GalleryType"`
	GalleryUrl   string       `xml:"GalleryURL"`
	PhotoDisplay string       `xml:"PhotoDisplay"`
	PictureUrls  []string     `xml:"PictureURL"`
}

 */

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

/*
通过API获得 item_id 对应的商品的 sku，查询数据库中SKU对应的商品是不是存在。如果存在就更新，不存在就插入新的数据
通过一个 item_id， 是不是只能获得一个 sku，还是一个 item_id 可以对应多个 sku ?

http://developer.ebay.com/devzone/xml/docs/Reference/eBay/GetItem.html

 */
func fetch_product_item_to_db(config models.EBayAPIConfigInfo, db_session *xorm.Session, item_id string) (product_id int64, product_sku string) {

	item_response := fetch_product_item(config, item_id)

	item_sku := item_response.Item.Sku
	var product models.ProductInfo
	has_product, err := db_session.Where(fmt.Sprintf("sku = '%s'", item_sku)).Get(&product)
	if err != nil {
		revel.ERROR.Println(err)
	}

	product.ItemId = item_response.Item.ItemId
	product.Sku = item_response.Item.Sku
	product.Name = item_response.Item.Title
	product.ImageUrl = core.NewNullString(item_response.Item.PictureDetails.GalleryUrl, true)
	product.ViewItemUrl = core.NewNullString(item_response.Item.ListingDetails.ViewItemUrl, true)

	if product.StatusId == 0 {
		product.StatusId = models.PRODUCT_STATUS_ZC
	}

	if has_product == false {
		// do insert
		_, err = db_session.Insert(&product)
		if err != nil {
			revel.ERROR.Println(err)
		}
	} else {
		// do update
		update_col_list := []string{"item_id", "sku", "name", "image_url"}
		_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
		if err != nil {
			revel.ERROR.Println(err)
		}
	}

	// 同步商品图片 begin

	// 插入商品图片，注意有时是更新图片不是插入图片
	// 如果这个商品的这个图片URL已经存在，就不做插入了

	// 获得这个商品的所有图片
	var product_image_list_in_db = make([]models.ProductImageInfo, 0)
	err = db_session.Where(fmt.Sprintf("product_id = %d", product.Id)).Find(&product_image_list_in_db)
	if err != nil {
		revel.ERROR.Println(err)
	}

	// 插入数据库中没有的图片
	var product_image_url_map_in_db = make(map[string]int)
	for _, product_image_in_db := range product_image_list_in_db {
		product_image_url_map_in_db[product_image_in_db.Url] = 1
	}

	for _, product_image_url_download := range item_response.Item.PictureDetails.PictureUrls {
		if _, found := product_image_url_map_in_db[product_image_url_download]; found == false {

			var product_image models.ProductImageInfo
			product_image.ProductId = product.Id
			product_image.ProductCode = product.Sku
			product_image.TypeId = models.PRODUCT_IMAGE_TYPE_CCTP
			product_image.ShopId = config.ShopId
			product_image.Url = product_image_url_download

			_, err := db_session.Insert(product_image)
			if err != nil {
				revel.ERROR.Println(err)
			}
		}
	}

	// 删除数据库中多余的图片
	var product_image_url_map_download = make(map[string]int)
	for _, product_image_url_download := range item_response.Item.PictureDetails.PictureUrls {
		product_image_url_map_download[product_image_url_download] = 1
	}

	for _, product_image_in_db := range product_image_list_in_db {
		if _, found := product_image_url_map_download[product_image_in_db.Url]; found == false {
			db_session.Id(product_image_in_db.Id).Delete(new(models.ProductImageInfo))
		}
	}

	// 同步商品图片 end

	return product.Id, product.Sku
}
