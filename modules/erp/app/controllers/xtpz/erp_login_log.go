package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpLoginLog struct {
	*revel.Controller
}

func (c ErpLoginLog) Index() revel.Result {
	return c.RenderTemplate("erp/xtpz/login_log/login_log_index.html")
}

func (c ErpLoginLog) ListData() revel.Result {
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

	login_log_list := make([]models.LoginLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&login_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.LoginLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  login_log_list,
		Total: count,
	})
}
