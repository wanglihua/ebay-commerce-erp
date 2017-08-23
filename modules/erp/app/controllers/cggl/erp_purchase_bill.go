package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/bill_service"
	"strconv"
	"time"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpPurchaseBill struct {
	*revel.Controller
}

func (c ErpPurchaseBill) Index() revel.Result {
	return c.RenderTemplate("erp/cggl/purchase_bill/purchase_bill_index.html")
}

type ErpPurchaseBillViewItem struct {
	kcgl_models.PurchaseBillInfo `xorm:"extends" json:"pb"`
	auth_models.UserInfo    `xorm:"extends" json:"u"`
	models.StatusInfo       `xorm:"extends" json:"s"`
}

func (c ErpPurchaseBill) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
	SELECT pb.*, u.*, s.*
	FROM erp_purchase_bill pb
	    INNER JOIN auth_user u ON pb.user_id = u.id
	    INNER JOIN erp_status s ON pb.status_id = s.id
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("pb.*, u.*, s.*")
	sql_builder.Count("count(pb.id) as bill_count")
	sql_builder.From(`erp_purchase_bill pb
    INNER JOIN auth_user u ON pb.user_id = u.id
    INNER JOIN erp_status s ON pb.status_id = s.id`)

	sql_builder.Where(filter)

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("pb.id ASC")
	}

	sql_builder.Limit(offset, limit)

	purchase_bill_list := make([]ErpPurchaseBillViewItem, 0)
	err := db_session.Sql(sql_builder.SqlList()).Find(&purchase_bill_list)
	core.HandleError(err)

	var bill_count struct {
		Count int64 `xorm:"'bill_count'"`
	}

	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&bill_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  purchase_bill_list,
		Total: bill_count.Count,
	})
}

type PurchaseBillDetailForm struct {
	PurchaseBill kcgl_models.PurchaseBillInfo `json:"purchase_bill"`
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

	var purchase_bill kcgl_models.PurchaseBillInfo
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
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YCJ

		affected, err = db_session.Insert(&purchase_bill)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(purchase_bill.Id).Update(&purchase_bill)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPurchaseBill) Delete() revel.Result {
	db_session := c.DbSession

	purchase_bill_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_id_list, "id_list")

	_, err := db_session.In("bill_id", purchase_bill_id_list).Delete(new(kcgl_models.PurchaseBillEntryInfo))
	core.HandleError(err)

	_, err = db_session.In("id", purchase_bill_id_list).Delete(new(kcgl_models.PurchaseBillInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据删除成功!"})
}
