package order_logistics_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_history_service"
)

func save_order_logistics(p_order *models.OrderInfo, match_result match_logistics_mode_result, unlock_express bool, login_name string, db_session *xorm.Session) models.OrderLogisticsInfo {
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

	return order_logistics
}
