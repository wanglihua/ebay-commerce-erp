package ebay

import (
	"encoding/xml"
)

type GetMemberMessagesResponse struct {
	Xmlns                 string                            `xml:"xmlns,attr"`
	Timestamp             string
	Ack                   string
	Version               string
	Build                 string
	Errors                Errors
	HasMoreItems          bool
	CorrelationId         string                           `xml:"CorrelationID"`
	HardExpirationWarning string
	Pagination            PaginationResult                 `xml:"PaginationResult"`
	MemberMessagesArray   MemberMessageArray               `xml:"MemberMessage"`
}

type MemberMessageArray struct {
	XMLName                xml.Name                  `xml:"MemberMessage"`
	MemberMessageExchanges []MemberMessageExchange   `xml:"MemberMessageExchange"`
}

type MemberMessageExchange struct {
	XMLName          xml.Name         `xml:"MemberMessageExchange"`
	Item             MsgPrdItem
	CreationDate     string
	LastModifiedDate string
	MessageStatus    string
	Response         string
	MessageMedia     MessageMedia     `xml:"MessageMedia"` //这里也不知道有问题不，是不是应该是个数组
	Question         Question
}

type MsgPrdItem struct {
	XMLName        xml.Name                  `xml:"Item"`
	ItemId         string                    `xml:"ItemID"`
	Title          string                    `xml:"Title"`
	ListingDetails MsgPrdItemListingDetails `xml:"ListingDetails"`
	PictureDetails MsgPrdItemPictureDetails `xml:"PictureDetails"`
}

type MsgPrdItemListingDetails struct {
	XMLName     xml.Name       `xml:"ListingDetails"`
	ViewItemUrl string         `xml:"ViewItemURL"`
}

type MsgPrdItemPictureDetails struct {
	XMLName      xml.Name     `xml:"PictureDetails"`
	GalleryType  string       `xml:"GalleryType"`
	GalleryUrl   string       `xml:"GalleryURL"`
	PhotoDisplay string       `xml:"PhotoDisplay"`
	PictureUrls  []string     `xml:"PictureURL"`
}

type MessageMedia struct {
	XMLName   xml.Name                `xml:"MessageMedia"`
	MediaName string
	MediaUrl  string                  `xml:"MediaURL"`
}

type Question struct {
	DisplayToPublic bool
	MessageId       string           `xml:"MessageID"`
	MessageMedia    MessageMedia     `xml:"MessageMedia"` //这里也不知道有问题不，是不是应该是个数组
	MessageType     string
	QuestionType    string
	RecipientId     string           `xml:"RecipientID"`
	SenderEmail     string
	SenderId        string           `xml:"SenderID"`
	Subject         string
	Body            string
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
