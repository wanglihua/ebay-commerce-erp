
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/track/models"
)

type TrackIssueComment struct {
	*revel.Controller
}

func (c TrackIssueComment) Index() revel.Result {
	return c.RenderTemplate("track/issue_comment/issue_comment_index.html")
}

func (c TrackIssueComment) ListData() revel.Result {
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

	issue_comment_list := make([]models.IssueCommentInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&issue_comment_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.IssueCommentInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_comment_list,
		Total: count,
	})
}

type IssueCommentDetailForm struct {
	IssueComment models.IssueCommentInfo `json:"issue_comment"`
}

func (f *IssueCommentDetailForm) IsCreate() bool {
	return f.IssueComment.Id == 0
}

func (f *IssueCommentDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.IssueComment.IssueId).Message("缺陷id不能为空！")

	validation.Required(f.IssueComment.Body).Message("正文不能为空！")
	if f.IssueComment.Body != "" {
		validation.MinSize(f.IssueComment.Body, 2).Message("正文长度不能小于2！")
	}
	if f.IssueComment.Body != "" {
		validation.MaxSize(f.IssueComment.Body, 4000).Message("正文长度不能大于4000！")
	}

	validation.Required(f.IssueComment.UserId).Message("提报人用户id不能为空！")

	validation.Required(f.IssueComment.UserCode).Message("提报人工号不能为空！")

	validation.Required(f.IssueComment.Create_time).Message("Create_time不能为空！")

	validation.Required(f.IssueComment.Update_time).Message("Update_time不能为空！")

	return validation.HasErrors() == false
}

func (c TrackIssueComment) DetailData() revel.Result {
	db_session := c.DbSession

	var issue_comment_id int64
	c.Params.Bind(&issue_comment_id, "id")

	var issue_comment models.IssueCommentInfo
	if issue_comment_id != 0 {
		has, err := db_session.Id(issue_comment_id).Get(&issue_comment)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的缺陷评论不存在！"})
		}
	}

	var detail_form IssueCommentDetailForm
	detail_form.IssueComment = issue_comment

	return c.RenderJson(detail_form)
}

func (c TrackIssueComment) Save() revel.Result {
	db_session := c.DbSession

	var detail_form IssueCommentDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	issue_comment := detail_form.IssueComment

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&issue_comment)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(issue_comment.Id).Update(&issue_comment)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c TrackIssueComment) Delete() revel.Result {
	db_session := c.DbSession

	issue_comment_id_list := make([]int64, 0)
	c.Params.Bind(&issue_comment_id_list, "id_list")

	affected, err := db_session.In("id", issue_comment_id_list).Delete(new(models.IssueCommentInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

