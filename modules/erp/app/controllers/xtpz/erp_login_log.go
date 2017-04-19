package controllers

import (
	"fmt"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpLoginLog struct {
	*revel.Controller
}

func (c ErpLoginLog) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LOGIN_LOG, c.Session)

	return c.RenderTemplate("erp/xtpz/login_log/login_log_index.html")
}

func (c ErpLoginLog) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LOGIN_LOG, c.Session)

	db_session := c.DbSession

	var login_name_query string
	var begin_time_query core.NullLongTime
	var end_time_query core.NullLongTime

	c.Params.Bind(&login_name_query, "login_name_query")
	c.Params.Bind(&begin_time_query, "begin_time_query")
	c.Params.Bind(&end_time_query, "end_time_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here
	if strings.TrimSpace(login_name_query) != "" {
		query = query.Where(fmt.Sprintf("login_name like '%%%s%%'", login_name_query))
	}

	if begin_time_query.Valid == true {
		query = query.Where(fmt.Sprintf("create_time >= '%s'", begin_time_query.Time.Format("2006-01-02 15:04:05")))
	}

	if end_time_query.Valid == true {
		query = query.Where(fmt.Sprintf("create_time <= '%s'", end_time_query.Time.Format("2006-01-02 15:04:05")))
	}

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
