package controllers

import (
	"bytes"
	"fmt"
	"math"
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
	"sort"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/modules/erp/models/ddlc_models"
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
	excel_export_company_elife_code   int = 400
)

func (c ErpOrderLogisticsCompanyExport) ExportOptionData() revel.Result {
	var db_session = c.DbSession
	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)

	var excel_export_company_list = [](map[string]interface{}){
		map[string]interface{}{
			"code": excel_export_company_eparcel_code,
			"name": "GIANT AU POST eParcel",
		},
		map[string]interface{}{
			"code": excel_export_company_elife_code,
			"name": "ELIFE AU POST eParcel",
		},
		map[string]interface{}{
			"code": excel_export_company_fastway_code,
			"name": "FASTWAY",
		},
		map[string]interface{}{
			"code": excel_export_company_toll_code,
			"name": "TOLL IPEC",
		},
	}

	return c.RenderJson(map[string]interface{}{
		"success":      true,
		"stock_list":   stock_list,
		"company_list": excel_export_company_list,
	})
}

type order_logistics_eparcel_view_item_t struct {
	OrderLogistics ddlc_models.OrderLogisticsInfo `xorm:"extends" json:"ol"`
	LogisticsMode  ddpz_models.LogisticsModeInfo  `xorm:"extends" json:"lm"`
	Order          ddlc_models.OrderInfo          `xorm:"extends" json:"o"`
	ShopAccount    xtpz_models.ShopAccountInfo    `xorm:"extends" json:"sa"`
	sku_qty_loc    string                         `xorm:"-"`
	loc_name       string                         `xorm:"-"`
}

type order_logistics_eparcel_view_item_slice []order_logistics_eparcel_view_item_t

func (s order_logistics_eparcel_view_item_slice) Len() int           { return len(s) }
func (s order_logistics_eparcel_view_item_slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s order_logistics_eparcel_view_item_slice) Less(i, j int) bool { return s[i].loc_name < s[j].loc_name }

type order_logistics_fastway_view_item_t struct {
	OrderLogistics ddlc_models.OrderLogisticsInfo `xorm:"extends" json:"ol"`
	Order          ddlc_models.OrderInfo          `xorm:"extends" json:"o"`
	LogisticsMode  ddpz_models.LogisticsModeInfo  `xorm:"extends" json:"lm"`
	ShopAccount    xtpz_models.ShopAccountInfo    `xorm:"extends" json:"sa"`
	sku_qty_loc    string                         `xorm:"-"`
	loc_name       string                         `xorm:"-"`
}

type order_logistics_fastway_view_item_slice []order_logistics_fastway_view_item_t

func (s order_logistics_fastway_view_item_slice) Len() int           { return len(s) }
func (s order_logistics_fastway_view_item_slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s order_logistics_fastway_view_item_slice) Less(i, j int) bool { return s[i].loc_name < s[j].loc_name }

type order_logistics_toll_view_item_t struct {
	OrderLogistics         ddlc_models.OrderLogisticsInfo       `xorm:"extends" json:"ol"`
	Order                  ddlc_models.OrderInfo                `xorm:"extends" json:"o"`
	LogisticsMode          ddpz_models.LogisticsModeInfo        `xorm:"extends" json:"lm"`
	ShopAccount            xtpz_models.ShopAccountInfo          `xorm:"extends" json:"sa"`
	sku_qty_loc            string                               `xorm:"-"`
	loc_name               string                               `xorm:"-"`
	qty                    int64                                `xorm:"-"`
	service_id             string                               `xorm:"-"`
	order_shipping_address ddlc_models.OrderShippingAddressInfo `xorm:"-"`
}

type order_logistics_toll_view_item_slice []order_logistics_toll_view_item_t

func (s order_logistics_toll_view_item_slice) Len() int           { return len(s) }
func (s order_logistics_toll_view_item_slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s order_logistics_toll_view_item_slice) Less(i, j int) bool { return s[i].loc_name < s[j].loc_name }

