package models

import (
	"matrix/core"
)

type MessageInfo struct {
	Id               int64             `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	ShopId           int64             `xorm:"bigint notnull index 'shop_id'" json:"shop_id"`
	MessageId        string            `xorm:"nvarchar(50) notnull index 'message_id'" json:"message_id"`
	ProductId        int64             `xorm:"bigint notnull index 'product_id'" json:"product_id"`
	ItemId           string            `xorm:"nvarchar(50) notnull index 'item_id'" json:"item_id"` // ebay系统的 ItemId，这个是产品的 item_id ??
	Sku              string            `xorm:"nvarchar(200) notnull index 'sku'" json:"sku"`
	CreationDate     core.NullLongTime `xorm:"datetime null 'creation_date'" json:"creation_date"`
	LastModifiedDate core.NullLongTime `xorm:"datetime null 'last_modified_date'" json:"last_modified_date"`
	MessageStatus    string            `xorm:"nvarchar(50) notnull 'message_status'" json:"message_status"` // Answered Unanswered
	DisplayToPublic  bool              `xorm:"bool notnull 'display_to_public'" json:"display_to_public"`
	MessageType      string            `xorm:"nvarchar(50) notnull 'message_type'" json:"message_type"`
	QuestionType     string            `xorm:"nvarchar(50) notnull 'question_type'" json:"question_type"`
	RecipientId      string            `xorm:"nvarchar(100) notnull 'recipient_id'" json:"recipient_id"`
	SenderEmail      string            `xorm:"nvarchar(200) notnull 'sender_email'" json:"sender_email"`
	SenderId         string            `xorm:"nvarchar(100) notnull 'sender_id'" json:"sender_id"`
	Subject          string            `xorm:"nvarchar(500) notnull 'subject'" json:"subject"`
	Body             string            `xorm:"nvarchar(5000) notnull 'body'" json:"body"`
	StatusId         int64             `xorm:"bigint notnull index 'status_id'" json:"status_id"`
	CreateTime       core.Time         `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime       core.Time         `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e MessageInfo) TableName() string {
	return TablePrefix + "message"
}

func (e MessageInfo) ModelDesc() string {
	return "verbose_name=消息,entity_json=message,route=message"
}

/*


Value	Description
Answered	The question has been answered at least once.
CustomCode	Reserved for future or internal use.
Unanswered	The question has not yet been answered.

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
