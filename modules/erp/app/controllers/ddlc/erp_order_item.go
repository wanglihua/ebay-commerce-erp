package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/inventory_service"
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

	var order_item models.OrderItemInfo
	_, err := db_session.Id(item_id).Get(&order_item)
	core.HandleError(err)

	var order_item_form OrderItemForm
	order_item_form.ItemId = order_item.Id
	order_item_form.OrderId = order_id
	order_item_form.Sku = order_item.ItemSku
	order_item_form.Qty = order_item.Quantity

	return c.RenderJson(map[string]interface{}{
		"success":           true,
		"order_item_detail": order_item_form,
	})
}

type OrderItemForm struct {
	ItemId  int64  `json:"item_id"` //这是 order_item 表的主键
	OrderId int64  `json:"order_id"`
	Sku     string `json:"sku"`
	Qty     int    `json:"qty"`
}

func (c ErpOrderItem) Save() revel.Result {

	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_VIEW_DETAIL, c.Session)

	var db_session = c.DbSession

	var order_item_form OrderItemForm
	err := json.Unmarshal(c.GetRequestBody(), &order_item_form)
	core.HandleError(err)

	// 判断订单当前状态是不是可以修改
	var order models.OrderInfo
	_, err = db_session.Id(order_item_form.OrderId).Get(&order)
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
	var order_logistics models.OrderLogisticsInfo
	_, err = db_session.Where("order_id = ?", order.Id).Get(&order_logistics)
	core.HandleError(err)

	var stock_id = order_logistics.StockId

	//获取 product_id
	var product models.ProductInfo
	_, err = db_session.Where(fmt.Sprintf("sku = '%s'", order_item_form.Sku)).Get(&product)
	core.HandleError(err)

	if product.Id == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的SKU不存在!",
		})
	}

	var order_item models.OrderItemInfo
	if order_item_form.ItemId != 0 {
		_, err = db_session.Id(order_item_form.ItemId).Get(&order_item)
		core.HandleError(err)
	}

	order_item.OrderId = order_item_form.OrderId
	order_item.TransactionPriceCurrency = "AUD"
	order_item.ItemProductId = product.Id
	order_item.ItemSku = order_item_form.Sku
	order_item.Quantity = order_item_form.Qty

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
		_, err = db_session.Id(order_item_form.ItemId).Cols(update_cols...).Update(&order_item)
		core.HandleError(err)
	}

	// 计算订单成本 AUD
	var order_cost_amount = order_service.GetOrderCostAmountSum(order.Id, db_session)

	revel.TRACE.Println("order_cost_amount", order_cost_amount)

	order_logistics.CostAmount = order_cost_amount
	_, err = db_session.Id(order_logistics.Id).Cols("cost_amount").Update(&order_logistics)
	core.HandleError(err)

	// 重新计算订单毛利
	var order_profit_au = order_service.CalculateOrderProfit(&order, &order_logistics)
	revel.TRACE.Println("order_profit_au", order_profit_au)
	var order_profit = order_profit_au * order.ExchangeRate

	order.ProfitAu = order_profit_au
	order.ProfitRmb = order_profit

	_, err = db_session.Id(order.Id).Cols("profit_au", "profit_rmb").Update(&order)
	core.HandleError(err)

	var order_item_list = make([]models.OrderItemInfo, 0)
	core.HandleError(db_session.Where("order_id = ?", order.Id).Find(&order_item_list))

	/*
	order_import_same_code_list[index].price_au = ((order_import.cost_price * (float64(order_import.Spsl)) / cost_price_total) * order_import.Skje) / float64(order_import.Spsl)
	order_import_same_code_list[index].price_rmb = order_import_same_code_list[index].price_au * aud_exchange_rate_setting
	order_import_same_code_list[index].gross_profit = (order_import_same_code_list[index].price_rmb - order_import_same_code_list[index].cost_price) * float64(order_import.Spsl)
	 */

	for _, oi := range order_item_list {
		// 价格分配
		oi.TransactionPrice = oi.CostPrice * float64(oi.Quantity) * order.AmountPaid / order_cost_amount / float64(oi.Quantity)
		oi.TransactionPriceCurrency = "AUD"
		oi.TransactionPriceRmb = oi.TransactionPrice * order.ExchangeRate

		// 利润分配
		oi.GrossProfit = oi.CostPrice * float64(oi.Quantity) * order.ProfitAu / order_cost_amount
		revel.TRACE.Println("oi.CostPrice * float64(oi.Quantity)", oi.CostPrice*float64(oi.Quantity))
		revel.TRACE.Println("order.ProfitAu", order.ProfitAu)
		revel.TRACE.Println("order_cost_amount", order_cost_amount)
		revel.TRACE.Println("oi.GrossProfit", oi.GrossProfit)

		_, err := db_session.Id(oi.Id).Cols("transaction_price", "transaction_price_currency", "transaction_price_rmb", "gross_profit").Update(&oi)
		core.HandleError(err)
	}

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
	var order_item models.OrderItemInfo
	_, err := db_session.Id(item_id).Get(&order_item)
	core.HandleError(err)

	var order models.OrderInfo
	_, err = db_session.Id(order_item.OrderId).Get(&order)
	core.HandleError(err)

	if order_service.Can_order_modify(&order) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "订单当前状态不能删除订单明细！",
		})
	}

	_, err = db_session.Id(item_id).Delete(new(models.OrderItemInfo))
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "订单明细删除成功！",
	})
}
