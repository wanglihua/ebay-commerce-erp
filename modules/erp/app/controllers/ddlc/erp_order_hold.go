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

func (c ErpOrderHold) PageData() revel.Result {
	var db_session = c.DbSession

	// hold类型列表
	var hold_type_list = make([]models.OrderHoldTypeInfo, 0)
	err := db_session.OrderBy("id ASC").Find(&hold_type_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"hold_type_id":   0,
		"hold_type_list": hold_type_list,
	})
}

func (c ErpOrderHold) ListData() revel.Result {
	db_session := c.DbSession

	var hold_type_id int64
	c.Params.Bind(&hold_type_id, "hold_type_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	filter = core.SqlAnd(filter, fmt.Sprintf("o.status_id = %d", models.ORDER_STATUS_HOLD))

	if hold_type_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.hold_type_id = %d", hold_type_id))
	}

	order_list, count := order_service.Get_order_data_list(db_session, filter, order, offset, limit)

	return c.RenderJson(core.GridResult{
		Data:  order_list,
		Total: count,
	})
}
