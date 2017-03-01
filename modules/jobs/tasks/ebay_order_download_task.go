package tasks

import (
	"github.com/revel/revel"
	"github.com/mithereal/go-ebay"
	"io/ioutil"
	"encoding/xml"
	"github.com/franela/goreq"
	"fmt"
	"github.com/go-xorm/xorm"
	"time"
	"matrix/modules/erp/models"
	"matrix/core"
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

type EBayOrderDownloadTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t EBayOrderDownloadTask) Run() {
	// get db driver and db url

	db_engine, err := xorm.NewEngine(t.db_driver, t.db_url)
	if err != nil {
		revel.ERROR.Println(err)
		return
	}

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	config := GetEBayConfig()

	has_more_orders := true
	page_number := 1
	for has_more_orders == true {
		response := t.fetch_orders(config, page_number)

		order_list := response.OrdersArray.Orders
		for _, order_download := range order_list {
			order_code := order_download.OrderID

			var order_in_db models.OrderInfo
			has, err := db_session.Where("code = ?", order_code).Get(&order_in_db)
			core.HandleError(err)

			if has == true {
				// update order
				_, err = db_session.Id(order_in_db.Id).Cols("type_id").Update(&order_in_db)
				core.HandleError(err)
			} else {
				// insert order
				order_in_db.Code = order_code
				order_in_db.TypeId = models.ORDER_TYPE_XZDD
				order_in_db.StatusId = models.ORDER_STATUS_YZF

				_, err = db_session.Insert(&order_in_db)
				core.HandleError(err)
			}

			_ = order_code
		}

		//todo: 进行订单数据库操作等

		//todo: 可以把返回的 时间戳 存入数据库，作为下一次订单下载请求的开始时间，或者是这个时间的稍微前面一点点

		// do something with the orders

		switch response.Ack {
		case "Success":
			revel.TRACE.Printf("eBay GetOrders Response: " + response.Ack + "\n")
			revel.TRACE.Printf("eBay GetOrders actual fetched orders count: %d", response.ReturnedOrderCountActual)
		case "Failure":
			revel.ERROR.Printf("eBay GetOrders Response: " + response.Ack + "\n")
		default:
			revel.ERROR.Printf("eBay GetOrders Response: " + response.Ack + "\n")
		}

		has_more_orders = response.HasMoreOrders
		page_number += 1
	}
}

func (t EBayOrderDownloadTask) fetch_orders(config ebay.Config, page_number int) ebay.GetOrdersRequestResponse {
	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetOrdersRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<NumberOfDays>1</NumberOfDays>
	<OrderRole>Seller</OrderRole>
	<OrderStatus>Active</OrderStatus>
	<Pagination>
		<EntriesPerPage>100</EntriesPerPage>
		<PageNumber>%d</PageNumber>
	</Pagination>
</GetOrdersRequest>
`
	req_body = fmt.Sprintf(req_body, config.Token, page_number)

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.EndPoint,
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		ShowDebug:   false,
	}

	req.AddHeader("X-EBAY-API-CALL-NAME", "GetOrders")
	req.AddHeader("X-EBAY-API-DEV-NAME", config.DevID)
	req.AddHeader("X-EBAY-API-CERT-NAME", config.CertID)
	req.AddHeader("X-EBAY-API-APP-NAME", config.AppID)
	req.AddHeader("X-EBAY-API-VERSION", config.Version)
	req.AddHeader("X-EBAY-API-COMPATIBILITY-LEVEL", config.Version)
	req.AddHeader("X-EBAY-API-REQUEST-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-RESPONSE-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-SITEID", config.SiteId)
	req.AddHeader("Accept", "application/xml,application/xhtml+xml")

	res, err := req.Do()

	if err != nil {
		revel.ERROR.Println("eBay GetOrders ERROR - processUrl -> req.Do: " + err.Error())
		return ebay.GetOrdersRequestResponse{}
	}

	data, err := ioutil.ReadAll(res.Body)

	revel.TRACE.Println(string(data))

	if err != nil {
		revel.ERROR.Println("eBay GetOrders ERROR - ioutil.ReadAll : " + err.Error())
		return ebay.GetOrdersRequestResponse{}
	}

	response := ebay.GetOrdersRequestResponse{}
	xml.Unmarshal(data, &response)

	orders := ebay.OrderArray{}
	xml.Unmarshal(data, &orders)

	response.OrdersArray = orders

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetOrders ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return ebay.GetOrdersRequestResponse{}
	}

	return response
}
