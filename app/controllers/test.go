package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/service/inventory_service"
)

type Test struct {
	*revel.Controller
}

func (c Test) Index() revel.Result {

	go inventory_service.Update_au_stock_cost()

	return c.RenderJson(map[string]interface{}{
		"success": true,
	})
}
