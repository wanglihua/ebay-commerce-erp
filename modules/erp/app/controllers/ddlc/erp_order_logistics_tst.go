package controllers

import (
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_logistics_service"
)

type ErpOrderLogisticsTest struct {
	*revel.Controller
}

func (c ErpOrderLogisticsTest) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_logistics_test/order_logistics_test_index.html")
}

func (c ErpOrderLogisticsTest) Match() revel.Result {

	var db_session = c.DbSession

	var order_codes string
	c.Params.Bind(&order_codes, "order_codes")

	var order_code_list = strings.Split(order_codes, ",")

	var order_list = make([]models.OrderInfo, 0)
	err := db_session.In("code", order_code_list).Find(&order_list)
	core.HandleError(err)

	var p_order_logistics = new(order_logistics_service.OrderLogistics)
	err = p_order_logistics.CalculateList(order_list, db_session)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "匹配物流完成！",
	})
}
