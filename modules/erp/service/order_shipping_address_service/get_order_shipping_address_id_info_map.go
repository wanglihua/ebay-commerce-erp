package order_shipping_address_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/core"
)

func Get_order_shipping_address_id_info_map(address_id_list []int64, db_session *xorm.Session) map[int64]ddlc_models.OrderShippingAddressInfo {
	var order_shipping_address_id_info_map = make(map[int64]ddlc_models.OrderShippingAddressInfo)

	if len(address_id_list) == 0 {
		return order_shipping_address_id_info_map
	}

	var address_list = make([]ddlc_models.OrderShippingAddressInfo, 0)
	err := db_session.In("id", address_id_list).Find(&address_list)
	core.HandleError(err)

	for _, address := range address_list {
		order_shipping_address_id_info_map[address.Id] = address
	}

	return order_shipping_address_id_info_map
}
