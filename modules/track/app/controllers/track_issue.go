package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
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
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	issue_list := make([]models.IssueInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&issue_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.IssueInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_list,
		Total: count,
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
	validation.Required(f.Issue.CateId).Message("类别id不能为空！")

	validation.Required(f.Issue.TypeId).Message("类型id不能为空！")

	validation.Required(f.Issue.Title).Message("标题不能为空！")
	if f.Issue.Title != "" {
		validation.MinSize(f.Issue.Title, 2).Message("标题长度不能小于2！")
	}
	if f.Issue.Title != "" {
		validation.MaxSize(f.Issue.Title, 100).Message("标题长度不能大于100！")
	}

	validation.Required(f.Issue.Body).Message("正文不能为空！")
	if f.Issue.Body != "" {
		validation.MinSize(f.Issue.Body, 2).Message("正文长度不能小于2！")
	}
	if f.Issue.Body != "" {
		validation.MaxSize(f.Issue.Body, 4000).Message("正文长度不能大于4000！")
	}

	validation.Required(f.Issue.SubmitUserId).Message("提报人用户id不能为空！")

	validation.Required(f.Issue.SubmitUserCode).Message("提报人工号不能为空！")

	validation.Required(f.Issue.EnterUserId).Message("录入人用户id不能为空！")

	validation.Required(f.Issue.EnterUserCode).Message("录入人工号不能为空！")

	validation.Required(f.Issue.StatusId).Message("状态id不能为空！")

	validation.Required(f.Issue.Create_time).Message("Create_time不能为空！")

	validation.Required(f.Issue.Update_time).Message("Update_time不能为空！")

	return validation.HasErrors() == false
}

func (c TrackIssue) Detail() revel.Result {
	var issue_id int64
	c.Params.Bind(&issue_id, "id")

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

	if detail_form.IsCreate() {
		_, err = db_session.Insert(&issue)
		core.HandleError(err)

		issue.StatusId = models.ISSUE_STATUS_YTJ
	} else {
		affected, err := db_session.Id(issue.Id).Update(&issue)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})
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
