package order_history_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
)

func Insert_order_history(order ddlc_models.OrderInfo, login_user_id int64, description string, db_session *xorm.Session) {
	var order_history ddlc_models.OrderHistoryInfo

	order_history.OrderId = order.Id
	order_history.OrderCode = order.Code
	order_history.StatusId = order.StatusId

	var order_status models.StatusInfo
	_, err := db_session.Id(order.StatusId).Get(&order_status)
	core.HandleError(err)

	order_history.StatusName = order_status.Name
	order_history.Description = description
	order_history.UserId = login_user_id

	_, err = db_session.Insert(&order_history)
	core.HandleError(err)
}
