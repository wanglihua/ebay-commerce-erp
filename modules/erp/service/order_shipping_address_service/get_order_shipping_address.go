package order_shipping_address_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ddlc_models"
	"fmt"
	"matrix/core"
)

/*
优先取得是 eBay 的送货地址

 */

func Get_order_shipping_address(order_id int64, db_session *xorm.Session) (ddlc_models.OrderShippingAddressInfo, ddlc_models.OrderShippingAddressInfo, ddlc_models.OrderShippingAddressInfo) {

	var shipping_address ddlc_models.OrderShippingAddressInfo // 取 ebay_shipping_address 中的一个， 优先取 paypal_shipping_address

	// 这里面取得是PayPal的地址  PayPal
	var paypal_shipping_address ddlc_models.OrderShippingAddressInfo
	has_paypal_address, err := db_session.Where(fmt.Sprintf("order_id = %d AND address_owner='PayPal'", order_id)).OrderBy("update_time DESC").Limit(1, 0).Get(&paypal_shipping_address)
	core.HandleError(err)
	if has_paypal_address == true {
		shipping_address = paypal_shipping_address
	}

	// 这里面取得是平台的地址  eBay
	var ebay_shipping_address ddlc_models.OrderShippingAddressInfo
	has_ebay_address, err := db_session.Where(fmt.Sprintf("order_id = %d AND address_owner='eBay'", order_id)).OrderBy("update_time DESC").Limit(1, 0).Get(&ebay_shipping_address)
	core.HandleError(err)
	if has_ebay_address == true {
		shipping_address = ebay_shipping_address
	}

	if has_ebay_address == false && has_paypal_address == false {
		_, err := db_session.Where(fmt.Sprintf("order_id = %d", order_id)).OrderBy("update_time DESC").Limit(1, 0).Get(&ebay_shipping_address)
		core.HandleError(err)

		ebay_shipping_address.AddressOwner = "eBay"
		shipping_address = ebay_shipping_address
	}

	return ebay_shipping_address, paypal_shipping_address, shipping_address
}
