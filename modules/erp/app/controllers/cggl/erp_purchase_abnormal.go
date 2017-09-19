package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/cggl_models"
	"fmt"
)

type ErpPurchaseAbnormal struct {
	*revel.Controller
}

func (c ErpPurchaseAbnormal) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_ABNORMAL, c.Session)

	return c.RenderTemplate("erp/cggl/purchase_abnormal/purchase_abnormal_index.html")
}

type PurchaseAbnormalViewItem struct {
	cggl_models.PurchaseBillEntryInfo    `xorm:"extends" json:"pbe"`
	cggl_models.PurchaseBillInfo         `xorm:"extends" json:"pb"`
	kcgl_models.ProductInfo              `xorm:"extends" json:"p"`
	cggl_models.PurchaseAbnormalTypeInfo `xorm:"extends" json:"at"`
	models.StatusInfo                    `xorm:"extends" json:"s"`
}

func (c ErpPurchaseAbnormal) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_ABNORMAL, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT pbe.*, pb.*, p.*, at.*, s.*
FROM erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_purchase_abnormal_type at ON pbe.abnormal_type_id = at.id
    INNER JOIN erp_status s ON pbe.status_id = s.id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("pbe.*, pb.*, p.*, at.*, s.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_purchase_abnormal_type at ON pbe.abnormal_type_id = at.id
    INNER JOIN erp_status s ON pbe.status_id = s.id`)


    filter = core.SqlAnd(filter, fmt.Sprintf("pbe.status_id "))

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("pbe.id ASC")
	}

	sb.Limit(offset, limit)

	purchase_bill_entry_list := make([]PurchaseAbnormalViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&purchase_bill_entry_list)
	core.HandleError(err)

	var purchase_bill_entry_count struct {
		Count int64 `xorm:"'entry_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&purchase_bill_entry_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  purchase_bill_entry_list,
		Total: purchase_bill_entry_count.Count,
	})
}

type PurchaseBillEntryAbnormalProcessDetailForm struct {
	PurchaseBill              cggl_models.PurchaseBillInfo                    `json:"purchase_bill"`
	PurchaseBillEntry         cggl_models.PurchaseBillEntryInfo               `json:"purchase_bill_entry"`
	ProductName               string                                          `json:"product_name"`
	AbnormalProcessResultList []cggl_models.PurchaseAbnormalProcessResultInfo `json:"abnormal_process_result_list"`
}

func (f *PurchaseBillEntryAbnormalProcessDetailForm) IsCreate() bool {
	return f.PurchaseBillEntry.Id == 0
}

func (f *PurchaseBillEntryAbnormalProcessDetailForm) Valid(validation *revel.Validation) bool {

	return validation.HasErrors() == false
}

func (c ErpPurchaseAbnormal) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_ABNORMAL, c.Session)

	db_session := c.DbSession

	var purchase_bill_entry_id int64
	c.Params.Bind(&purchase_bill_entry_id, "id")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	has, err := db_session.Id(purchase_bill_entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单明细不存在！"})
	}

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err = db_session.Id(purchase_bill_entry.BillId).Get(&purchase_bill)
	core.HandleError(err)

	var product kcgl_models.ProductInfo
	product_id := purchase_bill_entry.ProductId
	_, err = db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var abnormal_process_result_list = make([]cggl_models.PurchaseAbnormalProcessResultInfo, 0)
	core.HandleError(db_session.OrderBy("id ASC").Find(&abnormal_process_result_list))

	var detail_form PurchaseBillEntryAbnormalProcessDetailForm
	detail_form.PurchaseBill = purchase_bill
	detail_form.PurchaseBillEntry = purchase_bill_entry
	detail_form.AbnormalProcessResultList = abnormal_process_result_list
	detail_form.ProductName = product.Name

	return c.RenderJson(detail_form)
}
