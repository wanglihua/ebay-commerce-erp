package tasks

import (
	"github.com/revel/revel"
	"github.com/mithereal/go-ebay"
	"encoding/xml"
	"github.com/franela/goreq"
	"io/ioutil"
	"fmt"
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
	<MailMessageType>AskSellerQuestion</MailMessageType>
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

func (t EBayMessageDownloadTask) Run() {
	config := GetEBayConfig()

	response := t.fetch_messages(config)

	switch response.Ack {
	case "Success":
		revel.TRACE.Printf("eBay GetMemberMessages Response: " + response.Ack + "\n")
		revel.TRACE.Printf("eBay GetMemberMessages Response HasMoreItems: %t", response.HasMoreMessages)
	case "Failure":
		revel.ERROR.Printf("eBay GetMemberMessages Response: " + response.Ack + "\n")
	default:
		revel.ERROR.Printf("eBay GetMemberMessages Response: " + response.Ack + "\n")
	}
}

func (t EBayMessageDownloadTask) fetch_messages(config ebay.Config) ebay.GetMemberMessagesRequestResponse {
	req_body := `<?xml version="1.0" encoding="utf-8"?>
<GetMemberMessagesRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
	<MailMessageType>AskSellerQuestion</MailMessageType>
	<MessageStatus>Unanswered</MessageStatus>
	<StartCreationTime>%s</StartCreationTime>
	<EndCreationTime>%s</EndCreationTime>
	<Pagination>
		<EntriesPerPage>5</EntriesPerPage>
		<PageNumber>1</PageNumber>
	</Pagination>
</GetMemberMessagesRequest>
`

	req_body = fmt.Sprintf(req_body, config.Token, "2017-01-01T00:00:00", "2017-01-31T23:59:59")

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.EndPoint,
		Body:        req_body,
		ContentType: "application/xml; charset=utf-8",
		//UserAgent:   "go-ebay-fetch-orders",
		ShowDebug:   false,
	}

	req.AddHeader("X-EBAY-API-CALL-NAME", "GetMemberMessages")
	req.AddHeader("X-EBAY-API-DEV-NAME", config.DevID)
	req.AddHeader("X-EBAY-API-CERT-NAME", config.CertID)
	req.AddHeader("X-EBAY-API-APP-NAME", config.AppID)
	req.AddHeader("X-EBAY-API-VERSION", ebay.Version)
	req.AddHeader("X-EBAY-API-COMPATIBILITY-LEVEL", ebay.Version)
	req.AddHeader("X-EBAY-API-REQUEST-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-RESPONSE-ENCODING", "XML")
	req.AddHeader("X-EBAY-API-SITEID", config.SiteId)
	req.AddHeader("Accept", "application/xml,application/xhtml+xml")

	res, err := req.Do()

	var response ebay.GetMemberMessagesRequestResponse

	if err != nil {
		revel.ERROR.Println("ERROR - processUrl -> req.Do: " + err.Error())
		return response
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("ERROR - ioutil.ReadAll : " + err.Error())
		return response
	}

	xml.Unmarshal(data, &response)

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return response
	}

	return response
}
