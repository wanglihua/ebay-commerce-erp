package tasks

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
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
}

func (t EBayMessageDownloadTask) Name() string {
	return "消息下载"
}

func (t EBayMessageDownloadTask) Run() {
	defer func() {
		if err := recover(); err != nil {
			revel.ERROR.Printf("下载消息出现错误： %s\n", err)
		}
	}()

	job_time_begin := time.Now()

	// get db driver and db url

	if task_db_driver == "" {
		task_db_driver, _ = revel.Config.String("db.driver")
	}

	if task_db_url == "" {
		task_db_url, _ = revel.Config.String("db.url")
	}

	db_engine, err := xorm.NewEngine(task_db_driver, task_db_url)
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

	ebay_api_config_list := get_ebay_api_config_list(db_session)

	for _, config := range ebay_api_config_list {
		has_more_orders := true
		page_number := 1
		for has_more_orders == true {

			response := t.fetch_messages(config, page_number)

			// insert into database
			message_download_list := response.MemberMessagesArray.MemberMessageExchanges
			for _, message_download := range message_download_list {

				message_id := message_download.Question.MessageId

				var message_in_db models.MessageInfo
				message_in_db.ShopId = config.ShopId

				has_message, err := db_session.Where(fmt.Sprintf("message_id = '%s'", message_id)).Get(&message_in_db)
				if err != nil {
					revel.ERROR.Println(err)
				}
				message_in_db.MessageId = message_id

				// 如果这时产品不存在会怎样，下载产品
				var product_result struct {
					ProductId  int64 `xorm:"'product_id'"`
					ProductSku string `xorm:"'product_sku'"`
				}
				has_item_product, err := db_session.Sql(fmt.Sprintf("SELECT p.id AS product_id, p.sku AS product_sku FROM erp_product p WHERE p.item_id = '%s'", message_download.Item.ItemId)).Get(&product_result)
				if err != nil {
					revel.ERROR.Println(err)
				}

				product_id := product_result.ProductId
				product_sku := product_result.ProductSku
				if has_item_product == false {
					// download product
					product_id, product_sku, _ = fetch_product_item_to_db(config, db_session, message_download.Item.ItemId)
				}

				message_in_db.ProductId = product_id
				message_in_db.Sku = product_sku
				message_in_db.ItemId = message_download.Item.ItemId

				if message_download.CreationDate != "" {
					creation_date, err := time.Parse("2006-01-02T15:04:05.000Z", message_download.CreationDate)
					if err != nil {
						revel.ERROR.Println(err)
					}
					// 转成北京时间
					// modified_time = modified_time.Add(8 * time.Hour)
					creation_date = creation_date.Local()
					message_in_db.CreationDate = core.MakeNullLongTime(creation_date, true)
				}

				if message_download.LastModifiedDate != "" {
					last_modified_date, err := time.Parse("2006-01-02T15:04:05.000Z", message_download.LastModifiedDate)
					if err != nil {
						revel.ERROR.Println(err)
					}
					// 转成北京时间
					last_modified_date = last_modified_date.Local()
					message_in_db.LastModifiedDate = core.MakeNullLongTime(last_modified_date, true)
				}

				message_in_db.MessageStatus = message_download.MessageStatus // Answered Unanswered
				if strings.ToLower(message_download.MessageStatus) == strings.ToLower("Answered") {
					message_in_db.StatusId = models.MESSAGE_STATUS_ANSWERED
				} else if strings.ToLower(message_download.MessageStatus) == strings.ToLower("Unanswered") {
					message_in_db.StatusId = models.MESSAGE_STATUS_UNANSWERED
				}

				message_in_db.DisplayToPublic = message_download.Question.DisplayToPublic
				message_in_db.MessageType = message_download.Question.MessageType
				message_in_db.QuestionType = message_download.Question.QuestionType
				message_in_db.RecipientId = message_download.Question.RecipientId
				message_in_db.SenderEmail = message_download.Question.SenderEmail
				message_in_db.SenderId = message_download.Question.SenderId
				message_in_db.Subject = message_download.Question.Subject
				message_in_db.Body = message_download.Question.Body

				if has_message == true {
					// do update
					_, err = db_session.Id(message_in_db.Id).Update(&message_in_db)
					if err != nil {
						revel.ERROR.Println(err)
					}
				} else {
					// do insert
					_, err = db_session.Insert(&message_in_db)
					if err != nil {
						revel.ERROR.Println(err)
					}
				}
			}

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

	revel.TRACE.Println("eBay GetMemberMessages time duration: ", job_duration.Seconds(), "seconds")
}

func (t EBayMessageDownloadTask) fetch_messages(config models.EBayAPIConfigInfo, page_number int) ebay.GetMemberMessagesResponse {
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
	time_range_end := time.Now().UTC() //用的是 UTC 时间
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

	var response ebay.GetMemberMessagesResponse

	if err != nil {
		revel.ERROR.Println("eBay GetMemberMessages ERROR - processUrl -> req.Do: " + err.Error())
		return response
	}

	data, err := ioutil.ReadAll(res.Body)

	if revel.DevMode && log_response_xml_to_file {
		core.Force_write_to_file(fmt.Sprintf(`z:\test_message_response\%s.xml`, time.Now().Format("20060102150405")), string(data))
		//core.Force_write_to_file(fmt.Sprintf(`/home/wanglihua/ebay_sync/test_message_response/%s.xml`, time.Now().Format("20060102150405")), string(data))
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

/*

<?xml version="1.0" encoding="utf-8"?>
<GetMemberMessagesRequest xmlns="urn:ebay:apis:eBLBaseComponents">
  <RequesterCredentials>
    <eBayAuthToken>ABC...123</eBayAuthToken>
  </RequesterCredentials>
  <WarningLevel>High</WarningLevel>
  <MailMessageType>AskSellerQuestion</MailMessageType>
  <MessageStatus>Unanswered</MessageStatus>
  <StartCreationTime>2016-02-08T00:00:00.000Z</StartCreationTime>
  <EndCreationTime>2016-02-16T19:23:13.821Z</EndCreationTime>
  <Pagination>
    <EntriesPerPage>5</EntriesPerPage>
    <PageNumber>1</PageNumber>
  </Pagination>
</GetMemberMessagesRequest>


<GetMemberMessagesResponse xmlns="urn:ebay:apis:eBLBaseComponents">
  <Timestamp>2016-02-16T19:23:52.811Z</Timestamp>
  <Ack>Success</Ack>
  <Version>967</Version>
  <Build>E967_CORE_BUNDLED_10639415_R1</Build>
  <MemberMessage>
    <MemberMessageExchange>
      <Item>
        <ItemID>110043671232</ItemID>
        <ListingDetails>
          <StartTime>2016-02-16T18:52:41.000Z</StartTime>
          <EndTime>2016-02-23T18:52:41.000Z</EndTime>
          <ViewItemURL>http://cgi.sandbox.ebay.com/ws/eBayISAPI.dll?ViewItem&item=110043671232&category=377</ViewItemURL>
        </ListingDetails>
        <Seller>
          <UserID>magicalbookseller</UserID>
        </Seller>
        <SellingStatus>
          <CurrentPrice currencyID="USD">1.0</CurrentPrice>
        </SellingStatus>
        <TimeLeft>P6DT23H28M49S</TimeLeft>
        <Title>Harry Potter and the Philosopher's Stone</Title>
      </Item>
      <Question>
        <MessageType>AskSellerQuestion</MessageType>
        <QuestionType>CustomizedSubject</QuestionType>
        <DisplayToPublic>false</DisplayToPublic>
        <SenderID>bountifulbuyer</SenderID>
        <SenderEmail>BuyerEmailAddress</SenderEmail>
        <RecipientID>magicalbookseller</RecipientID>
        <Subject>
            bountifulbuyer has sent a question about :General question about this item for item #110043671232,
            ending on Feb-23-10 10:52:41 PST - Harry Potter and the Philosopher's Stone
        </Subject>
        <Body>
            I thought the first book in the Harry Potter series was titled "Harry Potter and the Sorcerer's Stone".
            Is this a different book?
        </Body>
        <MessageID>763286</MessageID>
      </Question>
      <MessageStatus>Unanswered</MessageStatus>
      <CreationDate>2016-02-16T19:22:26.000Z</CreationDate>
      <LastModifiedDate>2016-02-16T19:22:26.000Z</LastModifiedDate>
    </MemberMessageExchange>
  </MemberMessage>
  <PaginationResult>
    <TotalNumberOfPages>1</TotalNumberOfPages>
    <TotalNumberOfEntries>1</TotalNumberOfEntries>
  </PaginationResult>
  <HasMoreItems>false</HasMoreItems>
</GetMemberMessagesResponse>

 */
