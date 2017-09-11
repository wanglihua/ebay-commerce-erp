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
)

type ErpPurchaseFinance struct {
	*revel.Controller
}

func (c ErpPurchaseFinance) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

	return c.RenderTemplate("erp/cggl/purchase_finance/purchase_finance_index.html")
}

type ErpPurchaseFinanceViewItem struct {
	cggl_models.PurchaseBillInfo `xorm:"extends" json:"pb"`
	auth_models.UserInfo         `xorm:"extends" json:"u"`
	models.StatusInfo            `xorm:"extends" json:"s"`
}

func (c ErpPurchaseFinance) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_FINANCE, c.Session)

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

	filter = core.SqlAnd(filter, fmt.Sprintf("status_id = %d", models.PURCHASE_BILL_STATUS_YSP))

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

	// 设置为已下单状态
	for _, purchase_bill := range purchase_bill_list {
		purchase_bill.StatusId = models.PURCHASE_BILL_STATUS_CGZ
		_, err := db_session.Id(purchase_bill.Id).Cols("status_id").Update(&purchase_bill)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "采购单下单成功！",
	})
}
