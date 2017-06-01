package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/service/stock_pick_crew_service"
)

type Test struct {
	*revel.Controller
}

func (c Test) Index() revel.Result {
	var db_session = c.DbSession
	var web_session = c.Session

	var stock_id_str_list = stock_pick_crew_service.Get_login_user_stock_id_str_list(db_session, web_session)

	return c.RenderJson(stock_id_str_list)
}
