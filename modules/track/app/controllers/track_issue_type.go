package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/track/models"
)

type TrackIssueType struct {
	*revel.Controller
}

func (c TrackIssueType) Index() revel.Result {
	return c.RenderTemplate("track/issue_type/issue_type_index.html")
}

func (c TrackIssueType) ListData() revel.Result {
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

	issue_type_list := make([]models.IssueTypeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&issue_type_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.IssueTypeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_type_list,
		Total: count,
	})
}
