package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
)

type ErpOrderItem struct {
	*revel.Controller
}

func (c ErpOrderItem) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var item_id int64
	c.Params.Bind(&item_id, "item_id")

	var order_id int64
	c.Params.Bind(&order_id, "order_id")

	var db_session = c.DbSession

	var order ddlc_models.OrderInfo
	_, err := db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	var order_item ddlc_models.OrderItemInfo
	_, err = db_session.Id(item_id).Get(&order_item)
	core.HandleError(err)

	var order_item_form OrderItemForm
	order_item_form.ItemId = order_item.Id
	order_item_form.OrderId = order_id
	order_item_form.Sku = order_item.ItemSku
	order_item_form.Qty = order_item.Quantity
	order_item_form.Price = order_item.TransactionPrice
	order_item_form.IsOrderSgcj = order.TypeId == ddpz_models.ORDER_TYPE_SGCJ

	return c.RenderJson(map[string]interface{}{
		"success":           true,
		"order_item_detail": order_item_form,
	})
}

type OrderItemForm struct {
	ItemId      int64   `json:"item_id"` //这是 order_item 表的主键
	OrderId     int64   `json:"order_id"`
	Sku         string  `json:"sku"`
	Qty         int     `json:"qty"`
	Price       float64 `json:"price"`
	IsOrderSgcj bool    `json:"is_order_sgcj"` // 是否是手工创建的订单
}

func (c ErpOrderItem) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var db_session = c.DbSession

	var order_item_form OrderItemForm
	err := json.Unmarshal(c.GetRequestBody(), &order_item_form)
	core.HandleError(err)

	var order_id = order_item_form.OrderId

	// 判断订单当前状态是不是可以修改
	var order ddlc_models.OrderInfo
	_, err = db_session.Id(order_id).Get(&order)
	core.HandleError(err)

	if order_service.Can_order_modify(&order) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单当前状态不能修改订单明细！",
		})
	}

	c.Validation.Required(order_item_form.Sku).Message("SKU 不能为空！")
	c.Validation.Required(order_item_form.Qty).Message("数量 不能为空！")

	if c.Validation.HasErrors() {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	// 获得当前订单的仓库
	var order_logistics ddlc_models.OrderLogisticsInfo
	_, err = db_session.Where("order_id = ?", order.Id).Get(&order_logistics)
	core.HandleError(err)

	var stock_id = order_logistics.StockId

	//获取 product_id
	var product kcgl_models.ProductInfo
	_, err = db_session.Where(fmt.Sprintf("sku = '%s'", order_item_form.Sku)).Get(&product)
	core.HandleError(err)

	if product.Id == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的SKU不存在!",
		})
	}

	var order_item ddlc_models.OrderItemInfo
	if order_item_form.ItemId != 0 {
		_, err = db_session.Id(order_item_form.ItemId).Get(&order_item)
		core.HandleError(err)
	}

	var exchange_rate = order.ExchangeRate

	order_item.OrderId = order_item_form.OrderId
	order_item.TransactionPriceCurrency = "AUD"
	order_item.ItemProductId = product.Id
	order_item.ItemSku = order_item_form.Sku
	order_item.Quantity = order_item_form.Qty
	order_item.TransactionPrice = order_item_form.Price
	order_item.TransactionPriceRmb = order_item_form.Price * exchange_rate

	//获取成本价 AUD
	if stock_id == 0 {
		order_item.CostPrice = inventory_service.GetProductCnStockCostPrice(product.Id, db_session) / order.ExchangeRate
	} else {
		order_item.CostPrice = inventory_service.GetProductStockCostPrice(stock_id, product.Id, db_session) / order.ExchangeRate
	}

	if order_item_form.ItemId == 0 {
		_, err = db_session.Insert(&order_item)
		core.HandleError(err)
	} else {
		var update_cols = []string{"item_product_id", "item_sku", "quantity", "cost_price", "transaction_price_currency"}
		if order.TypeId == ddpz_models.ORDER_TYPE_SGCJ {
			update_cols = append(update_cols, "transaction_price", "transaction_price_rmb")
		}

		_, err = db_session.Id(order_item_form.ItemId).Cols(update_cols...).Update(&order_item)
		core.HandleError(err)
	}

	// 如果是 手工创建的的订单的，需要更新收款金额
	if order.TypeId == ddpz_models.ORDER_TYPE_SGCJ {
		type order_item_total_amount_t struct {
			TotalAmount float64 `xorm:"'total_amount'"`
		}

		/*
		SELECT sum(oi.transaction_price * oi.quantity) AS total_amount
		FROM erp_order_item oi
		WHERE oi.order_id = 0;
		*/
		var sql_item_total_amount = `
SELECT sum(oi.transaction_price * oi.quantity) AS total_amount
FROM erp_order_item oi
WHERE oi.order_id = %d;
`
		sql_item_total_amount = fmt.Sprintf(sql_item_total_amount, order_id)
		var order_item_total_amount order_item_total_amount_t
		_, err := db_session.Sql(sql_item_total_amount).Get(&order_item_total_amount)
		core.HandleError(err)

		order.AmountPaid = order_item_total_amount.TotalAmount
		order.AmountPaidRmb = order_item_total_amount.TotalAmount * order.ExchangeRate
		order.AmountPaidCurrency = "AUD"
		order.TotalAmount = order_item_total_amount.TotalAmount
		order.TotalAmountRmb = order_item_total_amount.TotalAmount * order.ExchangeRate
		order.TotalAmountCurrency = "AUD"
		order.AdjustmentAmount = order_item_total_amount.TotalAmount
		order.AdjustmentAmountRmb = order_item_total_amount.TotalAmount * order.ExchangeRate
		order.AdjustmentAmountCurrency = "AUD"

		var update_cols = []string{
			"amount_paid", "amount_paid_rmb", "amount_paid_currency",
			"total_amount", "total_amount_rmb", "total_amount_currency",
			"adjustment_amount", "adjustment_amount_rmb", "adjustment_amount_currency",
		}
		_, err = db_session.Id(order_id).Cols(update_cols...).Update(&order)
		core.HandleError(err)
	}

	order_service.Update_order_logistics_cost_amount_and_order_profit_and_order_item_price_profit(order, order_logistics, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单明细保存成功！",
	})
}

func (c ErpOrderItem) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var item_id int64
	c.Params.Bind(&item_id, "item_id")

	var db_session = c.DbSession

	//判断订单当前状态是不是可以修改
	var order_item ddlc_models.OrderItemInfo
	_, err := db_session.Id(item_id).Get(&order_item)
	core.HandleError(err)

	var order ddlc_models.OrderInfo
	_, err = db_session.Id(order_item.OrderId).Get(&order)
	core.HandleError(err)

	if order_service.Can_order_modify(&order) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单当前状态不能删除订单明细！",
		})
	}

	_, err = db_session.Id(item_id).Delete(new(ddlc_models.OrderItemInfo))
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单明细删除成功！",
	})
}
