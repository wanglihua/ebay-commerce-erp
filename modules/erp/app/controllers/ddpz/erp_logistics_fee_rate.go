package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpLogisticsFeeRate struct {
	*revel.Controller
}

func (c ErpLogisticsFeeRate) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_fee_rate/logistics_fee_rate_index.html")
}

func (c ErpLogisticsFeeRate) ListData() revel.Result {
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

	logistics_fee_rate_list := make([]models.LogisticsFeeRateInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&logistics_fee_rate_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.LogisticsFeeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_fee_rate_list,
		Total: count,
	})
}

type LogisticsFeeRateDetailForm struct {
	LogisticsFeeRate models.LogisticsFeeRateInfo `json:"logistics_fee_rate"`
}

func (f *LogisticsFeeRateDetailForm) IsCreate() bool {
	return f.LogisticsFeeRate.Id == 0
}

func (f *LogisticsFeeRateDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsFeeRate.ModeId).Message("物流公司不能为空！")

	validation.Required(f.LogisticsFeeRate.PostCode).Message("邮编不能为空！")
	if f.LogisticsFeeRate.PostCode != "" {
		validation.MinSize(f.LogisticsFeeRate.PostCode, 2).Message("邮编长度不能小于2！")
	}
	if f.LogisticsFeeRate.PostCode != "" {
		validation.MaxSize(f.LogisticsFeeRate.PostCode, 100).Message("邮编长度不能大于100！")
	}

	validation.Required(f.LogisticsFeeRate.Qbj).Message("起步价不能为空！")

	validation.Required(f.LogisticsFeeRate.Qbjzl).Message("起步价重量不能为空！")

	validation.Required(f.LogisticsFeeRate.Zdsf).Message("最低收费不能为空！")

	validation.Required(f.LogisticsFeeRate.Xzjt).Message("续重阶梯不能为空！")

	validation.Required(f.LogisticsFeeRate.Xzjtsf).Message("续重阶梯收费不能为空！")

	validation.Required(f.LogisticsFeeRate.Sfpydq).Message("是否为偏远地区不能为空！")

	validation.Required(f.LogisticsFeeRate.Pydqjsyf).Message("偏远地区加收运费不能为空！")

	if f.LogisticsFeeRate.Desc.Valid == true {
		validation.MaxSize(f.LogisticsFeeRate.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsFeeRate) DetailData() revel.Result {
	db_session := c.DbSession

	var logistics_fee_rate_id int64
	c.Params.Bind(&logistics_fee_rate_id, "id")

	var logistics_fee_rate models.LogisticsFeeRateInfo
	if logistics_fee_rate_id != 0 {
		has, err := db_session.Id(logistics_fee_rate_id).Get(&logistics_fee_rate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流费率不存在！"})
		}
	}

	var detail_form LogisticsFeeRateDetailForm
	detail_form.LogisticsFeeRate = logistics_fee_rate

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsFeeRate) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsFeeRateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_fee_rate := detail_form.LogisticsFeeRate

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&logistics_fee_rate)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(logistics_fee_rate.Id).Update(&logistics_fee_rate)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsFeeRate) Delete() revel.Result {
	db_session := c.DbSession

	logistics_fee_rate_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_fee_rate_id_list, "id_list")

	affected, err := db_session.In("id", logistics_fee_rate_id_list).Delete(new(models.LogisticsFeeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpLogisticsFeeRate) Import() revel.Result {
	return c.RenderEmpty()
}
