package ppgl_models

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

type OrderRefundReasonInfo struct {
	Id          int64     `xorm:"bigint notnull pk 'id'" json:"id"`
	Name        string    `xorm:"nvarchar(200) notnull index 'name'" json:"name"`
	Description string    `xorm:"nvarchar(500) 'description'" json:"description"`
	CreateTime  core.Time `xorm:"created notnull 'create_time'" json:"create_time"`
	UpdateTime  core.Time `xorm:"updated notnull 'update_time'" json:"update_time"`
}

func (e OrderRefundReasonInfo) TableName() string {
	return "erp_order_refund_reason"
}

func (e OrderRefundReasonInfo) ModelDesc() string {
	return "verbose_name=订单退款原因,entity_json=order_refund_reason,route=order/refund/reason"
}

func (e OrderRefundReasonInfo) InitData(db_session *xorm.Session) {

	var reason = new(OrderRefundReasonInfo)
	reason.Id = 1
	reason.Name = "质量问题"
	reason.Description = "质量问题"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 2
	reason.Name = "客户原因"
	reason.Description = "客户原因"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 3
	reason.Name = "订单漏发"
	reason.Description = "订单漏发"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 4
	reason.Name = "发错货"
	reason.Description = "发错货"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 5
	reason.Name = "改评价"
	reason.Description = "改评价"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 6
	reason.Name = "物流运输损坏"
	reason.Description = "物流运输损坏"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 7
	reason.Name = "物流丢包"
	reason.Description = "物流丢包"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 8
	reason.Name = "物流运送时间过长"
	reason.Description = "物流运送时间过长"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 9
	reason.Name = "包裹退回"
	reason.Description = "包裹退回"
	db_session.Insert(reason)

	reason = new(OrderRefundReasonInfo)
	reason.Id = 10
	reason.Name = "其他"
	reason.Description = "其他"
	db_session.Insert(reason)
}
