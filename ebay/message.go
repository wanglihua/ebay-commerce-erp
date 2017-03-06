package ebay

import (
	"encoding/xml"
)

type GetMemberMessagesRequestResponse struct {
	Xmlns                 string                            `xml:"xmlns,attr"`
	Timestamp             string
	Ack                   string
	Version               string
	Build                 string
	Errors                Errors
	HasMoreItems          bool
	CorrelationID         string
	HardExpirationWarning string
	Pagination            PaginationResult                 `xml:"PaginationResult"`
	MemberMessagesArray   MemberMessageArray               `xml:"MemberMessage"`
}

type MemberMessageArray struct {
	XMLName                xml.Name                  `xml:"MemberMessage"`
	MemberMessageExchanges []MemberMessageExchange   `xml:"MemberMessageExchange"`
}

type MemberMessageExchange struct {
	XMLName          xml.Name           `xml:"MemberMessageExchange"`
	Item             Item
	LastModifiedDate string
	MessageStatus    string
	Response         string
	MessageMedia    MessageMedia     `xml:"MessageMedia"`   //这里也不知道有问题不，是不是应该是个数组
	Question         Question
}

type MessageMedia struct {
	XMLName   xml.Name                  `xml:"MessageMedia"`
	MediaName string
	MediaURL  string
}

type Question struct {
	Body            string
	DisplayToPublic string
	MessageID       string
	MessageMedia   MessageMedia     `xml:"MessageMedia"`   //这里也不知道有问题不，是不是应该是个数组
	MessageType     string
	QuestionType    string
	RecipientID     string
	SenderEmail     string
	SenderID        string
	Subject         string
}
