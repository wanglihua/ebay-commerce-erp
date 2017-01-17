package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
)

type ErpProductImport struct {
	*revel.Controller
}

func (c ErpProductImport) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_import/product_import.html")
}

func (c ErpProductImport) ListData() revel.Result {
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

	product_import_list := make([]models.ProductImportInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_import_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductImportInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_import_list,
		Total: count,
	})
}

func (c ErpProductImport) Import() revel.Result {
	return c.RenderJson(core.JsonResult{Success: true, Message:  "Excel数据上传成功!"})
}

func (c ErpProductImport) Confirm() revel.Result {
	return c.RenderJson(core.JsonResult{Success: true, Message:  "商品数据导入成功!"})
}

func (c ErpProductImport) Clear() revel.Result {
	db_session := c.DbSession

	login_user := core.GetLoginUser(c.Session)
	login_user_id := login_user.UserId

	_, err := db_session.
	Where(fmt.Sprintf("%s = ?", models.ProductImportCols.UserId), login_user_id).
		Delete(new(models.ProductImportInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message:  "临时导入数据清空成功!"})
}
