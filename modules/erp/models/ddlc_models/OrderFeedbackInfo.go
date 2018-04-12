package ddlc_models

import (
	"matrix/core"
)

type OrderFeedbackInfo struct {
	Id          int64     `xorm:"bigint notnull pk autoincr 'id'" json:"id"`
	OrderId     int64     `xorm:"bigint notnull index 'order_id'" json:"order_id"`
	OrderItemId int64     `xorm:"bigint notnull index 'order_item_id'" json:"order_item_id"`
	CommentType int64     `xorm:"bigint notnull 'comment_type'" json:"comment_type"`
	CommentText string    `xorm:"varchar(200) notnull 'comment_text'" json:"comment_text"`
	CommentTime core.Time `xorm:"datetime notnull 'comment_time'" json:"comment_time"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderFeedbackInfo) TableName() string {
	return "erp_order_feedback"
}

func (e OrderFeedbackInfo) ModelDesc() string {
	return "verbose_name=订单反馈,entity_json=order_feedback,route=order/feedback"
}
