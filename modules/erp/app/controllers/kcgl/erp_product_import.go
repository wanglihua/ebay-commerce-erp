package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/import_batch_code_service"
	"matrix/modules/erp/service/product_import_batch_service"
	"matrix/modules/erp/service/status_service"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpProductImport struct {
	*revel.Controller
}

func (c ErpProductImport) Index() revel.Result {
	pass_service.Check_login_user_pass_or([]int64{pass_res.KCGL_PRODUCT_IMPORT, pass_res.KCGL_PRODUCT_EXPORT}, c.Session)

	var has_product_export_pass = pass_service.Is_login_user_has_pass(pass_res.KCGL_PRODUCT_EXPORT, c.Session)
	c.RenderArgs["has_product_export_pass"] = has_product_export_pass

	return c.RenderTemplate("erp/kcgl/product_import/product_import.html")
}

func (c ErpProductImport) PageData() revel.Result {
	pass_service.Check_login_user_pass_or([]int64{pass_res.KCGL_PRODUCT_IMPORT, pass_res.KCGL_PRODUCT_EXPORT}, c.Session)

	var db_session = c.DbSession

	// 获得导入的批次列表
	var batch_code_list = product_import_batch_service.Get_product_import_batch_code_list(db_session)

	return c.RenderJson(map[string]interface{}{
		"batch_code":      0,
		"batch_code_list": batch_code_list,
	})
}

func (c ErpProductImport) ListData() revel.Result {
	pass_service.Check_login_user_pass_or([]int64{pass_res.KCGL_PRODUCT_IMPORT, pass_res.KCGL_PRODUCT_EXPORT}, c.Session)

	db_session := c.DbSession

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var product_sku_query string
	c.Params.Bind(&product_sku_query, "product_sku_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_view_item_t struct {
		Product         kcgl_models.ProductInfo           `xorm:"extends" json:"prd"`
		ProductCate     kcgl_models.ProductCateInfo       `xorm:"extends" json:"cate"`
		Supplier        models.SupplierInfo          `xorm:"extends" json:"sup"`
		PackageMaterial models.PackageMaterialInfo   `xorm:"extends" json:"pm"`
		Status          models.StatusInfo            `xorm:"extends" json:"st"`
	}

	/*
SELECT prd.*, cate.*, sup.*, pm.*, st.*
FROM erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id
WHERE prd.sku LIKE '%abc%' AND prd.batch_code = 0
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("prd.*, cate.*, sup.*, pm.*, st.*")
	sql_builder.Count("count(prd.id) as product_count")
	sql_builder.From(`erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id`)

	//query extra filter here

	if strings.TrimSpace(product_sku_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_sku_query))
	}

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.batch_code = %d", batch_code))
	}

	sql_builder.Where(filter)

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("prd.id ASC")
	}

	sql_builder.Limit(offset, limit)

	product_list := make([]product_view_item_t, 0, limit)
	err := db_session.Sql(sql_builder.SqlList()).Find(&product_list)
	core.HandleError(err)

	var product_count struct {
		Count int64 `xorm:"'product_count'"`
	}
	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: product_count.Count,
	})
}

func (c ErpProductImport) Upload() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_IMPORT, c.Session)

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

	// 编号 SKU 中文名称  采购价（人民币元） 产品长度（cm） 产品宽度（cm） 产品高度（cm） 实际重量(kg) 状态
	type product_import_t struct {
		ord_num    string
		sku        string
		cn_name    string
		cost_price float64
		length     float64
		width      float64
		height     float64
		weight     float64
		status     string
	}

	product_import_list := make([]product_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
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

			var product_import product_import_t

			col_index := 0

			if len(row.Cells) < 9 { // 数据有 9 列
				break
			}

			// 编号 SKU 中文名称 采购价（人民币元） 产品长度（cm） 产品宽度（cm） 产品高度（cm） 实际重量(kg)

			// 编号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			product_import.ord_num = cell_val
			col_index++

			//SKU
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			//revel.TRACE.Println(cell_val)
			product_import.sku = cell_val
			col_index++

			// 中文名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			//revel.TRACE.Println(cell_val)
			product_import.cn_name = cell_val
			col_index++

			// 采购价（人民币元）
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			cost_price, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_import.cost_price = cost_price
			col_index++

			// length
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_length, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_import.length = product_length
			col_index++

			// width
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_width, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_import.width = product_width
			col_index++

			// height
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_height, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_import.height = product_height
			col_index++

			// weight
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			weight, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_import.weight = weight
			col_index++

			// 状态
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			//revel.TRACE.Println(cell_val)
			product_import.status = cell_val
			col_index++

			product_import_list = append(product_import_list, product_import)
		}
	}

	if len(product_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	var login_user = core.GetLoginUser(c.Session)
	var import_batch_code = import_batch_code_service.Get_import_batch_code(login_user.LoginName)

	var product_status_list = status_service.Get_status_list(db_session, models.STATUS_TYPE_PRODUCT)

	var get_product_status_id_by_name = func(status_name string) int64 {
		for _, status := range product_status_list {
			if status.Name == status_name {
				return status.Id
			}
		}

		return 0
	}

	has_zero_data := false
	for _, product_import := range product_import_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_import.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_import.sku
		product.CnName = product_import.cn_name
		product.Name = product_import.sku
		product.CostPrice = product_import.cost_price
		product.Length = product_import.length
		product.Width = product_import.width
		product.Height = product_import.height
		product.Weight = product_import.weight

		var status_name = product_import.status
		// 查看状态是否存在
		var status_id = get_product_status_id_by_name(status_name)

		if status_id == 0 {
			/*
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("序号 %s 的商品状态 %s 无效", product_import.ord_num, product_import.status),
			})
			*/

			panic(fmt.Sprintf("序号 %s 的商品状态 %s 无效", product_import.ord_num, product_import.status))
		}

		product.StatusId = status_id

		if product_import.length == 0 || product_import.width == 0 || product_import.height == 0 || product_import.weight == 0 {
			has_zero_data = true
		}

		if has == true {
			//允许更新商品，这样导入商品的属性有变时会方便好多

			// panic(fmt.Sprintf("序号 %s 的商品SKU %s 已存在", product_import.ord_num, product_import.sku))

			// 编号 SKU 中文名称 采购价（人民币元） 产品长度（cm） 产品宽度（cm） 产品高度（cm） 实际重量(kg)
			update_col_list := []string{"sku", "cn_name", "cost_price", "length", "width", "height", "weight"} //这里不更新英文名称
			_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
			core.HandleError(err)
		} else {
			product.BatchCode = import_batch_code
			_, err = db_session.Insert(&product)
			core.HandleError(err)
		}
	}

	message := "数据导入成功！"
	if has_zero_data == true {
		message = "数据导入成功, 部分长宽高数据为零！"
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": message,
	})
}

