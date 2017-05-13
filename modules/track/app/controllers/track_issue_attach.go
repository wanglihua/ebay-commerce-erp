
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/track/models"
)

type TrackIssueAttach struct {
	*revel.Controller
}

func (c TrackIssueAttach) Index() revel.Result {
	return c.RenderTemplate("track/issue_attach/issue_attach_index.html")
}

func (c TrackIssueAttach) ListData() revel.Result {
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

	issue_attach_list := make([]models.IssueAttachInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&issue_attach_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.IssueAttachInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_attach_list,
		Total: count,
	})
}

type IssueAttachDetailForm struct {
	IssueAttach models.IssueAttachInfo `json:"issue_attach"`
}

func (f *IssueAttachDetailForm) IsCreate() bool {
	return f.IssueAttach.Id == 0
}

func (f *IssueAttachDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.IssueAttach.IssueId).Message("缺陷id不能为空！")

	validation.Required(f.IssueAttach.Type).Message("类型不能为空！")

	validation.Required(f.IssueAttach.DirName).Message("路径名不能为空！")

	validation.Required(f.IssueAttach.FileName).Message("文件名不能为空！")

	validation.Required(f.IssueAttach.Status).Message("状态不能为空！")

	validation.Required(f.IssueAttach.Create_time).Message("Create_time不能为空！")

	validation.Required(f.IssueAttach.Update_time).Message("Update_time不能为空！")

	return validation.HasErrors() == false
}

func (c TrackIssueAttach) DetailData() revel.Result {
	db_session := c.DbSession

	var issue_attach_id int64
	c.Params.Bind(&issue_attach_id, "id")

	var issue_attach models.IssueAttachInfo
	if issue_attach_id != 0 {
		has, err := db_session.Id(issue_attach_id).Get(&issue_attach)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的缺陷附件不存在！"})
		}
	}

	var detail_form IssueAttachDetailForm
	detail_form.IssueAttach = issue_attach

	return c.RenderJson(detail_form)
}

func (c TrackIssueAttach) Save() revel.Result {
	db_session := c.DbSession

	var detail_form IssueAttachDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	issue_attach := detail_form.IssueAttach

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&issue_attach)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(issue_attach.Id).Update(&issue_attach)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c TrackIssueAttach) Delete() revel.Result {
	db_session := c.DbSession

	issue_attach_id_list := make([]int64, 0)
	c.Params.Bind(&issue_attach_id_list, "id_list")

	affected, err := db_session.In("id", issue_attach_id_list).Delete(new(models.IssueAttachInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