func (c ErpOrderLogisticsCompanyExport) ExportExcel() revel.Result {
	var db_session = c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var company_code int
	c.Params.Bind(&company_code, "company_code")

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var export_type int64
	c.Params.Bind(&export_type, "export_type")
	// 如果 export_type == 0， 是订单物流导出； 如果 export_type == 1 是待配货导出

	var order_id_list_str string
	if export_type == 1 {
		c.Params.Bind(&order_id_list_str, "order_id_list")
	}

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
	var stock sppz_models.StockInfo
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

	var set_cell_style = func(cell *xlsx.Cell) {
		style := cell.GetStyle()

		font := * xlsx.NewFont(11, "宋体")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyFont = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		cell.SetStyle(style)
	}

	var excel_template_logistics_export_dir = filepath.Join(revel.BasePath, "/static/excel-template/logistics_export")

	var xlsx_file *xlsx.File
	var export_excel_file_name string
	if company_code == excel_export_company_eparcel_code && (stock_code == excel_export_stock_mel_code || stock_code == excel_export_stock_bne_code ) {
		type ModeId struct {
			ModeId int64 `xorm:"'mode_id'" json:"mode_id"`
		}

		var mode_id_list = make([]ModeId, 0)
		sql_get_mode_id_list := `SELECT lm.id AS mode_id FROM erp_logistics_mode lm WHERE lm.stock_id = %d AND lm.name LIKE '%%Giant eParcel%%';`
		sql_get_mode_id_list = fmt.Sprintf(sql_get_mode_id_list, stock_id)
		err := db_session.Sql(sql_get_mode_id_list).Find(&mode_id_list)
		core.HandleError(err)

		var mode_id_str_list = make([]string, 0)
		for _, mode_id := range mode_id_list {
			mode_id_str_list = append(mode_id_str_list, fmt.Sprint(mode_id.ModeId))
		}

		/*
SELECT ol.*, lm.*, o.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_logistics_mode lm on lm.id = ol.mode_id
    INNER JOIN erp_order o ON o.id = ol.order_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE ol.mode_id IN (0) AND ol.stock_id = 0
		 */

		sql_get_logistics_list := `SELECT ol.*, lm.*, o.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_logistics_mode lm on lm.id = ol.mode_id
    INNER JOIN erp_order o ON o.id = ol.order_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE o.status_id = %d AND ol.mode_id IN (%s) AND ol.stock_id = %d`

		if export_type == 0 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_PHZ, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if export_type == 1 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_DPH, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if batch_code != 0 {
			sql_get_logistics_list += fmt.Sprintf(" AND o.batch_code = %d", batch_code)
		}

		if export_type == 1 { //待配货的导出
			sql_get_logistics_list = sql_get_logistics_list + fmt.Sprintf(" AND o.id IN(%s)", order_id_list_str)
		}

		var order_logistics_list = make([]order_logistics_eparcel_view_item_t, 0)
		err = db_session.Sql(sql_get_logistics_list).Find(&order_logistics_list)
		core.HandleError(err)

		for index, order_logistics := range order_logistics_list {
			order_logistics_list[index].sku_qty_loc, order_logistics_list[index].loc_name = get_order_sku_for_export(order_logistics.Order.Id, db_session)
		}

		var order_logistics_list_slice = order_logistics_eparcel_view_item_slice(order_logistics_list)
		sort.Sort(order_logistics_list_slice)

		order_logistics_list = []order_logistics_eparcel_view_item_t(order_logistics_list_slice)

		if stock_code == excel_export_stock_mel_code {
			xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "eparcel_mel.xlsx"))
			core.HandleError(err)
			export_excel_file_name = "GIANT eParcel MEL"

			sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
			for index, order_logistics := range order_logistics_list {
				var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
				// sheet.Cell(index+2, excel_col_label_to_index("C")).Value = "7D85"
				sheet.Cell(index+2, excel_col_label_to_index("C")).Value = order_logistics.LogisticsMode.SenderId
				sheet.Cell(index+2, excel_col_label_to_index("D")).Value = order_logistics.Order.FullTranId

				name1, name2 := split_long_string(order_shipping_address.Name, 35)
				sheet.Cell(index+2, excel_col_label_to_index("E")).Value = name1
				sheet.Cell(index+2, excel_col_label_to_index("F")).Value = name2

				var full_address = order_shipping_address.Street1 + "," + order_shipping_address.Street2
				address1, address2 := split_long_string(full_address, 35)
				sheet.Cell(index+2, excel_col_label_to_index("G")).Value = address1
				sheet.Cell(index+2, excel_col_label_to_index("H")).Value = address2
				sheet.Cell(index+2, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
				sheet.Cell(index+2, excel_col_label_to_index("L")).Value = get_state_short_name(order_shipping_address.Province)
				sheet.Cell(index+2, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
				//sheet.Cell(index+2, excel_col_label_to_index("N")).Value = order_shipping_address.CountryName
				sheet.Cell(index+2, excel_col_label_to_index("N")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("O")).Value = order_shipping_address.Phone
				sheet.Cell(index+2, excel_col_label_to_index("P")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("S")).Value = "A"
				sheet.Cell(index+2, excel_col_label_to_index("T")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("V")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("X")).Value = fmt.Sprintf("%s FROM:%s", order_logistics.Order.FullTranId, order_logistics.ShopAccount.Name)
				sheet.Cell(index+2, excel_col_label_to_index("Y")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("Z")).Value = order_logistics.sku_qty_loc
				sheet.Cell(index+2, excel_col_label_to_index("AA")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AC")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AD")).Value = "iepoch MEL (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AE")).Value = "54-58 Arkwright Drive"
				sheet.Cell(index+2, excel_col_label_to_index("AI")).Value = "DANDENONG SOUTH"
				sheet.Cell(index+2, excel_col_label_to_index("AJ")).Value = "VIC"
				sheet.Cell(index+2, excel_col_label_to_index("AK")).Value = "3175"
				sheet.Cell(index+2, excel_col_label_to_index("AL")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("AM")).Value = "iepoch MEL (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AN")).Value = "iepoch MEL (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
				sheet.Cell(index+2, excel_col_label_to_index("AY")).Value = "TRACKADV"
				sheet.Cell(index+2, excel_col_label_to_index("BI")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Weight)
				sheet.Cell(index+2, excel_col_label_to_index("BO")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("BP")).Value = "FALSE"
			}

		} else if stock_code == excel_export_stock_bne_code {
			xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "eparcel_bne.xlsx"))
			core.HandleError(err)
			export_excel_file_name = "GIANT eParcel BNE"

			sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
			for index, order_logistics := range order_logistics_list {
				var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
				// sheet.Cell(index+2, excel_col_label_to_index("C")).Value = "7D85"
				sheet.Cell(index+2, excel_col_label_to_index("C")).Value = order_logistics.LogisticsMode.SenderId
				sheet.Cell(index+2, excel_col_label_to_index("D")).Value = order_logistics.Order.FullTranId

				name1, name2 := split_long_string(order_shipping_address.Name, 35)
				sheet.Cell(index+2, excel_col_label_to_index("E")).Value = name1
				sheet.Cell(index+2, excel_col_label_to_index("F")).Value = name2

				var full_address = order_shipping_address.Street1 + "," + order_shipping_address.Street2
				address1, address2 := split_long_string(full_address, 35)
				sheet.Cell(index+2, excel_col_label_to_index("G")).Value = address1
				sheet.Cell(index+2, excel_col_label_to_index("H")).Value = address2
				sheet.Cell(index+2, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
				sheet.Cell(index+2, excel_col_label_to_index("L")).Value = get_state_short_name(order_shipping_address.Province)
				sheet.Cell(index+2, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
				//sheet.Cell(index+2, excel_col_label_to_index("N")).Value = order_shipping_address.CountryName
				sheet.Cell(index+2, excel_col_label_to_index("N")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("O")).Value = order_shipping_address.Phone
				sheet.Cell(index+2, excel_col_label_to_index("P")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("S")).Value = "A"
				sheet.Cell(index+2, excel_col_label_to_index("T")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("V")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("X")).Value = fmt.Sprintf("%s FROM:%s", order_logistics.Order.FullTranId, order_logistics.ShopAccount.Name)
				sheet.Cell(index+2, excel_col_label_to_index("Y")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("Z")).Value = order_logistics.sku_qty_loc
				sheet.Cell(index+2, excel_col_label_to_index("AA")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AC")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AD")).Value = "iepoch BNE (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AE")).Value = "Basement 2/30 Gow Street"
				sheet.Cell(index+2, excel_col_label_to_index("AI")).Value = "MOOROOKA"
				sheet.Cell(index+2, excel_col_label_to_index("AJ")).Value = "QLD"
				sheet.Cell(index+2, excel_col_label_to_index("AK")).Value = "4105"
				sheet.Cell(index+2, excel_col_label_to_index("AL")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("AM")).Value = "iepoch BNE (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AN")).Value = "iepoch BNE (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
				sheet.Cell(index+2, excel_col_label_to_index("AY")).Value = "TRACKADV"
				sheet.Cell(index+2, excel_col_label_to_index("BI")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Weight)
				sheet.Cell(index+2, excel_col_label_to_index("BO")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("BP")).Value = "FALSE"
			}
		}
	} else if company_code == excel_export_company_elife_code && (stock_code == excel_export_stock_mel_code || stock_code == excel_export_stock_bne_code ) {
		type ModeId struct {
			ModeId int64 `xorm:"'mode_id'" json:"mode_id"`
		}

		var mode_id_list = make([]ModeId, 0)
		sql_get_mode_id_list := `SELECT lm.id AS mode_id FROM erp_logistics_mode lm WHERE lm.stock_id = %d AND lm.name LIKE '%%Elife eParcel%%';`
		sql_get_mode_id_list = fmt.Sprintf(sql_get_mode_id_list, stock_id)
		err := db_session.Sql(sql_get_mode_id_list).Find(&mode_id_list)
		core.HandleError(err)

		var mode_id_str_list = make([]string, 0)
		for _, mode_id := range mode_id_list {
			mode_id_str_list = append(mode_id_str_list, fmt.Sprint(mode_id.ModeId))
		}

		/*
SELECT ol.*, lm.*, o.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    INNER JOIN erp_order o ON o.id = ol.order_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE ol.mode_id IN (0) AND ol.stock_id = 0
		 */

		sql_get_logistics_list := `SELECT ol.*, lm.*, o.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    INNER JOIN erp_order o ON o.id = ol.order_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE o.status_id = %d AND ol.mode_id IN (%s) AND ol.stock_id = %d`

		if export_type == 0 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_PHZ, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if export_type == 1 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_DPH, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if batch_code != 0 {
			sql_get_logistics_list += fmt.Sprintf(" AND o.batch_code = %d", batch_code)
		}

		if export_type == 1 { //待配货的导出
			sql_get_logistics_list = sql_get_logistics_list + fmt.Sprintf(" AND o.id IN(%s)", order_id_list_str)
		}

		var order_logistics_list = make([]order_logistics_eparcel_view_item_t, 0)
		err = db_session.Sql(sql_get_logistics_list).Find(&order_logistics_list)
		core.HandleError(err)

		for index, order_logistics := range order_logistics_list {
			order_logistics_list[index].sku_qty_loc, order_logistics_list[index].loc_name = get_order_sku_for_export(order_logistics.Order.Id, db_session)
		}

		var order_logistics_list_slice = order_logistics_eparcel_view_item_slice(order_logistics_list)
		sort.Sort(order_logistics_list_slice)

		order_logistics_list = []order_logistics_eparcel_view_item_t(order_logistics_list_slice)

		if stock_code == excel_export_stock_mel_code {
			xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "elife_ep_mel.xlsx"))
			core.HandleError(err)
			export_excel_file_name = "ELIFE eParcel MEL"

			sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
			for index, order_logistics := range order_logistics_list {
				var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
				// sheet.Cell(index+2, excel_col_label_to_index("C")).Value = "7D83"
				sheet.Cell(index+2, excel_col_label_to_index("C")).Value = order_logistics.LogisticsMode.SenderId
				sheet.Cell(index+2, excel_col_label_to_index("D")).Value = order_logistics.Order.FullTranId

				name1, name2 := split_long_string(order_shipping_address.Name, 35)
				sheet.Cell(index+2, excel_col_label_to_index("E")).Value = name1
				sheet.Cell(index+2, excel_col_label_to_index("F")).Value = name2

				var full_address = order_shipping_address.Street1 + "," + order_shipping_address.Street2
				address1, address2 := split_long_string(full_address, 35)
				sheet.Cell(index+2, excel_col_label_to_index("G")).Value = address1
				sheet.Cell(index+2, excel_col_label_to_index("H")).Value = address2
				sheet.Cell(index+2, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
				sheet.Cell(index+2, excel_col_label_to_index("L")).Value = get_state_short_name(order_shipping_address.Province)
				sheet.Cell(index+2, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
				//sheet.Cell(index+2, excel_col_label_to_index("N")).Value = order_shipping_address.CountryName
				sheet.Cell(index+2, excel_col_label_to_index("N")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("O")).Value = order_shipping_address.Phone
				sheet.Cell(index+2, excel_col_label_to_index("P")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("S")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("V")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("X")).Value = fmt.Sprintf("%s FROM:%s", order_logistics.Order.FullTranId, order_logistics.ShopAccount.Name)
				sheet.Cell(index+2, excel_col_label_to_index("Y")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("Z")).Value = order_logistics.sku_qty_loc
				sheet.Cell(index+2, excel_col_label_to_index("AA")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AC")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AD")).Value = "iepoch MEL (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AE")).Value = "54-58 Arkwright Drive"
				sheet.Cell(index+2, excel_col_label_to_index("AI")).Value = "DANDENONG SOUTH"
				sheet.Cell(index+2, excel_col_label_to_index("AJ")).Value = "VIC"
				sheet.Cell(index+2, excel_col_label_to_index("AK")).Value = "3175"
				sheet.Cell(index+2, excel_col_label_to_index("AL")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("AM")).Value = "iepoch MEL (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AN")).Value = "iepoch MEL (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
				sheet.Cell(index+2, excel_col_label_to_index("AY")).Value = "TRACKADV"
				//sheet.Cell(index+2, excel_col_label_to_index("BI")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Weight)
				sheet.Cell(index+2, excel_col_label_to_index("BV")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("BW")).Value = fmt.Sprintf("%0.2f", order_logistics.OrderLogistics.Weight)
				// sheet.Cell(index+2, excel_col_label_to_index("BX")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Length)
				// sheet.Cell(index+2, excel_col_label_to_index("BY")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Width)
				// sheet.Cell(index+2, excel_col_label_to_index("BZ")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Height)
				sheet.Cell(index+2, excel_col_label_to_index("CC")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("CD")).Value = "FALSE"
			}

		} else if stock_code == excel_export_stock_bne_code {
			xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "elife_ep_bne.xlsx"))
			core.HandleError(err)
			export_excel_file_name = "ELIFE eParcel BNE"

			sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
			for index, order_logistics := range order_logistics_list {
				var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
				// sheet.Cell(index+2, excel_col_label_to_index("C")).Value = "7D83"
				sheet.Cell(index+2, excel_col_label_to_index("C")).Value = order_logistics.LogisticsMode.SenderId
				sheet.Cell(index+2, excel_col_label_to_index("D")).Value = order_logistics.Order.FullTranId

				name1, name2 := split_long_string(order_shipping_address.Name, 35)
				sheet.Cell(index+2, excel_col_label_to_index("E")).Value = name1
				sheet.Cell(index+2, excel_col_label_to_index("F")).Value = name2

				var full_address = order_shipping_address.Street1 + "," + order_shipping_address.Street2
				address1, address2 := split_long_string(full_address, 35)
				sheet.Cell(index+2, excel_col_label_to_index("G")).Value = address1
				sheet.Cell(index+2, excel_col_label_to_index("H")).Value = address2
				sheet.Cell(index+2, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
				sheet.Cell(index+2, excel_col_label_to_index("L")).Value = get_state_short_name(order_shipping_address.Province)
				sheet.Cell(index+2, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
				//sheet.Cell(index+2, excel_col_label_to_index("N")).Value = order_shipping_address.CountryName
				sheet.Cell(index+2, excel_col_label_to_index("N")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("O")).Value = order_shipping_address.Phone
				sheet.Cell(index+2, excel_col_label_to_index("P")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("S")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("V")).Value = "N"
				sheet.Cell(index+2, excel_col_label_to_index("X")).Value = fmt.Sprintf("%s FROM:%s", order_logistics.Order.FullTranId, order_logistics.ShopAccount.Name)
				sheet.Cell(index+2, excel_col_label_to_index("Y")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("Z")).Value = order_logistics.sku_qty_loc
				sheet.Cell(index+2, excel_col_label_to_index("AA")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AC")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("AD")).Value = "iepoch BNE (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AE")).Value = "Basement 2/30 Gow Street"
				sheet.Cell(index+2, excel_col_label_to_index("AI")).Value = "MOOROOKA"
				sheet.Cell(index+2, excel_col_label_to_index("AJ")).Value = "QLD"
				sheet.Cell(index+2, excel_col_label_to_index("AK")).Value = "4105"
				sheet.Cell(index+2, excel_col_label_to_index("AL")).Value = "AU" // 固定
				sheet.Cell(index+2, excel_col_label_to_index("AM")).Value = "iepoch BNE (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AN")).Value = "iepoch BNE (3PL)"
				sheet.Cell(index+2, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
				sheet.Cell(index+2, excel_col_label_to_index("AY")).Value = "TRACKADV"
				// sheet.Cell(index+2, excel_col_label_to_index("BI")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Weight)
				sheet.Cell(index+2, excel_col_label_to_index("BV")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("BW")).Value = fmt.Sprintf("%0.2f", order_logistics.OrderLogistics.Weight)
				// sheet.Cell(index+2, excel_col_label_to_index("BX")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Length)
				// sheet.Cell(index+2, excel_col_label_to_index("BY")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Width)
				// sheet.Cell(index+2, excel_col_label_to_index("BZ")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Height)
				sheet.Cell(index+2, excel_col_label_to_index("CC")).Value = "FALSE"
				sheet.Cell(index+2, excel_col_label_to_index("CD")).Value = "FALSE"
			}
		}
	} else if company_code == excel_export_company_fastway_code && (stock_code == excel_export_stock_mel_code || stock_code == excel_export_stock_bne_code) {
		xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "fastway.xlsx"))
		core.HandleError(err)

		if stock_code == excel_export_stock_mel_code {
			export_excel_file_name = "FASTWAY MEL"
		} else if stock_code == excel_export_stock_bne_code {
			export_excel_file_name = "FASTWAY BNE"
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
			mode_id_str_list = append(mode_id_str_list, fmt.Sprint(mode_id.ModeId))
		}

		/*
SELECT ol.*, o.*, lm.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE ol.mode_id IN (0) AND ol.stock_id = 0
		 */

		sql_get_logistics_list := `SELECT ol.*, o.*, lm.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE  o.status_id = %d AND ol.mode_id IN (%s) AND ol.stock_id = %d`

		if export_type == 0 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_PHZ, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if export_type == 1 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_DPH, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if batch_code != 0 {
			sql_get_logistics_list += fmt.Sprintf(" AND o.batch_code = %d", batch_code)
		}

		if export_type == 1 { //待配货的导出
			sql_get_logistics_list = sql_get_logistics_list + fmt.Sprintf(" AND o.id IN(%s)", order_id_list_str)
		}

		var order_logistics_list = make([]order_logistics_fastway_view_item_t, 0)
		err = db_session.Sql(sql_get_logistics_list).Find(&order_logistics_list)
		core.HandleError(err)

		for index, order_logistics := range order_logistics_list {
			order_logistics_list[index].sku_qty_loc, order_logistics_list[index].loc_name = get_order_sku_for_export(order_logistics.Order.Id, db_session)
		}

		var order_logistics_list_slice = order_logistics_fastway_view_item_slice(order_logistics_list)
		sort.Sort(order_logistics_list_slice)

		order_logistics_list = []order_logistics_fastway_view_item_t(order_logistics_list_slice)

		sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
		for index, order_logistics := range order_logistics_list {
			var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
			sheet.Cell(index+1, excel_col_label_to_index("B")).Value = order_logistics.Order.FullTranId

			var c_cell = sheet.Cell(index+1, excel_col_label_to_index("C"))
			//sheet.Cell(index+1, excel_col_label_to_index("C")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Weight)

			if order_logistics.OrderLogistics.Weight < 1 {
				c_cell.SetFloat(1)
			} else {
				c_cell.SetFloat(math.Floor(order_logistics.OrderLogistics.Weight + 0.5))
			}
			set_cell_style(c_cell)

			//sheet.Cell(index+1, excel_col_label_to_index("D")).Value = "1"
			var d_cell = sheet.Cell(index+1, excel_col_label_to_index("D"))
			d_cell.SetInt(1)
			set_cell_style(d_cell)

			sheet.Cell(index+1, excel_col_label_to_index("H")).Value = order_shipping_address.Name
			var full_address = order_shipping_address.Street1 + "," + order_shipping_address.Street2
			address1, address2 := split_long_string(full_address, 25)
			sheet.Cell(index+1, excel_col_label_to_index("I")).Value = address1
			sheet.Cell(index+1, excel_col_label_to_index("J")).Value = address2
			sheet.Cell(index+1, excel_col_label_to_index("K")).Value = order_shipping_address.CityName
			sheet.Cell(index+1, excel_col_label_to_index("L")).Value = get_state_short_name(order_shipping_address.Province)
			sheet.Cell(index+1, excel_col_label_to_index("M")).Value = order_shipping_address.PostalCode
			sheet.Cell(index+1, excel_col_label_to_index("N")).Value = order_shipping_address.Name
			sheet.Cell(index+1, excel_col_label_to_index("O")).Value = order_logistics.Order.BuyerEmail
			sheet.Cell(index+1, excel_col_label_to_index("P")).Value = order_shipping_address.Phone
			sheet.Cell(index+1, excel_col_label_to_index("R")).Value = "ATL at safe place"
			sheet.Cell(index+1, excel_col_label_to_index("S")).Value = fmt.Sprintf("FROM:%s", order_logistics.ShopAccount.Name)
			sheet.Cell(index+1, excel_col_label_to_index("T")).Value = order_logistics.sku_qty_loc
			sheet.Cell(index+1, excel_col_label_to_index("U")).Value = get_fastway_packaging_code(order_logistics.LogisticsMode)
		}

	} else if company_code == excel_export_company_toll_code && (stock_code == excel_export_stock_mel_code || stock_code == excel_export_stock_bne_code) {
		xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "toll.xlsx"))
		core.HandleError(err)
		if stock_code == excel_export_stock_mel_code {
			export_excel_file_name = "TOLL MEL"
		} else if stock_code == excel_export_stock_bne_code {
			export_excel_file_name = "TOLL BNE"
		}

		type ModeId struct {
			ModeId int64 `xorm:"'mode_id'" json:"mode_id"`
		}

		var mode_id_list = make([]ModeId, 0)
		sql_get_mode_id_list := `SELECT lm.id AS mode_id FROM erp_logistics_mode lm WHERE lm.stock_id = %d AND lm.name LIKE '%%TOLL%%';`
		sql_get_mode_id_list = fmt.Sprintf(sql_get_mode_id_list, stock_id)
		err := db_session.Sql(sql_get_mode_id_list).Find(&mode_id_list)
		core.HandleError(err)

		var mode_id_str_list = make([]string, 0)
		for _, mode_id := range mode_id_list {
			mode_id_str_list = append(mode_id_str_list, fmt.Sprint(mode_id.ModeId))
		}

		/*
SELECT ol.*, o.*, lm.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE ol.mode_id IN (0) AND ol.stock_id = 0
		 */

		sql_get_logistics_list := `SELECT ol.*, o.*, lm.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_order o ON o.id = ol.order_id
    INNER JOIN erp_logistics_mode lm ON lm.id = ol.mode_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE o.status_id = %d AND ol.mode_id IN (%s) AND ol.stock_id = %d`

		if export_type == 0 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_PHZ, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if export_type == 1 {
			sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_DPH, strings.Join(mode_id_str_list, ","), stock_id)
		}

		if batch_code != 0 {
			sql_get_logistics_list += fmt.Sprintf(" AND o.batch_code = %d", batch_code)
		}

		if export_type == 1 { //待配货的导出
			sql_get_logistics_list = sql_get_logistics_list + fmt.Sprintf(" AND o.id IN(%s)", order_id_list_str)
		}

		var order_logistics_list_in_db = make([]order_logistics_toll_view_item_t, 0)
		err = db_session.Sql(sql_get_logistics_list).Find(&order_logistics_list_in_db)
		core.HandleError(err)

		var order_logistics_list = make([]order_logistics_toll_view_item_t, 0)
		for _, order_logistics := range order_logistics_list_in_db {
			var ol = order_logistics

			var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
			ol.order_shipping_address = order_shipping_address

			// TOLL 的 导出模板，B列 service ID。如果是墨尔本仓库，邮编第一个数字是“3”的话，填充2，否则填充4 。 如果是布里斯班仓库，邮编第一个数字是“4”的话，填充2，否则填充4 。
			//var service_id string = get_toll_service_id(stock_code, order_shipping_address)

			// 逻辑改了，现在不采用上面的方式了，而是从一个表里面取。只针对MEL的TOLL的 Service Id
			var service_id string = ""

			if stock_code == excel_export_stock_mel_code {
				service_id = get_mel_toll_service_id(order_shipping_address.PostalCode, db_session)
			}

			if stock_code == excel_export_stock_bne_code {
				service_id = get_toll_service_id(excel_export_stock_bne_code, order_shipping_address)
			}
			ol.service_id = service_id

			ol.sku_qty_loc, ol.loc_name, ol.qty = get_toll_order_sku_for_export(order_logistics.Order.Id, service_id, stock_code, db_session)

			order_logistics_list = append(order_logistics_list, ol)
		}

		var order_logistics_list_slice = order_logistics_toll_view_item_slice(order_logistics_list)
		sort.Sort(order_logistics_list_slice)

		order_logistics_list = []order_logistics_toll_view_item_t(order_logistics_list_slice)

		sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
		for index, order_logistics := range order_logistics_list {

			var service_id = order_logistics.service_id
			var order_shipping_address = order_logistics.order_shipping_address

			sheet.Cell(index+1, excel_col_label_to_index("A")).Value = order_logistics.LogisticsMode.SenderId
			sheet.Cell(index+1, excel_col_label_to_index("B")).Value = "IPEC"
			sheet.Cell(index+1, excel_col_label_to_index("C")).Value = service_id
			sheet.Cell(index+1, excel_col_label_to_index("D")).Value = order_shipping_address.Name
			sheet.Cell(index+1, excel_col_label_to_index("E")).Value = order_shipping_address.Name // 和 D　是一样的？

			var full_address = order_shipping_address.Street1 + "," + order_shipping_address.Street2
			address1, address2 := split_long_string(full_address, 25)
			sheet.Cell(index+1, excel_col_label_to_index("F")).Value = address1
			sheet.Cell(index+1, excel_col_label_to_index("G")).Value = address2
			sheet.Cell(index+1, excel_col_label_to_index("H")).Value = order_shipping_address.CityName
			sheet.Cell(index+1, excel_col_label_to_index("I")).Value = get_state_short_name(order_shipping_address.Province)
			sheet.Cell(index+1, excel_col_label_to_index("J")).Value = order_shipping_address.PostalCode
			sheet.Cell(index+1, excel_col_label_to_index("L")).Value = "CARTON"
			sheet.Cell(index+1, excel_col_label_to_index("M")).Value = get_toll_packaging_code(order_logistics.OrderLogistics) // 包裹类型要计算

			/*
			if service_id == "2" && stock_code == excel_export_stock_mel_code {
				// service_id 为 2 的不拆包
				sheet.Cell(index+1, excel_col_label_to_index("N")).Value = "1"
			} else {
				// service_id 不为 2 的要拆包
				sheet.Cell(index+1, excel_col_label_to_index("N")).Value = fmt.Sprint(order_logistics.qty)
			}
			*/

			// 后来又改为拆包的模式
			// if get_toll_packaged(service_id, stock_code) == false {
			// service_id 为 2 的不拆包
			// sheet.Cell(index+1, excel_col_label_to_index("N")).Value = "1"
			// } else {
			// service_id 不为 2 的要拆包
			sheet.Cell(index+1, excel_col_label_to_index("N")).Value = fmt.Sprint(order_logistics.qty)
			// }

			sheet.Cell(index+1, excel_col_label_to_index("O")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Length)
			sheet.Cell(index+1, excel_col_label_to_index("P")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Width)
			sheet.Cell(index+1, excel_col_label_to_index("Q")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Height/float64(order_logistics.qty))
			sheet.Cell(index+1, excel_col_label_to_index("R")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.DimWeight/float64(order_logistics.qty))

			if order_logistics.OrderLogistics.Weight/float64(order_logistics.qty) >= 1 {
				sheet.Cell(index+1, excel_col_label_to_index("S")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Weight/float64(order_logistics.qty))
			} else {
				sheet.Cell(index+1, excel_col_label_to_index("S")).Value = fmt.Sprintf("%0.1f", math.Ceil(order_logistics.OrderLogistics.Weight/float64(order_logistics.qty)))
			}

			sheet.Cell(index+1, excel_col_label_to_index("T")).Value = order_logistics.Order.FullTranId
			sheet.Cell(index+1, excel_col_label_to_index("U")).Value = fmt.Sprintf("FROM:%s", order_logistics.ShopAccount.Name)
			sheet.Cell(index+1, excel_col_label_to_index("V")).Value = "ATL at safe place"
			sheet.Cell(index+1, excel_col_label_to_index("W")).Value = order_logistics.sku_qty_loc
			sheet.Cell(index+1, excel_col_label_to_index("X")).Value = order_logistics.Order.FullTranId // 和 T 一样的？

			sheet.Cell(index+1, excel_col_label_to_index("Y")).Value = order_shipping_address.Phone
			sheet.Cell(index+1, excel_col_label_to_index("Z")).Value = order_logistics.Order.BuyerEmail
		}
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

