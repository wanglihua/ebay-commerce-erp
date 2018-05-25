package controllers

import (
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/order_service"
	"matrix/modules/erp/service/stock_service"
	"sort"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
	"strings"
)

type ErpOrderLogisticsExpressExport struct {
	*revel.Controller
}

type order_logistics_express_view_item_t struct {
	OrderLogistics ddlc_models.OrderLogisticsInfo `xorm:"extends" json:"ol"`
	LogisticsMode  ddpz_models.LogisticsModeInfo  `xorm:"extends" json:"lm"`
	Order          ddlc_models.OrderInfo          `xorm:"extends" json:"o"`
	ShopAccount    xtpz_models.ShopAccountInfo    `xorm:"extends" json:"sa"`
	sku_qty_loc    string                         `xorm:"-"`
	loc_name       string                         `xorm:"-"`
}

type order_logistics_express_view_item_slice []order_logistics_express_view_item_t

func (s order_logistics_express_view_item_slice) Len() int           { return len(s) }
func (s order_logistics_express_view_item_slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s order_logistics_express_view_item_slice) Less(i, j int) bool { return s[i].loc_name < s[j].loc_name }

func (c ErpOrderLogisticsExpressExport) ExportExcel() revel.Result {
	var db_session = c.DbSession

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var order_id_list = make([]int64, 0)
	c.Params.Bind(&order_id_list, "order_id_list")

	var order_id_str_list = make([]string, 0)
	for _, order_id := range order_id_list {
		order_id_str_list = append(order_id_str_list, fmt.Sprint(order_id))
	}

	var order_id_list_str = strings.Join(order_id_str_list, ",")
	revel.TRACE.Println("order_id_list_str", order_id_list_str)

	// 判断当前登录用户是否有选择的仓库的权限
	var mel_stock_id = stock_service.Get_mel_stock_id(db_session)
	var stock_list = stock_service.Get_login_user_stock_list(db_session, c.Session)
	var has_stock = false
	for _, stock := range stock_list {
		if mel_stock_id == stock.Id {
			has_stock = true
			break
		}
	}

	if has_stock == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "你没有导出墨尔本仓库物流的权限！",
		})
	}

	// 采用和 eparcel_mel 同样的模板
	// eparcel_mel.xlsx

	var excel_template_logistics_export_dir = filepath.Join(revel.BasePath, "/static/excel-template/logistics_export")

	var xlsx_file *xlsx.File

	/*
SELECT ol.*, lm.*, o.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_logistics_mode lm on lm.id = ol.mode_id
    INNER JOIN erp_order o ON o.id = ol.order_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE o.status_id = 0 AND o.hold_type_id = 0 AND ol.stock_id = 0;
	 */

	sql_get_logistics_list := `SELECT ol.*, lm.*, o.*, sa.*
FROM erp_order_logistics ol
    INNER JOIN erp_logistics_mode lm on lm.id = ol.mode_id
    INNER JOIN erp_order o ON o.id = ol.order_id
    LEFT JOIN erp_shop_account sa ON sa.id = o.shop_id
WHERE o.status_id = %d AND o.hold_type_id = %d AND ol.stock_id = %d`

	sql_get_logistics_list = fmt.Sprintf(sql_get_logistics_list, models.ORDER_STATUS_HOLD, ddpz_models.ORDER_HOLD_TYPE_EXPRESS, mel_stock_id)

	if batch_code != 0 {
		sql_get_logistics_list += fmt.Sprintf(" AND o.batch_code = %d", batch_code)
	}

	sql_get_logistics_list = sql_get_logistics_list + fmt.Sprintf(" AND o.id IN (%s)", order_id_list_str)

	var order_logistics_list = make([]order_logistics_express_view_item_t, 0)
	err := db_session.Sql(sql_get_logistics_list).Find(&order_logistics_list)
	core.HandleError(err)

	for index, order_logistics := range order_logistics_list {
		order_logistics_list[index].sku_qty_loc, order_logistics_list[index].loc_name = get_order_sku_for_export(order_logistics.Order.Id, db_session)
	}

	var order_logistics_list_slice = order_logistics_express_view_item_slice(order_logistics_list)
	sort.Sort(order_logistics_list_slice)

	order_logistics_list = []order_logistics_express_view_item_t(order_logistics_list_slice)

	// xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "eparcel_mel.xlsx"))
	xlsx_file, err = xlsx.OpenFile(filepath.Join(excel_template_logistics_export_dir, "elife_ep_mel_express.xlsx"))
	core.HandleError(err)

	sheet := xlsx_file.Sheets[0] // 模板里面应该至少有一个sheet
	for index, order_logistics := range order_logistics_list {
		var order_shipping_address = order_service.Get_order_shipping_address(order_logistics.Order.Id, db_session)
		/*
		sheet.Cell(index+2, excel_col_label_to_index("C")).Value = "7J85"
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
		sheet.Cell(index+2, excel_col_label_to_index("AE")).Value = "2 Venture Court"
		sheet.Cell(index+2, excel_col_label_to_index("AI")).Value = "DANDENONG SOUTH"
		sheet.Cell(index+2, excel_col_label_to_index("AJ")).Value = "VIC"
		sheet.Cell(index+2, excel_col_label_to_index("AK")).Value = "3175"
		sheet.Cell(index+2, excel_col_label_to_index("AL")).Value = "AU" // 固定
		sheet.Cell(index+2, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
		sheet.Cell(index+2, excel_col_label_to_index("AY")).Value = "TRACKADV"
		sheet.Cell(index+2, excel_col_label_to_index("BI")).Value = fmt.Sprintf("%0.1f", order_logistics.OrderLogistics.Weight)
		sheet.Cell(index+2, excel_col_label_to_index("BO")).Value = "FALSE"
		sheet.Cell(index+2, excel_col_label_to_index("BP")).Value = "FALSE"
		*/

		// sheet.Cell(index+2, excel_col_label_to_index("C")).Value = "7J83"
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
		//sheet.Cell(index+2, excel_col_label_to_index("AE")).Value = "2 Venture Court"
		sheet.Cell(index+2, excel_col_label_to_index("AE")).Value = "54-58 Arkwright Drive"
		sheet.Cell(index+2, excel_col_label_to_index("AI")).Value = "DANDENONG SOUTH"
		sheet.Cell(index+2, excel_col_label_to_index("AJ")).Value = "VIC"
		sheet.Cell(index+2, excel_col_label_to_index("AK")).Value = "3175"
		sheet.Cell(index+2, excel_col_label_to_index("AL")).Value = "AU" // 固定
		sheet.Cell(index+2, excel_col_label_to_index("AM")).Value = "iepoch MEL (3PL)"
		sheet.Cell(index+2, excel_col_label_to_index("AN")).Value = "iepoch MEL (3PL)"
		sheet.Cell(index+2, excel_col_label_to_index("AX")).Value = order_logistics.Order.BuyerEmail
		sheet.Cell(index+2, excel_col_label_to_index("AY")).Value = "TRACKADV"
		sheet.Cell(index+2, excel_col_label_to_index("BV")).Value = "FALSE"
		sheet.Cell(index+2, excel_col_label_to_index("BW")).Value = fmt.Sprintf("%0.2f", order_logistics.OrderLogistics.Weight)
		// sheet.Cell(index+2, excel_col_label_to_index("BX")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Length)
		// sheet.Cell(index+2, excel_col_label_to_index("BY")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Width)
		// sheet.Cell(index+2, excel_col_label_to_index("BZ")).Value = fmt.Sprintf("%0.0f", order_logistics.OrderLogistics.Height)
		sheet.Cell(index+2, excel_col_label_to_index("CC")).Value = "FALSE"
		sheet.Cell(index+2, excel_col_label_to_index("CD")).Value = "FALSE"
	}

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口
	file_name := "快递类.xlsx"

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}
