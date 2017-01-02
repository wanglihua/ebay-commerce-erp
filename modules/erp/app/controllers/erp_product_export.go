
package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductExport struct {
	*revel.Controller
}

func (c ErpProductExport) Index() revel.Result {
	return c.RenderTemplate("erp/product_export/product_export.html")
}

func (c ErpProductExport) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_list := make([]models.ProductInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: count,
	})
}

func (c ErpProductExport) Export() revel.Result {
	return c.RenderJson(core.JsonResult{Success: true, Message: "数据导出成功!"})
}
