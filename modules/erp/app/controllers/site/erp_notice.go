package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/xtgl_models"
)

type ErpNotice struct {
	*revel.Controller
}

func (c ErpNotice) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTGL_NOTICE, c.Session)

	return c.RenderTemplate("erp/site/notice/notice_index.html")
}

func (c ErpNotice) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTGL_NOTICE, c.Session)

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

	var affected int64
	if detail_form.IsCreate() {
		var login_user = core.GetLoginUser(c.Session)
		notice.UserId = login_user.UserId
		notice.LoginName = login_user.LoginName

		affected, err = db_session.Insert(&notice)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(notice.Id).Cols("title", "body").Update(&notice)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
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
