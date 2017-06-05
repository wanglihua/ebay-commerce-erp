package controllers

import (
	"github.com/revel/revel"
)

type ErpOrderLogisticsImport struct {
	*revel.Controller
}

func (c ErpOrderLogisticsImport) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_logistics_import/order_logistics_import_index.html")
}

// 导入运单号
func (c ErpOrderLogisticsImport) ImportTrackingNumber() revel.Result {
	return c.RenderEmpty()
}

// 导入实际运费
func (c ErpOrderLogisticsImport) ImportActualFee() revel.Result {
	return c.RenderEmpty()
}
