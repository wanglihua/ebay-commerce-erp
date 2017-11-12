package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/erp/service/purchase_bill_log_service"
	"matrix/modules/erp/service/purchase_bill_service"
	"math"
	"time"
	"strings"
)

type ErpPurchaseBillEntry struct {
	*revel.Controller
}

func (c ErpPurchaseBillEntry) Index() revel.Result {
	// pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)
	pass_service.Check_login_user_pass_or([]int64{pass_res.CGGL_PURCHASE_BILL, pass_res.CGGL_PURCHASE_FINANCE}, c.Session)

	var db_session = c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
	core.HandleError(err)

	c.RenderArgs["bill_code"] = purchase_bill.Code
	return c.RenderTemplate("erp/cggl/purchase_bill_entry/purchase_bill_entry_index.html")
}

type PurchaseBillEntryViewItem struct {
	cggl_models.PurchaseBillEntryInfo     `xorm:"extends" json:"pbe"`
	cggl_models.PurchaseBillEntryTypeInfo `xorm:"extends" json:"et"`
	kcgl_models.ProductInfo               `xorm:"extends" json:"p"`
	models.StatusInfo                     `xorm:"extends" json:"s"`
}

func (c ErpPurchaseBillEntry) ListData() revel.Result {
	// pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)
	pass_service.Check_login_user_pass_or([]int64{pass_res.CGGL_PURCHASE_BILL, pass_res.CGGL_PURCHASE_FINANCE}, c.Session)

	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	var purchase_bill cggl_models.PurchaseBillInfo
	if purchase_bill_id != 0 {
		has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
		}
	}

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT pbe.*, et.*, p.*, s.*
FROM erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill_entry_type et ON et.id = pbe.type_id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_status s ON pbe.status_id = s.id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("pbe.*, et.*, p.*, s.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_purchase_bill_entry pbe
    INNER JOIN erp_purchase_bill_entry_type et ON et.id = pbe.type_id
    INNER JOIN erp_product p ON pbe.product_id = p.id
    INNER JOIN erp_status s ON pbe.status_id = s.id`)

	sb.Where(core.SqlAnd(filter, fmt.Sprintf("bill_id = %d", purchase_bill_id)))

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("pbe.id ASC")
	}

	sb.Limit(offset, limit)

	purchase_bill_entry_list := make([]PurchaseBillEntryViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&purchase_bill_entry_list)
	core.HandleError(err)

	var purchase_bill_entry_count struct {
		Count int64 `xorm:"'entry_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&purchase_bill_entry_count)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"purchase_bill_code": func() string {
			if purchase_bill.Code == "" {
				return "请先单击采购单行"
			}
			return purchase_bill.Code
		}(),
		"list_data": core.GridResult{
			Data:  purchase_bill_entry_list,
			Total: purchase_bill_entry_count.Count,
		},
	})
}

type PurchaseBillEntryDetailForm struct {
	PurchaseBill      cggl_models.PurchaseBillInfo      `json:"purchase_bill"`
	PurchaseBillEntry cggl_models.PurchaseBillEntryInfo `json:"purchase_bill_entry"`
	ProductName       string                            `json:"product_name"`
}

func (f *PurchaseBillEntryDetailForm) IsCreate() bool {
	return f.PurchaseBillEntry.Id == 0
}

func (f *PurchaseBillEntryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.PurchaseBillEntry.BillId).Message("批次不能为空！")

	validation.Required(f.PurchaseBillEntry.ProductId).Message("商品不能为空！")

	validation.Required(f.PurchaseBillEntry.Qty).Message("数量不能为空！")

	validation.Required(f.PurchaseBillEntry.Price).Message("单价不能为空！")

	validation.Required(f.PurchaseBillEntry.Amount).Message("金额不能为空！")

	if f.PurchaseBillEntry.Desc.Valid == true {
		validation.MinSize(f.PurchaseBillEntry.Desc.String, 2).Message("备注长度不能小于2！")
	}
	if f.PurchaseBillEntry.Desc.Valid == true {
		validation.MaxSize(f.PurchaseBillEntry.Desc.String, 400).Message("备注长度不能大于400！")
	}

	return validation.HasErrors() == false
}

