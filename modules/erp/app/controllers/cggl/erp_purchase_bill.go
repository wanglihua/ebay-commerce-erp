package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/bill_service"
	"time"
)

type ErpPurchaseBill struct {
	*revel.Controller
}

func (c ErpPurchaseBill) Index() revel.Result {
	return c.RenderTemplate("erp/cggl/purchase_bill/purchase_bill_index.html")
}

func (c ErpPurchaseBill) ListData() revel.Result {
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

	purchase_bill_list := make([]models.PurchaseBillInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&purchase_bill_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.PurchaseBillInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  purchase_bill_list,
		Total: count,
	})
}

type PurchaseBillDetailForm struct {
	PurchaseBill models.PurchaseBillInfo `json:"purchase_bill"`
}

func (f *PurchaseBillDetailForm) IsCreate() bool {
	return f.PurchaseBill.Id == 0
}

func (f *PurchaseBillDetailForm) Valid(validation *revel.Validation) bool {

	if f.PurchaseBill.Desc.Valid == true {
		validation.MinSize(f.PurchaseBill.Desc.String, 2).Message("备注长度不能小于2！")
	}
	if f.PurchaseBill.Desc.Valid == true {
		validation.MaxSize(f.PurchaseBill.Desc.String, 400).Message("备注长度不能大于400！")
	}

	validation.Required(f.PurchaseBill.Date).Message("日期不能为空！")

	validation.Required(f.PurchaseBill.UserId).Message("操作人不能为空！")

	validation.Required(f.PurchaseBill.TransExpenses).Message("运费不能为空！")

	validation.Required(f.PurchaseBill.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPurchaseBill) DetailData() revel.Result {
	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "id")

	var purchase_bill models.PurchaseBillInfo
	if purchase_bill_id != 0 {
		has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
		}
	} else {
		purchase_bill.Code = bill_service.GetNextBillNo(db_session)
		purchase_bill.Date = core.Date(time.Now())
	}

	var detail_form PurchaseBillDetailForm
	detail_form.PurchaseBill = purchase_bill

	return c.RenderJson(detail_form)
}

func (c ErpPurchaseBill) Save() revel.Result {
	db_session := c.DbSession

	var detail_form PurchaseBillDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	purchase_bill := detail_form.PurchaseBill

	var affected int64
	if detail_form.IsCreate() {
		loginuser := core.GetLoginUser(c.Session)
		purchase_bill.UserId = loginuser.UserId
		purchase_bill.StatusId = models.PURCHASE_BILL_YCJ

		affected, err = db_session.Insert(&purchase_bill)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(purchase_bill.Id).Update(&purchase_bill)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPurchaseBill) Delete() revel.Result {
	db_session := c.DbSession

	purchase_bill_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_id_list, "id_list")

	affected, err := db_session.In("id", purchase_bill_id_list).Delete(new(models.PurchaseBillInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
