
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrderLog struct {
	*revel.Controller
}

func (c ErpOrderLog) Index() revel.Result {
	return c.RenderTemplate("erp/order_log/order_log_index.html")
}

func (c ErpOrderLog) ListData() revel.Result {
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

	order_log_list := make([]models.OrderLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.OrderLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_log_list,
		Total: count,
	})
}

type OrderLogDetailForm struct {
	OrderLog models.OrderLogInfo `json:"order_log"`
}

func (f *OrderLogDetailForm) IsCreate() bool {
	return f.OrderLog.Id == 0
}

func (f *OrderLogDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.OrderLog.OrderId).Message("订单Id不能为空！")

	validation.Required(f.OrderLog.OrderCode).Message("订单编号不能为空！")

	validation.Required(f.OrderLog.TransCode).Message("交易编号不能为空！")

	validation.Required(f.OrderLog.UserId).Message("用户Id不能为空！")

	validation.Required(f.OrderLog.LoginName).Message("操作工号不能为空！")

	validation.Required(f.OrderLog.Operation).Message("操作描述不能为空！")

	return validation.HasErrors() == false
}

func (c ErpOrderLog) DetailData() revel.Result {
	db_session := c.DbSession

	var order_log_id int64
	c.Params.Bind(&order_log_id, "id")

	var order_log models.OrderLogInfo
	if order_log_id != 0 {
		has, err := db_session.Id(order_log_id).Get(&order_log)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的订单修改日志不存在！"})
		}
	}

	var detail_form OrderLogDetailForm
	detail_form.OrderLog = order_log

	return c.RenderJson(detail_form)
}

func (c ErpOrderLog) Save() revel.Result {
	db_session := c.DbSession

	var detail_form OrderLogDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	order_log := detail_form.OrderLog

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&order_log)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(order_log.Id).Update(&order_log)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpOrderLog) Delete() revel.Result {
	db_session := c.DbSession

	order_log_id_list := make([]int64, 0)
	c.Params.Bind(&order_log_id_list, "id_list")

	affected, err := db_session.In("id", order_log_id_list).Delete(new(models.OrderLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

