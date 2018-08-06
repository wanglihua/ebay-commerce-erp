package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/service/ebay_token_status_check_service"
)

type Test struct {
	*revel.Controller
}

func (c Test) Index() revel.Result {

	ebay_token_status_check_service.Check_ebay_token_status()

	return c.RenderTemplate("test/test_index.html")
}
