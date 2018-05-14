package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/core"
	"matrix/modules/erp/models/ppgl_models"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_paypal_service"
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

/*
	REFUND_STATUS_YXJ = 7100 // 已新建
	REFUND_STATUS_YSH = 7200 // 已审核
	REFUND_STATUS_YCG = 7300 // 已成功
	REFUND_STATUS_YSB = 7400 // 已失败
	REFUND_STATUS_YZF = 7500 // 已作废
	REFUND_STATUS_YWJ = 7600 // 已完结
 */

func (c ErpRefund) Approve() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	var db_session = c.DbSession

	var order_refund_id int64
	c.Params.Bind(&order_refund_id, "order_refund_id")

	var order_refund ppgl_models.OrderRefundInfo
	_, err := db_session.Id(order_refund_id).Get(&order_refund)
	core.HandleError(err)

	if order_refund.StatusId != models.REFUND_STATUS_YXJ {
		return c.RenderJson(core.JsonResult{
			Success: false,
			Message: "订单退款当前不是已新建状态！",
		})
	}

	var refund_status_ysh models.StatusInfo
	_, err = db_session.Id(models.REFUND_STATUS_YSH).Get(&refund_status_ysh)
	core.HandleError(err)

	order_refund.StatusId = models.REFUND_STATUS_YSH
	order_refund.StatusName = refund_status_ysh.Name

	_, err = db_session.Id(order_refund_id).Update(&order_refund)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "退款审核成功！",
	})
}

func (c ErpRefund) PaypalRefund() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	var db_session = c.DbSession
	var order_refund_id int64

	c.Params.Bind(&order_refund_id, "order_refund_id")

	var result = order_paypal_service.Order_refund(order_refund_id, db_session)

	return c.RenderJson(result)

}

func (c ErpRefund) Complete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	var db_session = c.DbSession

	var order_refund_id int64
	c.Params.Bind(&order_refund_id, "order_refund_id")

	var order_refund ppgl_models.OrderRefundInfo
	_, err := db_session.Id(order_refund_id).Get(&order_refund)
	core.HandleError(err)

	if order_refund.MethodId == ppgl_models.PAYMENT_METHOD_BANK && order_refund.StatusId != models.REFUND_STATUS_YSH {
		return c.RenderJson(core.JsonResult{
			Success: false,
			Message: "订单退款退款方式是银行转账且当前不是已审核状态！",
		})
	}

	if order_refund.MethodId == ppgl_models.PAYMENT_METHOD_PALPAY && order_refund.StatusId != models.REFUND_STATUS_YCG {
		return c.RenderJson(core.JsonResult{
			Success: false,
			Message: "订单退款退款方式是Paypal且当前不是已成功状态！",
		})
	}

	var refund_status_ywj models.StatusInfo
	_, err = db_session.Id(models.REFUND_STATUS_YWJ).Get(&refund_status_ywj)
	core.HandleError(err)

	order_refund.StatusId = models.REFUND_STATUS_YWJ
	order_refund.StatusName = refund_status_ywj.Name

	_, err = db_session.Id(order_refund_id).Update(&order_refund)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "退款完结成功！",
	})
}

func (c ErpRefund) Cancel() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	var db_session = c.DbSession

	var order_refund_id int64
	c.Params.Bind(&order_refund_id, "order_refund_id")

	var order_refund ppgl_models.OrderRefundInfo
	_, err := db_session.Id(order_refund_id).Get(&order_refund)
	core.HandleError(err)

	var refund_status_yzf models.StatusInfo
	_, err = db_session.Id(models.REFUND_STATUS_YZF).Get(&refund_status_yzf)
	core.HandleError(err)

	order_refund.StatusId = models.REFUND_STATUS_YZF
	order_refund.StatusName = refund_status_yzf.Name

	_, err = db_session.Id(order_refund_id).Update(&order_refund)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "退款审核成功！",
	})
}
