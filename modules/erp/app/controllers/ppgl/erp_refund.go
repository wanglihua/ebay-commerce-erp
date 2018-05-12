package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/core"
	"matrix/modules/erp/models/ppgl_models"
	"matrix/modules/erp/models/xtpz_models"
)

type ErpRefund struct {
	*revel.Controller
}

func (c ErpRefund) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	var db_session = c.DbSession

	var paypal_account_list = make([]xtpz_models.PaypalAccountInfo, 0)
	core.HandleError(db_session.Find(&paypal_account_list))

	c.RenderArgs["account_list"] = paypal_account_list

	return c.RenderTemplate("erp/ppgl/refund/refund_index.html")
}

func (c ErpRefund) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	var db_session = c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var query = db_session.Where(filter)

	dataQuery := *query
	if order != "" {
		dataQuery = *dataQuery.OrderBy(order)
	} else {
		dataQuery = *dataQuery.Asc("id")
	}

	refundList := make([]ppgl_models.OrderRefundInfo, 0, limit)
	err := dataQuery.Limit(limit, offset).Find(&refundList)
	core.HandleError(err)

	countQuery := *query
	count, err := countQuery.Count(new(ppgl_models.OrderRefundInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  refundList,
		Total: count,
	})
}
