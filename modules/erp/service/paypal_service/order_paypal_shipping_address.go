package paypal_service

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/franela/goreq"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/ddlc_models"
)

const log_response_xml_to_file = true

func Get_order_paypal_shipping_address(db_session *xorm.Session, order_id int64) string {
	var order ddlc_models.OrderInfo
	_, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	item_id := "282369681238"
	ebay_order_id := "282369681238-1645636260018"

	var api_config xtpz_models.EBayAPIConfigInfo
	_, err = db_session.Where(fmt.Sprintf("shop_id = %d", order.ShopId)).Get(&api_config)
	core.HandleError(err)

	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetItemTransactionsRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<ItemID>%s</ItemID>
	<OrderLineItemID>%s</OrderLineItemID>
	<DetailLevel>ReturnAll</DetailLevel>
</GetItemTransactionsRequest>
`
	req_body = fmt.Sprintf(req_body, api_config.Token, item_id, ebay_order_id)

	req := goreq.Request{
		Method:      "POST",
		Uri:         api_config.GetEndPoint(),
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		ShowDebug:   revel.DevMode,
	}

	req.AddHeader("X-EBAY-API-CALL-NAME", "GetItemTransactions")
	req.AddHeader("X-EBAY-API-DEV-NAME", api_config.DevId)
	req.AddHeader("X-EBAY-API-CERT-NAME", api_config.CertId)
	req.AddHeader("X-EBAY-API-APP-NAME", api_config.AppId)
	req.AddHeader("X-EBAY-API-VERSION", api_config.GetApiVersion())
	req.AddHeader("X-EBAY-API-COMPATIBILITY-LEVEL", api_config.GetApiVersion())
	req.AddHeader("X-EBAY-API-REQUEST-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-RESPONSE-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-SITEID", api_config.GetSiteId())
	req.AddHeader("Accept", "application/xml,application/xhtml+xml")

	res, err := req.Do()

	if err != nil {
		revel.ERROR.Println("eBay GetItemTransactionsRequest ERROR - processUrl -> req.Do: " + err.Error())
		return ""
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("eBay GetItemTransactionsRequest ERROR - ioutil.ReadAll : " + err.Error())
		return ""
	}

	if revel.DevMode && log_response_xml_to_file {
		core.Force_write_to_file(fmt.Sprintf(`z:\test_order_response\item_transactions_%s.xml`, time.Now().Format("20060102150405")), string(data))
		//core.Force_write_to_file(fmt.Sprintf(`/home/wanglihua/ebay_sync/test_order_response/%s.xml`, time.Now().Format("20060102150405")), string(data))
	}

	/*
	response := ebay.GetOrdersResponse{}
	err = xml.Unmarshal(data, &response)
	if err != nil {
		revel.ERROR.Println(err)
	}

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetOrders ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return ebay.GetOrdersResponse{}
	}
	*/

	return ""
}
