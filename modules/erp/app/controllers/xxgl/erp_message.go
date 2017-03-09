package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpMessage struct {
	*revel.Controller
}

func (c ErpMessage) Index() revel.Result {
	return c.RenderTemplate("erp/xxgl/message/message_index.html")
}

func (c ErpMessage) ListData() revel.Result {
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

	message_list := make([]models.MessageInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&message_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.MessageInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  message_list,
		Total: count,
	})
}
