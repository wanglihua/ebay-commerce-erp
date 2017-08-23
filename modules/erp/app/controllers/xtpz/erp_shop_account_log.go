package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models/xtpz_models"
)

type ErpShopAccountLog struct {
	*revel.Controller
}

func (c ErpShopAccountLog) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	return c.RenderTemplate("erp/xtpz/shop_account_log/shop_account_log_index.html")
}

func (c ErpShopAccountLog) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	shop_account_log_list := make([]xtpz_models.ShopAccountLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&shop_account_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(xtpz_models.ShopAccountLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  shop_account_log_list,
		Total: count,
	})
}
