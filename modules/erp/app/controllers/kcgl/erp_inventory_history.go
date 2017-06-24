package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpInventoryHistory struct {
	*revel.Controller
}

func (c ErpInventoryHistory) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	var inventory_id int64
	c.Params.Bind(&inventory_id, "inventory_id")

	return c.RenderTemplate("erp/kcgl/inventory_history/inventory_history_index.html")
}

func (c ErpInventoryHistory) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	var inventory_id int64
	c.Params.Bind(&inventory_id, "inventory_id")

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

	inventory_history_list := make([]models.InventoryHistoryInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&inventory_history_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.InventoryHistoryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_history_list,
		Total: count,
	})
}
