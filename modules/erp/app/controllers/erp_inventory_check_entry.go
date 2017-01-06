
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpInventoryCheckEntry struct {
	*revel.Controller
}

func (c ErpInventoryCheckEntry) Index() revel.Result {
	return c.RenderTemplate("erp/inventory_check_entry/inventory_check_entry_index.html")
}

func (c ErpInventoryCheckEntry) ListData() revel.Result {
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

	inventory_check_entry_list := make([]models.InventoryCheckEntryInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&inventory_check_entry_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.InventoryCheckEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_check_entry_list,
		Total: count,
	})
}

type InventoryCheckEntryDetailForm struct {
	InventoryCheckEntry models.InventoryCheckEntryInfo `json:"check_entry"`
}

func (f *InventoryCheckEntryDetailForm) IsCreate() bool {
	return f.InventoryCheckEntry.Id == 0
}

func (f *InventoryCheckEntryDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.InventoryCheckEntry.CheckId).Message("盘点批次不能为空！")

	validation.Required(f.InventoryCheckEntry.StorageLocId).Message("库位不能为空！")

	validation.Required(f.InventoryCheckEntry.ProductId).Message("商品不能为空！")

	validation.Required(f.InventoryCheckEntry.PreQty).Message("盘点前库存不能为空！")

	validation.Required(f.InventoryCheckEntry.CheckQty).Message("盘点库存不能为空！")

	validation.Required(f.InventoryCheckEntry.QtyDiff).Message("库存差值不能为空！")

	return validation.HasErrors() == false
}

func (c ErpInventoryCheckEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var inventory_check_entry_id int64
	c.Params.Bind(&inventory_check_entry_id, "id")

	var inventory_check_entry models.InventoryCheckEntryInfo
	if inventory_check_entry_id != 0 {
		has, err := db_session.Id(inventory_check_entry_id).Get(&inventory_check_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的库存盘点明细不存在！"})
		}
	}

	var detail_form InventoryCheckEntryDetailForm
	detail_form.InventoryCheckEntry = inventory_check_entry

	return c.RenderJson(detail_form)
}

func (c ErpInventoryCheckEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form InventoryCheckEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	inventory_check_entry := detail_form.InventoryCheckEntry

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&inventory_check_entry)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(inventory_check_entry.Id).Update(&inventory_check_entry)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpInventoryCheckEntry) Delete() revel.Result {
	db_session := c.DbSession

	inventory_check_entry_id_list := make([]int64, 0)
	c.Params.Bind(&inventory_check_entry_id_list, "id_list")

	affected, err := db_session.In("id", inventory_check_entry_id_list).Delete(new(models.InventoryCheckEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

