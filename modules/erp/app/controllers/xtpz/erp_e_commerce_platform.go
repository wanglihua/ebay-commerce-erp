package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpECommercePlatform struct {
	*revel.Controller
}

func (c ErpECommercePlatform) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_ECOMMERCE_PLATFORM, c.Session)

	return c.RenderTemplate("erp/xtpz/e_commerce_platform/e_commerce_platform_index.html")
}

func (c ErpECommercePlatform) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_ECOMMERCE_PLATFORM, c.Session)

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

	e_commerce_platform_list := make([]models.ECommercePlatformInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&e_commerce_platform_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ECommercePlatformInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  e_commerce_platform_list,
		Total: count,
	})
}

/*

type ECommercePlatformDetailForm struct {
	ECommercePlatform models.ECommercePlatformInfo `json:"platform"`
}

func (f *ECommercePlatformDetailForm) IsCreate() bool {
	return f.ECommercePlatform.Id == 0
}

func (f *ECommercePlatformDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.ECommercePlatform.Name).Message("名称不能为空！")
	if f.ECommercePlatform.Name != "" {
		validation.MinSize(f.ECommercePlatform.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ECommercePlatform.Name != "" {
		validation.MaxSize(f.ECommercePlatform.Name, 100).Message("名称长度不能大于100！")
	}

	if f.ECommercePlatform.Desc.Valid == true {
		validation.MaxSize(f.ECommercePlatform.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpECommercePlatform) DetailData() revel.Result {
	db_session := c.DbSession

	var e_commerce_platform_id int64
	c.Params.Bind(&e_commerce_platform_id, "id")

	var e_commerce_platform models.ECommercePlatformInfo
	if e_commerce_platform_id != 0 {
		has, err := db_session.Id(e_commerce_platform_id).Get(&e_commerce_platform)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的电商平台不存在！"})
		}
	}

	var detail_form ECommercePlatformDetailForm
	detail_form.ECommercePlatform = e_commerce_platform

	return c.RenderJson(detail_form)
}

func (c ErpECommercePlatform) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ECommercePlatformDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	e_commerce_platform := detail_form.ECommercePlatform

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&e_commerce_platform)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(e_commerce_platform.Id).Update(&e_commerce_platform)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpECommercePlatform) Delete() revel.Result {
	db_session := c.DbSession

	e_commerce_platform_id_list := make([]int64, 0)
	c.Params.Bind(&e_commerce_platform_id_list, "id_list")

	affected, err := db_session.In("id", e_commerce_platform_id_list).Delete(new(models.ECommercePlatformInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

*/
