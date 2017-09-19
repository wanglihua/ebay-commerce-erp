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
)

type ErpPurchaseCheckin struct {
	*revel.Controller
}

func (c ErpPurchaseCheckin) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_CHECKIN, c.Session)

	var db_session = c.DbSession

	// 获得有明细状态是 采购中 的所有采购单的批次号
	/*
		SELECT DISTINCT pb.code AS bill_code
		FROM erp_purchase_bill    pb
		    INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
		WHERE pbe.status_id = 0;
	*/
	var sql_get_bill_code_list = `
	SELECT DISTINCT pb.code AS bill_code
	FROM erp_purchase_bill pb
	    INNER JOIN erp_purchase_bill_entry pbe ON pbe.bill_id = pb.id
	WHERE pbe.status_id = %d;`

	sql_get_bill_code_list = fmt.Sprintf(sql_get_bill_code_list, models.PURCHASE_BILL_STATUS_CGZ)

	type bill_code_struct struct {
		BillCode string `xorm:"'bill_code'"`
	}

	var bill_code_struct_list = make([]bill_code_struct, 0)
	core.HandleError(db_session.Sql(sql_get_bill_code_list).Find(&bill_code_struct_list))

	var bill_code_list = make([]string, 0)
	for _, bill_code_struct := range bill_code_struct_list {
		bill_code_list = append(bill_code_list, bill_code_struct.BillCode)
	}

	c.RenderArgs["bill_code_list"] = bill_code_list

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

	filter = core.SqlAnd(filter, fmt.Sprintf("pbe.status_id = %d", models.PURCHASE_BILL_STATUS_CGZ))

	if bill_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("pb.code = '%s'", bill_code_query))
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

	// 检查状态
	if bill_entry_in_db.StatusId != models.PURCHASE_BILL_STATUS_CGZ {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单明细状态当前不是采购中状态，不能入库！",
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
		var purchase_bill cggl_models.PurchaseBillInfo
		_, err = db_session.Id(bill_id).Get(&purchase_bill)
		core.HandleError(err)

		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YRK
		_, err = db_session.Id(purchase_bill.Id).Cols("status_id").Update(&purchase_bill)
		core.HandleError(err)
	}

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

	// 更新中国仓的库存数量
	for _, bill_entry := range bill_entry_list {
		inventory_service.AddCnInventoryQty(bill_entry.ProductId, bill_entry.CheckinQty, db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细入库成功！",
	})
}
