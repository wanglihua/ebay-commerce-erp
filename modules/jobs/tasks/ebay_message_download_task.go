package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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
http://developer.ebay.com/Devzone/XML/docs/Reference/eBay/GetMemberMessages.html

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
<GetMemberMessagesRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>AgAAAA**AQAAAA**aAAAAA**uVWPWA**nY+sHZ2PrBmdj6wVnY+sEZ2PrA2dj6wFk4GiD5aDpAydj6x9nY+seQ**pRMEAA**AAMAAA**VrbSv0o7iGoYiFvprWmTAsvp1TYnEDq74jjQOsqLWlwk5eS1y2Uyrq9aYWuWdWvRuM7/L0nUP1IFac3aBTY2aanmI+37H8fN0mUrSY7gFIwEdGx1pneZFwweLzc3l0yZcV1m22tWrSqRMQ6g+A31y8vxtjGHkSZkbeMJQJhMHLo5gMOjo2veZNacluDlzL2hKrs+Z7iIPuqxqwYEz2hDGe36BeXdBY5sDUvbZBdRVV4ws4Ta5L5Bal5tJO8tH/8WO0fXu5OKlw/BIGpgtbQc8Xe5MmhotfbznjfnIgrzxoSbsky6HZqpQ7AQh3UBZ/T7i2H5BdD9Rr3sMeXas6x7E8V0xsDqfQL1N8qAfKCOhDVP8Vxg/AOxP3dVMQaqaX0gBRp8CTnMIlYvHJKzW9URZ5vipjAJ+Fcx/IW5sHdiJFGHawHlfi0yPS3cYecLcPBUAKYqr1fBYqEZAQGsqId/PMVue8hwbl/gzaLus3WAOUPME4hFXMu2n9ADkYturlMa2TYRyHufNJmziHHUUb68bALZpTMhwH6QuG4ILSIluVDz9YcM7GuRGtQzBDhbTRdim6WGcjO7dkogo665c8QqSBIW2kT1+TTT24gXs8dUkv6Zea3Ihiko/bgh6mm2vCJ+LtMNqk+/vFHhL4nMhToDGX+0fFxu2U43aq3FJGo5GEYbn6MowwY/JzRHFKDkN2wzVKo7DJ2c1+faulxDGyVGhOMeEmfdZgL3iEd24JPm0mcvT1xA/UXpDYyCdhRb6I7+</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<MailMessageType>All</MailMessageType>
	<MessageStatus>Unanswered</MessageStatus>
	<StartCreationTime>2017-01-01T00:00:00</StartCreationTime>
	<EndCreationTime>2017-01-31T23:59:59</EndCreationTime>
	<Pagination>
		<EntriesPerPage>5</EntriesPerPage>
		<PageNumber>1</PageNumber>
	</Pagination>
</GetMemberMessagesRequest>

response body:
<?xml version="1.0" encoding="UTF-8"?>
<GetMemberMessagesResponse xmlns="urn:ebay:apis:eBLBaseComponents">
	<Timestamp>2017-02-01T07:21:45.555Z</Timestamp>
	<Ack>Success</Ack>
	<Version>989</Version>
	<Build>E989_CORE_API_18131074_R1</Build>
	<PaginationResult>
		<TotalNumberOfPages>0</TotalNumberOfPages>
		<TotalNumberOfEntries>0</TotalNumberOfEntries>
	</PaginationResult>
	<HasMoreItems>false</HasMoreItems>
</GetMemberMessagesResponse>

 */

type EBayMessageDownloadTask struct {
	db_driver string //mysql
	db_url    string //root:root@tcp(localhost:3306)/iot?charset=utf8
}

func (t EBayMessageDownloadTask) Run() {
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

	ebay_api_config_list := make([]models.EBayAPIConfigInfo, 0)

	/*
	SELECT cfg.* FROM erp_ebay_api_config cfg INNER JOIN erp_shop_account sa ON cfg.shop_id = sa.id
	*/
	err = db_session.Sql("SELECT cfg.* FROM erp_ebay_api_config cfg INNER JOIN erp_shop_account sa ON cfg.shop_id = sa.id").Find(&ebay_api_config_list)
	core.HandleError(err)

	for _, config := range ebay_api_config_list {
		has_more_orders := true
		page_number := 1
		for has_more_orders == true {

			response := t.fetch_messages(config, page_number)

			switch response.Ack {
			case "Success":
				revel.TRACE.Printf("eBay GetMemberMessages Response: " + response.Ack + "\n")
				revel.TRACE.Printf("eBay GetMemberMessages Response HasMoreItems: %t", response.HasMoreItems)
			case "Failure":
				revel.ERROR.Printf("eBay GetMemberMessages Response: " + response.Ack + "\n")
			default:
				revel.ERROR.Printf("eBay GetMemberMessages Response: " + response.Ack + "\n")
			}

			has_more_orders = response.HasMoreItems
			page_number += 1
		} // for has_more_orders == true
	} //for _, config := range ebay_api_config_list

	job_duration := time.Since(job_time_begin)

	revel.TRACE.Println("sync message time duration: ", job_duration.Seconds(), "seconds")
}

func (t EBayMessageDownloadTask) fetch_messages(config models.EBayAPIConfigInfo, page_number int) ebay.GetMemberMessagesRequestResponse {
	/*

	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetMemberMessagesRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<MailMessageType>All</MailMessageType>
	<MessageStatus>Answered</MessageStatus>
	<StartCreationTime>%s</StartCreationTime>
	<EndCreationTime>%s</EndCreationTime>
	<Pagination>
		<EntriesPerPage>5</EntriesPerPage>
		<PageNumber>%d</PageNumber>
	</Pagination>
</GetMemberMessagesRequest>
`
	req_body = fmt.Sprintf(req_body, config.Token, "2017-01-01T00:00:00", "2017-01-31T23:59:59", page_number)
	*/

	number_of_days := 1

	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetMemberMessagesRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<MailMessageType>All</MailMessageType>
	<StartCreationTime>%s</StartCreationTime>
	<EndCreationTime>%s</EndCreationTime>
	<Pagination>
		<EntriesPerPage>200</EntriesPerPage>
		<PageNumber>%d</PageNumber>
	</Pagination>
</GetMemberMessagesRequest>
`
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

	req.AddHeader("X-EBAY-API-CALL-NAME", "GetMemberMessages")
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

	var response ebay.GetMemberMessagesRequestResponse

	if err != nil {
		revel.ERROR.Println("eBay GetMemberMessages ERROR - processUrl -> req.Do: " + err.Error())
		return response
	}

	data, err := ioutil.ReadAll(res.Body)

	if revel.DevMode {
		core.Force_write_to_file(fmt.Sprintf(`z:\test_message_response\%s.xml`, time.Now().Format("20060102150405")), string(data))
	}

	if err != nil {
		revel.ERROR.Println("eBay GetMemberMessages ERROR - ioutil.ReadAll : " + err.Error())
		return response
	}

	err = xml.Unmarshal(data, &response)
	if err != nil {
		revel.ERROR.Println(err)
	}

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetMemberMessages ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return response
	}

	return response
}
