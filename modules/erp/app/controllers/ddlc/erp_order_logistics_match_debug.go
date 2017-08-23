package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/ddlc_models"
)

type ErpOrderLogisticsMatchDebug struct {
	*revel.Controller
}

func (c ErpOrderLogisticsMatchDebug) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_logistics_match_debug/order_logistics_match_debug_index.html")
}

type logistics_match_debug_view_item_t struct {
	OrderLogisticsMatchDebug ddlc_models.OrderLogisticsMatchDebugInfo `xorm:"extends" json:"md"`
	Stock                    sppz_models.StockInfo                    `xorm:"extends" json:"s"`
	OrderShippingAddress     ddlc_models.OrderShippingAddressInfo     `xorm:"extends" json:"osa"`
	LogisticsCompany         ddpz_models.LogisticsCompanyInfo         `xorm:"extends" json:"lc"`
	LogisticsMode            ddpz_models.LogisticsModeInfo            `xorm:"extends" json:"lm"`
	LogisticsFeeRate         ddpz_models.LogisticsFeeRateInfo         `xorm:"extends" json:"lfr"`
}

func (c ErpOrderLogisticsMatchDebug) ListData() revel.Result {

	var order_code_query string
	c.Params.Bind(&order_code_query, "order_code_query")

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT md.*, s.*, osa.*, lc.*, lm.*, lfr.*
FROM erp_order_logistics_match_debug md
    INNER JOIN erp_stock s ON s.id = md.stock_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = md.address_id
    INNER JOIN erp_logistics_company lc ON lc.id = md.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = md.mode_id
    LEFT JOIN erp_logistics_fee_rate lfr ON lfr.id = md.fee_rate_id
WHERE md.order_code = '2053968244015'

	var sb = core.NewSqlBuilder()
	sb.Select("md.*, s.*, osa.*, lc.*, lm.*, lfr.*")
	sb.Count("count(md.id) as match_debug_count")
	sb.From(`erp_order_logistics_match_debug md
    INNER JOIN erp_stock s ON s.id = md.stock_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = md.address_id
    INNER JOIN erp_logistics_company lc ON lc.id = md.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = md.mode_id
    LEFT JOIN erp_logistics_fee_rate lfr ON lfr.id = md.fee_rate_id`)

	if order_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("md.order_code = '%s'", order_code_query))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("md.id ASC")
	}
	*/

	var sb = core.NewSqlBuilder()
	sb.Select("id, order_code, address_postal_code, stock_name, company_name, mode_name, fee_rate_post_code, total_fee, total_cost, match_result, create_time")
	sb.Count("count(1) as match_debug_count")
	sb.From(`erp_order_logistics_match_debug`)

	if order_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("order_code = '%s'", order_code_query))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("id DESC")
	}

	sb.Limit(offset, limit)

	order_logistics_match_debug_list := make([]ddlc_models.OrderLogisticsMatchDebugInfo, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&order_logistics_match_debug_list)
	core.HandleError(err)

	var match_debug_count struct {
		Count int64 `xorm:"'match_debug_count'"`
	}

	if order_code_query != "" {
		_, err = db_session.Sql(sb.SqlCount()).Get(&match_debug_count)
		core.HandleError(err)
	} else {
		var sql_count = `SELECT id AS match_debug_count FROM erp_order_logistics_match_debug ORDER BY id DESC LIMIT 1;`
		_, err = db_session.Sql(sql_count).Get(&match_debug_count)
		core.HandleError(err)
	}

	return c.RenderJson(core.GridResult{
		Data:  order_logistics_match_debug_list,
		Total: match_debug_count.Count,
	})
}

func (c ErpOrderLogisticsMatchDebug) Detail(id int64) revel.Result {
	var db_session = c.DbSession

	/*
SELECT md.*, s.*, osa.*, lc.*, lm.*, lfr.*
FROM erp_order_logistics_match_debug md
    INNER JOIN erp_stock s ON s.id = md.stock_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = md.address_id
    INNER JOIN erp_logistics_company lc ON lc.id = md.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = md.mode_id
    LEFT JOIN erp_logistics_fee_rate lfr ON lfr.id = md.fee_rate_id
WHERE md.id = 1;
    */

	var sql_get = `SELECT md.*, s.*, osa.*, lc.*, lm.*, lfr.*
FROM erp_order_logistics_match_debug md
    INNER JOIN erp_stock s ON s.id = md.stock_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = md.address_id
    INNER JOIN erp_logistics_company lc ON lc.id = md.company_id
    INNER JOIN erp_logistics_mode lm ON lm.id = md.mode_id
    LEFT JOIN erp_logistics_fee_rate lfr ON lfr.id = md.fee_rate_id
WHERE md.id = %d;`
	sql_get = fmt.Sprintf(sql_get, id)

	var logistics_match_debug logistics_match_debug_view_item_t
	_, err := db_session.Sql(sql_get).Get(&logistics_match_debug)
	core.HandleError(err)

	c.RenderArgs["md"] = logistics_match_debug
	return c.RenderTemplate("erp/ddlc/order_logistics_match_debug/order_logistics_match_debug_detail.html")
}

func (c ErpOrderLogisticsMatchDebug) Clear() revel.Result {
	var db_session = c.DbSession

	_, err := db_session.Exec("DROP TABLE erp_order_logistics_match_debug;")
	core.HandleError(err)

	err = db_session.Sync2(new(ddlc_models.OrderLogisticsMatchDebugInfo))
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "清理成功!",
	})
}
