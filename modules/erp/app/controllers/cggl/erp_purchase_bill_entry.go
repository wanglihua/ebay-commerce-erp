package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
)

type ErpPurchaseBillEntry struct {
	*revel.Controller
}

func (c ErpPurchaseBillEntry) Index() revel.Result {
	return c.RenderTemplate("erp/cggl/purchase_bill_entry/purchase_bill_entry_index.html")
}

func (c ErpPurchaseBillEntry) ListData() revel.Result {
	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here
	query = query.Where(fmt.Sprintf("bill_id = %d", purchase_bill_id))

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	purchase_bill_entry_list := make([]models.PurchaseBillEntryInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&purchase_bill_entry_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.PurchaseBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  purchase_bill_entry_list,
		Total: count,
	})
}

type PurchaseBillEntryDetailForm struct {
	PurchaseBill      models.PurchaseBillInfo      `json:"purchase_bill"`
	PurchaseBillEntry models.PurchaseBillEntryInfo `json:"purchase_bill_entry"`
}

func (f *PurchaseBillEntryDetailForm) IsCreate() bool {
	return f.PurchaseBillEntry.Id == 0
}

func (f *PurchaseBillEntryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.PurchaseBillEntry.BillId).Message("批次不能为空！")

	validation.Required(f.PurchaseBillEntry.ProductId).Message("商品不能为空！")

	validation.Required(f.PurchaseBillEntry.Qty).Message("数量不能为空！")

	validation.Required(f.PurchaseBillEntry.Price).Message("单价不能为空！")

	validation.Required(f.PurchaseBillEntry.Amount).Message("金额不能为空！")

	if f.PurchaseBillEntry.Desc.Valid == true {
		validation.MinSize(f.PurchaseBillEntry.Desc.String, 2).Message("备注长度不能小于2！")
	}
	if f.PurchaseBillEntry.Desc.Valid == true {
		validation.MaxSize(f.PurchaseBillEntry.Desc.String, 400).Message("备注长度不能大于400！")
	}

	return validation.HasErrors() == false
}

func (c ErpPurchaseBillEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	if purchase_bill_id == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "请先选择采购单！"})
	}

	var purchase_bill models.PurchaseBillInfo
	has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
	}

	var purchase_bill_entry_id int64
	c.Params.Bind(&purchase_bill_entry_id, "id")

	var purchase_bill_entry models.PurchaseBillEntryInfo
	if purchase_bill_entry_id != 0 {
		has, err := db_session.Id(purchase_bill_entry_id).Get(&purchase_bill_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单明细不存在！"})
		}
	}

	var detail_form PurchaseBillEntryDetailForm
	detail_form.PurchaseBill = purchase_bill
	detail_form.PurchaseBillEntry = purchase_bill_entry

	return c.RenderJson(detail_form)
}

func (c ErpPurchaseBillEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form PurchaseBillEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	purchase_bill_entry := detail_form.PurchaseBillEntry
	purchase_bill_entry.BillId = detail_form.PurchaseBill.Id

	purchase_bill_entry.Amount = purchase_bill_entry.Price * float64(purchase_bill_entry.Qty)

	var affected int64
	if detail_form.IsCreate() {
		billId_count, err := db_session.Where("bill_id = ? and product_id = ?", purchase_bill_entry.BillId, purchase_bill_entry.ProductId).Count(&purchase_bill_entry)
		core.HandleError(err)
		if billId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，批次、商品已存在！" })
		}

		affected, err = db_session.Insert(&purchase_bill_entry)
		core.HandleError(err)
	} else {
		billId_count, err := db_session.Where("id <> ? and bill_id = ? and product_id = ?", purchase_bill_entry.Id, purchase_bill_entry.BillId, purchase_bill_entry.ProductId).Count(&purchase_bill_entry)
		core.HandleError(err)
		if billId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，批次、商品已存在！" })
		}

		affected, err = db_session.Id(purchase_bill_entry.Id).Update(&purchase_bill_entry)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPurchaseBillEntry) Delete() revel.Result {
	db_session := c.DbSession

	purchase_bill_entry_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_entry_id_list, "id_list")

	affected, err := db_session.In("id", purchase_bill_entry_id_list).Delete(new(models.PurchaseBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
