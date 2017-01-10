
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpLogisticsNumberLog struct {
	*revel.Controller
}

func (c ErpLogisticsNumberLog) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_number_log/logistics_number_log_index.html")
}

func (c ErpLogisticsNumberLog) ListData() revel.Result {
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

	logistics_number_log_list := make([]models.LogisticsNumberLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&logistics_number_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.LogisticsNumberLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_number_log_list,
		Total: count,
	})
}

type LogisticsNumberLogDetailForm struct {
	LogisticsNumberLog models.LogisticsNumberLogInfo `json:"number_log"`
}

func (f *LogisticsNumberLogDetailForm) IsCreate() bool {
	return f.LogisticsNumberLog.Id == 0
}

func (f *LogisticsNumberLogDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.LogisticsNumberLog.NumId).Message("运单号Id不能为空！")

	validation.Required(f.LogisticsNumberLog.ModeName).Message("货运方式不能为空！")

	validation.Required(f.LogisticsNumberLog.NumRange).Message("号段不能为空！")

	if f.LogisticsNumberLog.Desc.Valid == true {
		validation.MaxSize(f.LogisticsNumberLog.Desc.String, 500).Message("描述长度不能大于500！")
	}

	validation.Required(f.LogisticsNumberLog.UserId).Message("用户Id不能为空！")

	validation.Required(f.LogisticsNumberLog.LoginName).Message("用户工号不能为空！")

	return validation.HasErrors() == false
}

func (c ErpLogisticsNumberLog) DetailData() revel.Result {
	db_session := c.DbSession

	var logistics_number_log_id int64
	c.Params.Bind(&logistics_number_log_id, "id")

	var logistics_number_log models.LogisticsNumberLogInfo
	if logistics_number_log_id != 0 {
		has, err := db_session.Id(logistics_number_log_id).Get(&logistics_number_log)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的运单号操作日志不存在！"})
		}
	}

	var detail_form LogisticsNumberLogDetailForm
	detail_form.LogisticsNumberLog = logistics_number_log

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsNumberLog) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsNumberLogDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_number_log := detail_form.LogisticsNumberLog

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&logistics_number_log)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(logistics_number_log.Id).Update(&logistics_number_log)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsNumberLog) Delete() revel.Result {
	db_session := c.DbSession

	logistics_number_log_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_number_log_id_list, "id_list")

	affected, err := db_session.In("id", logistics_number_log_id_list).Delete(new(models.LogisticsNumberLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

