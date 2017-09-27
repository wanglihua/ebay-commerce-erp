package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"fmt"
	"matrix/modules/erp/models/cggl_models"
)

type ErpPurchaseBillLog struct {
	*revel.Controller
}

func (c ErpPurchaseBillLog) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	var db_session = c.DbSession

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	var purchase_bill cggl_models.PurchaseBillInfo
	_, err := db_session.Id(bill_id).Get(&purchase_bill)
	core.HandleError(err)

	c.RenderArgs["bill_code"] = purchase_bill.Code

	return c.RenderTemplate("erp/cggl/purchase_bill_log/purchase_bill_log_index.html")
}

func (c ErpPurchaseBillLog) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_PURCHASE_BILL, c.Session)

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	db_session := c.DbSession

	filter, _, offset, limit := core.GetGridRequestParam(c.Request)

	filter = core.SqlAnd(filter, fmt.Sprintf("bill_id = %d", bill_id))

	query := db_session.
	Where(filter)

	//query extra filter here

	data_query := *query
	data_query = *data_query.Desc("id")

	purchase_bill_log_list := make([]cggl_models.PurchaseBillLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&purchase_bill_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(cggl_models.PurchaseBillLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  purchase_bill_log_list,
		Total: count,
	})
}
