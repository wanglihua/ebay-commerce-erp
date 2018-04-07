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
	"unicode/utf8"
	"strings"
	"github.com/go-xorm/xorm"
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
		models.PURCHASE_BILL_STATUS_YTJ,
		models.PURCHASE_BILL_STATUS_YSP,
		models.PURCHASE_BILL_STATUS_CGZ,
		models.PURCHASE_BILL_STATUS_YRK,
		models.PURCHASE_BILL_STATUS_YWJ}).Find(&status_list))

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

	var sql_purchase_bill_to_pay = `
SELECT *
FROM erp_purchase_bill pb
WHERE pb.pay_days <> 0 AND pb.status_id = %d AND (datediff(DATE(date_add(pb.pay_time, INTERVAL pb.pay_days DAY)), DATE(now())) <= 7);`
	sql_purchase_bill_to_pay = fmt.Sprintf(sql_purchase_bill_to_pay, models.PURCHASE_BILL_STATUS_CGZ)

	var purchase_bill_to_pay_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.Sql(sql_purchase_bill_to_pay).Find(&purchase_bill_to_pay_list))

	c.RenderArgs["purchase_bill_to_pay_list"] = purchase_bill_to_pay_list
	c.RenderArgs["purchase_bill_to_pay_count"] = len(purchase_bill_to_pay_list)
	c.RenderArgs["purchase_bill_to_pay_max_index"] = len(purchase_bill_to_pay_list) - 1

	// 获取有明细的账期快到的采购单列表
	/*

SELECT pb.*
FROM erp_purchase_bill pb
	INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
WHERE pb.pay_days <> 0 AND pbe.status_id = 4600 AND (datediff(DATE(date_add(pbe.pay_time, INTERVAL pb.pay_days DAY)), DATE(now())) <= 7);

	*/

	var sql_purchase_bill_entry_to_pay = `
SELECT pb.*
FROM erp_purchase_bill pb
	INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
WHERE pb.pay_days <> 0 AND pbe.status_id = %d AND (datediff(DATE(date_add(pbe.pay_time, INTERVAL pb.pay_days DAY)), DATE(now())) <= 7);`
	sql_purchase_bill_entry_to_pay = fmt.Sprintf(sql_purchase_bill_entry_to_pay, models.PURCHASE_BILL_STATUS_CGZ)

	var purchase_bill_entry_to_pay_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.Sql(sql_purchase_bill_entry_to_pay).Find(&purchase_bill_entry_to_pay_list))

	c.RenderArgs["purchase_bill_entry_to_pay_list"] = purchase_bill_entry_to_pay_list
	c.RenderArgs["purchase_bill_entry_to_pay_count"] = len(purchase_bill_entry_to_pay_list)
	c.RenderArgs["purchase_bill_entry_to_pay_max_index"] = len(purchase_bill_entry_to_pay_list) - 1

	/*
SELECT group_concat(DISTINCT pb.code) AS code_list_str
FROM erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
WHERE pbe.status_id = 4500 AND pbe.type_id = 200;
	 */

	var sql_entry_to_purchase_code = `SELECT group_concat(DISTINCT pb.code) AS code_list_str
FROM erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
WHERE pbe.status_id = %d AND pbe.type_id = %d;`
	sql_entry_to_purchase_code = fmt.Sprintf(sql_entry_to_purchase_code, models.PURCHASE_BILL_STATUS_YSP, cggl_models.PURCHASE_BILL_TYPE_IN_TURN)
	var entry_purchase_code_strs struct {
		CodeListStr string `xorm:"'code_list_str'"`
	}
	_, err := db_session.Sql(sql_entry_to_purchase_code).Get(&entry_purchase_code_strs)
	core.HandleError(err)

	c.RenderArgs["entry_purchase_code_strs_len"] = int64(utf8.RuneCountInString(entry_purchase_code_strs.CodeListStr))
	c.RenderArgs["entry_purchase_code_strs"] = entry_purchase_code_strs.CodeListStr

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.Find(&user_list))

	c.RenderArgs["status_list"] = status_list
	c.RenderArgs["time_query_list"] = get_time_query_list()
	c.RenderArgs["user_query_id"] = 0
	c.RenderArgs["user_query_value"] = 0
	c.RenderArgs["user_query_list"] = get_user_query_list()
	c.RenderArgs["user_list"] = user_list

	return c.RenderTemplate("erp/cggl/purchase_finance/purchase_finance_index.html")
}

