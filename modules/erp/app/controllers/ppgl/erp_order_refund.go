package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/core"
	"matrix/modules/erp/models/ppgl_models"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/models/sppz_models"
	"fmt"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/kcgl_models"
	"math"
	"matrix/modules/erp/models"
	"time"
)

type ErpOrderRefund struct {
	*revel.Controller
}

func (c ErpOrderRefund) Refund() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPROVE, c.Session)

	/*
	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "refund.order_id")

	var amount float64
	c.Params.Bind(&amount, "refund.amount")

	var memo string
	c.Params.Bind(&memo, "refund.memo")

	var result = order_paypal_service.Order_refund(order_id, amount, memo, db_session)

	return c.RenderJson(result)
	*/

	return c.RenderJson(core.JsonResult{
		Success: false,
		Message: "先审核再退款！",
	})
}

func (c ErpOrderRefund) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var db_session = c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var refund_list = make([]ppgl_models.OrderRefundInfo, 0)
	err := db_session.Where("order_id = ?", order_id).Find(&refund_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"refund_list": refund_list,
	})
}

func (c ErpOrderRefund) ApplyShow(order_id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPLY, c.Session)
	var db_session = c.DbSession

	c.RenderArgs["order_id"] = order_id

	var order_refund_method_list = make([]ppgl_models.OrderRefundMethodInfo, 0)
	core.HandleError(db_session.Find(&order_refund_method_list))

	var order_refund_type_list = make([]ppgl_models.OrderRefundTypeInfo, 0)
	core.HandleError(db_session.Find(&order_refund_type_list))

	var order_refund_reason_list = make([]ppgl_models.OrderRefundReasonInfo, 0)
	core.HandleError(db_session.Find(&order_refund_reason_list))

	c.RenderArgs["order_refund_method_list"] = order_refund_method_list
	c.RenderArgs["order_refund_type_list"] = order_refund_type_list
	c.RenderArgs["order_refund_reason_list"] = order_refund_reason_list

	return c.RenderTemplate("erp/ppgl/order_refund/order_refund_apply.tpl")
}

