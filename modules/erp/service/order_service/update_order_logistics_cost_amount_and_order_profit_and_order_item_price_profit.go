package order_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
)

func Update_order_logistics_cost_amount_and_order_profit_and_order_item_price_profit(order ddlc_models.OrderInfo, order_logistics ddlc_models.OrderLogisticsInfo, db_session *xorm.Session) {
	// 计算订单成本 AUD
	var order_cost_amount = GetOrderCostAmountSum(order.Id, db_session)

	order_logistics.CostAmount = order_cost_amount
	_, err := db_session.Id(order_logistics.Id).Cols("cost_amount").Update(&order_logistics)
	core.HandleError(err)

	// 重新计算订单毛利
	var order_profit_au = CalculateOrderProfit(&order, &order_logistics)
	var order_profit = order_profit_au * order.ExchangeRate

	order.ProfitAu = order_profit_au
	order.ProfitRmb = order_profit

	_, err = db_session.Id(order.Id).Cols("profit_au", "profit_rmb").Update(&order)
	core.HandleError(err)

	var order_item_list = make([]ddlc_models.OrderItemInfo, 0)
	core.HandleError(db_session.Where("order_id = ?", order.Id).Find(&order_item_list))

	/*
	order_import_same_code_list[index].price_au = ((order_import.cost_price * (float64(order_import.Spsl)) / cost_price_total) * order_import.Skje) / float64(order_import.Spsl)
	order_import_same_code_list[index].price_rmb = order_import_same_code_list[index].price_au * aud_exchange_rate_setting
	order_import_same_code_list[index].gross_profit = (order_import_same_code_list[index].price_rmb - order_import_same_code_list[index].cost_price) * float64(order_import.Spsl)
	 */

	var order_item_update_cols = []string{"gross_profit"}

	// 手工创建的订单不自动计算价格和分配，其它的订单计算
	if order.TypeId != ddpz_models.ORDER_TYPE_SGCJ {
		order_item_update_cols = append(order_item_update_cols, "transaction_price", "transaction_price_currency", "transaction_price_rmb")
	}

	for _, oi := range order_item_list {
		// 手工创建的订单不自动计算价格和分配，其它的订单计算
		if order.TypeId != ddpz_models.ORDER_TYPE_SGCJ {
			// 价格分配
			oi.TransactionPrice = oi.CostPrice * float64(oi.Quantity) * order.AmountPaid / order_cost_amount / float64(oi.Quantity)
			oi.TransactionPriceCurrency = "AUD"
			oi.TransactionPriceRmb = oi.TransactionPrice * order.ExchangeRate
		}

		// 利润分配
		oi.GrossProfit = oi.CostPrice * float64(oi.Quantity) * order.ProfitAu / order_cost_amount

		_, err := db_session.Id(oi.Id).Cols(order_item_update_cols...).Update(&oi)
		core.HandleError(err)
	}
}
