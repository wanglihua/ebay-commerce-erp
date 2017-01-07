package controllers

import (
	"matrix/core"
	"github.com/revel/revel"
)

type ErpProductStorageLocBatchEdit struct {
	*revel.Controller
}

func (c ErpProductStorageLocBatchEdit) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/product_storage_loc_batch_edit/product_storage_loc_batch_edit_index.html")
}

func (c ErpProductStorageLocBatchEdit) Save() revel.Result {
	return c.RenderJson(core.JsonResult{Success: false, Message: "功能未实现！"})
}
