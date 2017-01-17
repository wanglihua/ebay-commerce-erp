package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderPhz struct {
	*revel.Controller
}

func (c ErpOrderPhz) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_phz/order_phz_index.html")
}

// /erp/order/phz
// erp_order_phz
// order_phz
// ErpOrderPhz
// erp_order_phz_index
