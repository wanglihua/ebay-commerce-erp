package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"strings"
)

type ErpCityPostcode struct {
	*revel.Controller
}

func (c ErpCityPostcode) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_CITY_POSTCODE, c.Session)

	c.RenderArgs["city_id_mel"] = ddpz_models.CITY_ID_MEL
	c.RenderArgs["city_id_bne"] = ddpz_models.CITY_ID_BNE
	c.RenderArgs["city_id_syd"] = ddpz_models.CITY_ID_SYD

	return c.RenderTemplate("erp/ddpz/city_postcode/city_postcode_index.html")
}

func (c ErpCityPostcode) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_CITY_POSTCODE, c.Session)

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

	city_postcode_list := make([]ddpz_models.CityPostcodeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&city_postcode_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ddpz_models.CityPostcodeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  city_postcode_list,
		Total: count,
	})
}

type CityPostcodeDetailForm struct {
	CityPostcode ddpz_models.CityPostcodeInfo `json:"city_postcode"`
}

func (f *CityPostcodeDetailForm) IsCreate() bool {
	return f.CityPostcode.Id == 0
}

func (f *CityPostcodeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.CityPostcode.CityId).Message("CityId不能为空！")

	validation.Required(f.CityPostcode.Postcode).Message("Postcode不能为空！")

	return validation.HasErrors() == false
}

func (c ErpCityPostcode) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_CITY_POSTCODE, c.Session)

	db_session := c.DbSession

	var city_postcode_id int64
	c.Params.Bind(&city_postcode_id, "id")

	var city_postcode ddpz_models.CityPostcodeInfo
	if city_postcode_id != 0 {
		has, err := db_session.Id(city_postcode_id).Get(&city_postcode)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的城市邮编不存在！"})
		}
	}

	var detail_form CityPostcodeDetailForm
	detail_form.CityPostcode = city_postcode

	return c.RenderJson(detail_form)
}

func (c ErpCityPostcode) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_CITY_POSTCODE, c.Session)

	db_session := c.DbSession

	var detail_form CityPostcodeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	city_postcode_from_ui := detail_form.CityPostcode

	if city_postcode_from_ui.CityId == ddpz_models.CITY_ID_MEL {
		city_postcode_from_ui.CityName = "MEL"
	}

	if city_postcode_from_ui.CityId == ddpz_models.CITY_ID_BNE {
		city_postcode_from_ui.CityName = "BNE"
	}

	if city_postcode_from_ui.CityId == ddpz_models.CITY_ID_SYD {
		city_postcode_from_ui.CityName = "SYD"
	}

	var city_id = city_postcode_from_ui.CityId
	var city_name = city_postcode_from_ui.CityName
	var postcode_list = strings.Split(city_postcode_from_ui.Postcode, ",")

	for _, postcode := range postcode_list {
		postcode = strings.TrimSpace(postcode)

		if postcode != "" {
			var city_postcode_in_db ddpz_models.CityPostcodeInfo
			has, err := db_session.Where("city_id = ? AND postcode = ?", city_id, postcode).Get(&city_postcode_in_db)
			core.HandleError(err)

			if has == false {
				var city_postcode ddpz_models.CityPostcodeInfo
				city_postcode.CityId = city_id
				city_postcode.CityName = city_name
				city_postcode.Postcode = postcode

				_, err = db_session.Insert(&city_postcode)
				core.HandleError(err)
			} else {
				city_postcode_in_db.CityName = city_name

				_, err = db_session.Id(city_postcode_in_db.Id).Update(&city_postcode_in_db)
				core.HandleError(err)
			}
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "邮编新增成功!"})
}

func (c ErpCityPostcode) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_CITY_POSTCODE, c.Session)

	db_session := c.DbSession

	city_postcode_id_list := make([]int64, 0)
	c.Params.Bind(&city_postcode_id_list, "id_list")

	affected, err := db_session.In("id", city_postcode_id_list).Delete(new(ddpz_models.CityPostcodeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
