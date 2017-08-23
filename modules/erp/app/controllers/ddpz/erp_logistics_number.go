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

// 弃用

type ErpLogisticsNumber struct {
	*revel.Controller
}

func (c ErpLogisticsNumber) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_number/logistics_number_index.html")
}

type ErpLogisticsNumberViewItem struct {
	ddpz_models.LogisticsNumberInfo `xorm:"extends" json:"n"`
	ddpz_models.LogisticsModeInfo   `xorm:"extends" json:"m"`
	models.StatusInfo          `xorm:"extends" json:"s"`
}

func (c ErpLogisticsNumber) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT n.*, m.*, s.*
		FROM erp_logistics_number n
		    INNER JOIN erp_logistics_mode m ON n.mode_id = m.id
		    INNER JOIN erp_status s ON n.status_id = s.id
	*/

	query := db_session.
		Select(" n.*, m.*, s.*").
		Table("erp_logistics_number").Alias("n").
		Join("INNER", []string{"erp_logistics_mode", "m"}, "n.mode_id = m.id").
		Join("INNER", []string{"erp_status", "s"}, "n.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("n.id")
	}

	logistics_number_list := make([]ErpLogisticsNumberViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&logistics_number_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpLogisticsNumberViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_number_list,
		Total: count,
	})
}

type LogisticsNumberDetailForm struct {
	LogisticsNumber   ddpz_models.LogisticsNumberInfo `json:"logistics_number"`
	LogisticsModeList []ddpz_models.LogisticsModeInfo `json:"logistics_mode_list"`
	StatusList        []models.StatusInfo        `json:"status_list"`
}

func (f *LogisticsNumberDetailForm) IsCreate() bool {
	return f.LogisticsNumber.Id == 0
}

func (f *LogisticsNumberDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsNumber.Num).Message("物流单号不能为空！")
	validation.MaxSize(f.LogisticsNumber.Num, 100).Message("物流单号长度不能大于100！")

	validation.Required(f.LogisticsNumber.ModeId).Message("物流方式不能为空！")

	if f.LogisticsNumber.Desc.Valid == true {
		validation.MaxSize(f.LogisticsNumber.Desc.String, 500).Message("详情长度不能大于500！")
	}

	validation.Required(f.LogisticsNumber.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpLogisticsNumber) DetailData() revel.Result {
	db_session := c.DbSession

	var logistics_number_id int64
	c.Params.Bind(&logistics_number_id, "id")

	var logistics_number ddpz_models.LogisticsNumberInfo
	if logistics_number_id != 0 {
		has, err := db_session.Id(logistics_number_id).Get(&logistics_number)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的运单号不规则号段不存在！"})
		}
	}

	var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.Find(&logistics_mode_list)
	core.HandleError(err)

	status_list := status_service.Get_status_list(db_session, models.STATUS_TYPE_LOGISTICS_NUMBER_BATCH)

	var detail_form LogisticsNumberDetailForm
	detail_form.LogisticsNumber = logistics_number
	detail_form.LogisticsModeList = logistics_mode_list
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsNumber) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsNumberDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_number := detail_form.LogisticsNumber

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&logistics_number)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(logistics_number.Id).Update(&logistics_number)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsNumber) Delete() revel.Result {
	db_session := c.DbSession

	logistics_number_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_number_id_list, "id_list")

	affected, err := db_session.In("id", logistics_number_id_list).Delete(new(ddpz_models.LogisticsNumberInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