func split_long_string(long_str string, letter_limit int) (str1 string, str2 string) {
	if letter_limit <= 0 {
		return long_str, ""
	}

	var long_str_rune_list = []rune(long_str)
	var long_str_rune_list_len = len(long_str_rune_list)

	if long_str_rune_list_len <= letter_limit {
		return long_str, ""
	}

	var space_rune = ([]rune(" "))[0]
	var index int
	for index = letter_limit; index > 0; index-- {
		if long_str_rune_list[index] == space_rune {
			break
		}
	}

	return strings.TrimSpace(string(long_str_rune_list[0:index])), strings.TrimSpace(string(long_str_rune_list[index:]))
}

func get_order_sku_for_export(order_id int64, db_session *xorm.Session) (order_sku_qty_loc_return string, loc_name string) {
	/*

	single product, single quantity

	1. 对于 SS， 显示 SKU 数量 [仓位号]
	2. 对于 SM， 显示 SKU DOG数量 [仓位号]
	3. 对于 MM, 显示 MM

	 */

	var order_sku_quantity_list = order_service.Get_order_sku_qty_loc_list(order_id, db_session)

	var order_sku_quantity_list_len = len(order_sku_quantity_list)
	if order_sku_quantity_list_len < 1 {
		return "", ""
	} else if order_sku_quantity_list_len == 1 {
		var order_sku_quantity = order_sku_quantity_list[0]
		if order_sku_quantity.Quantity == 1 {
			return fmt.Sprintf("%s %d [%s]", order_sku_quantity.Sku, order_sku_quantity.Quantity, order_sku_quantity.StorageLocName), order_sku_quantity.StorageLocName
		} else if order_sku_quantity.Quantity > 1 {
			return fmt.Sprintf("%s DOG%d [%s]", order_sku_quantity.Sku, order_sku_quantity.Quantity, order_sku_quantity.StorageLocName), order_sku_quantity.StorageLocName
		} else {
			return "", ""
		}
	} else { // if order_sku_quantity_list_len > 1 {
		return "MM", ""
	}
}

