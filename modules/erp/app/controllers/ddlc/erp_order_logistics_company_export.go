package controllers

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/stock_service"
)

type ErpOrderLogisticsCompanyExport struct {
	*revel.Controller
}

var (
	excel_export_stock_mel_code = 100
	excel_export_stock_bne_code = 200

	excel_export_company_eparcel_code int = 100
	excel_export_company_fastway_code int = 200
	excel_export_company_toll_code    int = 300
)

func (c ErpOrderLogisticsCompanyExport) ExportOptionData() revel.Result {
	var db_session = c.DbSession
	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)

	var excel_export_company_list = map[string]interface{}{
		fmt.Sprint(excel_export_company_eparcel_code): "AU POST eParcel",
		fmt.Sprint(excel_export_company_fastway_code): "FASTWAY,",
		fmt.Sprint(excel_export_company_toll_code):    "TOLL IPEC",
	}

	return c.RenderJson(map[string]interface{}{
		"success":      true,
		"stock_list":   stock_list,
		"company_list": excel_export_company_list,
	})
}

func (c ErpOrderLogisticsCompanyExport) ExportExcel() revel.Result {
	var db_session = c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var company_code int
	c.Params.Bind(&company_code, "company_code")

	// 判断当前登录用户是否有选择的仓库的权限
	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)
	var has_stock = false
	for _, stock := range stock_list {
		if stock_id == stock.Id {
			has_stock = true
			break
		}
	}

	if has_stock == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "你没有导出当前选择仓库物流的权限！",
		})
	}

	// eparcel_bne.xlsx
	// eparcel_mel.xlsx
	// fastway_bne.xlsx
	// fastway_mel.xlsx
	// toll_bne.xlsx
	// toll_mel.xlsx

	// the stock selected is mel or bne?
	var stock models.StockInfo
	_, err := db_session.Id(stock_id).Get(&stock)
	core.HandleError(err)

	var stock_code = 0 // 用于区分 mel 还是 bne
	if strings.ToLower(strings.TrimSpace(stock.ShortName.String)) == strings.ToLower("MEL") {
		stock_code = excel_export_stock_mel_code
	} else if strings.ToLower(strings.TrimSpace(stock.ShortName.String)) == strings.ToLower("BNE") {
		stock_code = excel_export_stock_bne_code
	} else {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "选择的仓库不能正确识别为墨尔本或布里斯班仓库！",
		})
	}

	var excel_template_logistics_export_dir = filepath.Join(revel.BasePath, "/static/excel-template/logistics_export")

	var xlsx_file *xlsx.File
	var export_excel_file_name string
	if company_code == excel_export_company_eparcel_code && (stock_code == excel_export_stock_mel_code || stock_code == excel_export_stock_bne_code ) {
		type ModeId struct {
			ModeId int64 `xorm:"'mode_id'" json:"mode_id"`
		}

		var mode_id_list = make([]ModeId, 0)
		sql_get_mode_id_list := `SELECT lm.id AS mode_id FROM erp_logistics_mode lm WHERE lm.stock_id = %d AND lm.name LIKE '%%eParcel%%';`
		sql_get_mode_id_list = fmt.Sprintf(sql_get_mode_id_list, stock_id)
		err := db_session.Sql(sql_get_mode_id_list).Find(&mode_id_list)
		core.HandleError(err)

		var mode_id_str_list = make([]string, 0)
		for _, mode_id := range mode_id_list {
			mode_id_str_list = append(mode_id_str_list, fmt.Sprint(mode_id))
		}

		/*
SELECT ol.*, o.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
WHERE ol.mode_id IN (0) AND ol.stock_id = 0
		 */

		type order_logistics_view_item_t struct {
			OrderLogistics models.OrderLogisticsInfo    `xorm:"extends" json:"ol"`
			Order          models.OrderInfo             `xorm:"extends" json:"o"`
		}

		sql_get_logistics_list := `SELECT ol.*, o.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
WHERE ol.mode_id IN (%s) AND ol.stock_id = %d`
		sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, strings.Join(mode_id_str_list, ","), stock_id)
		var order_logistics_list = make([]order_logistics_view_item_t, 0)
		err = db_session.Sql(sql_get_logistics_list).Find(&order_logistics_list)
		core.HandleError(err)

		if stock_code == excel_export_stock_mel_code {
			xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "eparcel_mel.xlsx"))
			core.HandleError(err)
			export_excel_file_name = "eParcel 墨尔本"

			sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
			for index, order_logistics := range order_logistics_list {
				var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
				sheet.Cell(index+1, excel_col_label_to_index("C")).Value = "7D85"
				sheet.Cell(index+1, excel_col_label_to_index("D")).Value = order_logistics.Order.FullTranId
				sheet.Cell(index+1, excel_col_label_to_index("E")).Value = order_shipping_address.Name

				var full_address = order_shipping_address.Street1 + order_shipping_address.Street2
				address1, address2 := split_address(full_address, 35)
				sheet.Cell(index+1, excel_col_label_to_index("G")).Value = address1
				sheet.Cell(index+1, excel_col_label_to_index("H")).Value = address2
				sheet.Cell(index+1, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
				sheet.Cell(index+1, excel_col_label_to_index("L")).Value = order_shipping_address.Province
				sheet.Cell(index+1, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
				sheet.Cell(index+1, excel_col_label_to_index("N")).Value = order_shipping_address.CountryName
				sheet.Cell(index+1, excel_col_label_to_index("O")).Value = order_shipping_address.Phone
				sheet.Cell(index+1, excel_col_label_to_index("P")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("S")).Value = "A"
				sheet.Cell(index+1, excel_col_label_to_index("T")).Value = "N"
				sheet.Cell(index+1, excel_col_label_to_index("V")).Value = "N"
				sheet.Cell(index+1, excel_col_label_to_index("X")).Value = order_logistics.Order.FullTranId
				sheet.Cell(index+1, excel_col_label_to_index("Y")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("Z")).Value = get_order_sku_for_export(order_logistics.Order.Id, db_session)
				sheet.Cell(index+1, excel_col_label_to_index("AA")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("AC")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("AD")).Value = "Iepoch MEL (3PL)"
				sheet.Cell(index+1, excel_col_label_to_index("AE")).Value = "2 Venture Court"
				sheet.Cell(index+1, excel_col_label_to_index("AI")).Value = "DANDENONG SOUTH"
				sheet.Cell(index+1, excel_col_label_to_index("AJ")).Value = "VIC"
				sheet.Cell(index+1, excel_col_label_to_index("AK")).Value = "3175"
				sheet.Cell(index+1, excel_col_label_to_index("AL")).Value = "AU"
				sheet.Cell(index+1, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
				sheet.Cell(index+1, excel_col_label_to_index("AY")).Value = "TRACKADV"
				sheet.Cell(index+1, excel_col_label_to_index("BI")).Value = fmt.Sprint(order_logistics.OrderLogistics.Weight)
				sheet.Cell(index+1, excel_col_label_to_index("BO")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("BP")).Value = "FALSE"
			}

		} else if stock_code == excel_export_stock_bne_code {
			xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "eparcel_bne.xlsx"))
			core.HandleError(err)
			export_excel_file_name = "eParcel 布里斯班"

			sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
			for index, order_logistics := range order_logistics_list {
				var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
				sheet.Cell(index+1, excel_col_label_to_index("C")).Value = "7D85"
				sheet.Cell(index+1, excel_col_label_to_index("D")).Value = order_logistics.Order.FullTranId
				sheet.Cell(index+1, excel_col_label_to_index("E")).Value = order_shipping_address.Name

				var full_address = order_shipping_address.Street1 + order_shipping_address.Street2
				address1, address2 := split_address(full_address, 35)
				sheet.Cell(index+1, excel_col_label_to_index("G")).Value = address1
				sheet.Cell(index+1, excel_col_label_to_index("H")).Value = address2
				sheet.Cell(index+1, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
				sheet.Cell(index+1, excel_col_label_to_index("L")).Value = order_shipping_address.Province
				sheet.Cell(index+1, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
				sheet.Cell(index+1, excel_col_label_to_index("N")).Value = order_shipping_address.CountryName
				sheet.Cell(index+1, excel_col_label_to_index("O")).Value = order_shipping_address.Phone
				sheet.Cell(index+1, excel_col_label_to_index("P")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("S")).Value = "A"
				sheet.Cell(index+1, excel_col_label_to_index("T")).Value = "N"
				sheet.Cell(index+1, excel_col_label_to_index("V")).Value = "N"
				sheet.Cell(index+1, excel_col_label_to_index("X")).Value = order_logistics.Order.FullTranId
				sheet.Cell(index+1, excel_col_label_to_index("Y")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("Z")).Value = get_order_sku_for_export(order_logistics.Order.Id, db_session)
				sheet.Cell(index+1, excel_col_label_to_index("AA")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("AC")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("AD")).Value = "Iepoch BNE (3PL)"
				sheet.Cell(index+1, excel_col_label_to_index("AE")).Value = "2/30 Gow Street"
				sheet.Cell(index+1, excel_col_label_to_index("AI")).Value = "MOOROOKA"
				sheet.Cell(index+1, excel_col_label_to_index("AJ")).Value = "QLD"
				sheet.Cell(index+1, excel_col_label_to_index("AK")).Value = "4105"
				sheet.Cell(index+1, excel_col_label_to_index("AL")).Value = "AU"
				sheet.Cell(index+1, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
				sheet.Cell(index+1, excel_col_label_to_index("AY")).Value = "TRACKADV"
				sheet.Cell(index+1, excel_col_label_to_index("BI")).Value = fmt.Sprint(order_logistics.OrderLogistics.Weight)
				sheet.Cell(index+1, excel_col_label_to_index("BO")).Value = "FALSE"
				sheet.Cell(index+1, excel_col_label_to_index("BP")).Value = "FALSE"
			}
		}
	} else if company_code == excel_export_company_fastway_code && (stock_code == excel_export_stock_mel_code || stock_code == excel_export_stock_bne_code) {
		xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "fastway.xlsx"))
		core.HandleError(err)

		if stock_code == excel_export_stock_mel_code {
			export_excel_file_name = "FASTWAY 墨尔本"
		} else if stock_code == excel_export_stock_bne_code {
			export_excel_file_name = "FASTWAY 布里斯班"
		}

		type ModeId struct {
			ModeId int64 `xorm:"'mode_id'" json:"mode_id"`
		}

		var mode_id_list = make([]ModeId, 0)
		sql_get_mode_id_list := `SELECT lm.id AS mode_id FROM erp_logistics_mode lm WHERE lm.stock_id = %d AND lm.name LIKE '%%FASTWAY%%';`
		sql_get_mode_id_list = fmt.Sprintf(sql_get_mode_id_list, stock_id)
		err := db_session.Sql(sql_get_mode_id_list).Find(&mode_id_list)
		core.HandleError(err)

		var mode_id_str_list = make([]string, 0)
		for _, mode_id := range mode_id_list {
			mode_id_str_list = append(mode_id_str_list, fmt.Sprint(mode_id))
		}

		/*
SELECT ol.*, o.*, lm.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
WHERE ol.mode_id IN (0) AND ol.stock_id = 0
		 */

		type order_logistics_view_item_t struct {
			OrderLogistics models.OrderLogisticsInfo    `xorm:"extends" json:"ol"`
			Order          models.OrderInfo             `xorm:"extends" json:"o"`
			LogisticsMode  models.LogisticsModeInfo     `xorm:"extends" json:"lm"`
		}

		sql_get_logistics_list := `SELECT ol.*, o.*, lm.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
WHERE ol.mode_id IN (%s) AND ol.stock_id = %d`
		sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, strings.Join(mode_id_str_list, ","), stock_id)
		var order_logistics_list = make([]order_logistics_view_item_t, 0)
		err = db_session.Sql(sql_get_logistics_list).Find(&order_logistics_list)
		core.HandleError(err)

		sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
		for index, order_logistics := range order_logistics_list {
			var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
			sheet.Cell(index+1, excel_col_label_to_index("B")).Value = order_logistics.Order.FullTranId
			sheet.Cell(index+1, excel_col_label_to_index("C")).Value = fmt.Sprint(order_logistics.OrderLogistics.Weight)
			sheet.Cell(index+1, excel_col_label_to_index("D")).Value = "1"

			sheet.Cell(index+1, excel_col_label_to_index("H")).Value = order_shipping_address.Name
			var full_address = order_shipping_address.Street1 + order_shipping_address.Street2
			address1, address2 := split_address(full_address, 25)
			sheet.Cell(index+1, excel_col_label_to_index("I")).Value = address1
			sheet.Cell(index+1, excel_col_label_to_index("J")).Value = address2
			sheet.Cell(index+1, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
			sheet.Cell(index+1, excel_col_label_to_index("L")).Value = order_shipping_address.Province
			sheet.Cell(index+1, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
			sheet.Cell(index+1, excel_col_label_to_index("N")).Value = order_shipping_address.Name
			sheet.Cell(index+1, excel_col_label_to_index("O")).Value = order_logistics.Order.BuyerEmail
			sheet.Cell(index+1, excel_col_label_to_index("P")).Value = order_shipping_address.Phone
			sheet.Cell(index+1, excel_col_label_to_index("R")).Value = "ATL at safe place"
			sheet.Cell(index+1, excel_col_label_to_index("S")).Value = order_logistics.Order.FullTranId
			sheet.Cell(index+1, excel_col_label_to_index("T")).Value = get_order_sku_for_export(order_logistics.Order.Id, db_session)
			sheet.Cell(index+1, excel_col_label_to_index("U")).Value = ""
		}

	} else if company_code == excel_export_company_toll_code && stock_code == excel_export_stock_mel_code {
		xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "toll_mel.xlsx"))
		core.HandleError(err)
		export_excel_file_name = "TOLL 墨尔本"

	} else if company_code == excel_export_company_toll_code && stock_code == excel_export_stock_bne_code {
		xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "toll_bne.xlsx"))
		core.HandleError(err)
		export_excel_file_name = "TOLL 布里斯班"

	} else {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "未识别的仓库和物流公司！",
		})
	}

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口
	file_name := export_excel_file_name + ".xlsx"

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

