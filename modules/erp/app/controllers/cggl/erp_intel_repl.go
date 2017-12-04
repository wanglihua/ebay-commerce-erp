package controllers

import (
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models"
	"strings"
	"matrix/modules/erp/service/bill_service"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/service/intel_repl_service"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/service/purchase_bill_service"
)

type ErpIntelRepl struct {
	*revel.Controller
}

func (c ErpIntelRepl) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	var db_session = c.DbSession

	/*
	var stock_list = make([]sppz_models.StockInfo, 0)
	core.HandleError(db_session.Where(fmt.Sprintf("stock_type_id = %d", sppz_models.STOCK_TYPE_AU)).Find(&stock_list))
	c.RenderArgs["stock_list"] = stock_list
	*/

	var status_list = make([]models.StatusInfo, 0)
	core.HandleError(db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_INTEL_REPL)).Find(&status_list))
	c.RenderArgs["status_list"] = status_list

	var product_status_list = make([]models.StatusInfo, 0)
	core.HandleError(db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT)).Find(&product_status_list))
	c.RenderArgs["product_status_list"] = product_status_list

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.OrderBy("login_name ASC").Find(&user_list))

	c.RenderArgs["sale_user_list"] = user_list
	c.RenderArgs["dev_user_list"] = user_list

	return c.RenderTemplate("erp/cggl/intel_repl/intel_repl_index.html")
}

type intel_repl_view_item_t struct {
	IntelRepl     cggl_models.IntelReplInfo `xorm:"extends" json:"ir"`
	Product       kcgl_models.ProductInfo   `xorm:"extends" json:"p"`
	ProductStatus models.StatusInfo         `xorm:"extends" json:"ps"`
	Status        models.StatusInfo         `xorm:"extends" json:"s"`
	Supplier      sppz_models.SupplierInfo  `xorm:"extends" json:"sup"`
}

func (c ErpIntelRepl) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	/*
sku_query:
min_profit_query:
max_profit_query:
status_id_query
	 */

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")
	sku_query = strings.TrimSpace(sku_query)

	var min_profit_query_str string
	c.Params.Bind(&min_profit_query_str, "min_profit_query")
	min_profit_query_str = strings.TrimSpace(min_profit_query_str)
	var min_profit_query float64
	var err error
	if min_profit_query_str != "" {
		min_profit_query, err = strconv.ParseFloat(min_profit_query_str, 64)
		if err != nil {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "最小利润率输入有误！",
			})
		}
	}

	var max_profit_query_str string
	c.Params.Bind(&max_profit_query_str, "max_profit_query")
	max_profit_query_str = strings.TrimSpace(max_profit_query_str)
	var max_profit_query float64

	if max_profit_query_str != "" {
		max_profit_query, err = strconv.ParseFloat(max_profit_query_str, 64)
		if err != nil {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "最大利润率输入有误！",
			})
		}
	}

	var product_status_id_query int64
	c.Params.Bind(&product_status_id_query, "product_status_id_query")

	var status_id_query int64
	c.Params.Bind(&status_id_query, "status_id_query")

	var sale_user_id_query int64
	c.Params.Bind(&sale_user_id_query, "sale_user_id_query")

	var dev_user_id_query int64
	c.Params.Bind(&dev_user_id_query, "dev_user_id_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*

SELECT ir.*, p.*, ps.*, s.*, sup.*
FROM erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_status ps ON p.status_id = ps.id
    INNER JOIN erp_status s ON ir.status_id = s.id
    LEFT JOIN erp_supplier sup ON p.supplier_id = sup.id;

SELECT count(ir.id) AS count
FROM erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_status ps ON p.status_id = ps.id
    INNER JOIN erp_status s ON ir.status_id = s.id
    LEFT JOIN erp_supplier sup ON p.supplier_id = sup.id;

	 */

	var sb = core.NewSqlBuilder()
	sb.Select("ir.*, p.*, ps.*, s.*, sup.*")
	sb.Count("count(ir.id) AS count")
	sb.From(`erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_status ps ON p.status_id = ps.id
    INNER JOIN erp_status s ON ir.status_id = s.id
	LEFT JOIN erp_supplier sup ON p.supplier_id = sup.id`)

	if sku_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.sku = '%s'", sku_query))
	}

	if min_profit_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.profit >= %f", min_profit_query))
	}

	if max_profit_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.profit <= %f", max_profit_query))
	}

	if product_status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.status_id = %d", product_status_id_query))
	}

	if status_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.status_id = %d", status_id_query))
	}

	if sale_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.sale_user_id = %d", sale_user_id_query))
	}

	if dev_user_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.dev_user_id = %d", dev_user_id_query))
	}

	sb.Where(filter)

	sb.Limit(offset, limit)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("ir.id DESC")
	}

	var intel_repl_list = make([]intel_repl_view_item_t, 0)
	core.HandleError(db_session.Sql(sb.SqlList()).Find(&intel_repl_list))

	var count core.Count
	_, err = db_session.Sql(sb.SqlCount()).Get(&count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  intel_repl_list,
		Total: count.Count,
	})
}

type IntelReplDetailForm struct {
	IntelRepl cggl_models.IntelReplInfo `json:"intel_repl"`
}

func (f *IntelReplDetailForm) IsCreate() bool {
	return f.IntelRepl.Id == 0
}

