package order_shipping_address_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func Get_orders_shipping_address_list(order_list []models.OrderInfo, db_session *xorm.Session) []models.OrderShippingAddressInfo {
	//Get_orders_shipping_address_list

	/*
	SELECT sa.* FROM erp_order_shipping_address sa WHERE sa.order_id IN (1, 2, 3);
	 */
	var order_id_str_list = make([]string, 0)
	for _, order := range order_list {
		order_id_str_list = append(order_id_str_list, fmt.Sprint(order.Id))
	}

	var sql_order_shipping_address_list = `SELECT sa.* FROM erp_order_shipping_address sa WHERE sa.order_id IN (%s)`
	sql_order_shipping_address_list = fmt.Sprintf(sql_order_shipping_address_list, strings.Join(order_id_str_list, ","))

	var order_shipping_address_list = make([]models.OrderShippingAddressInfo, 0)
	err := db_session.Sql(sql_order_shipping_address_list).Find(&order_shipping_address_list)
	core.HandleError(err)

	return order_shipping_address_list
}
