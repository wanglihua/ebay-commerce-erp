package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

type ErpProductStorageLocBatchEdit struct {
	*revel.Controller
}

func (c ErpProductStorageLocBatchEdit) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	return c.RenderTemplate("erp/sppz/product_storage_loc_batch_edit/product_storage_loc_batch_edit_index.html")
}

func (c ErpProductStorageLocBatchEdit) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STORAGE_LOC, c.Session)

	return c.RenderJson(core.JsonResult{Success: false, Message: "功能未实现！"})
}
