package controllers

import (
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/service/purchase_bill_service"

	"github.com/revel/revel"
	"fmt"
	"strings"
	"matrix/modules/erp/service/inventory_service"
	"encoding/json"
	"time"
	"matrix/modules/erp/service/purchase_bill_log_service"
	"math"
)

type ErpPurchaseCheckin struct {
	*revel.Controller
}

func (c ErpPurchaseCheckin) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_CHECKIN, c.Session)

	/*
	PURCHASE_BILL_STATUS_YCJ = 4400 // 采购单状态 已创建
	PURCHASE_BILL_STATUS_YTJ = 4450 // 采购单状态 已提交
	PURCHASE_BILL_STATUS_YSP = 4500 // 采购单状态 已审批
	PURCHASE_BILL_STATUS_CGZ = 4600 // 采购单状态 采购中
	PURCHASE_BILL_STATUS_YRK = 4700 // 采购单状态 已入库
	PURCHASE_BILL_STATUS_YC  = 4710 // 采购单状态 异常
	PURCHASE_BILL_STATUS_YCL = 4720 // 采购单状态 异常已处理
	 */

	var db_session = c.DbSession

	// 获得有明细状态是 采购中, 已入库， 异常，异常已处理 的所有采购单的批次号。 这时获得的采购单的状态只有两个 采购中，已入库

	// 批次号的列表
	/*
		SELECT DISTINCT pb.code AS bill_code, s.name AS bill_status
		FROM erp_purchase_bill pb
			INNER JOIN erp_status s ON s.id = pb.status_id
			INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
		WHERE pbe.status_id = 0;
	*/
	var sql_get_bill_code_list = `
	SELECT DISTINCT pb.code AS bill_code, s.name AS bill_status
	FROM erp_purchase_bill pb
		INNER JOIN erp_status s ON s.id = pb.status_id
		INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
	WHERE pbe.status_id = %d OR pbe.status_id = %d OR pbe.status_id = %d OR pbe.status_id = %d;`

	sql_get_bill_code_list = fmt.Sprintf(sql_get_bill_code_list,
		models.PURCHASE_BILL_STATUS_CGZ,
		models.PURCHASE_BILL_STATUS_YRK,
		models.PURCHASE_BILL_STATUS_YC,
		models.PURCHASE_BILL_STATUS_YCL)

	type bill_query_struct struct {
		BillCode   string `xorm:"'bill_code'"`
		BillStatus string `xorm:"'bill_status'"`
	}

	var bill_code_struct_list = make([]bill_query_struct, 0)
	core.HandleError(db_session.Sql(sql_get_bill_code_list).Find(&bill_code_struct_list))

	var bill_query_list = make([]bill_query_struct, 0)
	for _, bill_query_struct := range bill_code_struct_list {
		bill_query_list = append(bill_query_list, bill_query_struct)
	}

	// 状态的列表， 采购单明细的状态的列表
	var status_list = make([]models.StatusInfo, 0)
	core.HandleError(db_session.In("id", []int64{
		models.PURCHASE_BILL_STATUS_CGZ,
		models.PURCHASE_BILL_STATUS_YRK,
		models.PURCHASE_BILL_STATUS_YC,
		models.PURCHASE_BILL_STATUS_YCL}).Find(&status_list))

	// 获取账期快到的采购单列表
	/*

SELECT *
FROM erp_purchase_bill pb
WHERE pb.pay_days <> 0 AND pb.status_id = 4600 AND (datediff(DATE(date_add(pb.pay_time, INTERVAL pb.arrival_days DAY)), DATE(now())) <= 7);

	*/

	var sql_purchase_bill_to_arrival = `SELECT *
FROM erp_purchase_bill pb
WHERE pb.pay_days <> 0 AND pb.status_id = %d AND (datediff(DATE(date_add(pb.pay_time, INTERVAL pb.arrival_days DAY)), DATE(now())) <= 7);`
	sql_purchase_bill_to_arrival = fmt.Sprintf(sql_purchase_bill_to_arrival, models.PURCHASE_BILL_STATUS_CGZ)

	var purchase_bill_to_arrival_list = make([]cggl_models.PurchaseBillInfo, 0)
	core.HandleError(db_session.Sql(sql_purchase_bill_to_arrival).Find(&purchase_bill_to_arrival_list))

	c.RenderArgs["bill_query_list"] = bill_query_list
	c.RenderArgs["status_list"] = status_list
	c.RenderArgs["purchase_bill_to_arrival_list"] = purchase_bill_to_arrival_list
	c.RenderArgs["purchase_bill_to_arrival_count"] = len(purchase_bill_to_arrival_list)
	c.RenderArgs["purchase_bill_to_arrival_max_index"] = len(purchase_bill_to_arrival_list) - 1

	return c.RenderTemplate("erp/cggl/purchase_checkin/purchase_checkin_index.html")
}

