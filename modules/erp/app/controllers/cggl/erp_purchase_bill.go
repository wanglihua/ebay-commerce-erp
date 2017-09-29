package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/bill_service"
	"strconv"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"matrix/modules/erp/models/cggl_models"
	"strings"
	"matrix/modules/erp/service/purchase_bill_service"
	"matrix/modules/erp/models/sppz_models"
	"time"
	"matrix/modules/erp/service/purchase_bill_log_service"
)

type ErpPurchaseBill struct {
	*revel.Controller
}

func (c ErpPurchaseBill) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	var db_session = c.DbSession

	var status_list = make([]models.StatusInfo, 0)
	var status_id_list = []int64{
		models.PURCHASE_BILL_STATUS_YCJ,
		models.PURCHASE_BILL_STATUS_YTJ,
		models.PURCHASE_BILL_STATUS_YSP,
		models.PURCHASE_BILL_STATUS_CGZ,
		models.PURCHASE_BILL_STATUS_YRK}
	err := db_session.In("id", status_id_list).OrderBy("id ASC").Find(&status_list)
	core.HandleError(err)

	var user_has_approve_pass = pass_service.Is_login_user_has_pass(pass_res.CGGL_PURCHASE_APPROVE, c.Session)

	c.RenderArgs["status_list"] = status_list
	c.RenderArgs["user_has_approve_pass"] = user_has_approve_pass

	return c.RenderTemplate("erp/cggl/purchase_bill/purchase_bill_index.html")
}

type ErpPurchaseBillViewItem struct {
	cggl_models.PurchaseBillInfo `xorm:"extends" json:"pb"`
	sppz_models.SupplierInfo     `xorm:"extends" json:"sup"`
	auth_models.UserInfo         `xorm:"extends" json:"u"`
	models.StatusInfo            `xorm:"extends" json:"s"`
}

func (c ErpPurchaseBill) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	/*
		aoData.push({'name': 'status_id_query', 'value': this_vm.page.status_id_query});
		aoData.push({'name': 'code_query', 'value': this_vm.page.code_query});
	 */

	var status_id_query int64
	c.Params.Bind(&status_id_query, "status_id_query")

	var code_query string
	c.Params.Bind(&code_query, "code_query")
	code_query = strings.TrimSpace(code_query)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT pb.*, sup.*, u.*, s.*
FROM erp_purchase_bill pb
    INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
    INNER JOIN auth_user u ON pb.user_id = u.id
    INNER JOIN erp_status s ON pb.status_id = s.id
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("pb.*, sup.*, u.*, s.*")
	sql_builder.Count("count(pb.id) as bill_count")
	sql_builder.From(`erp_purchase_bill pb
	INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
    INNER JOIN auth_user u ON pb.user_id = u.id
    INNER JOIN erp_status s ON pb.status_id = s.id`)

	if status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("status_id = %d", status_id_query))
	}

	if code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("code = '%s'", code_query))
	}

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
	PurchaseBill cggl_models.PurchaseBillInfo `json:"purchase_bill"`
	SupplierName string                       `json:"supplier_name"`
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

	validation.Required(f.PurchaseBill.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPurchaseBill) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "id")

	var supplier_name string
	var purchase_bill cggl_models.PurchaseBillInfo
	if purchase_bill_id != 0 {
		has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
		}

		var supplier sppz_models.SupplierInfo
		_, err = db_session.Id(purchase_bill.SupplierId).Get(&supplier)
		core.HandleError(err)

	} else {
		purchase_bill.Code = bill_service.GetNextBillNo(db_session, c.Session)
	}

	var detail_form PurchaseBillDetailForm
	detail_form.PurchaseBill = purchase_bill
	detail_form.SupplierName = supplier_name

	return c.RenderJson(detail_form)
}

