
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/track/models"
)

type TrackIssueLog struct {
	*revel.Controller
}

func (c TrackIssueLog) Index() revel.Result {
	return c.RenderTemplate("track/issue_log/issue_log_index.html")
}

func (c TrackIssueLog) ListData() revel.Result {
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

	issue_log_list := make([]models.IssueLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&issue_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.IssueLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_log_list,
		Total: count,
	})
}

type IssueLogDetailForm struct {
	IssueLog models.IssueLogInfo `json:"issue_log"`
}

func (f *IssueLogDetailForm) IsCreate() bool {
	return f.IssueLog.Id == 0
}

func (f *IssueLogDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.IssueLog.IssueId).Message("缺陷id不能为空！")

	validation.Required(f.IssueLog.Content).Message("操作内容不能为空！")

	validation.Required(f.IssueLog.UserId).Message("录入人用户id不能为空！")

	validation.Required(f.IssueLog.UserCode).Message("录入人工号不能为空！")

	return validation.HasErrors() == false
}

func (c TrackIssueLog) DetailData() revel.Result {
	db_session := c.DbSession

	var issue_log_id int64
	c.Params.Bind(&issue_log_id, "id")

	var issue_log models.IssueLogInfo
	if issue_log_id != 0 {
		has, err := db_session.Id(issue_log_id).Get(&issue_log)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的缺陷日志不存在！"})
		}
	}

	var detail_form IssueLogDetailForm
	detail_form.IssueLog = issue_log

	return c.RenderJson(detail_form)
}

func (c TrackIssueLog) Save() revel.Result {
	db_session := c.DbSession

	var detail_form IssueLogDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	issue_log := detail_form.IssueLog

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&issue_log)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(issue_log.Id).Update(&issue_log)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c TrackIssueLog) Delete() revel.Result {
	db_session := c.DbSession

	issue_log_id_list := make([]int64, 0)
	c.Params.Bind(&issue_log_id_list, "id_list")

	affected, err := db_session.In("id", issue_log_id_list).Delete(new(models.IssueLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