func (c ErpProductImport) BatchDelete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_IMPORT, c.Session)

	var db_session = c.DbSession

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	/*
DELETE FROM erp_product WHERE batch_code = 0;
	 */

	var sql_delete_product = `DELETE FROM erp_product WHERE batch_code = %d`
	sql_delete_product = fmt.Sprintf(sql_delete_product, batch_code)
	_, err := db_session.Exec(sql_delete_product)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "删除成功！",
	})
}

func (c ErpProductImport) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_IMPORT, c.Session)

	var db_session = c.DbSession

	var product_id_list = make([]int64, 0)
	c.Params.Bind(&product_id_list, "id_list")

	_, err := db_session.In("id", product_id_list).Delete(new(kcgl_models.ProductInfo))
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "删除成功！",
	})
}

/*
func (c ErpProductImport) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_IMPORT, c.Session)

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

			if len(row.Cells) < 4 {
				break
			}

			var product_import models.ProductImportInfo

			//SKU
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
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
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
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
		"success": true,
		"message": "Excel数据上传成功！",
	})
}

func (c ErpProductImport) Confirm() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_IMPORT, c.Session)

	db_session := c.DbSession

	proc_sql := `
CREATE PROCEDURE erp_product_import_trans_v0001(p_user_id INT, p_product_status_id INT)
    BEGIN
        INSERT INTO erp_product (sku, name, cost_price, weight, status_id, create_time, update_time, version)
            SELECT pi.sku, pi.name, pi.cost_price, pi.weight, p_product_status_id, pi.import_time,
                pi.import_time, 1
            FROM erp_product_import pi
            WHERE pi.user_id = p_user_id;
        DELETE FROM erp_product_import WHERE user_id = p_user_id;
    END;
	`

	core.CreateProc("erp_product_import_trans_v0001", proc_sql, db_session)

	login_user := core.GetLoginUser(c.Session)
	login_user_id := login_user.UserId
	result, err := db_session.Exec("CALL erp_product_import_trans_v0001(?, ?);", login_user_id, models.PRODUCT_STATUS_ZC)
	core.HandleError(err)

	rows_affected, _ := result.RowsAffected()
	if rows_affected == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "没有商品导入数据!"})
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: "商品数据导入成功!"})
}

func (c ErpProductImport) Clear() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_PRODUCT_IMPORT, c.Session)

	db_session := c.DbSession

	login_user := core.GetLoginUser(c.Session)
	login_user_id := login_user.UserId

	_, err := db_session.
	Where(fmt.Sprintf("user_id = %d", login_user_id)).
		Delete(new(models.ProductImportInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "临时导入数据清空成功!"})
}
*/