// 把 A 转成 0， C 转成 2， 把 AA 转成 26， 把 BA 转成 52
func excel_col_label_to_index(row_label string) int {
	row_label = strings.ToUpper(row_label)
	var row_label_rune_list = []rune(row_label)
	var row_label_rune_list_len = len(row_label_rune_list)

	var row_index int = 0
	for i := 0; i < row_label_rune_list_len; i++ {
		row_index = row_index*26 + int(row_label_rune_list[i]-rune('A')) + 1
	}

	return row_index - 1
}

func split_address(address string, letter_limit int) (address1 string, address2 string) {
	if letter_limit <= 0 {
		return address, ""
	}

	var address_rune_list = []rune(address)
	var address_rune_list_len = len(address_rune_list)

	if address_rune_list_len <= letter_limit {
		return address, ""
	}

	var space_rune = ([]rune(" "))[0]
	var index int
	for index = letter_limit; index > 0; index-- {
		if address_rune_list[index] == space_rune {
			break
		}
	}

	return string(address_rune_list[0:index]), string(address_rune_list[index:])
}

func get_order_sku_for_export(order_id int64, db_session *xorm.Session) string {
	/*
	single product, single quantity

	1. 对于 SS， 显示 SKU 数量 [仓位号]
	2. 对于 SM， 显示 SKU DOG数量 [仓位号]
	3. 对于 MM, 显示 MM
	 */
	type OrderSkuQuantity struct {
		Sku            string `xorm:"'sku'" json:"sku"`
		Quantity       int64  `xorm:"'quantity'" json:"quantity"`
		StorageLocName string `xorm:"'storage_loc_name'" json:"storage_loc_name"`
	}

	var order_sku_quantity_list = make([]OrderSkuQuantity, 0)
	sql := `SELECT
    oi.item_sku  AS sku,
    oi.quantity  AS quantity,
    (SELECT sl.name
     FROM erp_inventory nvt
         INNER JOIN erp_storage_loc sl ON sl.id = nvt.storage_loc_id
     WHERE nvt.product_id = 0
     ORDER BY sl.update_time DESC
     LIMIT 0, 1) AS storage_loc_name
FROM erp_order_item oi
WHERE oi.order_id = %d;`
	sql = fmt.Sprintf(sql, order_id)
	err := db_session.Sql(sql).Find(&order_sku_quantity_list)
	core.HandleError(err)

	var order_sku_quantity_list_len = len(order_sku_quantity_list)
	if order_sku_quantity_list_len < 1 {
		return ""
	} else if order_sku_quantity_list_len == 1 {
		var order_sku_quantity = order_sku_quantity_list[0]
		if order_sku_quantity.Quantity == 1 {
			return fmt.Sprintf("%s %d [%s]", order_sku_quantity.Sku, order_sku_quantity.Quantity, order_sku_quantity.StorageLocName)
		} else if order_sku_quantity.Quantity > 1 {
			return fmt.Sprintf("%s DOG%d [%s]", order_sku_quantity.Sku, order_sku_quantity.Quantity, order_sku_quantity.StorageLocName)
		} else {
			return ""
		}
	} else { // if order_sku_quantity_list_len > 1 {
		return "MM"
	}
}