func (c ErpPurchaseBillEntry) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "purchase_bill_id")

	if purchase_bill_id == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "请先选择采购单！"})
	}

	var purchase_bill cggl_models.PurchaseBillInfo
	has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
	}

	var purchase_bill_entry_id int64
	c.Params.Bind(&purchase_bill_entry_id, "id")

	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	var product_name string
	if purchase_bill_entry_id != 0 {
		has, err := db_session.Id(purchase_bill_entry_id).Get(&purchase_bill_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单明细不存在！"})
		}

		var product kcgl_models.ProductInfo
		product_id := purchase_bill_entry.ProductId
		_, err = db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		product_name = product.CnName
	}

	var detail_form PurchaseBillEntryDetailForm
	detail_form.PurchaseBill = purchase_bill
	detail_form.PurchaseBillEntry = purchase_bill_entry
	detail_form.ProductName = product_name

	return c.RenderJson(detail_form)
}

func (c ErpPurchaseBillEntry) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	var detail_form PurchaseBillEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	purchase_bill_entry := detail_form.PurchaseBillEntry
	purchase_bill_entry.BillId = detail_form.PurchaseBill.Id

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err = db_session.Id(purchase_bill_entry.BillId).Get(&purchase_bill)
	core.HandleError(err)

	if purchase_bill_service.CanPurchaseBillModify(purchase_bill, c.Session) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单当前状态不能新增和修改明细！",
		})
	}

	var qty = purchase_bill_entry.MelQty + purchase_bill_entry.BneQty
	purchase_bill_entry.Qty = qty

	purchase_bill_entry.Amount = purchase_bill_entry.Price * float64(qty)

	if purchase_bill_entry.ProductId == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请先选择商品SKU！",
		})
	}

	if qty == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购数量不能为零！",
		})
	}

	if purchase_bill_entry.Price <= 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请输入采购价格！",
		})
	}

	var affected int64
	if detail_form.IsCreate() {
		billId_count, err := db_session.Where("bill_id = ? and product_id = ?", purchase_bill_entry.BillId, purchase_bill_entry.ProductId).Count(new(cggl_models.PurchaseBillEntryInfo))
		core.HandleError(err)
		if billId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，批次中该商品已存在！"})
		}

		purchase_bill_entry.StatusId = purchase_bill.StatusId // 设置成采购单的当前状态

		affected, err = db_session.Insert(&purchase_bill_entry)
		core.HandleError(err)

		purchase_bill_entry.DrawbackAmount = purchase_bill_service.CalculateBillEntryDrawbackAmount(purchase_bill_entry)

		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购明细新增", db_session, c.Session)
	} else {
		billId_count, err := db_session.Where("id <> ? and bill_id = ? and product_id = ?", purchase_bill_entry.Id, purchase_bill_entry.BillId, purchase_bill_entry.ProductId).Count(new(cggl_models.PurchaseBillEntryInfo))
		core.HandleError(err)
		if billId_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，批次、商品已存在！"})
		}

		purchase_bill_entry.DrawbackAmount = purchase_bill_service.CalculateBillEntryDrawbackAmount(purchase_bill_entry)

		affected, err = db_session.Id(purchase_bill_entry.Id).Update(&purchase_bill_entry)
		core.HandleError(err)

		purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购明细修改", db_session, c.Session)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	// 更新采购单总金额
	purchase_bill_service.UpdatePurchaseBillTotalAmount(purchase_bill.Id, db_session)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPurchaseBillEntry) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	purchase_bill_entry_id_list := make([]int64, 0)
	c.Params.Bind(&purchase_bill_entry_id_list, "id_list")

	if len(purchase_bill_entry_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择要删除的数据!",
		})
	}

	var entry_id = purchase_bill_entry_id_list[0]
	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err = db_session.Id(purchase_bill_entry.BillId).Get(&purchase_bill)
	core.HandleError(err)

	if purchase_bill_service.CanPurchaseBillModify(purchase_bill, c.Session) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单当前状态不能删除明细！",
		})
	}

	purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, "采购单明细删除", db_session, c.Session)

	affected, err := db_session.In("id", purchase_bill_entry_id_list).Delete(new(cggl_models.PurchaseBillEntryInfo))
	core.HandleError(err)

	// 更新采购单总金额
	purchase_bill_service.UpdatePurchaseBillTotalAmount(purchase_bill.Id, db_session)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpPurchaseBillEntry) ProductSelectListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type ErpProductViewItem struct {
		kcgl_models.ProductInfo         `xorm:"extends" json:"prd"`
		kcgl_models.ProductCateInfo     `xorm:"extends" json:"cate"`
		sppz_models.SupplierInfo        `xorm:"extends" json:"sup"`
		sppz_models.PackageMaterialInfo `xorm:"extends" json:"pm"`
		models.StatusInfo               `xorm:"extends" json:"st"`
	}

	/*
SELECT prd.*, cate.*, sup.*, pm.*, st.*
FROM erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("prd.*, cate.*, sup.*, pm.*, st.*")
	sql_builder.Count("count(prd.id) as product_count")
	sql_builder.From(`erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id`)

	sql_builder.Where(filter)

	//query extra filter here

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("prd.id ASC")
	}

	sql_builder.Limit(offset, limit)

	product_list := make([]ErpProductViewItem, 0, limit)
	err := db_session.Sql(sql_builder.SqlList()).Find(&product_list)
	core.HandleError(err)

	var product_count struct {
		Count int64 `xorm:"'product_count'"`
	}
	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: product_count.Count,
	})
}

func (c ErpPurchaseBillEntry) PurchaseInTurn() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	var db_session = c.DbSession

	var detail_form PurchaseBillEntryDetailForm
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

	// 检查类型
	if bill_entry_in_db.TypeId == cggl_models.PURCHASE_BILL_TYPE_IN_TURN {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单明细类型当前已经是分批采购类型，不能再进行分批采购！",
		})
	}

	// 检查状态
	if bill_entry_in_db.StatusId != models.PURCHASE_BILL_STATUS_CGZ {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "采购单明细状态当前不是采购中状态，不能进行分批采购！",
		})
	}

	if bill_entry_from_ui.CheckinQty > bill_entry_in_db.Qty { // 从界面上获取的分批采购的数量大于原采明细的数量
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "分批采购数量不能大于原采购数量！",
		})
	}

	var is_all_checkin = bill_entry_from_ui.CheckinQty == bill_entry_in_db.Qty

	// 检查是不是分批入库
	if is_all_checkin == true {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "全部数量采购，请使用正常采购流程而不是分批采购！",
		})
	}

	// 插入采购单历史，入库
	purchase_bill_log_service.Insert_purchase_bill_log(purchase_bill, fmt.Sprintf("采购单明细SKU:%s 分批采购，%d个", product.Sku, bill_entry_from_ui.CheckinQty), db_session, c.Session)

	var qty_before = bill_entry_in_db.Qty
	var mel_qty_before = bill_entry_in_db.MelQty
	// var bne_qty_before = bill_entry_in_db.BneQty
	var fitting_cost_before = bill_entry_in_db.FittingCost
	var other_cost_before = bill_entry_in_db.OtherCost
	var checkin_qty = bill_entry_from_ui.CheckinQty
	var checkin_rate = float64(checkin_qty) / float64(qty_before)
	var pay_amount = bill_entry_from_ui.PayAmount // 付款金额这里输入，财务下单

	bill_entry_in_db.Qty = qty_before - checkin_qty
	bill_entry_in_db.MelQty = bill_entry_in_db.MelQty - int64(math.Floor(float64(mel_qty_before)*checkin_rate+0.5)) // 四舍五入
	bill_entry_in_db.BneQty = bill_entry_in_db.Qty - bill_entry_in_db.MelQty
	bill_entry_in_db.Amount = bill_entry_in_db.Price * float64(bill_entry_in_db.Qty)

	var purchase_in_turn_desc = " 部分数量已分批采购"
	if strings.Contains(bill_entry_in_db.Desc.String, purchase_in_turn_desc) == false {
		bill_entry_in_db.Desc = core.MakeNullString(bill_entry_in_db.Desc.String+purchase_in_turn_desc, true)
	}
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
	new_bill_entry.PayAmount = pay_amount // 付款金额这里输入，财务下单
	new_bill_entry.Desc = core.MakeNullString(bill_entry_from_ui.Desc.String+fmt.Sprintf(" 分批采购，原数量 %d", qty_before), true)
	new_bill_entry.StatusId = models.PURCHASE_BILL_STATUS_YSP // 已审批
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

	// 更新尾款 // 付款由财务在财务管理维护, 更新尾款也放到那边 // 付款在采购管理里面做，扣除尾款在财务管理里面进行
	/*
	purchase_bill.FinalAmount = purchase_bill.FinalAmount - pay_amount
	_, err = db_session.Id(purchase_bill.Id).Cols("final_amount").Update(&purchase_bill)
	core.HandleError(err)
	*/

	// 计算中国仓成本价
	purchase_bill_service.CalculateBillEntryCnInventoryCost(purchase_bill.Id, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购明细分批采购成功！",
	})
}
