package order_history_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"matrix/core"
)

func Insert_order_history(order models.OrderInfo, description string, db_session *xorm.Session) {
	var order_history models.OrderHistoryInfo

	order_history.OrderId = order.Id
	order_history.OrderCode = order.Code
	order_history.StatusId = order.StatusId

	var order_status models.StatusInfo
	_, err := db_session.Id(order.StatusId).Get(&order_status)
	core.HandleError(err)

	order_history.StatusName = order_status.Name

	_, err = db_session.Insert(&order_history)
	core.HandleError(err)
}
