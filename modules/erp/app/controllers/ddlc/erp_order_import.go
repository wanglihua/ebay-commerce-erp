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
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/erp/service/import_batch_code_service"
	"matrix/modules/erp/service/order_import_batch_service"
	"matrix/modules/erp/service/order_logistics_service"
	"matrix/modules/erp/service/order_shipping_address_service"
	"matrix/modules/erp/service/order_history_service"
	"time"
)

type ErpOrderImport struct {
	*revel.Controller
}

func (c ErpOrderImport) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_IMPORT, c.Session)

	return c.RenderTemplate("erp/ddlc/order_import/order_import_index.html")
}

func (c ErpOrderImport) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_IMPORT, c.Session)

	var db_session = c.DbSession

	// 获得导入的批次列表
	var batch_code_list = order_import_batch_service.Get_order_import_batch_code_list(0, db_session)

	return c.RenderJson(map[string]interface{}{
		"batch_code":      0,
		"batch_code_list": batch_code_list,
	})
}

func (c ErpOrderImport) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_IMPORT, c.Session)

	db_session := c.DbSession

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var order_code_query string
	c.Params.Bind(&order_code_query, "order_code_query")
	order_code_query = strings.TrimSpace(order_code_query)

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	var buyer_name_query string
	c.Params.Bind(&buyer_name_query, "buyer_name_query")
	buyer_name_query = strings.TrimSpace(buyer_name_query)

	var buyer_email_query string
	c.Params.Bind(&buyer_email_query, "buyer_email_query")
	buyer_email_query = strings.TrimSpace(buyer_email_query)

	type order_import_view_item_t struct {
		models.OrderInfo              `xorm:"extends" json:"o"`
		models.OrderTypeInfo          `xorm:"extends" json:"ot"`
		models.OrderHoldTypeInfo      `xorm:"extends" json:"ht"`
		models.ECommercePlatformInfo  `xorm:"extends" json:"pf"`
		models.ShopAccountInfo        `xorm:"extends" json:"sa"`
		models.StatusInfo             `xorm:"extends" json:"s"`
	}

	/*
SELECT o.*, ot.*, ht.*, pf.*, sa.*, s.*
FROM erp_order o
    INNER JOIN erp_order_type ot ON ot.id = o.type_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_status s ON s.id = o.status_id
WHERE o.batch_code = 123 AND o.code LIKE '%abc%' AND o.buyer_user_id LIKE '%abc%' AND o.buyer_email LIKE '%abc%'
      AND o.id IN (SELECT oi.order_id FROM erp_order_item oi INNER JOIN erp_product p ON p.id = oi.item_product_id WHERE p.sku LIKE '%abc%');
	 */

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var sb = core.NewSqlBuilder()
	sb.Select("o.*, ot.*, ht.*, pf.*, sa.*, s.*")
	sb.Count("count(o.id) as order_count")
	sb.From(`erp_order o
    INNER JOIN erp_order_type ot ON ot.id = o.type_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_status s ON s.id = o.status_id`)

	filter = core.SqlAnd(filter, "o.batch_code IS NOT NULL AND o.batch_code <> 0")

	if batch_code != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.batch_code = %d", batch_code))
	}

	if order_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.code LIKE '%%%s%%'", order_code_query))
	}

	if buyer_name_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_user_id LIKE '%%%s%%'", buyer_name_query))
	}

	if buyer_email_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.buyer_email LIKE '%%%s%%'", buyer_email_query))
	}

	if product_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("o.id IN (SELECT oi.order_id FROM erp_order_item oi INNER JOIN erp_product p ON p.id = oi.item_product_id WHERE p.sku LIKE '%%%s%%')", product_code_query))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("o.id ASC")
	}

	sb.Limit(offset, limit)

	var order_import_list = make([]order_import_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&order_import_list)
	core.HandleError(err)

	var order_count struct {
		Count int64 `xorm:"order_count"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&order_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_import_list,
		Total: order_count.Count,
	})
}

// 上传 Excel 导入订单
func (c ErpOrderImport) Upload() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_IMPORT, c.Session)

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
		Yzbm         string  `json:"yzbm"`   // 0 邮政编码  address
		Yjdz         string  `json:"yjdz"`   // 1 邮寄地址  address
		Spbh         string  `json:"spbh"`   // 2 商品编号  order_item
		Spsl         int     `json:"spsl"`   // 3 商品数量  order_item
		Dpm          string  `json:"dpm"`    // 4 店铺名    order
		Ddbh         string  `json:"ddbh"`   // 5 订单编号  order
		Jydh         string  `json:"jydh"`   // 6 交易单号  order
		Skje         float64 `json:"skje"`   // 7 收款金额(澳币) order
		Yfsr         float64 `json:"yfsr"`   // 8 运费收入(澳币)  order
		Khzh         string  `json:"khzh"`   // 10 客户账号 order
		Khxm         string  `json:"khxm"`   // 11 客户姓名 address
		Khyx         string  `json:"khyx"`   // 12 客户邮箱 order
		Dh           string  `json:"dh"`     // 13 电话     address
		Guojia       string  `json:"guojia"` // 14 国家  address
		Ssdq         string  `json:"ssdq"`   // 15 所属地区  address
		Sscs         string  `json:"sscs"`   // 16 所属城市  address
		ErrorMessage string  `json:"error_message"`
		shop_id      int64
		shop_code    string
		product_id   int64
		price_rmb    float64
		price_au     float64
		cost_price   float64
		gross_profit float64
	}

	order_import_list := make([]order_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow

		revel.TRACE.Println("sheet.MaxRow", sheet.MaxRow)

		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			var cell *xlsx.Cell
			var cell_val string
			var err error

			// 邮政编码 邮寄地址 商品编号 商品数量 店铺名 订单编号 交易单号 收款金额(澳币) 运费收入(澳币) 客户账号 客户姓名 客户邮箱 电话 所属地区 所属城市
			var order_import order_import_t

			col_index := 0

			// 邮政编码
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Yzbm = cell_val
			col_index++

			// 邮寄地址
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Yjdz = cell_val
			col_index++

			// 商品编号
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Spbh = cell_val
			col_index++

			// 商品数量
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			spsl, spsl_err := strconv.ParseInt(cell_val, 10, 32)
			order_import.Spsl = int(spsl)
			col_index++

			// 店铺名
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Dpm = cell_val
			col_index++

			// 订单编号
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break; //如果订单编号为空了，就认为是数据结束了
			}
			order_import.Ddbh = cell_val
			col_index++

			// 商品数量 的错误提示，放在订单编号后面
			if spsl == 0 || spsl_err != nil {
				panic(fmt.Sprintf("订单编号 %s, 商品数量为零或者格式不对！", order_import.Ddbh))
			}

			// 交易单号
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Jydh = cell_val
			col_index++

			// 收款金额(澳币)
			cell = sheet.Cell(row_index, col_index)
			cell_val, _ = cell.String()
			cell_val = strings.TrimSpace(cell_val)
			if cell_val != "" {
				order_import.Skje, err = strconv.ParseFloat(cell_val, 64)
				if err != nil {
					panic(fmt.Sprintf("订单编号 %s， 收款金额(澳币)格式不对！", order_import.Ddbh))
				}
			}
			col_index++

			// 运费收入(澳币)
			cell = sheet.Cell(row_index, col_index)
			cell_val, _ = cell.String()
			cell_val = strings.TrimSpace(cell_val)
			if cell_val != "" {
				order_import.Yfsr, err = strconv.ParseFloat(cell_val, 64)
				if err != nil {
					panic(fmt.Sprintf("订单编号 %s， 运费收入(澳币)格式不对！", order_import.Ddbh))
				}
			}
			col_index++

			// 客户账号
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Khzh = cell_val
			col_index++

			// 客户姓名
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Khxm = cell_val
			col_index++

			// 客户邮箱
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Khyx = cell_val
			col_index++

			// 电话
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Dh = cell_val
			col_index++

			// 国家
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Guojia = cell_val
			col_index++

			// 所属地区
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Ssdq = cell_val
			col_index++

			// 所属城市
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			order_import.Sscs = cell_val
			col_index++

			order_import_list = append(order_import_list, order_import)
		}

		revel.TRACE.Println("sheet.MaxRow", sheet.MaxRow)
	}

	if len(order_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入的订单列表为空！",
		})
	}

	revel.TRACE.Println("len(order_import_list)", len(order_import_list))

	/*
	加载所有 店铺名称账户 和 商品SKU列表 ，方便接下来的验证，加快速度

SELECT p.id AS product_id, p.sku AS product_sku FROM erp_product p;
SELECT sa.id AS shop_id, sa.name AS shop_name FROM erp_shop_account sa;
	 */

	type product_id_and_sku_t struct {
		ProductId  int64  `xorm:"'product_id'"`
		ProductSku string `xorm:"'product_sku'"`
	}

	type shop_id_code_name_t struct {
		ShopId   int64  `xorm:"'shop_id'"`
		ShopCode string `xorm:"'shop_code'"`
		ShopName string `xorm:"'shop_name'"`
	}

	var sql_product = `SELECT p.id AS product_id, p.sku AS product_sku FROM erp_product p`
	var product_id_and_sku_list = make([]product_id_and_sku_t, 0)
	err = db_session.Sql(sql_product).Find(&product_id_and_sku_list)
	core.HandleError(err)

	var sql_shop = `SELECT sa.id AS shop_id, sa.code AS shop_code, sa.name AS shop_name FROM erp_shop_account sa;`
	var shop_id_and_name_list = make([]shop_id_code_name_t, 0)
	err = db_session.Sql(sql_shop).Find(&shop_id_and_name_list)
	core.HandleError(err)

	var is_shop_name_exist = func(shop_name string) (has_return bool, shop_id_return int64, shop_code_return string) {
		for _, shop_id_code_name := range shop_id_and_name_list {
			if strings.TrimSpace(shop_name) == strings.TrimSpace(shop_id_code_name.ShopName) {
				return true, shop_id_code_name.ShopId, shop_id_code_name.ShopCode
			}
		}
		return false, 0, ""
	}

	var is_product_sku_exist = func(product_sku string) (has_return bool, product_id_return int64) {
		for _, product_id_and_sku := range product_id_and_sku_list {
			if strings.TrimSpace(product_sku) == strings.TrimSpace(product_id_and_sku.ProductSku) {
				return true, product_id_and_sku.ProductId
			}
		}
		return false, 0
	}

	var order_import_error_list = make([]order_import_t, 0)

	// 错误检查
	for index, order_import := range order_import_list {
		// 检查导入的订单编号是否存在，订单编号检查放在最前面
		/*
		SELECT count(o.id) AS order_count FROM erp_order o WHERE o.code = '123'
		 */
		var order_count struct {
			Count int64 `xorm:"'order_count'"`
		}
		var sql_order_count = `SELECT count(o.id) AS order_count FROM erp_order o WHERE o.code = '%s'`
		sql_order_count = fmt.Sprintf(sql_order_count, order_import.Ddbh)
		_, err := db_session.Sql(sql_order_count).Get(&order_count)
		core.HandleError(err)
		if order_count.Count != 0 {
			// panic(fmt.Sprintf("订单编号 %s 数据库中已经存在！", order_import.Ddbh))
			order_import.ErrorMessage = fmt.Sprintf("订单编号 %s 数据库中已经存在！", order_import.Ddbh)
			order_import_error_list = append(order_import_error_list, order_import)
		}

		// 店铺是否存在
		has, shop_id, shop_code := is_shop_name_exist(order_import.Dpm)
		if has == false {
			//panic(fmt.Sprintf("导入文件中的店铺名 %s 不存在！", order_import.Dpm))
			order_import.ErrorMessage = fmt.Sprintf("编号为 %s 的订单中的店铺名 %s 不存在！", order_import.Ddbh, order_import.Dpm)
			order_import_error_list = append(order_import_error_list, order_import)
		}

		// 商品SKU是否存在
		has, product_id := is_product_sku_exist(order_import.Spbh)
		if has == false {
			/*
			感觉不自动插入商品，提示商品不存在更合适，因为还有商品导入的功能

			// 如果没有的话，插入商品
			var product models.ProductInfo
			product.Sku = order_import.Spbh
			product.CnName = order_import.Spbh
			product.Name = order_import.Spbh

			_, err := db_session.Insert(&product)
			core.HandleError(err)

			product_id = product.Id

			// 数据库里面新增加的也放到 刚刚预先加载的列表里面去
			var product_id_and_sku product_id_and_sku_t
			product_id_and_sku.ProductId = product_id
			product_id_and_sku.ProductSku = product.Sku

			product_id_and_sku_list = append(product_id_and_sku_list, product_id_and_sku)
			*/

			order_import.ErrorMessage = fmt.Sprintf("编号为 %s 的订单中的商品编号 %s 不存在！", order_import.Ddbh, order_import.Spbh)
			order_import_error_list = append(order_import_error_list, order_import)
		}

		order_import_list[index].shop_id = shop_id
		order_import_list[index].shop_code = shop_code
		order_import_list[index].product_id = product_id
	}

	if len(order_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": order_import_error_list,
		})
	}

	revel.TRACE.Println("order_import_list", len(order_import_list))
	// 合并订单，马帮对于多SKU的订单导出的时候，会根据SKU对订单进行拆解。比如一个订单有3个SKU，那么马帮订单导出的时候，会拆解成3条记录，这个时候需要根据同一订单编号在wollongzone对订单进行合并。
	// 马帮导出的拆解订单，每一行的收款都是总金额，所以这里金额不用累加
	var order_import_list_map = make(map[string][]order_import_t)

	for _, order_import := range order_import_list {
		order_import_same_code_list, found := order_import_list_map[order_import.Ddbh]
		if found == false {
			order_import_same_code_list = make([]order_import_t, 0)
		}
		order_import_same_code_list = append(order_import_same_code_list, order_import)
		order_import_list_map[order_import.Ddbh] = order_import_same_code_list
	}

	revel.TRACE.Println("order_import_list_map", len(order_import_list_map))

	var aud_exchange_rate_setting = exchange_rate_service.Get_aud_exchange_rate(db_session)

	// 计算售价，成本价，毛利
	for order_code := range order_import_list_map {
		var order_import_same_code_list = order_import_list_map[order_code]
		_ = order_import_same_code_list

		var cost_price_total float64 = 0
		for index, order_import := range order_import_same_code_list {
			var sql = `SELECT nvt.cost_price AS cost_price
FROM erp_inventory nvt
    INNER JOIN erp_stock st ON st.id = nvt.stock_id
WHERE st.stock_type_id = %d AND nvt.product_id = %d LIMIT 1;`
			sql = fmt.Sprintf(sql, models.STOCK_TYPE_CN, order_import.product_id)

			var product_cost_price struct {
				CostPrice float64 `xorm:"'cost_price'"`
			}
			_, err := db_session.Sql(sql).Get(&product_cost_price)
			core.HandleError(err)

			if product_cost_price.CostPrice != 0 {
				order_import_same_code_list[index].cost_price = product_cost_price.CostPrice
			} else {
				order_import_same_code_list[index].cost_price = 0.00001 // 是零的时候，随便给个数
			}

			cost_price_total = cost_price_total + order_import_same_code_list[index].cost_price*(float64(order_import.Spsl))
		}

		// 订单导入的时候，进行价格和利润的分配。
		// 现在是按照中国仓的成本价进行的分配和计算，不考虑物流等的费用；匹配物流后，会按照澳仓的成本进行分配和计算，考虑物流等的费用。
		for index, order_import := range order_import_same_code_list {
			order_import_same_code_list[index].price_au = ((order_import.cost_price * (float64(order_import.Spsl)) / cost_price_total) * order_import.Skje) / float64(order_import.Spsl)
			order_import_same_code_list[index].price_rmb = order_import_same_code_list[index].price_au * aud_exchange_rate_setting
			order_import_same_code_list[index].gross_profit = (order_import_same_code_list[index].price_rmb - order_import_same_code_list[index].cost_price) * float64(order_import.Spsl)
		}
	}

	// 例如 20170414212812
	var login_user = core.GetLoginUser(c.Session)
	var import_batch_code int64 = import_batch_code_service.Get_import_batch_code(login_user.LoginName)

	var order_code_count = 0
	for order_code := range order_import_list_map {
		order_code_count++
		revel.TRACE.Println("order_code_count", order_code_count)

		var order_import_same_code_list = order_import_list_map[order_code]

		// 订单
		var order models.OrderInfo

		if len(order_import_same_code_list) >= 1 {
			var order_import = order_import_same_code_list[0]

			order.ShopId = order_import.shop_id
			order.Code = order_import.Ddbh
			order.TranId = order_import.Jydh
			order.FullTranId = order_import.shop_code + order_import.Jydh

			var amount_au float64 = order_import.Skje
			var amount_currency = "AUD"                            // 澳币的也设置一下
			var amount_rmb = amount_au * aud_exchange_rate_setting // 汇率换算

			order.AmountPaidRmb = amount_rmb
			order.AmountPaidCurrency = amount_currency
			order.AmountPaid = amount_au

			order.TotalAmountRmb = amount_rmb
			order.TotalAmountCurrency = amount_currency
			order.TotalAmount = amount_au

			order.AdjustmentAmountRmb = amount_rmb
			order.AdjustmentAmountCurrency = amount_currency
			order.AdjustmentAmount = amount_au

			order.ExpressFee = order_import.Yfsr
			order.BuyerUserId = order_import.Khzh
			order.BuyerEmail = order_import.Khyx
			order.TypeId = models.ORDER_TYPE_DRDD // 导入订单
			order.StatusId = models.ORDER_STATUS_YZF
			order.BatchCode = import_batch_code // 导入批次

			order.ExchangeRate = aud_exchange_rate_setting

			order.PaidTime = core.MakeNullLongTime(time.Now(), true)

			_, err = db_session.Insert(&order)
			core.HandleError(err)

			var login_user = core.GetLoginUser(c.Session)
			order_history_service.Insert_order_history(order, login_user.UserId, fmt.Sprintf("订单导入新增订单，工号：%s", login_user.LoginName), db_session)

			// 地址
			var order_shipping_address models.OrderShippingAddressInfo
			order_shipping_address.OrderId = order.Id
			order_shipping_address.Name = order_import.Khxm
			order_shipping_address.PostalCode = order_import.Yzbm
			order_shipping_address.Street1 = order_import.Yjdz
			order_shipping_address.Phone = order_shipping_address_service.Normalize_phone(order_import.Dh) // 整理电话号码 "1234 5666"
			order_shipping_address.CountryName = order_import.Guojia
			order_shipping_address.Province = order_import.Ssdq
			order_shipping_address.CityName = order_import.Sscs
			order_shipping_address.AddressOwner = "eBay"
			_, err = db_session.Insert(&order_shipping_address)
			core.HandleError(err)

			// 回填 order address id
			order.AddressId = order_shipping_address.Id
			_, err = db_session.Id(order.Id).Cols("address_id").Update(&order)
			core.HandleError(err)

			// 订单详细
			var order_item models.OrderItemInfo
			order_item.OrderId = order.Id
			order_item.ItemProductId = order_import.product_id
			order_item.ItemSku = order_import.Spbh
			order_item.Quantity = order_import.Spsl
			order_item.TransactionPriceCurrency = "AUD"

			/*

			if order_item.Quantity != 0 {
				order_item.TransactionPrice = amount_au / float64(order_item.Quantity)
				order_item.TransactionPriceRmb = amount_rmb / float64(order_item.Quantity)
			}

			*/
			order_item.TransactionPrice = order_import.price_au
			order_item.TransactionPriceRmb = order_import.price_rmb
			order_item.CostPrice = order_import.cost_price
			order_item.GrossProfit = order_import.gross_profit

			//order_item.ProductCostPrice = order_import.Cpcb //  这是单个的产品成本， 后来改为成本从库存里面导入，这里就不导入了
			_, err = db_session.Insert(&order_item)
			core.HandleError(err)
		}

		if len(order_import_same_code_list) > 1 {
			for _, order_import := range order_import_same_code_list[1:] {
				// 由于马帮导出 BUG 的存在，这里的金额不用累加
				/*
				var item_amount_au = order_import.Skje
				// var item_amount_rmb = item_amount_au * aud_exchange_rate_setting

				var amount_au float64 = order.AmountPaid + item_amount_au // 这里累加
				var amount_rmb = amount_au * aud_exchange_rate_setting    // 汇率换算

				order.AmountPaidRmb = amount_rmb
				order.AmountPaid = amount_au

				order.TotalAmountRmb = amount_rmb
				order.TotalAmount = amount_au

				order.AdjustmentAmountRmb = amount_rmb
				order.AdjustmentAmount = amount_au

				order.ExpressFee = order.ExpressFee + order_import.Yfsr
				*/

				// 订单详细
				var order_item models.OrderItemInfo
				order_item.OrderId = order.Id
				order_item.ItemProductId = order_import.product_id
				order_item.ItemSku = order_import.Spbh
				order_item.Quantity = order_import.Spsl
				order_item.TransactionPriceCurrency = "AUD"

				/*

				if order_item.Quantity != 0 {
					order_item.TransactionPrice = item_amount_au / float64(order_item.Quantity)
					order_item.TransactionPriceRmb = item_amount_rmb  / float64(order_item.Quantity)
				}

				*/
				order_item.TransactionPrice = order_import.price_au
				order_item.TransactionPriceRmb = order_import.price_rmb
				order_item.CostPrice = order_import.cost_price
				order_item.GrossProfit = order_import.gross_profit

				_, err = db_session.Insert(&order_item)
				core.HandleError(err)
			}

			/*
			_, err = db_session.Id(order.Id).Cols(
				"amount_paid_rmb", "amount_paid",
				"total_amount_rmb", "total_amount",
				"adjustment_amount_rmb", "adjustment_amount",
				"express_fee").Update(&order)
			core.HandleError(err)
			*/
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpOrderImport) BatchDelete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_IMPORT, c.Session)

	var db_session = c.DbSession

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	/*
DELETE o.*, ol.*, oi.*, osa.*
FROM erp_order o
    LEFT JOIN erp_order_logistics ol ON ol.order_id = o.id
    INNER JOIN erp_order_item oi ON oi.order_id = o.id
    INNER JOIN erp_order_shipping_address osa ON osa.order_id = o.id
WHERE o.batch_code = 1;
	 */

	// 删除 采用多表一次性删除
	var sql_delete_orders = `DELETE o.*, ol.*, oi.*, osa.*
FROM erp_order o
    LEFT JOIN erp_order_logistics ol ON ol.order_id = o.id
    INNER JOIN erp_order_item oi ON oi.order_id = o.id
    INNER JOIN erp_order_shipping_address osa ON osa.order_id = o.id
WHERE o.batch_code = %d;`
	sql_delete_orders = fmt.Sprintf(sql_delete_orders, batch_code)
	_, err := db_session.Exec(sql_delete_orders)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "删除成功！",
	})
}

func (c ErpOrderImport) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_IMPORT, c.Session)

	var db_session = c.DbSession

	var order_id_list = make([]int64, 0)
	c.Params.Bind(&order_id_list, "id_list")

	if len(order_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择要删除的订单！",
		})
	}

	var order_id_str_list = core.Int_list_to_str_list(order_id_list)

	// 判断是不是导入的订单，如果不是导入的订单就不能删除
	var order_import_count struct {
		Count int64  `xorm:"'order_import_count'"`
	}
	/*
	SELECT count(o.batch_code) AS order_import_count FROM erp_order o WHERE o.id IN (2, 3, 4, 490, 491, 492);
	*/
	var sql_order_import_count = `SELECT count(o.batch_code) AS order_import_count FROM erp_order o WHERE o.batch_code IS NOT NULL AND o.batch_code <> 0 AND o.id IN (%s);`
	sql_order_import_count = fmt.Sprintf(sql_order_import_count, strings.Join(order_id_str_list, ","))
	_, err := db_session.Sql(sql_order_import_count).Get(&order_import_count)
	core.HandleError(err)

	// 有 order 不是在导入批次里面的
	if order_import_count.Count < int64(len(order_id_list)) {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "只能删除导入的订单！",
		})
	}

	// 删除 订单物流信息
	_, err = db_session.In("order_id", order_id_list).Delete(new(models.OrderLogisticsInfo))
	core.HandleError(err)

	// 删除 订单明细信息
	_, err = db_session.In("order_id", order_id_list).Delete(new(models.OrderItemInfo))
	core.HandleError(err)

	// 删除 订单收货地址信息
	_, err = db_session.In("order_id", order_id_list).Delete(new(models.OrderShippingAddressInfo))
	core.HandleError(err)

	// 删除 订单信息
	_, err = db_session.In("id", order_id_list).Delete(new(models.OrderInfo))
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "删除成功！",
	})
}

func (c ErpOrderImport) MatchLogistics() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDLC_ORDER_IMPORT, c.Session)

	var db_session = c.DbSession

	var batch_code int64
	c.Params.Bind(&batch_code, "batch_code")

	var p_order_logistics = new(order_logistics_service.OrderLogisticsService)
	var login_user = core.GetLoginUser(c.Session)
	p_order_logistics.LoginUser = login_user
	err := p_order_logistics.CalculateListImportBatch(batch_code, db_session)
	// core.HandleError(err)
	if err != nil {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "批次匹配物流完成！",
	})
}
