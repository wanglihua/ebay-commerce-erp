package controllers

import (
	"github.com/revel/revel"
)

type Log_order struct {
	*revel.Controller
}

func (c Log_order) Index() revel.Result {
	return c.RenderTemplate("log/order/order_index.html")
}
