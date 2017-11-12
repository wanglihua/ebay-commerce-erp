package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/cggl_models"
	"fmt"
	"time"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/service/purchase_bill_service"
	"matrix/modules/erp/service/purchase_bill_log_service"
)

type ErpPurchaseFinance struct {
	*revel.Controller
}

func (c ErpPurchaseFinance) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	// 状态的列表， 采购单的状态的列表
	var status_list = make([]models.StatusInfo, 0)
	core.HandleError(db_session.In("id", []int64{
		models.PURCHASE_BILL_STATUS_YSP,
		models.PURCHASE_BILL_STATUS_CGZ,
		models.PURCHASE_BILL_STATUS_YRK}).Find(&status_list))

	c.RenderArgs["status_list"] = status_list

	var pay_method_list = make([]sppz_models.SupplierPayMethodInfo, 0)
	core.HandleError(db_session.Find(&pay_method_list))

	c.RenderArgs["pay_method_list"] = pay_method_list

	// 获取账期快到的采购单列表
	/*

SELECT *
FROM erp_purchase_bill pb
WHERE pb.pay_days <> 0 AND pb.status_id = 4600 AND (datediff(DATE(date_add(pb.pay_time, INTERVAL pb.pay_days DAY)), DATE(now())) <= 7);

	*/

	var sql_purchase_bill_to_pay = `SELECT *
FROM erp_purchase_bill pb
WHERE pb.pay_days <> 0 AND pb.status_id = %d AND (datediff(DATE(date_add(pb.pay_time, INTERVAL pb.pay_days DAY)), DATE(now())) <= 7);`
	sql_purchase_bill_to_pay = fmt.Sprintf(sql_purchase_bill_to_pay, models.PURCHASE_BILL_STATUS_CGZ)

	var purchase_bill_to_pay_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.Sql(sql_purchase_bill_to_pay).Find(&purchase_bill_to_pay_list))

	c.RenderArgs["purchase_bill_to_pay_list"] = purchase_bill_to_pay_list
	c.RenderArgs["purchase_bill_to_pay_count"] = len(purchase_bill_to_pay_list)
	c.RenderArgs["purchase_bill_to_pay_max_index"] = len(purchase_bill_to_pay_list) - 1

	return c.RenderTemplate("erp/cggl/purchase_finance/purchase_finance_index.html")
}

type ErpPurchaseFinanceViewItem struct {
	PurchaseBill cggl_models.PurchaseBillInfo      `xorm:"extends" json:"pb"`
	Supplier     sppz_models.SupplierInfo          `xorm:"extends" json:"sup"`
	PayMethod    sppz_models.SupplierPayMethodInfo `xorm:"extends" json:"pm"`
	Status       models.StatusInfo                 `xorm:"extends" json:"s"`
	ApplyUser    auth_models.UserInfo              `xorm:"extends" json:"ayu"`
	ApproveUser  auth_models.UserInfo              `xorm:"extends" json:"avu"`
	PayUser      auth_models.UserInfo              `xorm:"extends" json:"pu"`
	CheckinUser  auth_models.UserInfo              `xorm:"extends" json:"cu"`
}

