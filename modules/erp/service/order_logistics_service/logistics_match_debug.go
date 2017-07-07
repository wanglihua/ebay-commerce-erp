package order_logistics_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func log_logistics_match_debug(match_result match_logistics_mode_result, result_desc string, match_debug string, db_session *xorm.Session) {
	var logistics_match_debug models.OrderLogisticsMatchDebugInfo

	logistics_match_debug.OrderId = match_result.order_id
	logistics_match_debug.OrderCode = match_result.order_code
	logistics_match_debug.AddressId = match_result.address_id
	logistics_match_debug.AddressPostalCode = match_result.address.PostalCode
	logistics_match_debug.CompanyId = match_result.company_id
	logistics_match_debug.CompanyName = match_result.company.Name
	logistics_match_debug.StockId = match_result.stock_id
	logistics_match_debug.StockName = match_result.stock.Name
	logistics_match_debug.ModeId = match_result.mode_id
	logistics_match_debug.ModeName = match_result.mode.Name
	logistics_match_debug.FeeRateId = match_result.fee_rate_id
	logistics_match_debug.FeeRatePostCode = match_result.fee_rate.Postcode
	logistics_match_debug.TotalFee = match_result.total_fee
	logistics_match_debug.TotalCost = match_result.total_amount
	logistics_match_debug.MatchResult = result_desc
	logistics_match_debug.MatchDebug = match_debug

	_, err := db_session.Insert(&logistics_match_debug)
	core.HandleError(err)
}
