package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/kcgl_models"
)

type ErpProductPostageCalWeight struct {
	*revel.Controller
}

func (c ErpProductPostageCalWeight) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_SALE_SETTING, c.Session)

	db_session := c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_postage_cal_weight_view_item_t struct {
		Product kcgl_models.ProductInfo `xorm:"extends" json:"p"`
		Status  models.StatusInfo       `xorm:"extends" json:"s"`
	}

	/*
		SELECT p.*, s.*
		FROM erp_product p
		    INNER JOIN erp_status s ON s.id = p.status_id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("p.*, s.*")
	sb.Count("count(p.id) AS product_count")
	sb.From(`erp_product p
		    INNER JOIN erp_status s ON s.id = p.status_id`)

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("p.id ASC")
	}

	sb.Limit(offset, limit)

	var product_list = make([]product_postage_cal_weight_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&product_list)
	core.HandleError(err)

	type product_count_t struct {
		Count int64 `xorm:"'product_count'"`
	}
	var product_count product_count_t
	_, err = db_session.Sql(sb.SqlCount()).Get(&product_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_list,
		Total: product_count.Count,
	})
}

type product_postage_cal_weight_detail_form_t struct {
	Product kcgl_models.ProductInfo `json:"product"`
}

func (c ErpProductPostageCalWeight) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_SALE_SETTING, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product kcgl_models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var detail_form product_postage_cal_weight_detail_form_t
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProductPostageCalWeight) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_SALE_SETTING, c.Session)

	db_session := c.DbSession

	var detail_form product_postage_cal_weight_detail_form_t
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	var product_id = detail_form.Product.Id

	var product_in_db kcgl_models.ProductInfo
	has, err := db_session.Id(product_id).Get(&product_in_db)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品不存在！",
		})
	}

	// var product_in_db_befor_modify = product_in_db

	product_in_db.PostageCalWeight = detail_form.Product.PostageCalWeight

	_, err = db_session.Id(product_id).Cols("postage_cal_weight").Update(&product_in_db)
	core.HandleError(err)

	// product_log_service.Insert_product_log(&detail_form.Product, &product_in_db_befor_modify, core.Db_operation_types.Update, db_session, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}

func (c ErpProductPostageCalWeight) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XSGL_PRODUCT_SALE_SETTING, c.Session)

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

	type product_postage_cal_weight_t struct {
		sku                string
		postage_cal_weight float64
	}

	product_postage_cal_weight_list := make([]product_postage_cal_weight_t, 0)

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
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_postage_cal_weight product_postage_cal_weight_t

			col_index := 0 //从第二列开始导入，略去第一列的编号

			//SKU
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			//revel.TRACE.Println(cell_val)
			product_postage_cal_weight.sku = cell_val
			col_index++

			// length
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			if strings.TrimSpace(cell_val) == "" {
				product_postage_cal_weight.postage_cal_weight = 0
			} else {
				postage_cal_weight, err := strconv.ParseFloat(cell_val, 64)
				if err != nil {
					panic("运费权重录入格式有错！")
				}
				product_postage_cal_weight.postage_cal_weight = postage_cal_weight
			}
			col_index++

			product_postage_cal_weight_list = append(product_postage_cal_weight_list, product_postage_cal_weight)
		}
	}

	if len(product_postage_cal_weight_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	for _, product_weight := range product_postage_cal_weight_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_weight.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_weight.sku
		product.PostageCalWeight = product_weight.postage_cal_weight

		if has == true {
			_, err = db_session.Id(product.Id).Cols("postage_cal_weight").Update(&product)
			core.HandleError(err)
		} else {
			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_weight.sku))
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}
