package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/ppgl_models"
	"matrix/modules/erp/models/xtpz_models"
	"fmt"
	"time"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpPaypalPayment struct {
	*revel.Controller
}

func (c ErpPaypalPayment) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_PAYPAL_PAYMENT, c.Session)

	var db_session = c.DbSession

	var paypal_account_list = make([]xtpz_models.PaypalAccountInfo, 0)
	core.HandleError(db_session.Find(&paypal_account_list))

	c.RenderArgs["account_list"] = paypal_account_list

	return c.RenderTemplate("erp/ppgl/paypal_payment/paypal_payment_index.html")
}

func (c ErpPaypalPayment) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_PAYPAL_PAYMENT, c.Session)

	db_session := c.DbSession

	var account_id_query int64
	c.Params.Bind(&account_id_query, "account_id_query")

	var query_begin_time string
	c.Params.Bind(&query_begin_time, "begin_time_query")
	query_begin_date, query_begin_date_err := core.Parse_time(query_begin_time)
	revel.TRACE.Println("query_begin_time", query_begin_time)

	var query_end_time string
	c.Params.Bind(&query_end_time, "end_time_query")
	query_end_date, query_end_date_err := core.Parse_time(query_end_time)

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	if account_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("account_id = %d", account_id_query))
	}

	if query_begin_date_err == nil {
		filter = core.SqlAnd(filter, fmt.Sprintf("pay_time >= '%04d-%02d-%02d'", query_begin_date.Year(), query_begin_date.Month(), query_begin_date.Day()))
	}

	if query_end_date_err == nil {
		query_end_date = query_end_date.Add(time.Hour * 24)
		filter = core.SqlAnd(filter, fmt.Sprintf("pay_time < '%04d-%02d-%02d'", query_end_date.Year(), query_end_date.Month(), query_end_date.Day()))
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
