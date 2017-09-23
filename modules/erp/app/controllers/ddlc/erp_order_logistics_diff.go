package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/setting"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
)

type ErpOrderLogisticsDiff struct {
	*revel.Controller
}

func (c ErpOrderLogisticsDiff) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_DIFF, c.Session)

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	return c.RenderTemplate("erp/ddlc/order_logistics_diff/order_logistics_diff_index.html")
}

func (c ErpOrderLogisticsDiff) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_DIFF, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
	})
}

func (c ErpOrderLogisticsDiff) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_LOGISTICS_DIFF, c.Session)

	var db_session = c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type logistics_diff_view_item_t struct {
		OrderLogisticsLogic  ddlc_models.OrderLogisticsLogicInfo  `xorm:"extends" json:"oll"`
		Order                ddlc_models.OrderInfo                `xorm:"extends" json:"o"`
		OrderShippingAddress ddlc_models.OrderShippingAddressInfo `xorm:"extends" json:"osa"`
		LogisticsMode        ddpz_models.LogisticsModeInfo        `xorm:"extends" json:"lm"`
		LogisticsModeSecond  ddpz_models.LogisticsModeInfo        `xorm:"extends" json:"lm_sec"`
	}

	/*
SELECT oll.*, o.*, osa.*, lm.*, lm_sec.*
FROM erp_order_logistics_logic oll
    INNER JOIN erp_order o ON o.id = oll.order_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = oll.address_id
    INNER JOIN erp_logistics_mode lm ON lm.id = oll.mode_id
    LEFT JOIN erp_logistics_mode lm_sec ON lm_sec.id = oll.mode_id_second
	 */

	sb := core.NewSqlBuilder()
	sb.Select("oll.*, o.*, osa.*, lm.*, lm_sec.*")
	sb.Count("count(oll.id) as count")
	sb.From(`erp_order_logistics_logic oll
    INNER JOIN erp_order o ON o.id = oll.order_id
    INNER JOIN erp_order_shipping_address osa ON osa.id = oll.address_id
    INNER JOIN erp_logistics_mode lm ON lm.id = oll.mode_id
    LEFT JOIN erp_logistics_mode lm_sec ON lm_sec.id = oll.mode_id_second`)

	sb.Where(filter)

	//query extra filter here

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("oll.id ASC")
	}

	sb.Limit(offset, limit)

	logistics_mode_list := make([]logistics_diff_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_mode_list)
	core.HandleError(err)

	var count core.Count
	_, err = db_session.Sql(sb.SqlCount()).Get(&count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_mode_list,
		Total: count.Count,
	})
}
