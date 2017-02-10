package controllers

import (
	"github.com/revel/revel"

	"fmt"
	"matrix/core"
	"matrix/modules/erp/models"
	"strings"
)

type ErpProductLog struct {
	*revel.Controller
}

func (c ErpProductLog) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/product_log/product_log_index.html")
}

func (c ErpProductLog) ListData() revel.Result {
	db_session := c.DbSession

	var product_code_query string
	var login_name_query string
	var begin_time_query core.NullLongTime
	var end_time_query core.NullLongTime

	c.Params.Bind(&product_code_query, "product_code_query")
	c.Params.Bind(&login_name_query, "login_name_query")
	c.Params.Bind(&begin_time_query, "begin_time_query")
	c.Params.Bind(&end_time_query, "end_time_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here
	if strings.TrimSpace(product_code_query) != "" {
		query = query.Where(fmt.Sprintf("product_code like '%%%s%%'", product_code_query))
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

	product_log_list := make([]models.ProductLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_log_list,
		Total: count,
	})
}