func (c ErpPurchaseFinance) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	db_session := c.DbSession

	var status_id_query int64
	c.Params.Bind(&status_id_query, "status_id_query")

	var pay_method_query int64
	c.Params.Bind(&pay_method_query, "pay_method_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT pb.*, sup.*, pm.*, s.*, ayu.*, avu.*, pu.*, cu.*
FROM erp_purchase_bill pb
    INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
    LEFT JOIN erp_supplier_pay_method pm ON pm.id = pb.pay_method_id
    INNER JOIN erp_status s ON pb.status_id = s.id
    INNER JOIN auth_user ayu ON pb.apply_user_id = ayu.id
    LEFT JOIN auth_user avu ON pb.approve_user_id = avu.id
    LEFT JOIN auth_user pu ON pb.pay_user_id = pu.id
    LEFT JOIN auth_user cu ON pb.checkin_user_id = cu.id;
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("pb.*, sup.*, pm.*, s.*, ayu.*, avu.*, pu.*, cu.*")
	sql_builder.Count("count(pb.id) as bill_count")
	sql_builder.From(`erp_purchase_bill pb
    INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
    LEFT JOIN erp_supplier_pay_method pm ON pm.id = pb.pay_method_id
    INNER JOIN erp_status s ON pb.status_id = s.id
    INNER JOIN auth_user ayu ON pb.apply_user_id = ayu.id
    LEFT JOIN auth_user avu ON pb.approve_user_id = avu.id
    LEFT JOIN auth_user pu ON pb.pay_user_id = pu.id
    LEFT JOIN auth_user cu ON pb.checkin_user_id = cu.id`)

	filter = core.SqlAnd(filter, fmt.Sprintf("pb.status_id = %d OR pb.status_id =%d OR pb.status_id =%d",
		models.PURCHASE_BILL_STATUS_YSP, models.PURCHASE_BILL_STATUS_CGZ, models.PURCHASE_BILL_STATUS_YRK))

	if status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.status_id = %d", status_id_query))
	}

	if pay_method_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.pay_method_id = %d", pay_method_query))
	}

	sql_builder.Where(filter)

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("pb.id ASC")
	}

	sql_builder.Limit(offset, limit)

	purchase_bill_list := make([]ErpPurchaseFinanceViewItem, 0)
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

func (c ErpPurchaseFinance) Purchase() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	purchase_bill_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_id_list, "id_list")

	var purchase_bill_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.In("id", purchase_bill_id_list).Find(&purchase_bill_list))

	// 检查状态
	for _, purchase_bill := range purchase_bill_list {
		if purchase_bill.StatusId != models.PURCHASE_BILL_STATUS_YSP {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("采购单 %s 当前状态不是已审批状态，不能下单！", purchase_bill.Code),
			})
		}
	}

	var login_user = core.GetLoginUser(c.Session)
	var pay_time = core.MakeNullTime(time.Now(), true)

	// 设置为已下单状态
	for _, purchase_bill := range purchase_bill_list {
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_CGZ
		purchase_bill.PayTime = pay_time
		purchase_bill.PayUserId = login_user.UserId
		_, err := db_session.Id(purchase_bill.Id).Cols("status_id", "pay_time", "pay_user_id").Update(&purchase_bill)
		core.HandleError(err)
	}

	for _, purchase_bill := range purchase_bill_list {
		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单下单", db_session, c.Session)
	}

	purchase_bill_service.SetEntriesStatus(purchase_bill_id_list, models.PURCHASE_BILL_STATUS_CGZ, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单下单成功！",
	})
}

func (c ErpPurchaseFinance) EntryList() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
	core.HandleError(err)

	c.RenderArgs["bill_code"] = purchase_bill.Code
	return c.RenderTemplate("erp/cggl/purchase_finance/purchase_finance_entry_list.html")
}

// 这个后来不使用了，付款金额在采购管理进行维护了
func (c ErpPurchaseFinance) EntryPayAmount() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	c.RenderArgs["entry"] = purchase_bill_entry

	return c.RenderTemplate("erp/cggl/purchase_finance/purchase_finance_entry_pay_amount.html")
}

// 这个后来不使用了，付款金额在采购管理进行维护了
func (c ErpPurchaseFinance) EntryPayAmountSave() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var pay_amount float64
	c.Params.Bind(&pay_amount, "pay_amount")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	purchase_bill_entry.PayAmount = pay_amount

	_, err = db_session.Id(entry_id).Cols("pay_amount").Update(&purchase_bill_entry)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细更新付款金额成功！",
	})
}

// 用于分批采购中的，采购明细审批
func (c ErpPurchaseFinance) PurchaseEntry() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	if purchase_bill_entry.StatusId != models.PURCHASE_BILL_STATUS_YSP || purchase_bill_entry.TypeId != cggl_models.PURCHASE_BILL_TYPE_IN_TURN {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "只有已审批的分批采购的明细单才能下单！",
		})
	}

	purchase_bill_entry.StatusId = models.PURCHASE_BILL_STATUS_CGZ

	_, err = db_session.Id(entry_id).Cols("status_id").Update(&purchase_bill_entry)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单明细分批采购下单成功！",
	})
}
