package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpProductCateLog struct {
	*revel.Controller
}

func (c ErpProductCateLog) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CATE, c.Session)

	return c.RenderTemplate("erp/sppz/product_cate_log/product_cate_log_index.html")
}

func (c ErpProductCateLog) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CATE, c.Session)

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

	product_cate_log_list := make([]models.ProductCateLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_cate_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductCateLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_cate_log_list,
		Total: count,
	})
}
