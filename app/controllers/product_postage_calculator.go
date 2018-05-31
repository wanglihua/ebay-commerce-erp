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

func (c ProductPostageCalculator) Index(id int64) revel.Result {
	var product_id = id
	c.RenderArgs["product_id"] = product_id
	return c.RenderTemplate("product_postage_calculator/product_postage_calculator_index.html")
}

func (c ProductPostageCalculator) Post() revel.Result {
	var db_session = c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "product_id")

	var postcode string
	c.Params.Bind(&postcode, "postcode")
	postcode = strings.TrimSpace(postcode)

	if postcode == "" {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "please input postcode",
		})
	}

	var err, fee = product_logistics_service.Match_product_postage_for_client(product_id, postcode, db_session)

	if err != nil {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Error，please contact customer service",
		})
	}

	if fee < 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Sorry, this area does not support",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"fee":     fmt.Sprintf("$%0.2f", fee),
	})
}