func (c ErpOrderRefund) ApplyShowData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPLY, c.Session)

	db_session := c.DbSession

	var order_id int64
	c.Params.Bind(&order_id, "id")

	revel.TRACE.Println("order_id:", order_id)

	var order ddlc_models.OrderInfo
	has_order, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)
	if has_order == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的订单不存在！",
		})
	}

	var status models.StatusInfo
	_, err = db_session.Id(order.StatusId).Get(&status)
	core.HandleError(err)

	var shop xtpz_models.ShopAccountInfo
	_, err = db_session.Id(order.ShopId).Get(&shop)
	core.HandleError(err)

	var order_type ddpz_models.OrderTypeInfo
	_, err = db_session.Id(order.TypeId).Get(&order_type)
	core.HandleError(err)

	var order_hold_type ddpz_models.OrderHoldTypeInfo
	_, err = db_session.Id(order.HoldTypeId).Get(&order_hold_type)
	core.HandleError(err)

	type order_item_view_t struct {
		OrderItem         ddlc_models.OrderItemInfo        `xorm:"extends" json:"oi"`
		Product           kcgl_models.ProductInfo          `xorm:"extends" json:"p"`
		ProductSotoreItem sppz_models.ProductStoreItemInfo `xorm:"extends" json:"psi"`
	}

	sql := `
SELECT oi.*, p.*, psi.*
FROM erp_order_item oi
    INNER JOIN erp_product p ON p.id = oi.item_product_id
    LEFT JOIN erp_product_store_item psi ON psi.id = oi.product_store_item_id
WHERE oi.order_id = %d;`
	sql = fmt.Sprintf(sql, order_id)

	order_item_list := make([]order_item_view_t, 0)
	err = db_session.Sql(sql).Find(&order_item_list)
	core.HandleError(err)

	var order_logistics ddlc_models.OrderLogisticsInfo
	_, err = db_session.Where(fmt.Sprintf("order_id = %d", order_id)).Get(&order_logistics)
	core.HandleError(err)

	var logistics_mode ddpz_models.LogisticsModeInfo
	if !(order.StatusId == models.ORDER_STATUS_HOLD && order.HoldTypeId.Valid == true &&
		(order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_ENVELOPE ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_LOGISTICS_NO_MATCH ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_INTERNATIONAL ||
			order.HoldTypeId.Int64 == ddpz_models.ORDER_HOLD_TYPE_PRODUCT_INFO_MISS)) {
		//如果是Hold状态且是信封Hold的就不显示匹配的物流，因为这些信息可以是前一次匹配时遗留的，后来指定了 Hold
		_, err = db_session.Id(order_logistics.ModeId).Get(&logistics_mode)
		core.HandleError(err)
	}

	var order_cost_amount = order_logistics.CostAmount * order.ExchangeRate

	// 物流费用

	var logistics_fee float64 = (order_logistics.TotalFee + order_logistics.PickFee) * order.ExchangeRate

	if order_logistics.ActualFee != 0 {
		logistics_fee = (order_logistics.ActualFee + order_logistics.PickFee) * order.ExchangeRate
	}

	// gst 消费税
	var order_gst float64 = order_service.Get_order_gst_fee(&order)

	// 重新按实际运费计算毛利
	var order_profit_au = order_service.CalculateOrderProfit(&order, &order_logistics)
	var order_profit = order_profit_au * order.ExchangeRate

	if math.Abs(order.ProfitAu-order_profit_au) > core.FLOAT_MIN_DIFF { // 之前许多订单没有按实际运费计算毛利的，这个机会计算一下
		order.ProfitAu = order_profit_au
		order.ProfitRmb = order_profit

		_, err := db_session.Id(order.Id).Cols("profit_au", "profit_rmb").Update(&order)
		core.HandleError(err)

		// 平均分配利润到订单项
		var order_item_list = make([]ddlc_models.OrderItemInfo, 0)
		core.HandleError(db_session.Where("order_id = ?", order.Id).Find(&order_item_list))

		var order_cost_amount = order_service.GetOrderCostAmountSum(order.Id, db_session)

		for _, oi := range order_item_list {
			// 利润分配
			oi.GrossProfit = oi.CostPrice * float64(oi.Quantity) * order.ProfitAu / order_cost_amount
			_, err := db_session.Id(oi.Id).Cols("gross_profit").Update(&oi)
			core.HandleError(err)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"order":             order,
		"status":            status,
		"shop":              shop,
		"order_type":        order_type,
		"order_hold_type":   order_hold_type,
		"order_item_list":   order_item_list,
		"order_logistics":   order_logistics,
		"order_cost_amount": order_cost_amount, // rmb
		"logistics_fee":     logistics_fee,     // rmb
		"order_gst":         order_gst,         // aud
		"order_profit":      order.ProfitRmb,   // rmb
	})
}

