package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/track/models"
)

type TrackIssueStatus struct {
	*revel.Controller
}

func (c TrackIssueStatus) Index() revel.Result {
	return c.RenderTemplate("track/issue_status/issue_status_index.html")
}

func (c TrackIssueStatus) ListData() revel.Result {
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

	issue_status_list := make([]models.IssueStatusInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&issue_status_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.IssueStatusInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_status_list,
		Total: count,
	})
}
