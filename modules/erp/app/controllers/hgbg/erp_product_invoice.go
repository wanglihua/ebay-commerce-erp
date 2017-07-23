package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpProductInvoice struct {
	*revel.Controller
}

func (c ErpProductInvoice) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	return c.RenderTemplate("erp/hgbg/product_invoice/product_invoice_index.html")
}

func (c ErpProductInvoice) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	db_session := c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type product_weight_view_item_t struct {
		Product models.ProductInfo `xorm:"extends" json:"p"`
		Status  models.StatusInfo  `xorm:"extends" json:"s"`
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

	var product_list = make([]product_weight_view_item_t, 0, limit)
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

type product_weight_detail_form_t struct {
	Product models.ProductInfo `json:"product"`
}

func (c ErpProductInvoice) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	db_session := c.DbSession

	var product_id int64
	c.Params.Bind(&product_id, "id")

	var product models.ProductInfo
	_, err := db_session.Id(product_id).Get(&product)
	core.HandleError(err)

	var detail_form product_weight_detail_form_t
	detail_form.Product = product

	return c.RenderJson(detail_form)
}

func (c ErpProductInvoice) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

	db_session := c.DbSession

	var detail_form product_weight_detail_form_t
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	var product_id = detail_form.Product.Id

	var product_in_db models.ProductInfo
	has, err := db_session.Id(product_id).Get(&product_in_db)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品不存在！",
		})
	}

	product_in_db.NeedInvoice = detail_form.Product.NeedInvoice

	updata_col_List := []string{"need_invoice"}
	_, err = db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}

func (c ErpProductInvoice) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_PRODUCT_INVOICE, c.Session)

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

	type product_invoice_t struct {
		sku          string
		need_invoice int64
	}

	product_invoice_list := make([]product_invoice_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow

		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_weight product_invoice_t

			col_index := 0 //从第二列开始导入，略去第一列的编号

			//SKU
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}

			//revel.TRACE.Println(cell_val)
			product_weight.sku = cell_val
			col_index++

			// length
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)

			if strings.TrimSpace(cell_val) == "是" {
				product_weight.need_invoice = 1
			} else {
				product_weight.need_invoice = 0
			}

			col_index++

			product_invoice_list = append(product_invoice_list, product_weight)
		}
	}

	if len(product_invoice_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	for _, product_invoice := range product_invoice_list {
		var product models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_invoice.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_invoice.sku
		product.NeedInvoice = product_invoice.need_invoice

		if has == true {
			_, err = db_session.Id(product.Id).Cols("need_invoice").Update(&product)
			core.HandleError(err)
		} else {
			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_invoice.sku))
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}
