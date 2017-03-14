package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/service/message_service"
)

type ErpMessage struct {
	*revel.Controller
}

func (c ErpMessage) Index() revel.Result {
	return c.RenderTemplate("erp/xxgl/message/message_index.html")
}

func (c ErpMessage) ListData() revel.Result {
	db_session := c.DbSession

	message_list, msg_count := message_service.Get_message_list_data(db_session, c.Request)

	return c.RenderJson(core.GridResult{
		Data:  message_list,
		Total: msg_count,
	})
}