type ErpPurchaseFinanceViewItem struct {
	PurchaseBill  cggl_models.PurchaseBillInfo      `xorm:"extends" json:"pb"`
	Supplier      sppz_models.SupplierInfo          `xorm:"extends" json:"sup"`
	PayMethod     sppz_models.SupplierPayMethodInfo `xorm:"extends" json:"pm"`
	Status        models.StatusInfo                 `xorm:"extends" json:"s"`
	ApplyUser     auth_models.UserInfo              `xorm:"extends" json:"ayu"`
	ApproveUser   auth_models.UserInfo              `xorm:"extends" json:"avu"`
	PayUser       auth_models.UserInfo              `xorm:"extends" json:"pu"`
	CheckinUser   auth_models.UserInfo              `xorm:"extends" json:"cu"`
	PayExpireDays int64                             `xorm:"'pay_expire_days'" json:"pay_expire_days"`
}

func (c ErpPurchaseFinance) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	db_session := c.DbSession

	var status_id_query int64
	c.Params.Bind(&status_id_query, "status_id_query")

	var pay_method_query int64
	c.Params.Bind(&pay_method_query, "pay_method_query")

	var code_query string
	c.Params.Bind(&code_query, "code_query")
	code_query = strings.TrimSpace(code_query)

	// supplier_query
	var supplier_query string
	c.Params.Bind(&supplier_query, "supplier_query")
	supplier_query = strings.TrimSpace(supplier_query)

	var user_query_id int64
	c.Params.Bind(&user_query_id, "user_query_id")

	var user_query_value int64
	c.Params.Bind(&user_query_value, "user_query_value")

	// time_query_id
	var time_query_id int64
	c.Params.Bind(&time_query_id, "time_query_id")

	var query_time_begin_str string
	c.Params.Bind(&query_time_begin_str, "query_time_begin")
	query_time_begin, query_time_begin_err := core.Parse_time(query_time_begin_str)

	var query_time_end_str string
	c.Params.Bind(&query_time_end_str, "query_time_end")
	query_time_end, query_time_end_err := core.Parse_time(query_time_end_str)

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")
	sku_query = strings.TrimSpace(sku_query)

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
	PURCHASE_BILL_STATUS_YCJ = 4400 // 采购单状态 已创建
	PURCHASE_BILL_STATUS_YTJ = 4450 // 采购单状态 已提交
	PURCHASE_BILL_STATUS_YSP = 4500 // 采购单状态 已审批
	PURCHASE_BILL_STATUS_CGZ = 4600 // 采购单状态 采购中
	PURCHASE_BILL_STATUS_YRK = 4700 // 采购单状态 已入库
	PURCHASE_BILL_STATUS_YC  = 4710 // 采购单状态 异常
	PURCHASE_BILL_STATUS_YCL = 4720 // 采购单状态 异常已处理
 	*/

	/*
SELECT pb.*, sup.*, pm.*, s.*, ayu.*, avu.*, pu.*, cu.*, pb.pay_time, if(pb.status_id <> 4600, 500000, datediff(DATE(date_add(ifnull(pb.pay_time, '3000-01-01'), INTERVAL pb.pay_days DAY)), DATE(now()))) AS pay_expire_days
FROM erp_purchase_bill pb
    INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
    LEFT JOIN erp_supplier_pay_method pm ON pm.id = pb.pay_method_id
    INNER JOIN erp_status s ON pb.status_id = s.id
    INNER JOIN auth_user ayu ON pb.apply_user_id = ayu.id
    LEFT JOIN auth_user avu ON pb.approve_user_id = avu.id
    LEFT JOIN auth_user pu ON pb.pay_user_id = pu.id
    LEFT JOIN auth_user cu ON pb.checkin_user_id = cu.id
ORDER BY pay_expire_days DESC;
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("pb.*, sup.*, pm.*, s.*, ayu.*, avu.*, pu.*, cu.*, pb.pay_time, if(pb.status_id <> 4600, 500000, datediff(DATE(date_add(ifnull(pb.pay_time, '3000-01-01'), INTERVAL pb.pay_days DAY)), DATE(now()))) AS pay_expire_days")
	sql_builder.Count("count(pb.id) as bill_count")
	sql_builder.From(`erp_purchase_bill pb
    INNER JOIN erp_supplier sup ON sup.id = pb.supplier_id
    LEFT JOIN erp_supplier_pay_method pm ON pm.id = pb.pay_method_id
    INNER JOIN erp_status s ON pb.status_id = s.id
    INNER JOIN auth_user ayu ON pb.apply_user_id = ayu.id
    LEFT JOIN auth_user avu ON pb.approve_user_id = avu.id
    LEFT JOIN auth_user pu ON pb.pay_user_id = pu.id
    LEFT JOIN auth_user cu ON pb.checkin_user_id = cu.id`)

	filter = core.SqlAnd(filter, fmt.Sprintf("pb.status_id = %d OR pb.status_id =%d OR pb.status_id =%d OR pb.status_id =%d",
		models.PURCHASE_BILL_STATUS_YSP, models.PURCHASE_BILL_STATUS_CGZ, models.PURCHASE_BILL_STATUS_YRK, models.PURCHASE_BILL_STATUS_YWJ))

	filter = core.SqlOr(filter, fmt.Sprintf(`pb.id IN (
		SELECT pbe.bill_id
		FROM erp_purchase_bill_entry pbe
		WHERE pbe.status_id = %d)`,
		models.PURCHASE_BILL_STATUS_YSP))

	if status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.status_id = %d", status_id_query))
	}

	if pay_method_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.pay_method_id = %d", pay_method_query))
	}

	if code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.code = '%s'", code_query))
	}

	if supplier_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("sup.name LIKE '%%%s%%'", supplier_query))
	}

	if user_query_id == USER_QUERY_APPLY {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.apply_user_id = %d", user_query_value))
	} else if user_query_id == USER_QUERY_PAY {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.pay_user_id = %d", user_query_value))
	} else if user_query_id == USER_QUERY_CHECKIN {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.checkin_user_id = %d", user_query_value))
	}

	if query_time_begin_err == nil {
		if time_query_id == TIME_QUERY_APPLY_TIME {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.apply_time >= '%s'", query_time_begin.Format("2006-01-02 15:04")))
		} else if time_query_id == TIME_QUERY_PAY_TIME {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.pay_time >= '%s'", query_time_begin.Format("2006-01-02 15:04")))
		} else if time_query_id == TIME_QUERY_CHECKIN_TIME {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.checkin_time >= '%s'", query_time_begin.Format("2006-01-02 15:04")))
		} /* else {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.apply_time >= '%s'", query_time_begin.Format("2006-01-02 15:04")))
		} */
	}

	if query_time_end_err == nil {

		if time_query_id == TIME_QUERY_APPLY_TIME {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.apply_time <= '%s'", query_time_end.Format("2006-01-02 15:04")))
		} else if time_query_id == TIME_QUERY_PAY_TIME {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.pay_time <= '%s'", query_time_end.Format("2006-01-02 15:04")))
		} else if time_query_id == TIME_QUERY_CHECKIN_TIME {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.checkin_time <= '%s'", query_time_end.Format("2006-01-02 15:04")))
		} /* else {
			filter = core.SqlAnd(filter, fmt.Sprintf("pb.apply_time <= '%s'", query_time_end.Format("2006-01-02 15:04")))
		} */
	}

	if sku_query != "" {
		/*
		SELECT *
		FROM erp_purchase_bill pb
		WHERE pb.id IN (
			SELECT pbe.bill_id
			FROM erp_purchase_bill_entry pbe
			WHERE pbe.product_id IN (
				SELECT p.id
				FROM erp_product p
				WHERE p.sku = '124'
			)
		)
		 */
		var sku_filter = `pb.id IN (
			SELECT pbe.bill_id
			FROM erp_purchase_bill_entry pbe
			WHERE pbe.product_id IN (
				SELECT p.id
				FROM erp_product p
				WHERE p.sku = '%s'
			)
		)`
		sku_filter = fmt.Sprintf(sku_filter, sku_query)
		filter = core.SqlAnd(filter, sku_filter)
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

	var bill_id = purchase_bill_entry.BillId
	purchase_bill_service.UpdatePurchaseBillAmount(bill_id, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细更新付款金额成功！",
	})
}

