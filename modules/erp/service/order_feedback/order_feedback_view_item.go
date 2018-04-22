package order_feedback

import (
	"matrix/modules/erp/models/ddlc_models"
)

type OrderFeedbackViewItem struct {
	OrderFeedback ddlc_models.OrderFeedbackInfo `xorm:"extends" json:"of"`
	OrderItem     ddlc_models.OrderItemInfo     `xorm:"extends" json:"oi"`
}
