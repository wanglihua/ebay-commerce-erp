package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
)

type ErpProductBatchEdit struct {
	*revel.Controller
}

func (c ErpProductBatchEdit) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/product_batch_edit/product_batch_edit_index.html")
}

func (c ErpProductBatchEdit) Save() revel.Result {
	return c.RenderJson(core.JsonResult{Success: false, Message: "功能未实现！"})
}