func (c ErpPurchaseBill) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

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
		purchase_bill.ApplyUserId = loginuser.UserId
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YCJ

		affected, err = db_session.Insert(&purchase_bill)
		core.HandleError(err)

		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单新建", db_session, c.Session)
	} else {
		var purchase_bill_in_db cggl_models.PurchaseBillInfo
		_, err := db_session.Id(purchase_bill.Id).Get(&purchase_bill_in_db)
		core.HandleError(err)

		if purchase_bill_service.CanPurchaseBillModify(purchase_bill_in_db, c.Session) == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "采购单当前状态不能被修改！"})
		}

		affected, err = db_session.Id(purchase_bill.Id).Update(&purchase_bill)
		core.HandleError(err)

		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单修改", db_session, c.Session)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPurchaseBill) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	purchase_bill_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_id_list, "id_list")

	var purchase_bill_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.In("id", purchase_bill_id_list).Find(&purchase_bill_list))

	for _, purchase_bill := range purchase_bill_list {
		if purchase_bill_service.CanPurchaseBillModify(purchase_bill, c.Session) == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "采购单当前状态不能删除！",
			})
		}
	}

	for _, purchase_bill := range purchase_bill_list {
		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单删除", db_session, c.Session)
	}

	_, err := db_session.In("bill_id", purchase_bill_id_list).Delete(new(cggl_models.PurchaseBillEntryInfo))
	core.HandleError(err)

	_, err = db_session.In("id", purchase_bill_id_list).Delete(new(cggl_models.PurchaseBillInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据删除成功!"})
}

func (c ErpPurchaseBill) Submit() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	var db_session = c.DbSession

	purchase_bill_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_id_list, "id_list")

	var purchase_bill_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.In("id", purchase_bill_id_list).Find(&purchase_bill_list))

	for _, purchase_bill := range purchase_bill_list {
		count, err := db_session.Where(fmt.Sprintf("bill_id = %d", purchase_bill.Id)).Count(new(cggl_models.PurchaseBillEntryInfo))
		core.HandleError(err)

		if count == 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("采购单 %s 当前明细列表为空，不能进行提交！", purchase_bill.Code),
			})
		}
	}

	// 检查状态
	for _, purchase_bill := range purchase_bill_list {
		if purchase_bill.StatusId != models.PURCHASE_BILL_STATUS_YCJ {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("采购单 %s 当前状态不是已创建状态，不能进行提交！", purchase_bill.Code),
			})
		}
	}

	// 设置为 已提交 状态
	for _, purchase_bill := range purchase_bill_list {
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YTJ
		_, err := db_session.Id(purchase_bill.Id).Cols("status_id").Update(&purchase_bill)
		core.HandleError(err)
	}

	purchase_bill_service.SetEntriesStatus(purchase_bill_id_list, models.PURCHASE_BILL_STATUS_YTJ, db_session)

	for _, purchase_bill := range purchase_bill_list {
		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单提交", db_session, c.Session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单提交成功！",
	})
}

func (c ErpPurchaseBill) Approve() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_APPROVE, c.Session)

	var db_session = c.DbSession

	purchase_bill_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_id_list, "id_list")

	var purchase_bill_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.In("id", purchase_bill_id_list).Find(&purchase_bill_list))

	// 检查状态
	for _, purchase_bill := range purchase_bill_list {
		if purchase_bill.StatusId != models.PURCHASE_BILL_STATUS_YTJ {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("采购单 %s 当前状态不是已提交状态，不能进行审批！", purchase_bill.Code),
			})
		}
	}

	var purchase_time = core.MakeNullTime(time.Now(), true)

	// 设置为已审批状态
	for _, purchase_bill := range purchase_bill_list {
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YSP
		purchase_bill.PurchaseTime = purchase_time // 已审批时填入 采购时间
		_, err := db_session.Id(purchase_bill.Id).Cols("status_id", "purchase_time").Update(&purchase_bill)
		core.HandleError(err)
	}

	var bill_entry_list = make([]cggl_models.PurchaseBillEntryInfo, 0)
	core.HandleError(db_session.In("bill_id", purchase_bill_id_list).Find(&bill_entry_list))

	for _, bill_entry := range bill_entry_list {
		bill_entry.StatusId = models.PURCHASE_BILL_STATUS_YSP
		bill_entry.PurchaseTime = purchase_time

		_, err := db_session.Id(bill_entry.Id).Cols("status_id", "purchase_time").Update(&bill_entry)
		core.HandleError(err)
	}

	for _, purchase_bill := range purchase_bill_list {
		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单审批", db_session, c.Session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单审批成功！",
	})
}

func (c ErpPurchaseBill) SupplierSelectListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type supplier_view_item_t struct {
		Supplier sppz_models.SupplierInfo `xorm:"extends" json:"s"`
		Region   sppz_models.RegionInfo   `xorm:"extends" json:"r"`
	}

	/*
SELECT s.*, r.*
FROM erp_supplier s
    LEFT JOIN erp_region r ON s.region_id = r.id
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("s.*, r.*")
	sb.Count("count(s.id) AS supplier_count")
	sb.From(`erp_supplier s
    LEFT JOIN erp_region r ON s.region_id = r.id`)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("o.id ASC")
	}

	sb.Limit(offset, limit)

	supplier_list := make([]supplier_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&supplier_list)
	core.HandleError(err)

	var supplier_count struct {
		Count int64 `xorm:"'supplier_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&supplier_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  supplier_list,
		Total: supplier_count.Count,
	})
}
