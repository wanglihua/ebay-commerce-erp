package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/ppgl_models"
)

type ErpPaypalPayment struct {
	*revel.Controller
}

func (c ErpPaypalPayment) Index() revel.Result {
	return c.RenderTemplate("erp/ppgl/paypal_payment/paypal_payment_index.html")
}

func (c ErpPaypalPayment) ListData() revel.Result {
	
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	query := db_session.
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	paypal_payment_list := make([]ppgl_models.PaypalPaymentInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&paypal_payment_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ppgl_models.PaypalPaymentInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  paypal_payment_list,
		Total: count,
	})
}
