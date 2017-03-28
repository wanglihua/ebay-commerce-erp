package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductWeight struct {
	*revel.Controller
}

func (c ErpProductWeight) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_weight/product_weight_index.html")
}

func (c ErpProductWeight) ListData() revel.Result {
	db_session := c.DbSession
	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var filter_type int64
	c.Params.Bind(&filter_type, "filter_type")

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

	if filter_type == 1 { //有零值数据的记录
		filter_has_zero_data := core.SqlOr("length = 0", "width = 0", "height = 0", "weight = 0")
		filter = core.SqlAnd(filter, filter_has_zero_data)
	}

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

func (c ErpProductWeight) DetailData() revel.Result {
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

func (c ErpProductWeight) Save() revel.Result {
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

	product_in_db.Length = detail_form.Product.Length
	product_in_db.Width = detail_form.Product.Width
	product_in_db.Height = detail_form.Product.Height
	product_in_db.Weight = detail_form.Product.Weight

	updata_col_List := []string{"length", "width", "height", "weight"}
	affected_count, err := db_session.Id(product_id).Cols(updata_col_List...).Update(&product_in_db)
	core.HandleError(err)

	if affected_count == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "商品长宽高重量数据未能更新！",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "保存成功！",
	})
}

func (c ErpProductWeight) Import() revel.Result {
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

	type product_weight_t struct {
		sku    string
		length float64
		width  float64
		height float64
		weight float64
	}

	product_weight_list := make([]product_weight_t, 0)

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

			var product_weight product_weight_t

			col_index := 1 //从第二列开始导入，略去第一列的编号

			//SKU
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_weight.sku = cell_val
			col_index++

			// length
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_length, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.length = product_length
			col_index++

			// width
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_width, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.width = product_width
			col_index++

			// height
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			product_height, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.height = product_height
			col_index++

			// weight
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			//revel.TRACE.Println(cell_val)
			weight, err := strconv.ParseFloat(cell_val, 64)
			core.HandleError(err)
			product_weight.weight = weight
			col_index++

			product_weight_list = append(product_weight_list, product_weight)
		}
	}

	if len(product_weight_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	has_zero_data := false
	for _, product_weight := range product_weight_list {
		revel.TRACE.Println(product_weight)

		var product models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_weight.sku)).Get(&product)
		core.HandleError(err)

		product.Sku = product_weight.sku
		product.Length = product_weight.length
		product.Width = product_weight.width
		product.Height = product_weight.height
		product.Weight = product_weight.weight

		if product_weight.length == 0 || product_weight.width == 0 || product_weight.height == 0 || product_weight.weight == 0 {
			has_zero_data = true
		}

		if has == true {
			update_col_list := []string{"length", "width", "height", "weight"}
			_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
			core.HandleError(err)
		} else {
			product.CnName = product.Sku
			product.Name = product.Sku
			product.StatusId = models.PRODUCT_STATUS_ZC
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
