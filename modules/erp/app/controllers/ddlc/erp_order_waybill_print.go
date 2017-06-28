package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/service/order_service"
	"matrix/core"
	"matrix/modules/erp/service/stock_service"
)

type ErpOrderWaybillPrint struct {
	*revel.Controller
}

func (c ErpOrderWaybillPrint) Index() revel.Result {
	return c.RenderEmpty()
}

func (c ErpOrderWaybillPrint) Preview() revel.Result {
	var db_session = c.DbSession

	var order_id_list = make([]int64, 0)
	c.Params.Bind(&order_id_list, "order_id_list")

	revel.TRACE.Println(order_id_list)

	var mel_stock_id = stock_service.Get_mel_stock_id(db_session)
	var bne_stock_id = stock_service.Get_bne_stock_id(db_session)

	var print_order_list = make([]Print_order_t, 0)
	for _, order_id := range order_id_list {
		var print_order Print_order_t

		var order models.OrderInfo
		_, err := db_session.Id(order_id).Get(&order)
		core.HandleError(err)
		print_order.Order = order

		var order_shipping_address models.OrderShippingAddressInfo
		_, err = db_session.Id(order.AddressId).Get(&order_shipping_address)
		core.HandleError(err)
		print_order.Address = order_shipping_address

		var logistics_company models.LogisticsCompanyInfo
		_, err = db_session.Id(order.CompanyId).Get(&logistics_company)
		core.HandleError(err)
		print_order.LogisticsCompanyName = logistics_company.Name

		print_order.SkuQtyLocList, print_order.PcsCount = get_order_sku_qty_loc_list_for_print(order_id, db_session)

		if order.StockId == mel_stock_id {
			print_order.StockAddressLine1 = "2 Venture Court"
			print_order.StockAddressLine2 = "DANDENONG SOUTH, VIC 3175"
		}

		if order.StockId == bne_stock_id {
			print_order.StockAddressLine1 = "Basement 2/30 Gow Street"
			print_order.StockAddressLine2 = "MOOROOKA, QLD 4105"
		}

		print_order_list = append(print_order_list, print_order)
	}

	c.RenderArgs["print_order_list"] = print_order_list
	return c.RenderTemplate("erp/ddlc/order_waybill_print/order_waybill_print_preview.html")
}

const (
	ORDER_BY_SKU = 100 // 按SKU排序
	ORDER_BY_SL  = 200 // 按库位排序
)

func (c ErpOrderWaybillPrint) SelectOrderByDialog() revel.Result {
	return c.RenderEmpty()
}

func (c ErpOrderWaybillPrint) Print() revel.Result {
	var db_session = c.DbSession

	var order_id_list = make([]int64, 0)
	c.Params.Bind(&order_id_list, "order_id_list")

	revel.TRACE.Println(order_id_list)

	var mel_stock_id = stock_service.Get_mel_stock_id(db_session)
	var bne_stock_id = stock_service.Get_bne_stock_id(db_session)

	var print_order_list = make([]Print_order_t, 0)
	for _, order_id := range order_id_list {
		var print_order Print_order_t

		var order models.OrderInfo
		_, err := db_session.Id(order_id).Get(&order)
		core.HandleError(err)
		print_order.Order = order

		var order_shipping_address models.OrderShippingAddressInfo
		_, err = db_session.Id(order.AddressId).Get(&order_shipping_address)
		core.HandleError(err)
		print_order.Address = order_shipping_address

		var logistics_company models.LogisticsCompanyInfo
		_, err = db_session.Id(order.CompanyId).Get(&logistics_company)
		core.HandleError(err)
		print_order.LogisticsCompanyName = logistics_company.Name

		print_order.SkuQtyLocList, print_order.PcsCount = get_order_sku_qty_loc_list_for_print(order_id, db_session)

		if order.StockId == mel_stock_id {
			print_order.StockAddressLine1 = "2 Venture Court"
			print_order.StockAddressLine2 = "DANDENONG SOUTH, VIC 3175"
		}

		if order.StockId == bne_stock_id {
			print_order.StockAddressLine1 = "Basement 2/30 Gow Street"
			print_order.StockAddressLine2 = "MOOROOKA, QLD 4105"
		}

		print_order_list = append(print_order_list, print_order)
	}

	c.RenderArgs["print_order_list"] = print_order_list
	return c.RenderTemplate("erp/ddlc/order_waybill_print/order_waybill_print.html")
}

type Print_order_t struct {
	Order   models.OrderInfo
	Address models.OrderShippingAddressInfo
	//OrderLogistics       models.OrderLogisticsInfo
	LogisticsCompanyName string
	SkuQtyLocList        []string
	PcsCount             int64
	StockAddressLine1    string
	StockAddressLine2    string
}

func get_order_sku_qty_loc_list_for_print(order_id int64, db_session *xorm.Session) (sku_qty_loc_list_return []string, pcs_count_return int64) {
	/*
		数量*SKU_仓位号
	 */

	var order_sku_qty_loc_list = order_service.Get_order_sku_qty_loc_list(order_id, db_session)

	var order_item_list = make([]string, 0)

	var pcs_count int64 = 0
	for _, order_sku_qty_loc := range order_sku_qty_loc_list {
		order_item_list = append(order_item_list, fmt.Sprintf("%d*%s_%s", order_sku_qty_loc.Quantity, order_sku_qty_loc.Sku, order_sku_qty_loc.StorageLocName))
		pcs_count = pcs_count + order_sku_qty_loc.Quantity
	}

	return order_item_list, pcs_count
}
