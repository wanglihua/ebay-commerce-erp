package controllers

import (
	"fmt"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductBindLog struct {
	*revel.Controller
}

func (c ErpProductBindLog) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/product_bind_log/product_bind_log_index.html")
}

func (c ErpProductBindLog) ListData() revel.Result {
	db_session := c.DbSession

	var product_bind_code_query string
	var login_name_query string
	var begin_time_query core.NullLongTime
	var end_time_query core.NullLongTime

	c.Params.Bind(&product_bind_code_query, "product_bind_code_query")
	c.Params.Bind(&login_name_query, "login_name_query")
	c.Params.Bind(&begin_time_query, "begin_time_query")
	c.Params.Bind(&end_time_query, "end_time_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here
	if strings.TrimSpace(product_bind_code_query) != "" {
		query = query.Where(fmt.Sprintf("product_bind_code like '%%%s%%'", product_bind_code_query))
	}

	if strings.TrimSpace(login_name_query) != "" {
		query = query.Where(fmt.Sprintf("login_name like '%%%s%%'", login_name_query))
	}

	if begin_time_query.Valid == true {
		query = query.Where(fmt.Sprintf("create_time >= '%s'", begin_time_query.Time.Format("2006-01-02 15:04:05")))
	}

	if end_time_query.Valid == true {
		query = query.Where(fmt.Sprintf("create_time <= '%s'", end_time_query.Time.Format("2006-01-02 15:04:05")))
	}

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_bind_log_list := make([]models.ProductBindLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_bind_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductBindLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_bind_log_list,
		Total: count,
	})
}
