package controllers

import "github.com/revel/revel"

type ErpOrderWaybillPrint struct {
	*revel.Controller
}

func (c ErpOrderWaybillPrint) Index() revel.Result {
	return c.RenderEmpty()
}

func (c ErpOrderWaybillPrint) Preview() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_waybill_print/order_waybill_print_preview.html")
}
