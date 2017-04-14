package controllers

import (
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrderImport struct {
	*revel.Controller
}

func (c ErpOrderImport) Index() revel.Result {
	return c.RenderTemplate("erp/ddlc/order_import/order_import_index.html")
}

func (c ErpOrderImport) ListData() revel.Result {
	return c.RenderEmpty()
}

// 上传 Excel 导入订单
func (c ErpOrderImport) Upload() revel.Result {
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

	type order_import_t struct {
		Yzbm string  // 0 邮政编码
		Yjdz string  // 1 邮寄地址
		Spbh string  // 2 商品编号
		Spsl int64   // 3 商品数量
		Dpm  string  // 4 店铺名
		Ddbh string  // 5 订单编号
		Jydh string  // 6 交易单号
		Skje float64 // 7 收款金额(RMB)
		Yfsr float64 // 8 运费收入
		Cpcb float64 // 9 产品成本
		Khzh string  // 10 客户账号
		Khxm string  // 11 客户姓名
		Khyx string  // 12 客户邮箱
		Dh   string  // 13 电话
		Ssdq string  // 14 所属地区
		Sscs string  // 15 所属城市
	}

	order_import_list := make([]order_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
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

			// 邮政编码 邮寄地址 商品编号 商品数量 店铺名 订单编号 交易单号 收款金额(RMB) 运费收入 产品成本 客户账号 客户姓名 客户邮箱 电话 所属地区 所属城市
			var order_import order_import_t

			col_index := 0

			if len(row.Cells) == 0 {
				break
			}

			// 邮政编码
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			order_import.Yzbm = cell_val
			col_index++

			// 邮寄地址
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Yjdz = cell_val
			col_index++

			// 商品编号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Spbh = cell_val
			col_index++

			// 商品数量
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Spsl, _ = strconv.ParseInt(cell_val, 10, 64)
			col_index++

			// 店铺名
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Dpm = cell_val
			col_index++

			// 订单编号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Ddbh = cell_val
			col_index++

			// 交易单号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Jydh = cell_val
			col_index++

			// 收款金额(RMB)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Skje, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 运费收入
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Yfsr, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 产品成本
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Cpcb, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 客户账号
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Khzh = cell_val
			col_index++

			// 客户姓名
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Khxm = cell_val
			col_index++

			// 客户邮箱
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Khyx = cell_val
			col_index++

			// 电话
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Dh = cell_val
			col_index++

			// 所属地区
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Ssdq = cell_val
			col_index++

			// 所属城市
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Sscs = cell_val
			col_index++

			order_import_list = append(order_import_list, order_import)
		}
	}

	if len(order_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入的订单列表为空！",
		})
	}

	/*
	加载所有 店铺名称账户 和 商品SKU列表 ，方便接下来的验证，加快速度

SELECT p.id AS product_id, p.sku AS product_sku FROM erp_product p;
SELECT sa.id AS shop_id, sa.name AS shop_name FROM erp_shop_account sa;
	 */

	type product_id_and_sku_t struct {
		ProductId  int64  `xorm:"'product_id'"`
		ProductSku string `xorm:"'product_sku'"`
	}

	type shop_id_and_name_t struct {
		ShopId   int64  `xorm:"'shop_id'"`
		ShopName string `xorm:"'shop_name'"`
	}

	var sql_product = `SELECT p.id AS product_id, p.sku AS product_sku FROM erp_product p`
	var product_id_and_sku_list = make([]product_id_and_sku_t, 0)
	err = db_session.Sql(sql_product).Find(&product_id_and_sku_list)
	core.HandleError(err)

	var sql_shop = `SELECT sa.id AS shop_id, sa.name AS shop_name FROM erp_shop_account sa`
	var shop_id_and_name_list = make([]shop_id_and_name_t, 0)
	err = db_session.Sql(sql_shop).Find(&shop_id_and_name_list)
	core.HandleError(err)

	var is_shop_name_exist = func(shop_name string) (has_return bool, shop_id_return int64) {
		for _, shop_id_and_name := range shop_id_and_name_list {
			if strings.TrimSpace(shop_name) == strings.TrimSpace(shop_id_and_name.ShopName) {
				return true, shop_id_and_name.ShopId
			}
		}
		return false, 0
	}

	var is_product_sku_exist = func(product_sku string) (has_return bool, product_id_return int64) {
		for _, product_id_and_sku := range product_id_and_sku_list {
			if strings.TrimSpace(product_sku) == strings.TrimSpace(product_id_and_sku.ProductSku) {
				return true, product_id_and_sku.ProductId
			}
		}
		return false, 0
	}

	for _, order_import := range order_import_list {
		// 店铺是否存在
		has, shop_id := is_shop_name_exist(order_import.Dpm)
		if has == false {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "导入文件中的店铺名不存在！",
			})
		}

		// 商品SKU是否存在
		has, product_id := is_product_sku_exist(order_import.Spbh)
		if has == false {
			// 如果没有的话，插入商品
			var product models.ProductInfo
			product.Sku = order_import.Spbh
			product.CnName = order_import.Spbh

			_, err := db_session.Insert(&product)
			core.HandleError(err)

			product_id = product.Id
		}

		_ = product_id

		var order models.OrderInfo

		order.ShopId = shop_id

		_, err = db_session.Insert(&order)
		core.HandleError(err)

	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}
