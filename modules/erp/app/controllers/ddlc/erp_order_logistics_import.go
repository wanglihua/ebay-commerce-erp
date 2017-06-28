package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrderLogisticsImport struct {
	*revel.Controller
}

func (c ErpOrderLogisticsImport) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_logistics_import/order_logistics_import_index.html")
}

// 导入运单号
func (c ErpOrderLogisticsImport) ImportTrackingNumber() revel.Result {
	db_session := c.DbSession

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type tracking_number_import_t struct {
		OrderNumber     string  `json:"order_number"`
		TrackingNumber  string  `json:"tracking_number"`
		ErrorMessage    string  `json:"error_message"`
		order_logistics models.OrderLogisticsInfo
	}

	tracking_number_import_list := make([]tracking_number_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow

		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			var cell *xlsx.Cell
			var cell_val string
			var err error

			var tracking_number_import tracking_number_import_t

			col_index := 0

			// 订单代码
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)

			if cell_val == "" {
				break
			}

			tracking_number_import.OrderNumber = cell_val
			col_index++

			// 运单号
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			tracking_number_import.TrackingNumber = cell_val
			col_index++

			tracking_number_import_list = append(tracking_number_import_list, tracking_number_import)
		}
	}

	if len(tracking_number_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入的列表为空！",
		})
	}

	var tracking_number_import_error_list = make([]tracking_number_import_t, 0)

	// 错误检查
	for index, order_import := range tracking_number_import_list {
		/*
			SELECT ol.* FROM erp_order o
			    INNER JOIN erp_order_logistics ol ON ol.order_id = o.id
			WHERE o.full_tran_id = ''
		*/

		var order_logistics models.OrderLogisticsInfo
		var sql_get_order_logistics = `SELECT ol.* FROM erp_order o
			    INNER JOIN erp_order_logistics ol ON ol.order_id = o.id
			WHERE o.full_tran_id = '%s'`
		sql_get_order_logistics = fmt.Sprintf(sql_get_order_logistics, order_import.OrderNumber)
		has, err := db_session.Sql(sql_get_order_logistics).Get(&order_logistics)
		core.HandleError(err)

		tracking_number_import_list[index].order_logistics = order_logistics

		if has == false {
			tracking_number_import_list[index].ErrorMessage = fmt.Sprintf("订单代码 %s 对应的订单物流不存在！", order_import.OrderNumber)
			tracking_number_import_error_list = append(tracking_number_import_error_list, tracking_number_import_list[index])
		}
	}

	if len(tracking_number_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": tracking_number_import_error_list,
		})
	}

	for _, order_import := range tracking_number_import_list {
		var order_logistics = order_import.order_logistics
		order_logistics.TrackingNum = order_import.TrackingNumber
		_, err := db_session.Id(order_logistics.Id).Cols("tracking_num").Update(&order_logistics)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

// 导入实际运费
func (c ErpOrderLogisticsImport) ImportActualFee() revel.Result {
	db_session := c.DbSession

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type actual_fee_import_t struct {
		OrderNumber     string   `json:"order_number"`
		ActualFee       float64  `json:"actual_fee"`
		ErrorMessage    string   `json:"error_message"`
		ActualFeeStr    string   `json:"actual_fee_str"`
		order_logistics models.OrderLogisticsInfo
	}

	actual_fee_import_list := make([]actual_fee_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow

		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			var cell *xlsx.Cell
			var cell_val string
			var err error

			var tracking_number_import actual_fee_import_t

			col_index := 0

			// 订单代码
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)

			if cell_val == "" {
				break
			}

			tracking_number_import.OrderNumber = cell_val
			col_index++

			// 运单号
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			tracking_number_import.ActualFeeStr = cell_val
			col_index++

			actual_fee_import_list = append(actual_fee_import_list, tracking_number_import)
		}
	}

	if len(actual_fee_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入的列表为空！",
		})
	}

	var order_import_error_list = make([]actual_fee_import_t, 0)

	// 错误检查
	for index, actual_fee_import := range actual_fee_import_list {
		/*
			SELECT ol.* FROM erp_order o
			    INNER JOIN erp_order_logistics ol ON ol.order_id = o.id
			WHERE o.full_tran_id = ''
		*/

		var order_logistics models.OrderLogisticsInfo
		var sql_get_order_logistics = `SELECT ol.* FROM erp_order o
			    INNER JOIN erp_order_logistics ol ON ol.order_id = o.id
			WHERE o.full_tran_id = '%s'`
		sql_get_order_logistics = fmt.Sprintf(sql_get_order_logistics, actual_fee_import.OrderNumber)
		has, err := db_session.Sql(sql_get_order_logistics).Get(&order_logistics)
		core.HandleError(err)

		actual_fee_import_list[index].order_logistics = order_logistics

		if has == false {
			actual_fee_import_list[index].ErrorMessage = fmt.Sprintf("订单代码 %s 对应的订单物流不存在！", actual_fee_import.OrderNumber)
			order_import_error_list = append(order_import_error_list, actual_fee_import_list[index])
		}

		actual_fee_import.ActualFee, err = strconv.ParseFloat(actual_fee_import.ActualFeeStr, 64)
		if err != nil {
			actual_fee_import_list[index].ErrorMessage = fmt.Sprintf("订单代码 %s 对应实际运费 %s 的格式不对！", actual_fee_import.OrderNumber, actual_fee_import.ActualFeeStr)
			order_import_error_list = append(order_import_error_list, actual_fee_import_list[index])
		}
	}

	if len(order_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": order_import_error_list,
		})
	}

	for _, order_import := range actual_fee_import_list {
		var order_logistics = order_import.order_logistics
		order_logistics.ActualFee = order_import.ActualFee
		_, err := db_session.Id(order_logistics.Id).Cols("tracking_num").Update(&order_logistics)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}
