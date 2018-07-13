package service

import (
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models"
)

type MessageOrderItem struct {
	OrderItem ddlc_models.OrderItemInfo `xorm:"extends" json:"oi"`
	Order     ddlc_models.OrderInfo     `xorm:"extends" json:"o"`
	Status    models.StatusInfo         `xorm:"extends" json:"s"`
}
