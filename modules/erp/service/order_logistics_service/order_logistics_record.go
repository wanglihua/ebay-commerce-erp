package order_logistics_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_history_service"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/pick_fee_service"
)

func save_order_logistics(p_ol *OrderLogisticsService, p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, match_result match_logistics_mode_result, unlock_express bool, login_name string, db_session *xorm.Session) models.OrderLogisticsInfo {
	// 匹配成功，记录匹配成功数据
	var order_logistics models.OrderLogisticsInfo
	has, err := db_session.Where(fmt.Sprintf("order_id = %d", match_result.order_id)).Get(&order_logistics)
	core.HandleError(err)

	order_logistics.OrderId = match_result.order_id
	order_logistics.AddressId = match_result.address_id
	order_logistics.CompanyId = match_result.company_id
	order_logistics.StockId = match_result.stock_id
	order_logistics.ModeId = match_result.mode_id
	order_logistics.FeeRateId = match_result.fee_rate_id
	order_logistics.Length = match_result.length
	order_logistics.Width = match_result.width
	order_logistics.Height = match_result.height
	order_logistics.DimWeight = match_result.dim_weight
	order_logistics.Sfkltjz = match_result.sfkltjz
	order_logistics.Weight = match_result.actual_weight
	order_logistics.TotalFee = match_result.total_fee
	order_logistics.PickFee = match_result.pick_fee
	order_logistics.CostAmount = match_result.cost_amount
	order_logistics.Ptfl = p_ol.logistics_setting.Ptfl
	order_logistics.Jyfl = p_ol.logistics_setting.Jyfl

	if has == true {
		_, err := db_session.Id(order_logistics.Id).Update(&order_logistics)
		core.HandleError(err)
	} else {
		_, err := db_session.Insert(&order_logistics)
		core.HandleError(err)
	}

	var update_col_list = []string{"status_id", "stock_id", "address_id", "company_id", "mode_id", "hold_type_id"}

	if match_result.has_stock == true {
		if match_result.mode.TypeId == models.LOGISTICS_MODE_TYPE_EXPRESS && unlock_express == false { // 匹配出来是快递，并且当前不解锁快递，就可以设置为快递 hold
			// 把订单状态修改为 Hold
			p_order.HoldTypeId = core.MakeNullInt(models.ORDER_HOLD_TYPE_EXPRESS, true)
			p_order.StatusId = models.ORDER_STATUS_HOLD // 修改订单状态为 hold 状态

			if login_name == "" {
				login_name = "自动任务"
			}
			order_history_service.Insert_order_history(*p_order, fmt.Sprintf("匹配物流订单Hold 快递类，工号：%s", login_name), db_session)
		} else {
			// 把订单状态修改为 待配货
			p_order.HoldTypeId = core.MakeNullInt(0, false)
			p_order.StatusId = models.ORDER_STATUS_DPH

			if login_name == "" {
				login_name = "自动任务"
			}
			order_history_service.Insert_order_history(*p_order, fmt.Sprintf("匹配物流订单待配货，工号：%s", login_name), db_session)
		}
	}

	// 填入匹配过程中的最后的 stock_id 和 address_id
	p_order.StockId = match_result.stock_id
	p_order.AddressId = match_result.address_id
	p_order.CompanyId = match_result.company_id
	p_order.ModeId = match_result.mode_id
	_, err = db_session.Id(p_order.Id).Cols(update_col_list...).Update(p_order)
	core.HandleError(err)

	// 保存成本价和分配利润
	var order_profit = order_service.CalculateOrderProfit(p_order, &order_logistics)
	var cost_price_total float64 = 0
	var product_info_list = p_order_logistics_data.order_product_info_list
	for _, product_info := range product_info_list {
		cost_price_total = cost_price_total + product_info.CostPrice
	}

	if cost_price_total != 0 {
		for _, product_info := range product_info_list {
			var item_profit = product_info.CostPrice * order_profit / cost_price_total

			var order_item models.OrderItemInfo
			order_item.CostPrice = product_info.CostPrice
			order_item.GrossProfit = item_profit
			_, err = db_session.Id(product_info.ItemId).Cols("cost_price", "gross_profit").Update(&order_item)
			core.HandleError(err)
		}
	}

	return order_logistics
}

func save_special_hold_order_logistics(p_ol *OrderLogisticsService, p_order *models.OrderInfo, p_order_logistics_data *order_logistics_data_t, envelope_fee_id int64, total_fee float64, db_session *xorm.Session) {

	var order_logistics models.OrderLogisticsInfo
	has, err := db_session.Where("order_id = ?", p_order.Id).Get(&order_logistics)
	core.HandleError(err)

	var whole_length, whole_width, whole_height, whole_weight float64 = get_order_products_whole_logistics_data(p_order_logistics_data)

	order_logistics.OrderId = p_order.Id
	order_logistics.AddressId = p_order_logistics_data.order_shipping_address.Id
	order_logistics.CompanyId = 0
	order_logistics.StockId = p_ol.mel_stock_id //设置发货仓库为墨尔本仓库
	order_logistics.ModeId = 0
	order_logistics.FeeRateId = 0
	order_logistics.EnvelopeFeeId = envelope_fee_id
	order_logistics.Length = whole_length
	order_logistics.Width = whole_width
	order_logistics.Height = whole_height
	order_logistics.Weight = whole_weight
	order_logistics.DimWeight = 0
	order_logistics.Sfkltjz = 0
	order_logistics.Weight = whole_weight
	order_logistics.Sfkltjz = 0
	order_logistics.TotalFee = total_fee
	order_logistics.PickFee = pick_fee_service.Get_product_pick_fee(whole_weight, p_ol.pick_fee_setting)
	_, order_logistics.CostAmount = check_stock_has_order_products_and_get_cost_amount(p_ol, p_ol.mel_stock_id, p_order_logistics_data)
	order_logistics.Ptfl = p_ol.logistics_setting.Ptfl
	order_logistics.Jyfl = p_ol.logistics_setting.Jyfl

	if has == false {
		_, err = db_session.Insert(&order_logistics)
		core.HandleError(err)
	} else {
		_, err = db_session.Id(order_logistics.Id).Update(&order_logistics)
		core.HandleError(err)
	}

	p_order.StockId = p_ol.mel_stock_id
	p_order.AddressId = p_order_logistics_data.order_shipping_address.Id
	p_order.CompanyId = 0
	p_order.ModeId = 0

	var update_col_list = []string{"stock_id", "address_id", "company_id", "mode_id"}
	_, err = db_session.Id(p_order.Id).Cols(update_col_list...).Update(p_order)
	core.HandleError(err)

	var product_info_list = p_order_logistics_data.order_product_info_list
	for _, product_info := range product_info_list {
		p_ol.pending_keeper.add(order_logistics.StockId, product_info.ProductId, product_info.Qty)
	}
}