// 返回 false 表示 不拆包， 返回 true 表示 拆包
func get_toll_packaged(service_id string, stock_code int) bool {
	if service_id == "2" && stock_code == excel_export_stock_mel_code {
		// service_id 为 2 的 MEL的 不拆包
		// sheet.Cell(index+1, excel_col_label_to_index("N")).Value = "1"
		return false
	} else {
		// service_id 不为 2 的要拆包
		// sheet.Cell(index+1, excel_col_label_to_index("N")).Value = fmt.Sprint(order_logistics.qty)
		return true
	}
}

func get_toll_order_sku_for_export(order_id int64, service_id string, stock_code int, db_session *xorm.Session) (order_sku_qty_loc_return string, loc_name_return string, qty_return int64) {
	var order_sku_quantity_list = order_service.Get_order_sku_qty_loc_list(order_id, db_session)

	var order_sku_quantity_list_len = len(order_sku_quantity_list)
	if order_sku_quantity_list_len < 1 {
		return "", "", 0
	} else if order_sku_quantity_list_len == 1 {
		var order_sku_quantity = order_sku_quantity_list[0]
		if order_sku_quantity.Quantity == 1 {
			return fmt.Sprintf("%s %d [%s]", order_sku_quantity.Sku, order_sku_quantity.Quantity, order_sku_quantity.StorageLocName), order_sku_quantity.StorageLocName, order_sku_quantity.Quantity
		} else if order_sku_quantity.Quantity > 1 {
			// 后来又改为都需要拆包
			//if get_toll_packaged(service_id, stock_code) == false {
			// service_id 为 2 的不拆包
			//	return fmt.Sprintf("%s DOG%d [%s]", order_sku_quantity.Sku, order_sku_quantity.Quantity, order_sku_quantity.StorageLocName), order_sku_quantity.StorageLocName, order_sku_quantity.Quantity
			//} else {
			// service_id 不为 2 的要拆包
			return fmt.Sprintf("%s DOG [%s]", order_sku_quantity.Sku, order_sku_quantity.StorageLocName), order_sku_quantity.StorageLocName, order_sku_quantity.Quantity
			//}

		} else {
			return "", "", 0
		}
	} else { // if order_sku_quantity_list_len > 1 {
		var total_qty int64 = 0
		for _, sku_qty := range order_sku_quantity_list {
			total_qty = total_qty + sku_qty.Quantity
		}
		return "MM", "", total_qty
	}
}

