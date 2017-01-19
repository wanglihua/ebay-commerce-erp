package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
)

// 弃用

type ErpLogisticsNumberBatch struct {
	*revel.Controller
}

func (c ErpLogisticsNumberBatch) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_number_batch/logistics_number_batch_index.html")
}

type ErpLogisticsNumberBatchViewItem struct {
	models.LogisticsNumberBatchInfo `xorm:"extends" json:"b"`
	models.LogisticsModeInfo        `xorm:"extends" json:"m"`
	models.StatusInfo               `xorm:"extends" json:"s"`
}

func (c ErpLogisticsNumberBatch) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT b.*, m.*, s.*
		FROM erp_logistics_number_batch b
		    INNER JOIN erp_logistics_mode m ON b.mode_id = m.id
		    INNER JOIN erp_status s ON b.status_id = s.id
	*/

	query := db_session.
		Select(" b.*, m.*, s.*").
		Table("erp_logistics_number_batch").Alias("b").
		Join("INNER", []string{"erp_logistics_mode", "m"}, "b.mode_id = m.id").
		Join("INNER", []string{"erp_status", "s"}, "b.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("b.id")
	}

	logistics_number_batch_list := make([]ErpLogisticsNumberBatchViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&logistics_number_batch_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpLogisticsNumberBatchViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_number_batch_list,
		Total: count,
	})
}

type LogisticsNumberBatchDetailForm struct {
	LogisticsNumberBatch models.LogisticsNumberBatchInfo `json:"number_batch"`
	LogisticsModeList    []models.LogisticsModeInfo      `json:"logistics_mode_list"`
	StatusList           []models.StatusInfo             `json:"status_list"`
}

func (f *LogisticsNumberBatchDetailForm) IsCreate() bool {
	return f.LogisticsNumberBatch.Id == 0
}

func (f *LogisticsNumberBatchDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsNumberBatch.ModeId).Message("物流方式不能为空！")

	validation.Required(f.LogisticsNumberBatch.Prefix).Message("单号前缀不能为空！")

	validation.Required(f.LogisticsNumberBatch.Suffix).Message("单号后缀不能为空！")

	validation.Required(f.LogisticsNumberBatch.BeginNum).Message("开始号段不能为空！")

	validation.Required(f.LogisticsNumberBatch.EndNum).Message("结束号段不能为空！")

	validation.Required(f.LogisticsNumberBatch.CurNum).Message("已使用号段不能为空！")

	validation.Required(f.LogisticsNumberBatch.Remind).Message("剩余号段不能为空！")

	validation.Required(f.LogisticsNumberBatch.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpLogisticsNumberBatch) DetailData() revel.Result {
	db_session := c.DbSession

	var logistics_number_batch_id int64
	c.Params.Bind(&logistics_number_batch_id, "id")

	var logistics_number_batch models.LogisticsNumberBatchInfo
	if logistics_number_batch_id != 0 {
		has, err := db_session.Id(logistics_number_batch_id).Get(&logistics_number_batch)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的订单运单号不存在！"})
		}
	}

	var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
	err := db_session.Find(&logistics_mode_list)
	core.HandleError(err)

	status_list := status_service.GetStatusList(db_session, models.STATUS_TYPE_LOGISTICS_NUMBER_BATCH)

	var detail_form LogisticsNumberBatchDetailForm
	detail_form.LogisticsNumberBatch = logistics_number_batch
	detail_form.LogisticsModeList = logistics_mode_list
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsNumberBatch) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsNumberBatchDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_number_batch := detail_form.LogisticsNumberBatch

	logistics_number_batch.Remind = logistics_number_batch.EndNum - logistics_number_batch.BeginNum + 1

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&logistics_number_batch)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(logistics_number_batch.Id).Update(&logistics_number_batch)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsNumberBatch) Delete() revel.Result {
	db_session := c.DbSession

	logistics_number_batch_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_number_batch_id_list, "id_list")

	affected, err := db_session.In("id", logistics_number_batch_id_list).Delete(new(models.LogisticsNumberBatchInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
