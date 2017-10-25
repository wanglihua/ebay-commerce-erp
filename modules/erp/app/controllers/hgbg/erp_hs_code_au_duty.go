package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/hgbg_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpHsCodeAuDuty struct {
	*revel.Controller
}

func (c ErpHsCodeAuDuty) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_AU_DUTY, c.Session)

	return c.RenderTemplate("erp/hgbg/hs_code_au_duty/hs_code_au_duty_index.html")
}

func (c ErpHsCodeAuDuty) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_AU_DUTY, c.Session)

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

	hs_code_au_duty_list := make([]hgbg_models.HsCodeAuDutyInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&hs_code_au_duty_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(hgbg_models.HsCodeAuDutyInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  hs_code_au_duty_list,
		Total: count,
	})
}

type HsCodeAuDutyDetailForm struct {
	HsCodeAuDuty hgbg_models.HsCodeAuDutyInfo `json:"hs_code_au"`
}

func (f *HsCodeAuDutyDetailForm) IsCreate() bool {
	return f.HsCodeAuDuty.Id == 0
}

func (f *HsCodeAuDutyDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.HsCodeAuDuty.HsCode).Message("海关编码不能为空！")
	if f.HsCodeAuDuty.HsCode != "" {
		validation.MinSize(f.HsCodeAuDuty.HsCode, 1).Message("海关编码长度不能小于1！")
	}
	if f.HsCodeAuDuty.HsCode != "" {
		validation.MaxSize(f.HsCodeAuDuty.HsCode, 10).Message("海关编码长度不能大于10！")
	}

	validation.Required(f.HsCodeAuDuty.CnName).Message("中文品名不能为空！")

	validation.Required(f.HsCodeAuDuty.EnName).Message("英文品名不能为空！")

	validation.Required(f.HsCodeAuDuty.DutyRate).Message("关税税率不能为空！")

	return validation.HasErrors() == false
}

func (c ErpHsCodeAuDuty) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_AU_DUTY, c.Session)

	db_session := c.DbSession

	var hs_code_au_duty_id int64
	c.Params.Bind(&hs_code_au_duty_id, "id")

	var hs_code_au_duty hgbg_models.HsCodeAuDutyInfo
	if hs_code_au_duty_id != 0 {
		has, err := db_session.Id(hs_code_au_duty_id).Get(&hs_code_au_duty)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的澳洲海关编码不存在！"})
		}
	}

	var detail_form HsCodeAuDutyDetailForm
	detail_form.HsCodeAuDuty = hs_code_au_duty

	return c.RenderJson(detail_form)
}

func (c ErpHsCodeAuDuty) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_AU_DUTY, c.Session)

	db_session := c.DbSession

	var detail_form HsCodeAuDutyDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	hs_code_au_duty := detail_form.HsCodeAuDuty

	var affected int64
	if detail_form.IsCreate() {
		var hs_code hgbg_models.HsCodeAuDutyInfo
		has, err := db_session.Where("en_name = ?", hs_code_au_duty.EnName).Get(&hs_code)
		core.HandleError(err)

		if has == true {
			return c.RenderJson(core.JsonResult{Success: false, Message: "海关申报品名（英文名）应唯一！"})
		}

		affected, err = db_session.Insert(&hs_code_au_duty)
		core.HandleError(err)
	} else {
		var hs_code hgbg_models.HsCodeAuDutyInfo
		has, err := db_session.Where("id <> ? AND en_name = ?", hs_code_au_duty.Id, hs_code_au_duty.EnName).Get(&hs_code)
		core.HandleError(err)

		if has == true {
			return c.RenderJson(core.JsonResult{Success: false, Message: "海关申报品名（英文名）应唯一！"})
		}

		affected, err = db_session.Id(hs_code_au_duty.Id).Update(&hs_code_au_duty)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpHsCodeAuDuty) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_AU_DUTY, c.Session)

	db_session := c.DbSession

	hs_code_au_duty_id_list := make([]int64, 0)
	c.Params.Bind(&hs_code_au_duty_id_list, "id_list")

	affected, err := db_session.In("id", hs_code_au_duty_id_list).Delete(new(hgbg_models.HsCodeAuDutyInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
