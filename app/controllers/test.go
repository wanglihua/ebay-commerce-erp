package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_logistics_service"
)

type Test struct {
	*revel.Controller
}

func (c Test) Index() revel.Result {
	var db_session = c.DbSession

	var order models.OrderInfo
	_, err := db_session.Where(fmt.Sprintf("code = '%s'", "18797")).Get(&order)
	core.HandleError(err)

	order.StatusId = models.ORDER_STATUS_YZF
	_, err = db_session.Id(order.Id).Cols("status_id").Update(&order)
	core.HandleError(err)

	var p_ol = new(order_logistics_service.OrderLogistics)
	err = p_ol.CalculateList([]models.OrderInfo{order}, db_session)
	core.HandleError(err)

	var order_logistics models.OrderLogisticsInfo
	_, err = db_session.Where(fmt.Sprintf("order_id = %d", order.Id)).Get(&order_logistics)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"length": order_logistics.Length,
		"width":  order_logistics.Width,
		"height": order_logistics.Height,
		"weight": order_logistics.Weight,
	})
}
