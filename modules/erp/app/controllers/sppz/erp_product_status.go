package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpProductStatus struct {
	*revel.Controller
}

func (c ErpProductStatus) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_STATUS, c.Session)

	return c.RenderTemplate("erp/sppz/product_status/product_status_index.html")
}

func (c ErpProductStatus) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_STATUS, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here
	query = query.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT))

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	status_list := make([]models.StatusInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&status_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.StatusInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  status_list,
		Total: count,
	})
}
