package controllers

import (
	"github.com/revel/revel"

	"fmt"
	"matrix/core"
	"matrix/modules/erp/models"
	"github.com/tealeg/xlsx"
	"strconv"
)

type ErpProductImport struct {
	*revel.Controller
}

func (c ErpProductImport) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_import/product_import.html")
}

func (c ErpProductImport) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here
	login_user := core.GetLoginUser(c.Session)
	login_user_id := login_user.UserId
	query = query.Where(fmt.Sprintf("%s = ?", models.ProductImportCols.UserId), login_user_id)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_import_list := make([]models.ProductImportInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_import_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductImportInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_import_list,
		Total: count,
	})
}

func (c ErpProductImport) Import() revel.Result {
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

	product_import_list := make([]models.ProductImportInfo, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	sheet_count := len(xlsx_file.Sheets)
	//revel.TRACE.Printf("sheet_count: %d\n", sheet_count)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//revel.TRACE.Printf("max_row: %d\n", max_row)
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error
			col_index := 0

			var product_import models.ProductImportInfo

			//SKU
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_import.Sku = cell_val
			col_index++

			//名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_import.Name = cell_val
			col_index++

			//英文名
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_import.EnglishName = cell_val
			col_index++

			//成本价(元)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			cost_price, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_import.CostPrice = cost_price
			col_index++

			//重量(kg)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			weight, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_import.Weight = weight
			col_index++

			product_import_list = append(product_import_list, product_import)
		}
	}

	if len(product_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success":  false,
			"message":  "Excel文件中导入商品列表为空！",
		})
	}

	login_user := core.GetLoginUser(c.Session)
	login_user_id := login_user.UserId

	for _, product_import := range product_import_list {
		product_import.UserId = login_user_id
		_, err := db_session.Insert(&product_import)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success":  true,
		"message":  "Excel数据上传成功！",
	})
}

func (c ErpProductImport) Confirm() revel.Result {
	db_session := c.DbSession

	sql := `
DELIMITER //
CREATE PROCEDURE erp_trans_import_product(IN login_user_id INT)
    BEGIN
        INSERT INTO erp_product (sku, name, english_name, cost_price, weight, create_time, update_time, version)
            SELECT pi.sku, pi.name, pi.english_name, pi.cost_price, pi.weight, pi.import_time, pi.import_time, 1
            FROM erp_product_import pi
            WHERE pi.user_id = login_user_id;
        DELETE FROM erp_product_import WHERE user_id = login_user_id;
    END
//
DELIMITER ;
	`
	//login_user := core.GetLoginUser(c.Session)
	//login_user_id := login_user.UserId

	//_, err := db_session.Exec(sql, login_user_id, login_user_id)
	_, err := db_session.Exec(sql)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "商品数据导入成功!"})
}

func (c ErpProductImport) Clear() revel.Result {
	db_session := c.DbSession

	login_user := core.GetLoginUser(c.Session)
	login_user_id := login_user.UserId

	_, err := db_session.
	Where(fmt.Sprintf("%s = ?", models.ProductImportCols.UserId), login_user_id).
		Delete(new(models.ProductImportInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "临时导入数据清空成功!"})
}
