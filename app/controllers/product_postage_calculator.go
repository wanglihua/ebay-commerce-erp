package controllers

import (
	"github.com/revel/revel"
	"strings"
	"matrix/modules/erp/service/product_logistics_service"
	"fmt"
)

type ProductPostageCalculator struct {
	*revel.Controller
}

func (c ProductPostageCalculator) Index() revel.Result {
	var sku string
	c.Params.Bind(&sku, "sku")
	sku = strings.TrimSpace(sku)

	c.RenderArgs["sku"] = sku

	return c.RenderTemplate("product_postage_calculator/product_postage_calculator_index.html")
}

func (c ProductPostageCalculator) Post() revel.Result {
	var db_session = c.DbSession

	var sku string
	c.Params.Bind(&sku, "sku")
	sku = strings.TrimSpace(sku)

	var postcode string
	c.Params.Bind(&postcode, "postcode")
	postcode = strings.TrimSpace(postcode)

	var err, fee = product_logistics_service.Match_product_postage(sku, postcode, db_session)

	if err != nil {
		c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "something is wrong!",
		})
	}

	if fee < 0 {
		c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "non delivery",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"fee":     fmt.Sprintf("%0.2f", fee),
	})
}