type PurchaseCheckinViewItem struct {
	cggl_models.PurchaseBillEntryInfo `xorm:"extends" json:"pbe"`
	cggl_models.PurchaseBillInfo      `xorm:"extends" json:"pb"`
	kcgl_models.ProductInfo           `xorm:"extends" json:"p"`
	models.StatusInfo                 `xorm:"extends" json:"s"`
}

func (c ErpPurchaseCheckin) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_CHECKIN, c.Session)

	var bill_code_query string
	c.Params.Bind(&bill_code_query, "bill_code_query")
	bill_code_query = strings.TrimSpace(bill_code_query)

	var status_id_query int64
	c.Params.Bind(&status_id_query, "status_id_query")

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT pbe.*, pb.*, p.*, s.*
		FROM erp_purchase_bill_entry pbe
		    INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
		    INNER JOIN erp_product p ON pbe.product_id = p.id
		    INNER JOIN erp_status s ON pbe.status_id = s.id
	*/

	sb := core.NewSqlBuilder()
	sb.Select("pbe.*, pb.*, p.*, s.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_purchase_bill_entry pbe
	INNER JOIN erp_purchase_bill pb ON pbe.bill_id = pb.id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_status s ON pbe.status_id = s.id`)

	filter = core.SqlAnd(filter, fmt.Sprintf("pbe.status_id = %d OR pbe.status_id = %d OR pbe.status_id = %d ",
		models.PURCHASE_BILL_STATUS_CGZ, models.PURCHASE_BILL_STATUS_YRK, models.PURCHASE_BILL_STATUS_YC))

	if bill_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.code = '%s'", bill_code_query))
	}

	if status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.status_id = %d", status_id_query))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("pbe.id ASC")
	}

	sb.Limit(offset, limit)

	purchase_bill_entry_list := make([]PurchaseCheckinViewItem, 0, limit)
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

type PurchaseBillEntryCheckinDetailForm struct {
	PurchaseBill      cggl_models.PurchaseBillInfo           `json:"purchase_bill"`
	PurchaseBillEntry cggl_models.PurchaseBillEntryInfo      `json:"purchase_bill_entry"`
	ProductName       string                                 `json:"product_name"`
	AbnormalTypeList  []cggl_models.PurchaseAbnormalTypeInfo `json:"abnormal_type_list"`
}

func (f *PurchaseBillEntryCheckinDetailForm) IsCreate() bool {
	return f.PurchaseBillEntry.Id == 0
}

func (f *PurchaseBillEntryCheckinDetailForm) Valid(validation *revel.Validation) bool {

	/*
	if f.PurchaseBillEntry.CheckinQty <= 0 || f.PurchaseBillEntry.CheckinQty > f.PurchaseBillEntry.Qty {
		validation.Error("入库数量输入有误！")
	}
	*/

	// 也有可能订了 10 个，到货了 12 个，最终入库的是 12个
	if f.PurchaseBillEntry.CheckinQty <= 0 {
		validation.Error("入库数量输入有误！")
	}

	return validation.HasErrors() == false
}

func (c ErpPurchaseCheckin) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_CHECKIN, c.Session)

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

	var abnormal_type_list = make([]cggl_models.PurchaseAbnormalTypeInfo, 0)
	core.HandleError(db_session.OrderBy("id ASC").Find(&abnormal_type_list))

	var detail_form PurchaseBillEntryCheckinDetailForm
	detail_form.PurchaseBill = purchase_bill
	detail_form.PurchaseBillEntry = purchase_bill_entry
	detail_form.PurchaseBillEntry.CheckinQty = purchase_bill_entry.Qty
	detail_form.AbnormalTypeList = abnormal_type_list
	detail_form.PurchaseBillEntry.AbnormalTypeId = cggl_models.ABNORMAL_TYPE_ZC
	detail_form.ProductName = product.Name

	return c.RenderJson(detail_form)
}

func (c ErpPurchaseCheckin) Checkin() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_CHECKIN, c.Session)

	var db_session = c.DbSession

	var detail_form PurchaseBillEntryCheckinDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var bill_entry_from_ui = detail_form.PurchaseBillEntry

	var bill_entry_in_db cggl_models.PurchaseBillEntryInfo
	_, err = db_session.Id(bill_entry_from_ui.Id).Get(&bill_entry_in_db)
	core.HandleError(err)

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err = db_session.Id(bill_entry_in_db.BillId).Get(&purchase_bill)
	core.HandleError(err)

	var product kcgl_models.ProductInfo
	_, err = db_session.Id(bill_entry_in_db.ProductId).Get(&product)
	core.HandleError(err)

	// 检查状态
	if bill_entry_in_db.StatusId != models.PURCHASE_BILL_STATUS_CGZ {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单明细状态当前不是采购中状态，不能入库！",
		})
	}

	if bill_entry_from_ui.CheckinQty > bill_entry_in_db.Qty {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "入库数量不能大于采购数量！",
		})
	}

	var is_abnormal = bill_entry_from_ui.CheckinQty < bill_entry_in_db.Qty

	// 检查是不是选择了异常类型
	if is_abnormal == true && bill_entry_from_ui.AbnormalTypeId == cggl_models.ABNORMAL_TYPE_ZC {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择异常类型！",
		})
	}

	// 更新明细的状态
	bill_entry_in_db.StatusId = models.PURCHASE_BILL_STATUS_YRK
	if is_abnormal == true {
		bill_entry_in_db.StatusId = models.PURCHASE_BILL_STATUS_YC
		bill_entry_in_db.AbnormalProcessResultId = cggl_models.ABNORMAL_PROCESS_RESULT_WCL

		// 插入采购单历史 异常
		// TODO： 可以写明 多少数量入库，多少数量异常
		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, fmt.Sprintf("采购单明细SKU:%s 异常", product.Sku), db_session, c.Session)
	} else {
		// 插入采购单历史，入库
		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, fmt.Sprintf("采购单明细SKU:%s 入库", product.Sku), db_session, c.Session)
	}

	bill_entry_in_db.CheckinQty = bill_entry_from_ui.CheckinQty
	bill_entry_in_db.AbnormalTypeId = bill_entry_from_ui.AbnormalTypeId
	bill_entry_in_db.CheckinTime = core.MakeNullTime(time.Now(), true)
	_, err = db_session.Id(bill_entry_in_db.Id).Cols("status_id", "abnormal_process_result_id", "checkin_qty", "abnormal_type_id", "checkin_time").Update(&bill_entry_in_db)
	core.HandleError(err)

	// 查看同一个采购单的明细是不是都是 已入库， 异常， 异常已处理 的这个状态了
	var bill_id = bill_entry_in_db.BillId
	count, err := db_session.
	Where(fmt.Sprintf("bill_id = %d AND status_id <> %d AND status_id <> %d AND status_id <> %d",
		bill_id,
		models.PURCHASE_BILL_STATUS_YRK,
		models.PURCHASE_BILL_STATUS_YC,
		models.PURCHASE_BILL_STATUS_YCL)).
		Count(new(cggl_models.PurchaseBillEntryInfo))
	core.HandleError(err)

	// 如果都是当前设置的这个状态了， 就更新采购单的状态
	if count == 0 {

		var login_user = core.GetLoginUser(c.Session)

		var purchase_bill cggl_models.PurchaseBillInfo
		_, err = db_session.Id(bill_id).Get(&purchase_bill)
		core.HandleError(err)

		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YRK
		purchase_bill.CheckinUserId = login_user.UserId
		purchase_bill.CheckinTime = core.MakeNullTime(time.Now(), true)
		_, err = db_session.Id(purchase_bill.Id).Cols("status_id", "checkin_user_id", "checkin_time").Update(&purchase_bill)
		core.HandleError(err)
	}

	// 计算中国仓成本价
	purchase_bill_service.CalculateBillEntryCnInventoryCost(purchase_bill.Id, db_session)

	// 更新中国仓的库存数量
	inventory_service.AddCnInventoryQty(bill_entry_from_ui.ProductId, bill_entry_from_ui.CheckinQty, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细入库成功！",
	})
}

func (c ErpPurchaseCheckin) CheckinList() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_CHECKIN, c.Session)

	var db_session = c.DbSession

	entry_id_list := make([]int64, 0)
	c.Params.Bind(&entry_id_list, "id_list")

	var bill_entry_list = make([]cggl_models.PurchaseBillEntryInfo, 0)
	core.HandleError(db_session.In("id", entry_id_list).Find(&bill_entry_list))

	// 检查状态
	for _, bill_entry := range bill_entry_list {
		if bill_entry.StatusId != models.PURCHASE_BILL_STATUS_CGZ {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "选中的采购单明细当前状态不全是采购中状态，不能入库！",
			})
		}
	}

	// 设置为已入库状态。更新明细的状态，如果一个采购单下面的所有明细都已入库状态，把采购单也置为已入库状态
	for _, bill_entry := range bill_entry_list {
		// 更新明细的状态
		bill_entry.StatusId = models.PURCHASE_BILL_STATUS_YRK
		bill_entry.CheckinQty = bill_entry.Qty                   // 批量入库的时候只能这样默认了
		bill_entry.AbnormalTypeId = cggl_models.ABNORMAL_TYPE_ZC // 批量入库的时候只能这样默认了
		bill_entry.CheckinTime = core.MakeNullTime(time.Now(), true)
		_, err := db_session.Id(bill_entry.Id).Cols("status_id", "checkin_qty", "abnormal_type_id", "checkin_time").Update(bill_entry)
		core.HandleError(err)

		purchase_bill_service.UpdateBillStatus(bill_entry, models.PURCHASE_BILL_STATUS_YRK, db_session)
	}

	for _, bill_entry := range bill_entry_list {
		// 计算中国仓成本价
		purchase_bill_service.CalculateBillEntryCnInventoryCost(bill_entry.BillId, db_session)

		// 更新中国仓的库存数量
		inventory_service.AddCnInventoryQty(bill_entry.ProductId, bill_entry.CheckinQty, db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细入库成功！",
	})
}

func (c ErpPurchaseCheckin) CheckinInTurn() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_CHECKIN, c.Session)

	var db_session = c.DbSession

	var detail_form PurchaseBillEntryCheckinDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var bill_entry_from_ui = detail_form.PurchaseBillEntry

	var bill_entry_in_db cggl_models.PurchaseBillEntryInfo
	_, err = db_session.Id(bill_entry_from_ui.Id).Get(&bill_entry_in_db)
	core.HandleError(err)

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err = db_session.Id(bill_entry_in_db.BillId).Get(&purchase_bill)
	core.HandleError(err)

	var product kcgl_models.ProductInfo
	_, err = db_session.Id(bill_entry_in_db.ProductId).Get(&product)
	core.HandleError(err)

	// 检查状态
	if bill_entry_in_db.StatusId != models.PURCHASE_BILL_STATUS_CGZ {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单明细状态当前不是采购中状态，不能入库！",
		})
	}

	if bill_entry_from_ui.CheckinQty > bill_entry_in_db.Qty {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "入库数量不能大于采购数量！",
		})
	}

	var is_all_checkin = bill_entry_from_ui.CheckinQty == bill_entry_in_db.Qty

	// 检查是不是分批入库
	if is_all_checkin == true {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "全部数量入库，请使用正常入库而不是分批入库！",
		})
	}

	// 插入采购单历史，入库
	purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, fmt.Sprintf("采购单明细SKU:%s 分批入库，%d个", product.Sku, bill_entry_from_ui.CheckinQty), db_session, c.Session)

	var qty_before = bill_entry_in_db.Qty
	var mel_qty_before = bill_entry_in_db.MelQty
	// var bne_qty_before = bill_entry_in_db.BneQty
	var fitting_cost_before = bill_entry_in_db.FittingCost
	var other_cost_before = bill_entry_in_db.OtherCost
	var checkin_qty = bill_entry_from_ui.CheckinQty
	var checkin_rate = float64(checkin_qty) / float64(qty_before)

	bill_entry_in_db.Qty = qty_before - checkin_qty
	bill_entry_in_db.MelQty = bill_entry_in_db.MelQty - int64(math.Floor(float64(mel_qty_before)*checkin_rate+0.5)) // 四舍五入
	bill_entry_in_db.BneQty = bill_entry_in_db.Qty - bill_entry_in_db.MelQty
	bill_entry_in_db.Amount = bill_entry_in_db.Price * float64(bill_entry_in_db.Qty)
	bill_entry_in_db.Desc = core.MakeNullString(bill_entry_in_db.Desc.String+" 部分数量已分批入库", true)
	bill_entry_in_db.FittingCost = (float64(bill_entry_in_db.Qty) / float64(qty_before)) * fitting_cost_before
	bill_entry_in_db.OtherCost = (float64(bill_entry_in_db.Qty) / float64(qty_before)) * other_cost_before
	_, err = db_session.Id(bill_entry_in_db.Id).Cols("qty", "mel_qty", "bne_qty", "amount", "desc", "fitting_cost", "other_cost").Update(&bill_entry_in_db)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	var new_bill_entry cggl_models.PurchaseBillEntryInfo

	new_bill_entry.BillId = bill_entry_in_db.BillId
	new_bill_entry.ProductId = bill_entry_in_db.ProductId
	new_bill_entry.Qty = checkin_qty
	new_bill_entry.MelQty = mel_qty_before - bill_entry_in_db.MelQty
	new_bill_entry.BneQty = checkin_qty - new_bill_entry.MelQty
	new_bill_entry.Price = bill_entry_in_db.Price
	new_bill_entry.Amount = new_bill_entry.Price * float64(new_bill_entry.Qty)
	new_bill_entry.Desc = core.MakeNullString(bill_entry_from_ui.Desc.String+fmt.Sprintf(" 分批入库，原数量 %d", qty_before), true)
	new_bill_entry.StatusId = models.PURCHASE_BILL_STATUS_YRK
	new_bill_entry.CheckinQty = checkin_qty
	new_bill_entry.CheckinTime = core.MakeNullTime(time.Now(), true)
	new_bill_entry.InventoryCost = bill_entry_in_db.InventoryCost
	new_bill_entry.CheckinUserId = login_user.UserId
	new_bill_entry.NeedInvoice = bill_entry_in_db.NeedInvoice
	new_bill_entry.TaxPoint = bill_entry_in_db.TaxPoint
	new_bill_entry.DrawbackPoint = bill_entry_in_db.DrawbackPoint
	new_bill_entry.DrawbackAmount = bill_entry_in_db.DrawbackAmount
	new_bill_entry.FittingCost = (float64(checkin_qty) / float64(qty_before)) * fitting_cost_before
	new_bill_entry.OtherCost = (float64(checkin_qty) / float64(qty_before)) * other_cost_before

	_, err = db_session.Insert(&new_bill_entry)
	core.HandleError(err)

	purchase_bill.FinalAmount = purchase_bill.FinalAmount - new_bill_entry.Amount
	_, err = db_session.Id(purchase_bill.Id).Cols("final_amount").Update(&purchase_bill)
	core.HandleError(err)

	// 计算中国仓成本价
	purchase_bill_service.CalculateBillEntryCnInventoryCost(purchase_bill.Id, db_session)

	// 更新中国仓的库存数量
	inventory_service.AddCnInventoryQty(bill_entry_from_ui.ProductId, bill_entry_from_ui.CheckinQty, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细分批入库成功！",
	})
}
