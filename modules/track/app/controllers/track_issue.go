package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
	"matrix/app/routes"
	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/track/models"
)

type TrackIssue struct {
	*revel.Controller
}

func (c TrackIssue) Index() revel.Result {
	return c.RenderTemplate("track/issue/issue_index.html")
}

func (c TrackIssue) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type issue_view_item_t struct {
		Issue       models.IssueInfo       `xorm:"extends" json:"ti"`
		IssueType   models.IssueTypeInfo   `xorm:"extends" json:"tit"`
		IssueCate   models.IssueCateInfo   `xorm:"extends" json:"tic"`
		IssueStatus models.IssueStatusInfo `xorm:"extends" json:"tis"`
	}

	/*
SELECT ti.*, tit.*, tic.*, tis.*
FROM track_issue ti
    INNER JOIN track_issue_type tit ON tit.id = ti.type_id
    LEFT JOIN track_issue_cate tic ON tic.id = ti.cate_id
    INNER JOIN track_issue_status tis ON tis.id = ti.status_id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("ti.*, tit.*, tic.*, tis.*")
	sb.Count("count(ti.id) as issue_count")
	sb.From(`track_issue ti
    INNER JOIN track_issue_type tit ON tit.id = ti.type_id
    LEFT JOIN track_issue_cate tic ON tic.id = ti.cate_id
    INNER JOIN track_issue_status tis ON tis.id = ti.status_id`)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("ti.id DESC")
	}

	sb.Limit(offset, limit)

	issue_list := make([]issue_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&issue_list)
	core.HandleError(err)

	var issue_count struct {
		Count int64  `xorm:"'issue_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&issue_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_list,
		Total: issue_count.Count,
	})
}

type IssueDetailForm struct {
	Issue         models.IssueInfo       `json:"issue"`
	IssueTypeList []models.IssueTypeInfo `json:"issue_type_list"`
	IssueCateList []models.IssueCateInfo `json:"issue_cate_list"`
	UserList      []auth_models.UserInfo `json:"user_list"`
}

func (f *IssueDetailForm) IsCreate() bool {
	return f.Issue.Id == 0
}

func (f *IssueDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Issue.TypeId).Message("类型id不能为空！")

	validation.Required(f.Issue.Title).Message("标题不能为空！")
	validation.MinSize(f.Issue.Title, 2).Message("标题长度不能小于2！")
	validation.MaxSize(f.Issue.Title, 100).Message("标题长度不能大于100！")

	validation.Required(f.Issue.Body).Message("正文不能为空！")
	validation.MinSize(f.Issue.Body, 2).Message("正文长度不能小于2！")
	//validation.MaxSize(f.Issue.Body, 4000).Message("正文长度不能大于4000！")

	validation.Required(f.Issue.SubmitUserId).Message("提报人不能为空！")
	validation.Required(f.Issue.EnterUserId).Message("录入人不能为空！")

	return validation.HasErrors() == false
}

func (c TrackIssue) Add() revel.Result {
	c.RenderArgs["issue_id"] = 0
	return c.RenderTemplate("track/issue/issue_detail.html")
}

func (c TrackIssue) Detail(issue_id int64) revel.Result {
	c.RenderArgs["issue_id"] = issue_id
	return c.RenderTemplate("track/issue/issue_detail.html")
}

func (c TrackIssue) DetailData() revel.Result {
	db_session := c.DbSession

	var issue_id int64
	c.Params.Bind(&issue_id, "id")

	var issue models.IssueInfo
	if issue_id != 0 {
		has, err := db_session.Id(issue_id).Get(&issue)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的软件缺陷不存在！"})
		}
	}

	if issue_id == 0 {
		var login_user = core.GetLoginUser(c.Session)
		issue.SubmitUserId = login_user.UserId

		issue.TypeId = models.ISSUE_TYPE_BUG
	}

	var issue_type_list = make([]models.IssueTypeInfo, 0)
	err := db_session.Find(&issue_type_list)
	core.HandleError(err)

	var issue_cate_list = make([]models.IssueCateInfo, 0)
	err = db_session.Find(&issue_cate_list)
	core.HandleError(err)

	var user_list = make([]auth_models.UserInfo, 0)
	err = db_session.Find(&user_list)
	core.HandleError(err)

	var detail_form IssueDetailForm
	detail_form.Issue = issue
	detail_form.IssueTypeList = issue_type_list
	detail_form.IssueCateList = issue_cate_list
	detail_form.UserList = user_list

	return c.RenderJson(detail_form)
}

func (c TrackIssue) Save() revel.Result {
	db_session := c.DbSession

	var detail_form IssueDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	issue := detail_form.Issue

	var submit_user auth_models.UserInfo
	_, err = db_session.Id(issue.SubmitUserId).Get(&submit_user)
	core.HandleError(err)
	issue.SubmitUserName = submit_user.NickName

	if detail_form.IsCreate() {
		issue.StatusId = models.ISSUE_STATUS_YTJ

		var login_user = core.GetLoginUser(c.Session)
		issue.EnterUserId = login_user.UserId
		issue.EnterUserName = login_user.NickName

		_, err = db_session.Insert(&issue)
		core.HandleError(err)

	} else {
		affected, err := db_session.Id(issue.Id).Update(&issue)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success":  true,
		"message":  "数据保存成功!",
		"issue_id": issue.Id,
	})
}

func (c TrackIssue) Delete() revel.Result {
	db_session := c.DbSession

	issue_id_list := make([]int64, 0)
	c.Params.Bind(&issue_id_list, "id_list")

	if len(issue_id_list) == 0 || len(issue_id_list) > 1 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择一条数据！",
		})
	}

	var issue_id = issue_id_list[0]

	var issue models.IssueInfo
	has, err := db_session.Id(issue_id).Get(&issue)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("id为 %d 的数据不存在！", issue_id),
		})
	}

	if issue.StatusId != models.ISSUE_STATUS_YTJ {
		if has == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "只有 已提交 状态的数据才能被删除！",
			})
		}
	}

	_, err = db_session.Id(issue_id).Delete(new(models.IssueInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据删除成功!"})
}

func (c TrackIssue) AddRedirect() revel.Result {
	var added_issue_id int64
	c.Params.Bind(&added_issue_id, "added_issue_id")

	return c.Redirect(routes.TrackIssue.Detail(added_issue_id))
}

func (c TrackIssue) DeleteRedirect() revel.Result {
	return c.Redirect(routes.TrackIssue.Add())
}
