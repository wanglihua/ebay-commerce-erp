package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models/ddpz_models"
)

type ErpEnvelopeFeeRate struct {
	*revel.Controller
}

func (c ErpEnvelopeFeeRate) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ENVELOPE_FEE_RATE, c.Session)

	return c.RenderTemplate("erp/ddpz/envelope_fee_rate/envelope_fee_rate_index.html")
}

func (c ErpEnvelopeFeeRate) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ENVELOPE_FEE_RATE, c.Session)

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

	envelope_fee_rate_list := make([]ddpz_models.EnvelopeFeeRateInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&envelope_fee_rate_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ddpz_models.EnvelopeFeeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  envelope_fee_rate_list,
		Total: count,
	})
}

type EnvelopeFeeRateDetailForm struct {
	EnvelopeFeeRate ddpz_models.EnvelopeFeeRateInfo `json:"envelope_fee_rate"`
}

func (f *EnvelopeFeeRateDetailForm) IsCreate() bool {
	return f.EnvelopeFeeRate.Id == 0
}

func (f *EnvelopeFeeRateDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.EnvelopeFeeRate.Name).Message("信封名称不能为空！")
	validation.Required(f.EnvelopeFeeRate.Length).Message("信封长度（cm）不能为空！")
	validation.Required(f.EnvelopeFeeRate.Width).Message("信封宽度（cm）不能为空！")
	validation.Required(f.EnvelopeFeeRate.Height).Message("信封高度（cm）不能为空！")
	validation.Required(f.EnvelopeFeeRate.Weight).Message("信封重量（kg）不能为空！")
	validation.Required(f.EnvelopeFeeRate.Fee).Message("费用（澳币）不能为空！")

	return validation.HasErrors() == false
}

func (c ErpEnvelopeFeeRate) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ENVELOPE_FEE_RATE, c.Session)

	db_session := c.DbSession

	var envelope_fee_rate_id int64
	c.Params.Bind(&envelope_fee_rate_id, "id")

	var envelope_fee_rate ddpz_models.EnvelopeFeeRateInfo
	if envelope_fee_rate_id != 0 {
		has, err := db_session.Id(envelope_fee_rate_id).Get(&envelope_fee_rate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的信封费率不存在！"})
		}
	}

	var detail_form EnvelopeFeeRateDetailForm
	detail_form.EnvelopeFeeRate = envelope_fee_rate

	return c.RenderJson(detail_form)
}

func (c ErpEnvelopeFeeRate) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ENVELOPE_FEE_RATE, c.Session)

	db_session := c.DbSession

	var detail_form EnvelopeFeeRateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	envelope_fee_rate := detail_form.EnvelopeFeeRate

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&envelope_fee_rate)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(envelope_fee_rate.Id).Update(&envelope_fee_rate)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpEnvelopeFeeRate) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ENVELOPE_FEE_RATE, c.Session)

	db_session := c.DbSession

	envelope_fee_rate_id_list := make([]int64, 0)
	c.Params.Bind(&envelope_fee_rate_id_list, "id_list")

	affected, err := db_session.In("id", envelope_fee_rate_id_list).Delete(new(ddpz_models.EnvelopeFeeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