func (c ErpPurchaseFinance) EntryPrice() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	c.RenderArgs["entry"] = purchase_bill_entry

	return c.RenderTemplate("erp/cggl/purchase_finance/purchase_finance_entry_price.html")
}

func (c ErpPurchaseFinance) EntryPriceSave() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	c.entry_price_save(db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细商品单价成功！",
	})
}

func (c ErpPurchaseFinance) entry_price_save(db_session *xorm.Session) {
	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var price float64
	c.Params.Bind(&price, "price")

	var qty int64
	c.Params.Bind(&qty, "qty")

	var freight_charge float64
	c.Params.Bind(&freight_charge, "freight_charge")

	var pay_amount float64
	c.Params.Bind(&pay_amount, "pay_amount")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	purchase_bill_entry.Price = price
	purchase_bill_entry.Qty = qty
	purchase_bill_entry.Amount = purchase_bill_entry.Price * float64(purchase_bill_entry.Qty)
	purchase_bill_entry.FreightCharge = freight_charge
	purchase_bill_entry.PayAmount = pay_amount
	purchase_bill_entry.DrawbackAmount = purchase_bill_service.CalculateBillEntryDrawbackAmount(purchase_bill_entry)

	_, err = db_session.Id(entry_id).Cols("price", "qty", "amount", "freight_charge", "pay_amount", "drawback_amount").Update(&purchase_bill_entry)
	core.HandleError(err)

	var bill_id = purchase_bill_entry.BillId
	purchase_bill_service.UpdatePurchaseBillAmount(bill_id, db_session)
}