func (c ErpOrderRefund) Apply() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_ORDER_REFUND_APPLY, c.Session)

	var web_session = c.Session
	var db_session = c.DbSession

	var refund_order_id int64
	c.Params.Bind(&refund_order_id, "refund_order_id")

	var refund_product_id int64
	c.Params.Bind(&refund_product_id, "refund_product_id")

	var refund_type_id int64
	c.Params.Bind(&refund_type_id, "refund_type_id")

	var refund_reason_id int64
	c.Params.Bind(&refund_reason_id, "refund_reason_id")

	var refund_method_id int64
	c.Params.Bind(&refund_method_id, "refund_method_id")

	var refund_amount float64
	c.Params.Bind(&refund_amount, "refund_amount")

	var refund_description string
	c.Params.Bind(&refund_description, "refund_description")

	var order_refund_in_db ppgl_models.OrderRefundInfo
	has, err := db_session.Where("order_id = ? AND status_id <> ?", refund_order_id, models.REFUND_STATUS_YZF).Get(&order_refund_in_db)
	core.HandleError(err)

	if has {
		return c.RenderJson(core.JsonResult{
			Success: false,
			Message: "订单已存在退款信息！",
		})
	}

	var refund_order ddlc_models.OrderInfo
	_, err = db_session.Id(refund_order_id).Get(&refund_order)
	core.HandleError(err)

	// 退款金额不能大于付款金额
	if refund_amount > refund_order.AmountPaid {
		return c.RenderJson(core.JsonResult{
			Success: false,
			Message: "退款金额不能大于订单付款金额！",
		})
	}

	var shop xtpz_models.ShopAccountInfo
	_, err = db_session.Id(refund_order.ShopId).Get(&shop)
	core.HandleError(err)

	// Payment
	var sql_get_order_payment_transaction = `
SELECT oet.*
FROM erp_order_external_transaction oet
    INNER JOIN erp_paypal_payment pp ON pp.id = oet.payment_id
WHERE pp.type = 'Payment' AND oet.order_id = %d
LIMIT 1;`
	sql_get_order_payment_transaction = fmt.Sprintf(sql_get_order_payment_transaction, refund_order_id)

	var order_external_transaction ddlc_models.OrderExternalTransactionInfo
	has, err = db_session.Sql(sql_get_order_payment_transaction).Get(&order_external_transaction)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(core.JsonResult{
			Success: false,
			Message: "未找到订单关联付款信息！",
		})
	}

	var refund_product kcgl_models.ProductInfo
	_, err = db_session.Id(refund_product_id).Get(&refund_product)
	core.HandleError(err)

	var refund_type ppgl_models.OrderRefundTypeInfo
	_, err = db_session.Id(refund_type_id).Get(&refund_type)
	core.HandleError(err)

	var refund_reason ppgl_models.OrderRefundReasonInfo
	_, err = db_session.Id(refund_reason_id).Get(&refund_reason)
	core.HandleError(err)

	var refund_method ppgl_models.OrderRefundMethodInfo
	_, err = db_session.Id(refund_method_id).Get(&refund_method)
	core.HandleError(err)

	var refund_status_yxj models.StatusInfo
	_, err = db_session.Id(models.REFUND_STATUS_YXJ).Get(&refund_status_yxj)
	core.HandleError(err)

	var order_refund ppgl_models.OrderRefundInfo
	order_refund.ShopId = shop.Id
	order_refund.ShopName = shop.Name
	order_refund.OrderId = refund_order.Id
	order_refund.OrderCode = refund_order.Code
	order_refund.OrderAmount = refund_order.AmountPaid
	order_refund.PayTransId = order_external_transaction.TransactionId
	order_refund.RefundAmount = refund_amount
	order_refund.TypeId = refund_type.Id
	order_refund.TypeName = refund_type.Name
	order_refund.ReasonId = refund_reason.Id
	order_refund.Reason = refund_reason.Name
	order_refund.MethodId = refund_method.Id
	order_refund.MethodName = refund_method.Name
	order_refund.Description = refund_description
	order_refund.ProductId = refund_product.Id
	order_refund.ProductCode = refund_product.Sku
	order_refund.ProductName = refund_product.CnName
	order_refund.SalesUserId = refund_product.SaleUserId
	order_refund.SalesUserName = refund_product.SaleUser

	var login_user = core.GetLoginUser(web_session)
	order_refund.ApplyUserId = login_user.UserId
	order_refund.ApplyUserName = fmt.Sprintf("%s(%s)", login_user.LoginName, login_user.NickName)

	order_refund.ApplyTime = core.LongTime(time.Now())
	order_refund.StatusId = refund_status_yxj.Id
	order_refund.StatusName = refund_status_yxj.Name

	_, err = db_session.Insert(&order_refund)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "退款申请成功！",
	})
}
