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
	"encoding/json"
	"matrix/modules/erp/service/inventory_service"
	"strings"
	"matrix/modules/erp/service/purchase_bill_log_service"
)

type ErpPurchaseAbnormal struct {
	*revel.Controller
}

func (c ErpPurchaseAbnormal) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_ABNORMAL, c.Session)

	var db_session = c.DbSession

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
	WHERE pbe.status_id = %d OR pbe.status_id = %d;`

	sql_get_bill_code_list = fmt.Sprintf(sql_get_bill_code_list,
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
		models.PURCHASE_BILL_STATUS_YC,
		models.PURCHASE_BILL_STATUS_YCL}).Find(&status_list))

	c.RenderArgs["bill_query_list"] = bill_query_list
	c.RenderArgs["status_list"] = status_list

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

	var bill_code_query string
	c.Params.Bind(&bill_code_query, "bill_code_query")
	bill_code_query = strings.TrimSpace(bill_code_query)

	var status_id_query int64
	c.Params.Bind(&status_id_query, "status_id_query")

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

	filter = core.SqlAnd(filter, fmt.Sprintf("pbe.status_id = %d OR pbe.status_id = %d", models.PURCHASE_BILL_STATUS_YC, models.PURCHASE_BILL_STATUS_YCL)) // 异常 和 异常已处理

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

func (c ErpPurchaseAbnormal) Process() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_ABNORMAL, c.Session)

	var db_session = c.DbSession

	var detail_form PurchaseBillEntryAbnormalProcessDetailForm
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
	if bill_entry_in_db.StatusId != models.PURCHASE_BILL_STATUS_YC {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单明细状态当前不是异常状态，不用进行处理！",
		})
	}

	// 更新明细的状态
	bill_entry_in_db.StatusId = models.PURCHASE_BILL_STATUS_YCL
	bill_entry_in_db.AbnormalProcessResultId = bill_entry_from_ui.AbnormalProcessResultId
	bill_entry_in_db.Desc = bill_entry_from_ui.Desc

	if bill_entry_from_ui.AbnormalProcessResultId == cggl_models.ABNORMAL_PROCESS_RESULT_WCL {
		_, err = db_session.Id(bill_entry_in_db.Id).Cols("desc").Update(&bill_entry_in_db)
		core.HandleError(err)

		return c.RenderJson(map[string]interface{}{
			"success": true,
			"message": "采购单异常未处理！",
		})
	}

	if bill_entry_from_ui.AbnormalProcessResultId == cggl_models.ABNORMAL_PROCESS_RESULT_YCXRK { // 已重新入库
		// 更新中国仓的库存数量
		inventory_service.AddCnInventoryQty(bill_entry_from_ui.ProductId, bill_entry_in_db.Qty-bill_entry_in_db.CheckinQty, db_session)
		bill_entry_in_db.CheckinQty = bill_entry_in_db.Qty
	}

	bill_entry_in_db.CheckinQty = bill_entry_from_ui.CheckinQty
	_, err = db_session.Id(bill_entry_in_db.Id).Cols("status_id", "abnormal_process_result_id", "desc", "checkin_qty").Update(&bill_entry_in_db)
	core.HandleError(err)

	purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, fmt.Sprintf("采购单明细SKU:%s 入库", product.Sku), db_session, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单异常处理成功！",
	})
}
