package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpHsCodeCnReportElem struct {
	*revel.Controller
}

func (c ErpHsCodeCnReportElem) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM, c.Session)

	return c.RenderTemplate("erp/hgbg/hs_code_cn_report_elem/hs_code_cn_report_elem_index.html")
}

func (c ErpHsCodeCnReportElem) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM, c.Session)

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

	hs_code_cn_report_elem_list := make([]models.HsCodeCnReportElemInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&hs_code_cn_report_elem_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.HsCodeCnReportElemInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  hs_code_cn_report_elem_list,
		Total: count,
	})
}

type HsCodeCnReportElemDetailForm struct {
	HsCodeCnReportElem models.HsCodeCnReportElemInfo `json:"hs_code_cn"`
}

func (f *HsCodeCnReportElemDetailForm) IsCreate() bool {
	return f.HsCodeCnReportElem.Id == 0
}

func (f *HsCodeCnReportElemDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.HsCodeCnReportElem.HsCode).Message("海关编码不能为空！")
	if f.HsCodeCnReportElem.HsCode != "" {
		validation.MinSize(f.HsCodeCnReportElem.HsCode, 1).Message("海关编码长度不能小于1！")
	}
	if f.HsCodeCnReportElem.HsCode != "" {
		validation.MaxSize(f.HsCodeCnReportElem.HsCode, 10).Message("海关编码长度不能大于10！")
	}

	validation.Required(f.HsCodeCnReportElem.CnName).Message("中文品名不能为空！")

	validation.Required(f.HsCodeCnReportElem.EnName).Message("英文品名不能为空！")

	return validation.HasErrors() == false
}

func (c ErpHsCodeCnReportElem) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM, c.Session)

	db_session := c.DbSession

	var hs_code_cn_report_elem_id int64
	c.Params.Bind(&hs_code_cn_report_elem_id, "id")

	var hs_code_cn_report_elem models.HsCodeCnReportElemInfo
	if hs_code_cn_report_elem_id != 0 {
		has, err := db_session.Id(hs_code_cn_report_elem_id).Get(&hs_code_cn_report_elem)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的中国海关编码不存在！"})
		}
	}

	var detail_form HsCodeCnReportElemDetailForm
	detail_form.HsCodeCnReportElem = hs_code_cn_report_elem

	return c.RenderJson(detail_form)
}

func (c ErpHsCodeCnReportElem) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM, c.Session)

	db_session := c.DbSession

	var detail_form HsCodeCnReportElemDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	hs_code_cn_report_elem := detail_form.HsCodeCnReportElem

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&hs_code_cn_report_elem)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(hs_code_cn_report_elem.Id).Update(&hs_code_cn_report_elem)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpHsCodeCnReportElem) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM, c.Session)

	db_session := c.DbSession

	hs_code_cn_report_elem_id_list := make([]int64, 0)
	c.Params.Bind(&hs_code_cn_report_elem_id_list, "id_list")

	affected, err := db_session.In("id", hs_code_cn_report_elem_id_list).Delete(new(models.HsCodeCnReportElemInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
