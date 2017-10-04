package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/xtgl_models"
	"html/template"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/service/notice_read_service"
)

type ErpNotice struct {
	*revel.Controller
}

func (c ErpNotice) GetNoticeDataForHomePage() revel.Result {
	// 这里没有权限限制

	var db_session = c.DbSession

	/*
SELECT n.*, u.*
FROM erp_notice n
    INNER JOIN auth_user u ON n.user_id = u.id
ORDER BY n.list_top DESC, n.create_time DESC
LIMIT 50;
	 */

	type notice_view_item_t struct {
		Notice xtgl_models.NoticeInfo `xorm:"extends" json:"n"`
		User   auth_models.UserInfo   `xorm:"extends" json:"u"`
	}

	var sql = `
SELECT n.*, u.*
FROM erp_notice n
    INNER JOIN auth_user u ON n.user_id = u.id
ORDER BY n.list_top DESC, n.create_time DESC
LIMIT 50;
	`
	var notice_list = make([]notice_view_item_t, 0)
	core.HandleError(db_session.Sql(sql).Find(&notice_list))

	var login_user = core.GetLoginUser(c.Session)
	var notice_read_list = make([]xtgl_models.NoticeReadInfo, 0)
	core.HandleError(db_session.Where("user_id = ?", login_user.UserId).Find(&notice_read_list))

	var is_notice_read = func(notice_id int64) bool {
		for _, notice_read := range notice_read_list {
			if notice_read.NoticeId == notice_id {
				return true
			}
		}

		return false
	}

	for index, notice := range notice_list {
		notice_list[index].Notice.IsRead = is_notice_read(notice.Notice.Id)
	}

	return c.RenderJson(map[string]interface{}{
		"notice_list": notice_list,
	})
}

func (c ErpNotice) GetNoticeListDataForHomePage() revel.Result {
	// 这里没有权限限制

	var db_session = c.DbSession

	_, _, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT n.*, u.*
FROM erp_notice n
    INNER JOIN auth_user u ON n.user_id = u.id
ORDER BY n.list_top DESC, n.create_time DESC
LIMIT 50;
	 */

	type notice_view_item_t struct {
		Notice xtgl_models.NoticeInfo `xorm:"extends" json:"n"`
		User   auth_models.UserInfo   `xorm:"extends" json:"u"`
	}

	var sb = core.NewSqlBuilder()
	sb.Select("n.*, u.*")
	sb.Count("count(n.id) AS count")
	sb.From(`erp_notice n
    INNER JOIN auth_user u ON n.user_id = u.id`)
	sb.OrderBy(`n.list_top DESC, n.create_time DESC`)
	sb.Limit(offset, limit)

	var notice_list = make([]notice_view_item_t, 0)
	core.HandleError(db_session.Sql(sb.SqlList()).Find(&notice_list))

	var count core.Count
	_, err := db_session.Sql(sb.SqlCount()).Get(&count)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)
	var notice_read_list = make([]xtgl_models.NoticeReadInfo, 0)
	core.HandleError(db_session.Where("user_id = ?", login_user.UserId).Find(&notice_read_list))

	var is_notice_read = func(notice_id int64) bool {
		for _, notice_read := range notice_read_list {
			if notice_read.NoticeId == notice_id {
				return true
			}
		}

		return false
	}

	for index, notice := range notice_list {
		notice_list[index].Notice.IsRead = is_notice_read(notice.Notice.Id)
	}

	return c.RenderJson(core.GridResult{
		Data:  notice_list,
		Total: count.Count,
	})
}

func (c ErpNotice) ViewDetail(id int64) revel.Result {
	var notice_id = id

	var db_session = c.DbSession

	var notice xtgl_models.NoticeInfo
	_, err := db_session.Id(notice_id).Get(&notice)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)
	var login_user_id = login_user.UserId

	// 将当前用户的
	notice_read_service.MarkNoticeUserRead(login_user_id, notice_id, db_session)

	var notice_read_user_list = notice_read_service.GetNoticeReadUserList(notice_id, db_session)

	c.RenderArgs["notice"] = notice
	c.RenderArgs["notice_body"] = template.HTML(notice.Body)
	c.RenderArgs["notice_read_user_list"] = notice_read_user_list
	c.RenderArgs["notice_read_user_list_max_index"] = len(notice_read_user_list) - 1
	return c.RenderTemplate("erp/xtgl/notice/notice_view_detail.html")
}

func (c ErpNotice) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTGL_NOTICE, c.Session)

	return c.RenderTemplate("erp/xtgl/notice/notice_index.html")
}

func (c ErpNotice) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTGL_NOTICE, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(" list_top desc, " + order)
	} else {
		data_query = *data_query.OrderBy("list_top desc, id desc")
	}

	notice_list := make([]xtgl_models.NoticeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&notice_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(xtgl_models.NoticeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  notice_list,
		Total: count,
	})
}

type NoticeDetailForm struct {
	Notice xtgl_models.NoticeInfo `json:"notice"`
}

func (f *NoticeDetailForm) IsCreate() bool {
	return f.Notice.Id == 0
}

func (f *NoticeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Notice.Title).Message("标题不能为空！")
	validation.Required(f.Notice.Body).Message("内容不能为空！")

	return validation.HasErrors() == false
}

func (c ErpNotice) Detail(id int64) revel.Result {
	var notice_id = id

	c.RenderArgs["notice_id"] = notice_id
	return c.RenderTemplate("erp/xtgl/notice/notice_detail.html")
}

func (c ErpNotice) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTGL_NOTICE, c.Session)

	db_session := c.DbSession

	var notice_id int64
	c.Params.Bind(&notice_id, "id")

	var notice xtgl_models.NoticeInfo
	if notice_id != 0 {
		has, err := db_session.Id(notice_id).Get(&notice)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的通知公告不存在！"})
		}
	}

	var detail_form NoticeDetailForm
	detail_form.Notice = notice

	return c.RenderJson(detail_form)
}

func (c ErpNotice) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTGL_NOTICE, c.Session)

	db_session := c.DbSession

	var detail_form NoticeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	notice := detail_form.Notice

	var is_create = detail_form.IsCreate()

	var affected int64
	if is_create == true {
		var login_user = core.GetLoginUser(c.Session)
		notice.UserId = login_user.UserId
		notice.LoginName = login_user.LoginName

		affected, err = db_session.Insert(&notice)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(notice.Id).Cols("title", "body", "list_top", "important").Update(&notice)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	//return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
	return c.RenderJson(map[string]interface{}{
		"success":   true,
		"is_create": is_create,
		"notice_id": notice.Id,
		"message":   "数据保存成功！",
	})
}

func (c ErpNotice) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTGL_NOTICE, c.Session)

	db_session := c.DbSession

	notice_id_list := make([]int64, 0)
	c.Params.Bind(&notice_id_list, "id_list")

	affected, err := db_session.In("id", notice_id_list).Delete(new(xtgl_models.NoticeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