// 用于分批采购中的，采购明细审批
func (c ErpPurchaseFinance) PurchaseEntry() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	c.entry_price_save(db_session)

	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	if purchase_bill_entry.StatusId != models.PURCHASE_BILL_STATUS_YSP {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "只有已审批的采购的明细单才能下单！",
		})
	}

	purchase_bill_entry.StatusId = models.PURCHASE_BILL_STATUS_CGZ

	var pay_time = core.MakeNullTime(time.Now(), true)
	purchase_bill_entry.PayTime = pay_time

	_, err = db_session.Id(entry_id).Cols("status_id", "pay_time").Update(&purchase_bill_entry)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单明细采购下单成功！",
	})
}

func (c ErpPurchaseFinance) PayFinalAmount() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	var entry_id int64
	c.Params.Bind(&entry_id, "entry_id")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	purchase_bill_entry.FinalAmountStatusId = core.MakeNullInt(models.PURCHASE_BILL_ENTRY_FINAL_AMOUNT_STATUS_YZF, true)

	_, err = db_session.Id(entry_id).Cols("final_amount_status_id").Update(&purchase_bill_entry)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单明细支付尾款成功！",
	})
}

func (c ErpPurchaseFinance) Complete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	var db_session = c.DbSession

	purchase_bill_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_id_list, "id_list")

	var purchase_bill_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.In("id", purchase_bill_id_list).Find(&purchase_bill_list))

	// 检查状态
	for _, purchase_bill := range purchase_bill_list {
		if purchase_bill.StatusId != models.PURCHASE_BILL_STATUS_YRK {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("采购单 %s 当前状态不是已入库状态，不能完结！", purchase_bill.Code),
			})
		}
	}

	// 设置为已完结状态
	for _, purchase_bill := range purchase_bill_list {
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YWJ
		_, err := db_session.Id(purchase_bill.Id).Cols("status_id").Update(&purchase_bill)
		core.HandleError(err)
	}

	var bill_entry_list = make([]cggl_models.PurchaseBillEntryInfo, 0)
	core.HandleError(db_session.In("bill_id", purchase_bill_id_list).Find(&bill_entry_list))

	for _, bill_entry := range bill_entry_list {
		bill_entry.StatusId = models.PURCHASE_BILL_STATUS_YWJ
		_, err := db_session.Id(bill_entry.Id).Cols("status_id").Update(&bill_entry)
		core.HandleError(err)
	}

	for _, purchase_bill := range purchase_bill_list {
		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单完结", db_session, c.Session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单完结成功！",
	})
}
