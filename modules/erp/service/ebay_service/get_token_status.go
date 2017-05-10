package ebay_service

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"

	"github.com/franela/goreq"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/ebay"
	"matrix/modules/erp/models"
)

/*

get token status

http://developer.ebay.com/devzone/xml/docs/Reference/eBay/GetTokenStatus.html

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
<GetTokenStatusRequest xmlns="urn:ebay:apis:eBLBaseComponents">
  <!-- (No call-specific Input fields) -->

  <!-- Standard Input Fields -->
  <ErrorLanguage> string </ErrorLanguage>
  <MessageID> string </MessageID>
  <Version> string </Version>
  <WarningLevel> WarningLevelCodeType </WarningLevel>
</GetTokenStatusRequest>


response body:

<?xml version="1.0" encoding="utf-8"?>
<GetTokenStatusResponse xmlns="urn:ebay:apis:eBLBaseComponents">
  <!-- Call-specific Output Fields -->
  <TokenStatus> TokenStatusType
    <EIASToken> string </EIASToken>
    <ExpirationTime> dateTime </ExpirationTime>
    <RevocationTime> dateTime </RevocationTime>
    <Status> TokenStatusCodeType </Status>
  </TokenStatus>
  <!-- Standard Output Fields -->
  <Ack> AckCodeType </Ack>
  <Build> string </Build>
  <CorrelationID> string </CorrelationID>
  <Errors> ErrorType
    <ErrorClassification> ErrorClassificationCodeType </ErrorClassification>
    <ErrorCode> token </ErrorCode>
    <ErrorParameters ParamID="string"> ErrorParameterType
      <Value> string </Value>
    </ErrorParameters>
    <!-- ... more ErrorParameters nodes allowed here ... -->
    <LongMessage> string </LongMessage>
    <SeverityCode> SeverityCodeType </SeverityCode>
    <ShortMessage> string </ShortMessage>
  </Errors>
  <!-- ... more Errors nodes allowed here ... -->
  <HardExpirationWarning> string </HardExpirationWarning>
  <Timestamp> dateTime </Timestamp>
  <Version> string </Version>
</GetTokenStatusResponse>

 */

func Get_token_status(config models.EBayAPIConfigInfo) ebay.GetTokenStatusResponse {

	var req_body = `<?xml version="1.0" encoding="utf-8"?>
<GetTokenStatusRequest xmlns="urn:ebay:apis:eBLBaseComponents">
	<RequesterCredentials>
		<eBayAuthToken>%s</eBayAuthToken>
	</RequesterCredentials>
	<ErrorLanguage>zh_CN</ErrorLanguage>
	<WarningLevel>High</WarningLevel>
</GetTokenStatusRequest>
`
	req_body = fmt.Sprintf(req_body, config.Token)

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
		revel.ERROR.Println("eBay GetTokenStatus ERROR - processUrl -> req.Do: " + err.Error())
		return ebay.GetTokenStatusResponse{}
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("eBay GetTokenStatus ERROR - ioutil.ReadAll : " + err.Error())
		return ebay.GetTokenStatusResponse{}
	}

	response := ebay.GetTokenStatusResponse{}
	err = xml.Unmarshal(data, &response)
	core.HandleError(err)

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("eBay GetTokenStatus ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return ebay.GetTokenStatusResponse{}
	}

	return response
}
