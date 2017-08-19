package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

type ErpSupplierExcel struct {
	*revel.Controller
}

func (c ErpSupplierExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

	return c.RenderEmpty()
}

func (c ErpSupplierExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_SUPPLIER, c.Session)

	return c.RenderEmpty()
}
