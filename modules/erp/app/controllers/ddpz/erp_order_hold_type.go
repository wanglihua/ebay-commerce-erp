package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/setting"
)

type ErpOrderHoldType struct {
	*revel.Controller
}

func (c ErpOrderHoldType) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ORDER_HOLD_TYPE, c.Session)

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	return c.RenderTemplate("erp/ddpz/order_hold_type/order_hold_type_index.html")
}

func (c ErpOrderHoldType) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ORDER_HOLD_TYPE, c.Session)

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

	order_hold_type_list := make([]models.OrderHoldTypeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_hold_type_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.OrderHoldTypeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_hold_type_list,
		Total: count,
	})
}

type OrderHoldTypeDetailForm struct {
	OrderHoldType models.OrderHoldTypeInfo `json:"order_hold_type"`
}

func (f *OrderHoldTypeDetailForm) IsCreate() bool {
	return f.OrderHoldType.Id == 0
}

func (f *OrderHoldTypeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.OrderHoldType.Name).Message("名称不能为空！")

	validation.Required(f.OrderHoldType.Desc).Message("备注不能为空！")

	return validation.HasErrors() == false
}

func (c ErpOrderHoldType) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ORDER_HOLD_TYPE, c.Session)

	db_session := c.DbSession

	var order_hold_type_id int64
	c.Params.Bind(&order_hold_type_id, "id")

	var order_hold_type models.OrderHoldTypeInfo
	if order_hold_type_id != 0 {
		has, err := db_session.Id(order_hold_type_id).Get(&order_hold_type)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的订单Hold类型不存在！"})
		}
	}

	var detail_form OrderHoldTypeDetailForm
	detail_form.OrderHoldType = order_hold_type

	return c.RenderJson(detail_form)
}

func (c ErpOrderHoldType) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_ORDER_HOLD_TYPE, c.Session)

	db_session := c.DbSession

	var detail_form OrderHoldTypeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	order_hold_type := detail_form.OrderHoldType

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where("name = ?", order_hold_type.Name).Count(&order_hold_type)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&order_hold_type)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", order_hold_type.Id, order_hold_type.Name).Count(&order_hold_type)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Id(order_hold_type.Id).Update(&order_hold_type)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}
