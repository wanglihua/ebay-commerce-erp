package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/bill_service"
	"strconv"
	"time"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"matrix/modules/erp/models/cggl_models"
	"strings"
	"matrix/modules/erp/service/purchase_bill_service"
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

	c.RenderArgs["status_list"] = status_list

	return c.RenderTemplate("erp/cggl/purchase_bill/purchase_bill_index.html")
}

type ErpPurchaseBillViewItem struct {
	cggl_models.PurchaseBillInfo `xorm:"extends" json:"pb"`
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
	SELECT pb.*, u.*, s.*
	FROM erp_purchase_bill pb
	    INNER JOIN auth_user u ON pb.user_id = u.id
	    INNER JOIN erp_status s ON pb.status_id = s.id
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("pb.*, u.*, s.*")
	sql_builder.Count("count(pb.id) as bill_count")
	sql_builder.From(`erp_purchase_bill pb
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

	validation.Required(f.PurchaseBill.Date).Message("日期不能为空！")

	validation.Required(f.PurchaseBill.UserId).Message("操作人不能为空！")

	validation.Required(f.PurchaseBill.TransExpenses).Message("运费不能为空！")

	validation.Required(f.PurchaseBill.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPurchaseBill) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	db_session := c.DbSession

	var purchase_bill_id int64
	c.Params.Bind(&purchase_bill_id, "id")

	var purchase_bill cggl_models.PurchaseBillInfo
	if purchase_bill_id != 0 {
		has, err := db_session.Id(purchase_bill_id).Get(&purchase_bill)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的采购单不存在！"})
		}
	} else {
		purchase_bill.Code = bill_service.GetNextBillNo(db_session)
		purchase_bill.Date = core.Date(time.Now())
	}

	var detail_form PurchaseBillDetailForm
	detail_form.PurchaseBill = purchase_bill

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
		purchase_bill.UserId = loginuser.UserId
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YCJ

		affected, err = db_session.Insert(&purchase_bill)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(purchase_bill.Id).Update(&purchase_bill)
		core.HandleError(err)

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

	purchase_bill_service.SetEntriesStatus(purchase_bill_id_list, models.PURCHASE_BILL_STATUS_YCJ, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单提交成功！",
	})
}

func (c ErpPurchaseBill) Approve() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session) // 应该有个专门的审批权限

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

	// 设置为已审批状态
	for _, purchase_bill := range purchase_bill_list {
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_YSP
		_, err := db_session.Id(purchase_bill.Id).Cols("status_id").Update(&purchase_bill)
		core.HandleError(err)
	}

	purchase_bill_service.SetEntriesStatus(purchase_bill_id_list, models.PURCHASE_BILL_STATUS_YCJ, db_session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单审批成功！",
	})
}
