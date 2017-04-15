package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
)

type ErpOrderHold struct {
	*revel.Controller
}

func (c ErpOrderHold) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_hold/order_hold_index.html")
}

func (c ErpOrderHold) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_HOLD))
	// filter = core.SqlAnd(filter, fmt.Sprintf("hold_type_id = %d", models.asdfasdfasdfa))

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