func get_fastway_packaging_code(logistics_mode ddpz_models.LogisticsModeInfo) string {
	/*
		FASTWAY
		字段    代码
		Label    1
		A2       4
		A3       5
		A4       6
		A5       7
		Small    16
		Medium   17
		Large    18
	*/

	if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "Label") {
		return "1"
	} else if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "A2") {
		return "4"
	} else if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "A3") {
		return "5"
	} else if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "A4") {
		return "6"
	} else if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "A5") {
		return "7"
	} else if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "Small") {
		return "16"
	} else if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "Medium") {
		return "17"
	} else if strings.Contains(logistics_mode.Name, "FASTWAY") && strings.Contains(logistics_mode.Name, "Large") {
		return "18"
	} else {
		return ""
	}
}

func get_toll_packaging_code(order_logistics ddlc_models.OrderLogisticsInfo) string {

	// 临时改成都返回 C
	return "C"

	/*
		包裹类型    代表
		C         TOLL < 25KG
		LC        TOLL > 25KG
		VL        TOLL > 50KG
		XL        TOLL > 100KG
		这个是打折后的体积重 与 实际重 取大值 判断
	 */

	/*
   var weight = math.Max(order_logistics.DimWeight, order_logistics.Weight)
   if weight <= 25 {
	   return "C"
   } else if weight > 25 && weight <= 50 {
	   return "LC"
   } else if weight > 50 && weight <= 100 {
	   return "VL"
   } else if weight > 100 {
	   return "XL"
   } else {
	   return ""
   }
	*/
}

