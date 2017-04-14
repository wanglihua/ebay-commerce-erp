package controllers

import "github.com/revel/revel"

type ErpOrderImport struct {
	*revel.Controller
}

func (c ErpOrderImport) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_import/order_import_index.html")
}

func (c ErpOrderImport) ListData() revel.Result {
	return c.RenderEmpty()
}

// 上传 Excel 导入订单
func (c ErpOrderImport) Upload() revel.Result {
	return c.RenderEmpty()
}
