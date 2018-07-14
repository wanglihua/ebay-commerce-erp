package service

import (
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/xtpz_models"
)

type MessageOrderItem struct {
	OrderItem   ddlc_models.OrderItemInfo   `xorm:"extends"`
	Order       ddlc_models.OrderInfo       `xorm:"extends"`
	Status      models.StatusInfo           `xorm:"extends"`
	ShopAccount xtpz_models.ShopAccountInfo `xorm:"extends"`
	OrdNum      string                      `xorm:"-"`
	Price       string                      `xorm:"-"`
	Amount      string                      `xorm:"-"`
	PaidTime    string                      `xorm:"-"`
	ShippedTime string                      `xorm:"-"`
}
