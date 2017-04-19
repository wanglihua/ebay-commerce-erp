package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpExchangeRateLog struct {
	*revel.Controller
}

func (c ErpExchangeRateLog) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EXCHANGE_RATE, c.Session)

	return c.RenderTemplate("erp/xtpz/exchange_rate_log/exchange_rate_log_index.html")
}

func (c ErpExchangeRateLog) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EXCHANGE_RATE, c.Session)

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

	exchange_rate_log_list := make([]models.ExchangeRateLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&exchange_rate_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ExchangeRateLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  exchange_rate_log_list,
		Total: count,
	})
}
