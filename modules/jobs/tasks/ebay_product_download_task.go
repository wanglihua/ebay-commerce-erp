package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"runtime/debug"
	"time"

	"github.com/franela/goreq"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/ebay"
	"matrix/modules/erp/models"
)

/*
get member messages document:
http://developer.ebay.com/Devzone/XML/docs/Reference/eBay/GetSellerList.html

test tools:
http://developer.ebay.com/DevZone/build-test/test-tool/?index=0&api=trading&call=GetMemberMessages&variation=xml#

request header:
X-EBAY-API-CALL-NAME:GetMemberMessages
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
<GetSellerListRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>AgAAAA**AQAAAA**aAAAAA**uVWPWA**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4GiD5aDpAydj6x9nY+seQ**pRMEAA**AAMAAA**VrbSv0o7iGoYiFvprWmTAsvp1TYnEDq74jjQOsqLWlwk5eS1y2Uyrq9aYWuWdWvRuM7/L0nUP1IFac3aBTY2aanmI+37H8fN0mUrSY7gFIwEdGx1pneZFwweLzc3l0yZcV1m22tWrSqRMQ6g+A31y8vxtjGHkSZkbeMJQJhMHLo5gMOjo2veZNacluDlzL2hKrs+Z7iIPuqxqwYEz2hDGe36BeXdBY5sDUvbZBdRVV4ws4Ta5L5Bal5tJO8tH/8WO0fXu5OKlw/BIGpgtbQc8Xe5MmhotfbznjfnIgrzxoSbsky6HZqpQ7AQh3UBZ/T7i2H5BdD9Rr3sMeXas6x7E8V0xsDqfQL1N8qAfKCOhDVP8Vxg/AOxP3dVMQaqaX0gBRp8CTnMIlYvHJKzW9URZ5vipjAJ+Fcx/IW5sHdiJFGHawHlfi0yPS3cYecLcPBUAKYqr1fBYqEZAQGsqId/PMVue8hwbl/gzaLus3WAOUPME4hFXMu2n9ADkYturlMa2TYRyHufNJmziHHUUb68bALZpTMhwH6QuG4ILSIluVDz9YcM7GuRGtQzBDhbTRdim6WGcjO7dkogo665c8QqSBIW2kT1+TTT24gXs8dUkv6Zea3Ihiko/bgh6mm2vCJ+LtMNqk+/vFHhL4nMhToDGX+0fFxu2U43aq3FJGo5GEYbn6MowwY/JzRHFKDkN2wzVKo7DJ2c1+faulxDGyVGhOMeEmfdZgL3iEd24JPm0mcvT1xA/UXpDYyCdhRb6I7+</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<StartTimeFrom>dateTime</StartTimeFrom>
	<StartTimeTo>dateTime</StartTimeTo>
	<EndTimeFrom>dateTime</EndTimeFrom>
	<EndTimeTo>dateTime</EndTimeTo>
	<Pagination>
		<EntriesPerPage>200</EntriesPerPage>
		<PageNumber>1</PageNumber>
	</Pagination>
</GetSellerListRequest>

 */

type EBayProductDownloadTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t EBayProductDownloadTask) Name() string {
	return "商品下载"
}

func (t EBayProductDownloadTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("下载商品出现错误： %s\n", err)
			revel.ERROR.Println(string(debug.Stack()))
		}
	}()

	job_time_begin := time.Now()

	// get db driver and db url

	if t.db_driver == "" {
		t.db_driver, _ = revel.Config.String("db.driver")
	}

	if t.db_url == "" {
		t.db_url, _ = revel.Config.String("db.url")
	}

	db_engine, err := xorm.NewEngine(t.db_driver, t.db_url)
	core.HandleError(err)

	db_engine.TZLocation = time.Local
	db_engine.SetMaxIdleConns(3)
	db_engine.SetMaxOpenConns(100)

	db_session := db_engine.NewSession()
	db_session = db_session.NoAutoCondition(true)

	defer func() {
		db_session.Close()
	}()

	ebay_api_config_list := get_ebay_api_config_list(db_session)

	for _, config := range ebay_api_config_list {
		has_more_orders := true
		page_number := 1
		for has_more_orders == true {

			list_response := t.fetch_product_list(config.APIConfig, page_number)

			switch list_response.Ack {
			case "Success":
				revel.TRACE.Printf("eBay GetSellerList Response: " + list_response.Ack + "\n")
				revel.TRACE.Printf("eBay GetSellerList Response HasMoreItems: %t", list_response.HasMoreItems)
			case "Failure":
				revel.ERROR.Printf("eBay GetSellerList Response: " + list_response.Ack + "\n")
			default:
				revel.ERROR.Printf("eBay GetSellerList Response: " + list_response.Ack + "\n")
			}

			has_more_orders = list_response.HasMoreItems
			page_number += 1

			//获取产品的详细信息
			for _, list_item := range list_response.ItemArray.Items {

				/*
				SELECT count(psi.id) AS product_item_count FROM erp_product_store_item psi WHERE psi.shop_id = 1 AND psi.item_id = '123'
				 */
				sql := `SELECT count(psi.id) AS product_item_count FROM erp_product_store_item psi WHERE psi.shop_id = %d AND psi.item_id = '%s'`
				sql = fmt.Sprintf(sql, config.APIConfig.ShopId, list_item.ItemId)

				var product_item_count struct {
					Count int64 `xorm:"'product_item_count'"`
				}

				_, err := db_session.Sql(sql).Get(&product_item_count)
				core.HandleError(err)

				if product_item_count.Count == 0 {
					fetch_product_item_to_db(config.APIConfig, db_session, list_item.ItemId)
				}
			}
		} // for has_more_orders == true
	} //for _, config := range ebay_api_config_list

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("eBay GetSellerList time duration: ", job_duration.Seconds(), "seconds")
}

func (t EBayProductDownloadTask) fetch_product_list(config models.EBayAPIConfigInfo, page_number int) ebay.GetSellerListResponse {

	number_of_days := 120 - 1

	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetSellerListRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<StartTimeFrom>%s</StartTimeFrom>
	<StartTimeTo>%s</StartTimeTo>
	<Pagination>
		<EntriesPerPage>200</EntriesPerPage>
		<PageNumber>%d</PageNumber>
	</Pagination>
</GetSellerListRequest>`

	time_range_end := time.Now().UTC()
	time_range_begin := time_range_end.Add(- time.Duration(number_of_days) * 24 * time.Hour)

	req_body = fmt.Sprintf(req_body, config.Token, time_range_begin.Format("2006-01-02T15:04:05"), time_range_end.Format("2006-01-02T15:04:05"), page_number)

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.GetEndPoint(),
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		ShowDebug:   revel.DevMode,
	}

	req.AddHeader("X-EBAY-API-CALL-NAME", "GetSellerList")
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

	var response ebay.GetSellerListResponse

	if err != nil {
		revel.ERROR.Println("eBay GetSellerList ERROR - processUrl -> req.Do: " + err.Error())
		return response
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("eBay GetSellerList ERROR - ioutil.ReadAll : " + err.Error())
		return response
	}

	if revel.DevMode && log_response_xml_to_file {
		core.Force_write_to_file(fmt.Sprintf(`z:\test_product_response\%s.xml`, time.Now().Format("20060102150405")), string(data))
		//core.Force_write_to_file(fmt.Sprintf(`/home/wanglihua/ebay_sync/test_product_response/%s.xml`, time.Now().Format("20060102150405")), string(data))
	}

	err = xml.Unmarshal(data, &response)
	core.HandleError(err)

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetSellerList ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return response
	}

	return response
}