// 主要用于 BNE 的 service id
func get_toll_service_id(stock_code int, order_shipping_address ddlc_models.OrderShippingAddressInfo) string {
	var service_id string = ""
	var postalcode_rune_list = []rune(order_shipping_address.PostalCode)
	if len(postalcode_rune_list) > 0 {
		if stock_code == excel_export_stock_mel_code {
			if postalcode_rune_list[0] == ([]rune("3"))[0] {
				service_id = "2"
			} else {
				service_id = "4"
			}
		} else if stock_code == excel_export_stock_bne_code {
			/*
			if postalcode_rune_list[0] == ([]rune("4"))[0] {
				service_id = "2"
			} else {
				service_id = "4"
			}
			*/

			/*
				布里斯班的TOLL的导出模板，第二列
				刘兴宗_浙江澳中欧易事多 2017-06-07 12:12:19
				这里全部固定为4
			 */

			service_id = "4"
		}
	}

	return service_id
}

func get_state_short_name(state_name string) string {
	/*
前三个字母		缩写
new    nsw	    NSW
vic             VIC
que    qld	    QLD
tas             TAS
sou    sa	    SA
wes    wa	   WA
nor    nt	   NT
aus    act	  ACT
	 */

	state_name = strings.ToLower(strings.TrimSpace(state_name))

	if strings.HasPrefix(state_name, "new") || strings.HasPrefix(state_name, "nsw") {
		return "NSW"
	} else if strings.HasPrefix(state_name, "vic") {
		return "VIC"
	} else if strings.HasPrefix(state_name, "que") || strings.HasPrefix(state_name, "qld") {
		return "QLD"
	} else if strings.HasPrefix(state_name, "tas") {
		return "TAS"
	} else if strings.HasPrefix(state_name, "sou") || strings.HasPrefix(state_name, "sa") {
		return "SA"
	} else if strings.HasPrefix(state_name, "wes") || strings.HasPrefix(state_name, "wa") {
		return "WA"
	} else if strings.HasPrefix(state_name, "nor") || strings.HasPrefix(state_name, "nt") {
		return "NT"
	} else if strings.HasPrefix(state_name, "aus") || strings.HasPrefix(state_name, "act") {
		return "ACT"
	} else {
		return ""
	}
}

func get_mel_toll_service_id(postcode string, db_session *xorm.Session) string {
	var mel_toll_service_id ddpz_models.MelTollSeviceIdInfo
	postcode_without_prefix_zero := strings.TrimPrefix(postcode, "0")
	_, err := db_session.Where(fmt.Sprintf("postcode = '%s' OR postcode = '%s'", postcode, postcode_without_prefix_zero)).Get(&mel_toll_service_id)
	core.HandleError(err)

	return mel_toll_service_id.ServiceId
}
