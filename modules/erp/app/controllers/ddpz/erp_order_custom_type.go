package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
	"matrix/modules/erp/models/ddpz_models"
)

type ErpOrderCustomType struct {
	*revel.Controller
}

func (c ErpOrderCustomType) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/order_custom_type/order_custom_type_index.html")
}

type ErpOrderCustomTypeViewItem struct {
	ddpz_models.OrderCustomTypeInfo `xorm:"extends" json:"t"`
	models.StatusInfo          `xorm:"extends" json:"s"`
}

func (c ErpOrderCustomType) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT t.*, s.*
		FROM erp_order_custom_type t
		    INNER JOIN erp_status s ON t.status_id = s.id
	*/

	query := db_session.
		Select("t.*, s.*").
		Table("erp_order_custom_type").Alias("t").
		Join("INNER", []string{"erp_status", "s"}, "t.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	order_custom_type_list := make([]ErpOrderCustomTypeViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_custom_type_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpOrderCustomTypeViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_custom_type_list,
		Total: count,
	})
}

type OrderCustomTypeDetailForm struct {
	OrderCustomType ddpz_models.OrderCustomTypeInfo `json:"order_type"`
	StatusList      []models.StatusInfo        `json:"status_list"`
}

func (f *OrderCustomTypeDetailForm) IsCreate() bool {
	return f.OrderCustomType.Id == 0
}

func (f *OrderCustomTypeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.OrderCustomType.Code).Message("编号不能为空！")

	validation.Required(f.OrderCustomType.Name).Message("名称不能为空！")

	validation.Required(f.OrderCustomType.Color).Message("颜色不能为空！")

	if f.OrderCustomType.Desc.Valid == true {
		validation.MaxSize(f.OrderCustomType.Desc.String, 500).Message("描述长度不能大于500！")
	}

	validation.Required(f.OrderCustomType.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpOrderCustomType) DetailData() revel.Result {
	db_session := c.DbSession

	var order_custom_type_id int64
	c.Params.Bind(&order_custom_type_id, "id")

	var order_custom_type ddpz_models.OrderCustomTypeInfo
	if order_custom_type_id != 0 {
		has, err := db_session.Id(order_custom_type_id).Get(&order_custom_type)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的订单自定义类型不存在！"})
		}
	}

	status_list := status_service.Get_status_list(db_session, models.STATUS_TYPE_ORDER_CUSTOM_TYPE)

	var detail_form OrderCustomTypeDetailForm
	detail_form.OrderCustomType = order_custom_type
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpOrderCustomType) Save() revel.Result {
	db_session := c.DbSession

	var detail_form OrderCustomTypeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	order_custom_type := detail_form.OrderCustomType

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&order_custom_type)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(order_custom_type.Id).Update(&order_custom_type)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpOrderCustomType) Delete() revel.Result {
	db_session := c.DbSession

	order_custom_type_id_list := make([]int64, 0)
	c.Params.Bind(&order_custom_type_id_list, "id_list")

	affected, err := db_session.In("id", order_custom_type_id_list).Delete(new(ddpz_models.OrderCustomTypeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
