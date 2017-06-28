package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/service/order_service"
	"matrix/core"
	"matrix/modules/erp/service/stock_service"
	"sort"
)

type ErpOrderWaybillPrint struct {
	*revel.Controller
}

func (c ErpOrderWaybillPrint) Index() revel.Result {
	return c.RenderEmpty()
}

func (c ErpOrderWaybillPrint) SelectPrintOrder() revel.Result {
	return c.RenderTemplate("")
}

func (c ErpOrderWaybillPrint) Preview() revel.Result {
	var db_session = c.DbSession

	var order_id_list = make([]int64, 0)
	c.Params.Bind(&order_id_list, "order_id_list")

	var order_by_code int64
	c.Params.Bind(&order_by_code, "order_by_code")

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

		var order_print_item_detail = get_order_sku_qty_loc_list_for_print(order_id, db_session)
		print_order.SkuQtyLocList = order_print_item_detail.sku_qty_loc_list
		print_order.PcsCount = order_print_item_detail.pcs_count
		print_order.first_sku = order_print_item_detail.first_sku
		print_order.first_storage_loc = order_print_item_detail.first_storage_loc

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

	if order_by_code == ORDER_BY_SKU {
		var print_order_list_order_by_sku = print_order_order_by_sku_slice_t(print_order_list)
		sort.Sort(print_order_list_order_by_sku)
		print_order_list = []Print_order_t(print_order_list_order_by_sku)
	}

	if order_by_code == ORDER_BY_SL {
		var print_order_list_order_by_sl = print_order_order_by_sl_slice_t(print_order_list)
		sort.Sort(print_order_list_order_by_sl)
		print_order_list = []Print_order_t(print_order_list_order_by_sl)
	}

	c.RenderArgs["print_order_list"] = print_order_list
	return c.RenderTemplate("erp/ddlc/order_waybill_print/order_waybill_print_preview.html")
}

const (
	ORDER_BY_SKU int64 = 100 // 按SKU排序
	ORDER_BY_SL  int64 = 200 // 按库位排序
)

func (c ErpOrderWaybillPrint) SelectOrderByDialog() revel.Result {
	return c.RenderEmpty()
}

func (c ErpOrderWaybillPrint) Print() revel.Result {
	var db_session = c.DbSession

	var order_id_list = make([]int64, 0)
	c.Params.Bind(&order_id_list, "order_id_list")

	var order_by_code int64
	c.Params.Bind(&order_by_code, "order_by_code")

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

		var order_print_item_detail = get_order_sku_qty_loc_list_for_print(order_id, db_session)
		print_order.SkuQtyLocList = order_print_item_detail.sku_qty_loc_list
		print_order.PcsCount = order_print_item_detail.pcs_count
		print_order.first_sku = order_print_item_detail.first_sku
		print_order.first_storage_loc = order_print_item_detail.first_storage_loc

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

	if order_by_code == ORDER_BY_SKU {
		var print_order_list_order_by_sku = print_order_order_by_sku_slice_t(print_order_list)
		sort.Sort(print_order_list_order_by_sku)
		print_order_list = []Print_order_t(print_order_list_order_by_sku)
	}

	if order_by_code == ORDER_BY_SL {
		var print_order_list_order_by_sl = print_order_order_by_sl_slice_t(print_order_list)
		sort.Sort(print_order_list_order_by_sl)
		print_order_list = []Print_order_t(print_order_list_order_by_sl)
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
	first_sku            string
	first_storage_loc    string
}

type print_order_order_by_sku_slice_t []Print_order_t

func (s print_order_order_by_sku_slice_t) Len() int {
	return len(s)
}

func (s print_order_order_by_sku_slice_t) Less(i int, j int) bool {
	return s[i].first_sku < s[j].first_sku
}

func (s print_order_order_by_sku_slice_t) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type print_order_order_by_sl_slice_t []Print_order_t

func (s print_order_order_by_sl_slice_t) Len() int {
	return len(s)
}

func (s print_order_order_by_sl_slice_t) Less(i int, j int) bool {
	return s[i].first_storage_loc < s[j].first_storage_loc
}

func (s print_order_order_by_sl_slice_t) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

type order_print_item_detail_t struct {
	sku_qty_loc_list  []string
	pcs_count         int64
	first_sku         string
	first_storage_loc string
}

func get_order_sku_qty_loc_list_for_print(order_id int64, db_session *xorm.Session) order_print_item_detail_t {
	/*
		数量*SKU_仓位号
	 */

	var order_sku_qty_loc_list = order_service.Get_order_sku_qty_loc_list(order_id, db_session)

	var first_sku string
	var first_storage_loc string

	var order_item_list = make([]string, 0)
	var pcs_count int64 = 0
	for _, order_sku_qty_loc := range order_sku_qty_loc_list {
		order_item_list = append(order_item_list, fmt.Sprintf("%d*%s_%s", order_sku_qty_loc.Quantity, order_sku_qty_loc.Sku, order_sku_qty_loc.StorageLocName))
		pcs_count = pcs_count + order_sku_qty_loc.Quantity

		if first_sku == "" {
			first_sku = order_sku_qty_loc.Sku
		}

		if first_storage_loc == "" {
			first_storage_loc = order_sku_qty_loc.StorageLocName
		}
	}

	var order_print_item_detail order_print_item_detail_t
	order_print_item_detail.sku_qty_loc_list = order_item_list
	order_print_item_detail.pcs_count = pcs_count
	order_print_item_detail.first_sku = first_sku
	order_print_item_detail.first_storage_loc = first_storage_loc

	return order_print_item_detail
}
