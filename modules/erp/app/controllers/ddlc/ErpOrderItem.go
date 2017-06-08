package controllers

import "github.com/revel/revel"

type ErpOrderItem struct {
	*revel.Controller
}

func (c ErpOrderItem) DetailData() revel.Result {
	return c.RenderEmpty()
}

func (c ErpOrderItem) Save() revel.Result {
	return c.RenderEmpty()
}