func (f *IntelReplDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.IntelRepl.ProductId).Message("商品不能为空！")

	validation.Required(f.IntelRepl.Qty).Message("数量不能为空！")

	validation.Required(f.IntelRepl.Profit).Message("毛利率不能为空！")

	validation.Required(f.IntelRepl.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

/*
func (c ErpIntelRepl) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	var intel_repl_id int64
	c.Params.Bind(&intel_repl_id, "id")

	var intel_repl cggl_models.IntelReplInfo
	if intel_repl_id != 0 {
		has, err := db_session.Id(intel_repl_id).Get(&intel_repl)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的智能补货不存在！"})
		}
	}

	var detail_form IntelReplDetailForm
	detail_form.IntelRepl = intel_repl

	return c.RenderJson(detail_form)
}

func (c ErpIntelRepl) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	var detail_form IntelReplDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	intel_repl := detail_form.IntelRepl

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&intel_repl)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(intel_repl.Id).Update(&intel_repl)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}
*/

func (c ErpIntelRepl) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	intel_repl_id_list := make([]int64, 0)
	c.Params.Bind(&intel_repl_id_list, "id_list")

	affected, err := db_session.In("id", intel_repl_id_list).Delete(new(cggl_models.IntelReplInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpIntelRepl) Purchase() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	var db_session = c.DbSession

	var intel_repl_id_list = make([]int64, 0)
	c.Params.Bind(&intel_repl_id_list, "id_list")

	var intel_repl_list = make([]cggl_models.IntelReplInfo, 0)
	core.HandleError(db_session.In("id", intel_repl_id_list).Find(&intel_repl_list))

	if len(intel_repl_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择智能补货条目！",
		})
	}

	for _, intel_repl := range intel_repl_list {
		if intel_repl.StatusId == models.INTEL_REPL_STATUS_YCN {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "已采购状态的智能补货条目不能再次采购！",
			})
		}
	}

	/*
		SELECT DISTINCT ifnull(p.supplier_id, 0) AS supplier_id
		FROM erp_product p
		WHERE p.id IN (23879, 23876, 17920);
	 */

	var product_id_str_list = make([]string, 0)
	for _, intel_repl := range intel_repl_list {
		product_id_str_list = append(product_id_str_list, fmt.Sprint(intel_repl.ProductId))
	}

	type supplier_id_t struct {
		SupplierId int64 `xorm:"'supplier_id'"`
	}

	var supplier_id_list = make([]supplier_id_t, 0)
	var sql_supplier_list = `SELECT DISTINCT ifnull(p.supplier_id, 0) AS supplier_id
FROM erp_product p
WHERE p.id IN (%s);`
	sql_supplier_list = fmt.Sprintf(sql_supplier_list, strings.Join(product_id_str_list, ","))

	err := db_session.Sql(sql_supplier_list).Find(&supplier_id_list)
	core.HandleError(err)

	if len(supplier_id_list) != 1 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品的供应商不是同一供应商的，不能生成一个采购单！",
		})
	}

	var supplier_id = supplier_id_list[0].SupplierId

	// 供应商不能为空
	if supplier_id == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单的供应商不能为空，请先设置SKU供应商！",
		})
	}

	// var mel_stock_id = stock_service.Get_mel_stock_id(db_session)
	// var bne_stock_id = stock_service.Get_bne_stock_id(db_session)

	var login_user = core.GetLoginUser(c.Session)

	// 生成采购单
	var purchase_bill cggl_models.PurchaseBillInfo
	purchase_bill.Code = bill_service.GetNextBillNo(bill_service.ENTITY_ID_PURCHASE_BILL, db_session, c.Session)
	purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YCJ
	purchase_bill.ApplyUserId = login_user.UserId
	purchase_bill.SupplierId = supplier_id

	_, err = db_session.Insert(&purchase_bill)
	core.HandleError(err)

	// 生成采购单明细，这个循环过程中会自动实现 MEL，BNE 仓库条目的合并
	for _, intel_repl := range intel_repl_list {
		var purchase_bill_entry cggl_models.PurchaseBillEntryInfo

		purchase_bill_entry.BillId = purchase_bill.Id
		purchase_bill_entry.ProductId = intel_repl.ProductId

		purchase_bill_entry.MelQty = intel_repl.MelQty
		purchase_bill_entry.BneQty = intel_repl.BneQty
		purchase_bill_entry.Qty = intel_repl.Qty

		purchase_bill_entry.StatusId = purchase_bill.StatusId
		purchase_bill_entry.TypeId = cggl_models.PURCHASE_BILL_TYPE_NORMAL

		// 获得商品的建议采购价
		purchase_bill_entry.Price = purchase_bill_service.Get_propose_purchase_price(intel_repl.ProductId, db_session)

		// 计算金额
		purchase_bill_entry.Amount = purchase_bill_entry.Price * float64(purchase_bill_entry.Qty)

		var product kcgl_models.ProductInfo
		_, err = db_session.Id(intel_repl.ProductId).Get(&product)
		core.HandleError(err)

		purchase_bill_entry.NeedInvoice = product.NeedInvoice
		purchase_bill_entry.TaxPoint = product.TaxPoint
		purchase_bill_entry.DrawbackPoint = product.DrawbackPoint
		purchase_bill_entry.DrawbackAmount = purchase_bill_service.CalculateBillEntryDrawbackAmount(purchase_bill_entry)

		// 保存采购明细
		_, err = db_session.Insert(&purchase_bill_entry)
		core.HandleError(err)

		intel_repl.StatusId = models.INTEL_REPL_STATUS_YCN
		_, err = db_session.Id(intel_repl.Id).Cols("status_id").Update(&intel_repl)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("生成采购单 %s 成功！", purchase_bill.Code),
	})
}

func (c ErpIntelRepl) Calculate() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	var db_session = c.DbSession

	intel_repl_service.Calculate_intel_repl(db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "智能补货计算完成！",
	})
}
