package order_logistics_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
)

func Get_order_logistics_id_info_map(order_id_list []int64, db_session *xorm.Session) map[int64]ddlc_models.OrderLogisticsInfo {
	var order_logistics_id_info_map = make(map[int64]ddlc_models.OrderLogisticsInfo)
	if len(order_id_list) == 0 {
		return order_logistics_id_info_map
	}

	var order_logistics_list = make([]ddlc_models.OrderLogisticsInfo, 0)
	err := db_session.In("order_id", order_id_list).Find(&order_logistics_list)
	core.HandleError(err)

	for _, order_logistics := range order_logistics_list {
		order_logistics_id_info_map[order_logistics.Id] = order_logistics
	}

	return order_logistics_id_info_map
}
