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

func (c ErpPurchaseCheckin) Checkin() revel.Result {
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

		/*
			bill_entry.StatusId = models.PURCHASE_BILL_STATUS_YRK
			_, err := db_session.Id(bill_entry.Id).Cols("status_id").Update(&bill_entry)
			core.HandleError(err)
		*/

		purchase_bill_service.SetEntryStatusAndUpdateBillStatus(bill_entry, models.PURCHASE_BILL_STATUS_YRK, db_session)
	}

	// 更新中国仓的库存数量
	for _, bill_entry := range bill_entry_list {
		inventory_service.AddCnInventoryQty(bill_entry.ProductId, bill_entry.Qty, db_session)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单下单成功！",
	})
}
