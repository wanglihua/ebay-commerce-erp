
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpInventoryCheck struct {
	*revel.Controller
}

func (c ErpInventoryCheck) Index() revel.Result {
	return c.RenderTemplate("erp/inventory_check/inventory_check_index.html")
}

func (c ErpInventoryCheck) ListData() revel.Result {
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

	inventory_check_list := make([]models.InventoryCheckInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&inventory_check_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.InventoryCheckInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_check_list,
		Total: count,
	})
}

type InventoryCheckDetailForm struct {
	InventoryCheck models.InventoryCheckInfo `json:"inventory_check"`
}

func (f *InventoryCheckDetailForm) IsCreate() bool {
	return f.InventoryCheck.Id == 0
}

func (f *InventoryCheckDetailForm) Valid(validation *revel.Validation) bool { 

	validation.Required(f.InventoryCheck.StockId).Message("仓库不能为空！")

	if f.InventoryCheck.Desc.Valid == true {
		validation.MinSize(f.InventoryCheck.Desc.String, 2).Message("备注长度不能小于2！")
	}
	if f.InventoryCheck.Desc.Valid == true {
		validation.MaxSize(f.InventoryCheck.Desc.String, 400).Message("备注长度不能大于400！")
	}

	validation.Required(f.InventoryCheck.CheckDate).Message("盘点日期不能为空！")

	validation.Required(f.InventoryCheck.CheckUserId).Message("盘点人不能为空！")


	validation.Required(f.InventoryCheck.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpInventoryCheck) DetailData() revel.Result {
	db_session := c.DbSession

	var inventory_check_id int64
	c.Params.Bind(&inventory_check_id, "id")

	var inventory_check models.InventoryCheckInfo
	if inventory_check_id != 0 {
		has, err := db_session.Id(inventory_check_id).Get(&inventory_check)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的库存盘点不存在！"})
		}
	}

	var detail_form InventoryCheckDetailForm
	detail_form.InventoryCheck = inventory_check

	return c.RenderJson(detail_form)
}

func (c ErpInventoryCheck) Save() revel.Result {
	db_session := c.DbSession

	var detail_form InventoryCheckDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	inventory_check := detail_form.InventoryCheck

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&inventory_check)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(inventory_check.Id).Update(&inventory_check)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpInventoryCheck) Delete() revel.Result {
	db_session := c.DbSession

	inventory_check_id_list := make([]int64, 0)
	c.Params.Bind(&inventory_check_id_list, "id_list")

	affected, err := db_session.In("id", inventory_check_id_list).Delete(new(models.InventoryCheckInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

