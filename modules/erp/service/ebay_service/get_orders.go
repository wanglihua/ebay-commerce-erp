package ebay_service

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/franela/goreq"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/ebay"
	"matrix/modules/erp/models"
)

/*

get orders document:
http://developer.ebay.com/devzone/xml/docs/Reference/eBay/GetOrders.html

test tools:
http://developer.ebay.com/DevZone/build-test/test-tool/?index=0&api=trading&call=GetOrders&variation=xml#

request header:
X-EBAY-API-CALL-NAME:GetOrders
X-EBAY-API-DEV-NAME:***REMOVED***
X-EBAY-API-CERT-NAME:***REMOVED***
X-EBAY-API-APP-NAME:***REMOVED***
X-EBAY-API-VERSION:967
X-EBAY-API-COMPATIBILITY-LEVEL:967
X-EBAY-API-REQUEST-ENCODING:XML
X-EBAY-API-RESPONSE-ENCODING:XML
X-EBAY-API-SITEID:0
Accept:application/xml,application/xhtml+xml

request body:
<?xml version="1.0" encoding="utf-8"?>
<GetOrdersRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>AgAAAA**AQAAAA**aAAAAA**uVWPWA**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4GiD5aDpAydj6x9nY+seQ**pRMEAA**AAMAAA**VrbSv0o7iGoYiFvprWmTAsvp1TYnEDq74jjQOsqLWlwk5eS1y2Uyrq9aYWuWdWvRuM7/L0nUP1IFac3aBTY2aanmI+37H8fN0mUrSY7gFIwEdGx1pneZFwweLzc3l0yZcV1m22tWrSqRMQ6g+A31y8vxtjGHkSZkbeMJQJhMHLo5gMOjo2veZNacluDlzL2hKrs+Z7iIPuqxqwYEz2hDGe36BeXdBY5sDUvbZBdRVV4ws4Ta5L5Bal5tJO8tH/8WO0fXu5OKlw/BIGpgtbQc8Xe5MmhotfbznjfnIgrzxoSbsky6HZqpQ7AQh3UBZ/T7i2H5BdD9Rr3sMeXas6x7E8V0xsDqfQL1N8qAfKCOhDVP8Vxg/AOxP3dVMQaqaX0gBRp8CTnMIlYvHJKzW9URZ5vipjAJ+Fcx/IW5sHdiJFGHawHlfi0yPS3cYecLcPBUAKYqr1fBYqEZAQGsqId/PMVue8hwbl/gzaLus3WAOUPME4hFXMu2n9ADkYturlMa2TYRyHufNJmziHHUUb68bALZpTMhwH6QuG4ILSIluVDz9YcM7GuRGtQzBDhbTRdim6WGcjO7dkogo665c8QqSBIW2kT1+TTT24gXs8dUkv6Zea3Ihiko/bgh6mm2vCJ+LtMNqk+/vFHhL4nMhToDGX+0fFxu2U43aq3FJGo5GEYbn6MowwY/JzRHFKDkN2wzVKo7DJ2c1+faulxDGyVGhOMeEmfdZgL3iEd24JPm0mcvT1xA/UXpDYyCdhRb6I7+</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<CreateTimeFrom>2017-01-01T00:00:00</CreateTimeFrom>
	<CreateTimeTo>2017-02-01T23:59:59</CreateTimeTo>
	<OrderRole>Seller</OrderRole>
	<OrderStatus>Active</OrderStatus>
</GetOrdersRequest>

response body:
<?xml version="1.0" encoding="UTF-8"?>
<GetOrdersResponse xmlns="urn:ebay:apis:eBLBaseComponents">
	<Timestamp>2017-01-31T12:57:53.266Z</Timestamp>
	<Ack>Success</Ack>
	<Version>989</Version>
	<Build>E989_CORE_API_18131074_R1</Build>
	<PaginationResult>
		<TotalNumberOfPages>0</TotalNumberOfPages>
		<TotalNumberOfEntries>0</TotalNumberOfEntries>
	</PaginationResult>
	<HasMoreOrders>false</HasMoreOrders>
	<OrderArray/>
	<OrdersPerPage>100</OrdersPerPage>
	<PageNumber>1</PageNumber>
	<ReturnedOrderCountActual>0</ReturnedOrderCountActual>
</GetOrdersResponse>

 */

func Get_orders(config models.EBayAPIConfigInfo, page_number int) ebay.GetOrdersResponse {

	number_of_days := 1

	/*
		The field is used to retrieve orders that are in a specific state. If this field is used with a date filter, only orders that satisfy both the date range and the OrderStatus value are retrieved.

		For eBay orders, this field's value can be set to 'Active', 'Completed', 'Canceled' or 'Inactive' to retrieve orders in these states. The 'Shipped' value is only applicable for Half.com orders.
	 */

	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetOrdersRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<NumberOfDays>%d</NumberOfDays>
	<OrderRole>Seller</OrderRole>
	<OrderStatus>All</OrderStatus>
	<Pagination>
		<EntriesPerPage>100</EntriesPerPage>
		<PageNumber>%d</PageNumber>
	</Pagination>
</GetOrdersRequest>
`

	req_body = fmt.Sprintf(req_body, config.Token, number_of_days, page_number)

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.GetEndPoint(),
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		ShowDebug:   revel.DevMode,
	}

	add_request_header(&req, config)

	res, err := req.Do()

	if err != nil {
		revel.ERROR.Println("eBay GetOrders ERROR - processUrl -> req.Do: " + err.Error())
		return ebay.GetOrdersResponse{}
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("eBay GetOrders ERROR - ioutil.ReadAll : " + err.Error())
		return ebay.GetOrdersResponse{}
	}

	if revel.DevMode && log_response_xml_to_file {
		core.Force_write_to_file(fmt.Sprintf(`z:\test_order_response\%s.xml`, time.Now().Format("20060102150405")), string(data))
		//core.Force_write_to_file(fmt.Sprintf(`/home/wanglihua/ebay_sync/test_order_response/%s.xml`, time.Now().Format("20060102150405")), string(data))
	}

	response := ebay.GetOrdersResponse{}
	err = xml.Unmarshal(data, &response)
	core.HandleError(err)

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetOrders ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return ebay.GetOrdersResponse{}
	}

	return response
}
