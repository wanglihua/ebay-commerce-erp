package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/sppz_models"
	"fmt"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models"
	"strings"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/service/bill_service"
)

type ErpIntelRepl struct {
	*revel.Controller
}

func (c ErpIntelRepl) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	var db_session = c.DbSession

	var stock_list = make([]sppz_models.StockInfo, 0)
	core.HandleError(db_session.Where(fmt.Sprintf("stock_type_id = %d", sppz_models.STOCK_TYPE_AU)).Find(&stock_list))

	c.RenderArgs["stock_list"] = stock_list

	return c.RenderTemplate("erp/cggl/intel_repl/intel_repl_index.html")
}

type intel_repl_view_item_t struct {
	IntelRepl cggl_models.IntelReplInfo `xorm:"extends" json:"ir"`
	Product   kcgl_models.ProductInfo   `xorm:"extends" json:"p"`
	Stock     sppz_models.StockInfo     `xorm:"extends" json:"st"`
	Status    models.StatusInfo         `xorm:"extends" json:"s"`
}

func (c ErpIntelRepl) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	/*
sku_query:
stock_id_query:0
min_profit_query:
max_profit_query:
	 */

	var sku_query string
	c.Params.Bind(&sku_query, "sku_query")
	sku_query = strings.TrimSpace(sku_query)

	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

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

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*

SELECT ir.*, p.*, st.*, s.*
FROM erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_stock st on ir.stock_id = st.id
    INNER JOIN erp_status s ON ir.status_id = s.id;

SELECT count(ir.id) AS count
FROM erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_stock st on ir.stock_id = st.id
    INNER JOIN erp_status s ON ir.status_id = s.id;

	 */

	var sb = core.NewSqlBuilder()
	sb.Select("ir.*, p.*, st.*, s.*")
	sb.Count("count(ir.id) AS count")
	sb.From(`erp_intel_repl ir
    INNER JOIN erp_product p ON ir.product_id = p.id
    INNER JOIN erp_stock st on ir.stock_id = st.id
    INNER JOIN erp_status s ON ir.status_id = s.id`)

	if sku_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("p.sku = '%s'", sku_query))
	}

	if stock_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.stock_id = %d", stock_id_query))
	}

	if min_profit_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.profit >= %f", min_profit_query))
	}

	if max_profit_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("ir.profit <= %f", max_profit_query))
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

	validation.Required(f.IntelRepl.StockId).Message("仓库不能为空！")

	validation.Required(f.IntelRepl.Qty).Message("数量不能为空！")

	validation.Required(f.IntelRepl.Profit).Message("毛利率不能为空！")

	validation.Required(f.IntelRepl.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

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

	var mel_stock_id = stock_service.Get_mel_stock_id(db_session)
	var bne_stock_id = stock_service.Get_bne_stock_id(db_session)

	var login_user = core.GetLoginUser(c.Session)

	// 生成采购单
	var purchase_bill cggl_models.PurchaseBillInfo
	purchase_bill.Code = bill_service.GetNextBillNo(bill_service.ENTITY_ID_PURCHASE_BILL, db_session, c.Session)
	purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YCJ
	purchase_bill.ApplyUserId = login_user.UserId

	_, err := db_session.Insert(&purchase_bill)
	core.HandleError(err)

	// 生成采购单明细
	for _, intel_repl := range intel_repl_list {
		var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
		has, err := db_session.Where(fmt.Sprintf("bill_id = %d AND product_id = %d", purchase_bill.Id, intel_repl.ProductId)).Get(&purchase_bill_entry)
		core.HandleError(err)

		purchase_bill_entry.BillId = purchase_bill.Id
		purchase_bill_entry.ProductId = intel_repl.ProductId

		if intel_repl.StockId == mel_stock_id {
			purchase_bill_entry.MelQty = intel_repl.Qty
		}

		if intel_repl.StockId == bne_stock_id {
			purchase_bill_entry.BneQty = intel_repl.Qty
		}

		purchase_bill_entry.Qty = purchase_bill_entry.MelQty + purchase_bill_entry.BneQty
		purchase_bill_entry.StatusId = purchase_bill.StatusId

		if has == true {
			var col_list = []string{"bill_id", "product_id", "mel_qty", "bne_qty", "qty", "status_id"}
			_, err = db_session.Id(purchase_bill_entry.Id).Cols(col_list...).Update(&purchase_bill_entry)
			core.HandleError(err)
		} else {
			_, err = db_session.Insert(&purchase_bill_entry)
			core.HandleError(err)
		}

		intel_repl.StatusId = models.INTEL_REPL_STATUS_YCG
		_, err = db_session.Id(intel_repl.Id).Cols("status_id").Update(&intel_repl)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("生成采购单 %s 成功！", purchase_bill.Code),
	})
}
