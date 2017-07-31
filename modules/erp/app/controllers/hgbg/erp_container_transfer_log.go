package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models"
	"fmt"
)

type ErpContainerTransferLog struct {
	*revel.Controller
}

func (c ErpContainerTransferLog) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	var db_session = c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var ct models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&ct)
	core.HandleError(err)

	c.RenderArgs["batch_code"] = ct.BatchCode

	return c.RenderTemplate("erp/hgbg/container_transfer_log/container_transfer_log_index.html")
}

func (c ErpContainerTransferLog) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	db_session := c.DbSession

	filter, _, offset, limit := core.GetGridRequestParam(c.Request)

	filter = core.SqlAnd(filter, fmt.Sprintf("transfer_id = %d", transfer_id))

	query := db_session.
	Where(filter)

	//query extra filter here

	data_query := *query
	data_query = *data_query.Desc("id")

	container_transfer_log_list := make([]models.ContainerTransferLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&container_transfer_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ContainerTransferLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  container_transfer_log_list,
		Total: count,
	})
}
