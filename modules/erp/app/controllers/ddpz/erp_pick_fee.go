
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpPickFee struct {
	*revel.Controller
}

func (c ErpPickFee) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/pick_fee/pick_fee_index.html")
}

func (c ErpPickFee) ListData() revel.Result {
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

	pick_fee_list := make([]models.PickFeeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&pick_fee_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.PickFeeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  pick_fee_list,
		Total: count,
	})
}

type PickFeeDetailForm struct {
	PickFee models.PickFeeInfo `json:"pick_fee"`
}

func (f *PickFeeDetailForm) IsCreate() bool {
	return f.PickFee.Id == 0
}

func (f *PickFeeDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.PickFee.Fee1).Message("5Kg以下不能为空！")

	validation.Required(f.PickFee.Fee2).Message("5.1KG-10KG不能为空！")

	validation.Required(f.PickFee.Fee3).Message("10.1KG-25KG不能为空！")

	validation.Required(f.PickFee.Fee4).Message("25.1KG-50KG不能为空！")

	validation.Required(f.PickFee.Fee5).Message("50.1KG以上不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPickFee) DetailData() revel.Result {
	db_session := c.DbSession

	var pick_fee_id int64
	c.Params.Bind(&pick_fee_id, "id")

	var pick_fee models.PickFeeInfo
	if pick_fee_id != 0 {
		has, err := db_session.Id(pick_fee_id).Get(&pick_fee)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的拣货费用不存在！"})
		}
	}

	var detail_form PickFeeDetailForm
	detail_form.PickFee = pick_fee

	return c.RenderJson(detail_form)
}

func (c ErpPickFee) Save() revel.Result {
	db_session := c.DbSession

	var detail_form PickFeeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	pick_fee := detail_form.PickFee

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&pick_fee)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(pick_fee.Id).Update(&pick_fee)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPickFee) Delete() revel.Result {
	db_session := c.DbSession

	pick_fee_id_list := make([]int64, 0)
	c.Params.Bind(&pick_fee_id_list, "id_list")

	affected, err := db_session.In("id", pick_fee_id_list).Delete(new(models.PickFeeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

