package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/ppgl_models"
	"matrix/modules/erp/models/xtpz_models"
	"fmt"
)

type ErpPaypalPayment struct {
	*revel.Controller
}

func (c ErpPaypalPayment) Index() revel.Result {
	var db_session = c.DbSession

	var paypal_account_list = make([]xtpz_models.PaypalAccountInfo, 0)
	core.HandleError(db_session.Find(&paypal_account_list))

	c.RenderArgs["account_list"] = paypal_account_list

	return c.RenderTemplate("erp/ppgl/paypal_payment/paypal_payment_index.html")
}

func (c ErpPaypalPayment) ListData() revel.Result {

	db_session := c.DbSession

	var account_id int64
	c.Params.Bind(&account_id, "account_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if account_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("account_id = %d", account_id))
	}

	query := db_session.
		Where(filter)

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
