package order_combine_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/core"
)

func insert_combine_order_address_from_first_order(combine_order_id int64, first_order_id int64, db_session *xorm.Session) {

	//复制第一个订单的地址，作为合并后的订单的地址
	var order_ship_address_list = make([]ddlc_models.OrderShippingAddressInfo, 0)
	err := db_session.Where("order_id = ?", first_order_id).Find(&order_ship_address_list)
	core.HandleError(err)

	for _, order_ship_address := range order_ship_address_list {
		order_ship_address.Id = 0
		order_ship_address.OrderId = combine_order_id

		_, err := db_session.Insert(&order_ship_address)
		core.HandleError(err)
	}
}
