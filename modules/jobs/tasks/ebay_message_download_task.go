package tasks

import (
	"github.com/revel/revel"
	"github.com/mithereal/go-ebay"
	"encoding/xml"
	"github.com/franela/goreq"
	"io/ioutil"
)

type EBayMessageDownloadTask struct {
}

func (t EBayMessageDownloadTask) Run() {
	config := GetEBayConfig()

	credentials := ebay.RequesterCredentials{
		RequestToken: config.Token,
	}

	member_messages_request := ebay.GetMemberMessagesRequest{
		Xmlns:             "urn:ebay:apis:eBLBaseComponents",
		DisplayToPublic:   "",
		EndCreationTime:   "2017-01-20T00:00:00",
		//ItemID:            *ItemID,
		MailMessageType:      "",
		MessageStatus:        "",
		SenderID:             "",
		StartCreationTime:    "2017-01-20T00:00:00",
		ErrorLanguage:        "",
		MessageID:            "",
		OutputSelector:       "",
		Version:              ebay.Version, // 这里是933，原来的实例中是935
		RequesterCredentials: credentials,
	}

	member_messages_xml, err := xml.Marshal(member_messages_request)

	if err != nil {
		revel.ERROR.Println("ERROR - xml.Marshal : " + err.Error())
		return
	}

	xml_header := []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>")

	body := append(xml_header, member_messages_xml...)

	req := goreq.Request{
		Method:      "POST",
		Uri:         config.EndPoint,
		Body:        body,
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

	if err != nil {
		revel.ERROR.Println("ERROR - processUrl -> req.Do: " + err.Error())
		return
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		revel.ERROR.Println("ERROR - ioutil.ReadAll : " + err.Error())
		return
	}

	var response ebay.GetMemberMessagesRequestResponse

	xml.Unmarshal(data, &response)

	e := response.Errors

	if e.ShortMessage != "" {
		revel.ERROR.Println("ERROR - " + e.ErrorCode + " : " + e.LongMessage)
		return
	}

	// dump response
}
