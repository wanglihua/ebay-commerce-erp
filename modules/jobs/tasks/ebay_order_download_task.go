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
	"os"
	"strconv"
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
	job_time_begin := time.Now()

	// get db driver and db url

	if t.db_driver == "" {
		t.db_driver, _ = revel.Config.String("db.driver")
	}

	if t.db_url == "" {
		t.db_url, _ = revel.Config.String("db.url")
	}

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
			if err != nil {
				revel.ERROR.Println(err)
			}

			revel.TRACE.Println("OrderStatus: ", order_download.OrderStatus)

			order_in_db.TypeId = models.ORDER_TYPE_XZDD
			order_in_db.EBayStatus = order_download.OrderStatus
			order_in_db.StatusId = models.ORDER_STATUS_YZF

			update_col_list := make([]string, 0)
			update_col_list = append(update_col_list, "type_id")
			update_col_list = append(update_col_list, "ebay_status")
			update_col_list = append(update_col_list, "status_id")

			if order_download.CreatedTime != "" {
				created_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.CreatedTime)
				if err != nil {
					revel.ERROR.Println(err)
				}
				// 转成北京时间
				created_time = created_time.Add(8 * time.Hour)
				order_in_db.CreatedTime = core.NewNullLongTime(created_time, true)
				update_col_list = append(update_col_list, "created_time")
			}

			if order_download.ModifiedTime != "" {
				modified_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.ModifiedTime)
				if err != nil {
					revel.ERROR.Println(err)
				}
				// 转成北京时间
				modified_time = modified_time.Add(8 * time.Hour)
				order_in_db.ModifiedTime = core.NewNullLongTime(modified_time, true)
				update_col_list = append(update_col_list, "modified_time")
			}

			if order_download.PaidTime != "" {
				paid_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.PaidTime)
				if err != nil {
					revel.ERROR.Println(err)
				}
				// 转成北京时间
				paid_time = paid_time.Add(8 * time.Hour)
				order_in_db.PaidTime = core.NewNullLongTime(paid_time, true)
				update_col_list = append(update_col_list, "paid_time")
			}

			if order_download.ShippedTime != "" {
				shipped_time, err := time.Parse("2006-01-02T15:04:05.000Z", order_download.ShippedTime)
				if err != nil {
					revel.ERROR.Println(err)
				}
				// 转成北京时间
				shipped_time = shipped_time.Add(8 * time.Hour)
				order_in_db.ShippedTime = core.NewNullLongTime(shipped_time, true)
				update_col_list = append(update_col_list, "shipped_time")
			}

			order_in_db.AdjustmentAmountCurrency = order_download.AdjustmentAmount.Currency
			update_col_list = append(update_col_list, "adjustment_amount_currency")
			order_in_db.AdjustmentAmount, err = strconv.ParseFloat(order_download.AdjustmentAmount.Value, 64)
			if err != nil {
				revel.ERROR.Println(err)
			}
			update_col_list = append(update_col_list, "adjustment_amount")

			order_in_db.AmountPaidCurrency = order_download.AmountPaid.Currency
			update_col_list = append(update_col_list, "amount_paid_currency")
			order_in_db.AmountPaid, err = strconv.ParseFloat(order_download.AmountPaid.Value, 64)
			if err != nil {
				revel.ERROR.Println(err)
			}
			update_col_list = append(update_col_list, "amount_paid")

			order_in_db.AmountSavedCurrency = order_download.AmountSaved.Currency
			update_col_list = append(update_col_list, "amount_saved_currency")
			order_in_db.AmountSaved, err = strconv.ParseFloat(order_download.AmountSaved.Value, 64)
			if err != nil {
				revel.ERROR.Println(err)
			}
			update_col_list = append(update_col_list, "amount_saved")

			if has == true {
				// update order
				_, err = db_session.Id(order_in_db.Id).Cols(update_col_list...).Update(&order_in_db)
				if err != nil {
					revel.ERROR.Println(err)
				}
			} else {
				// insert order
				order_in_db.Code = order_code

				_, err = db_session.Insert(&order_in_db)
				if err != nil {
					revel.ERROR.Println(err)
				}
			}

			_ = order_code
		}

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

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("sync order time duration: ", job_duration.Seconds(), "seconds")
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
	<OrderStatus>All</OrderStatus>
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

	//revel.TRACE.Println(string(data))
	t.force_write_to_file(fmt.Sprintf(`z:\test_order_response\%s.xml`, time.Now().Format("20060102150405")), string(data))

	if err != nil {
		revel.ERROR.Println("eBay GetOrders ERROR - ioutil.ReadAll : " + err.Error())
		return ebay.GetOrdersRequestResponse{}
	}

	response := ebay.GetOrdersRequestResponse{}
	xml.Unmarshal(data, &response)

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetOrders ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return ebay.GetOrdersRequestResponse{}
	}

	return response
}

func (t EBayOrderDownloadTask) force_write_to_file(file_name string, content string) {

	is_file_exist := func(file_name string) bool {
		var exist = true
		if _, err := os.Stat(file_name); os.IsNotExist(err) {
			exist = false
		}
		return exist
	}

	if is_file_exist(file_name) {
		err := os.Remove(file_name)
		core.HandleError(err)
	}

	var file *os.File
	var err error

	file, err = os.Create(file_name) //创建文件
	if err != nil {
		revel.ERROR.Println(err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		revel.ERROR.Println(err)
	}

	file.Sync()
	file.Close()
}
