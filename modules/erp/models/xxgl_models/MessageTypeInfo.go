package xxgl_models

import "github.com/go-xorm/xorm"

const (
	MESSAGE_TYPE_BEFORE_SALE = 1001
	MESSAGE_TYPE_AFTER_SALE  = 1002
	MESSAGE_TYPE_FROM_EBAY   = 1003
	MESSAGE_TYPE_OFFER       = 1004
	MESSAGE_TYPE_REQUEST     = 1005
	MESSAGE_TYPE_RETURN      = 1006
	MESSAGE_TYPE_PLATFORM    = 1008
	MESSAGE_TYPE_SALE        = 2001 // 售前售后消息
)

type MessageTypeInfo struct {
	Id          int64  `xorm:"bigint notnull pk  'id'" json:"id"`
	BigTypeId   int64  `xorm:"bigint notnull index 'big_type_id'" json:"big_type_id"`
	BigTypeName string `xorm:"nvarchar(100) notnull 'big_type_name'" json:"big_type_name"`
	TypeName    string `xorm:"nvarchar(100) notnull 'type_name'" json:"type_name"`
	TypeDesc    string `xorm:"nvarchar(300) notnull 'type_desc'" json:"type_desc"`
}

func (e MessageTypeInfo) TableName() string {
	return "erp_message_type"
}

func (e MessageTypeInfo) ModelDesc() string {
	return "verbose_name=消息类型,entity_json=message_type,route=message/type"
}

func (e MessageTypeInfo) InitData(db_session *xorm.Session) {
	var message_type = new(MessageTypeInfo)
	message_type.Id = MESSAGE_TYPE_BEFORE_SALE
	message_type.BigTypeId = MESSAGE_TYPE_BEFORE_SALE
	message_type.BigTypeName = "售前消息"
	message_type.TypeName = "售前消息"
	message_type.TypeDesc = "售前消息"
	db_session.Insert(message_type)

	message_type = new(MessageTypeInfo)
	message_type.Id = MESSAGE_TYPE_AFTER_SALE
	message_type.BigTypeId = MESSAGE_TYPE_AFTER_SALE
	message_type.BigTypeName = "售后消息"
	message_type.TypeName = "售后消息"
	message_type.TypeDesc = "售后消息"
	db_session.Insert(message_type)

	message_type = new(MessageTypeInfo)
	message_type.Id = MESSAGE_TYPE_FROM_EBAY
	message_type.BigTypeId = MESSAGE_TYPE_PLATFORM
	message_type.BigTypeName = "平台消息"
	message_type.TypeName = "eBay渠道"
	message_type.TypeDesc = "eBay渠道"
	db_session.Insert(message_type)

	message_type = new(MessageTypeInfo)
	message_type.Id = MESSAGE_TYPE_OFFER
	message_type.BigTypeId = MESSAGE_TYPE_PLATFORM
	message_type.BigTypeName = "平台消息"
	message_type.TypeName = "offer"
	message_type.TypeDesc = "offer"
	db_session.Insert(message_type)

	message_type = new(MessageTypeInfo)
	message_type.Id = MESSAGE_TYPE_REQUEST
	message_type.BigTypeId = MESSAGE_TYPE_PLATFORM
	message_type.BigTypeName = "平台消息"
	message_type.TypeName = "request"
	message_type.TypeDesc = "request"
	db_session.Insert(message_type)

	message_type = new(MessageTypeInfo)
	message_type.Id = MESSAGE_TYPE_RETURN
	message_type.BigTypeId = MESSAGE_TYPE_PLATFORM
	message_type.BigTypeName = "平台消息"
	message_type.TypeName = "return"
	message_type.TypeDesc = "return"
	db_session.Insert(message_type)

	message_type = new(MessageTypeInfo)
	message_type.Id = MESSAGE_TYPE_SALE
	message_type.BigTypeId = MESSAGE_TYPE_SALE
	message_type.BigTypeName = "售前售后消息"
	message_type.TypeName = "售前售后消息"
	message_type.TypeDesc = "售前售后消息"
	db_session.Insert(message_type)
}
