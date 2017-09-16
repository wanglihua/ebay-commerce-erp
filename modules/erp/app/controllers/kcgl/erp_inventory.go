package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/product_service"
	"matrix/modules/erp/service/stock_pick_crew_service"
	"matrix/modules/erp/service/stock_service"
	"matrix/modules/erp/service/storage_loc_service"
	"matrix/modules/erp/service/inventory_service"
	"matrix/modules/erp/service/inventory_history_service"
	"matrix/modules/erp/models/kcgl_models"
	"matrix/modules/erp/models/sppz_models"
)

type ErpInventory struct {
	*revel.Controller
}

const (
	PRODUCT_ONE_TO_ONE_STOCK_PRODUCT       = 100
	PRODUCT_ONE_TO_ONE_STORAGE_LOC_PRODUCT = 200
)

func (c ErpInventory) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	var has_sku_loc_res = pass_service.Is_login_user_has_pass(pass_res.KCGL_INVENTORY_SKU_LOC, c.Session)
	//var has_sku_qty_res = pass_service.Is_login_user_has_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY_OVER, c.Session) || pass_service.Is_login_user_has_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY_DIFF, c.Session)

	//var has_sku_qty_res = pass_service.Is_login_user_has_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY, c.Session)
	var has_sku_price_res = pass_service.Is_login_user_has_pass(pass_res.KCGL_INVENTORY_SKU_PRICE, c.Session)

	var has_sku_qty_import_res = pass_service.Is_login_user_has_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY, c.Session)
	var has_sku_qty_export_res = pass_service.Is_login_user_has_pass(pass_res.KCGL_INVENTORY_EXPORT_QTY, c.Session)

	c.RenderArgs["has_sku_loc_res"] = has_sku_loc_res
	//c.RenderArgs["has_sku_qty_res"] = has_sku_qty_res
	c.RenderArgs["has_sku_qty_import_res"] = has_sku_qty_import_res
	c.RenderArgs["has_sku_qty_export_res"] = has_sku_qty_export_res
	c.RenderArgs["has_sku_price_res"] = has_sku_price_res
	c.RenderArgs["product_one_to_one_stock_product"] = PRODUCT_ONE_TO_ONE_STOCK_PRODUCT
	c.RenderArgs["product_one_to_one_storage_loc_product"] = PRODUCT_ONE_TO_ONE_STORAGE_LOC_PRODUCT

	return c.RenderTemplate("erp/kcgl/inventory/inventory_index.html")
}

func (c ErpInventory) PageData() revel.Result {
	var db_session = c.DbSession
	var stock_list = make([]sppz_models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"stock_list":            stock_list,
		"stock_id_query":        0,
		"product_code_query":    "",
		"cost_price_zero_query": 0,
		"product_one2one_query": 0,
	})
}

type InventoryViewItem struct {
	Inventory  kcgl_models.InventoryInfo  `xorm:"extends" json:"ivt"`
	Stock      sppz_models.StockInfo      `xorm:"extends" json:"st"`
	StorageLoc sppz_models.StorageLocInfo `xorm:"extends" json:"loc"`
	Product    kcgl_models.ProductInfo    `xorm:"extends" json:"prd"`
}

func (c ErpInventory) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	var cost_price_zero_query int64
	c.Params.Bind(&cost_price_zero_query, "cost_price_zero_query")

	//product_one2one_query
	var product_one2one_query int64
	c.Params.Bind(&product_one2one_query, "product_one2one_query")

	/*
		SELECT ivt.*, st.*, loc.*, prd.*
		FROM erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	var sb = core.NewSqlBuilder()
	sb.Select(`ivt.*, st.*, loc.*, prd.*`)
	sb.Count(`count(ivt.id) AS ivt_count`)
	sb.From(`erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id`)

	if stock_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("st.id = %d", stock_id_query))
	}

	if product_code_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_code_query))
	}

	var stock_id_str_list = stock_pick_crew_service.Get_login_user_stock_id_str_list(db_session, c.Session)
	if len(stock_id_str_list) != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("st.id IN (%s)", strings.Join(stock_id_str_list, ",")))
	} else {
		filter = core.SqlAnd(filter, "st.id IN (0)")
	}

	if cost_price_zero_query != 0 {
		filter = core.SqlAnd(filter, "ivt.cost_price = 0")
	}

	if product_one2one_query != 0 {

		// 只显示 仓库商品 重复记录
		if product_one2one_query == PRODUCT_ONE_TO_ONE_STOCK_PRODUCT {
			filter = core.SqlAnd(filter, `
ivt.id IN (
    SELECT id
    FROM erp_inventory n
        INNER JOIN (
                       SELECT stock_id, product_id
                       FROM erp_inventory
                       GROUP BY stock_id, product_id
                       HAVING count(*) > 1
                   ) r ON n.stock_id = r.stock_id AND n.product_id = r.product_id)`)
		}

		// 只显示 仓位商品 重复记录
		if product_one2one_query == PRODUCT_ONE_TO_ONE_STORAGE_LOC_PRODUCT {
			filter = core.SqlAnd(filter, fmt.Sprintf(`
ivt.storage_loc_id IN (
	SELECT nvt.storage_loc_id
	FROM erp_inventory nvt
		INNER JOIN erp_storage_loc sl ON nvt.storage_loc_id = sl.id
		INNER JOIN erp_stock s ON sl.stock_id = s.id
	WHERE s.stock_type_id = %d AND nvt.storage_loc_id <> 0
	GROUP BY nvt.storage_loc_id
	HAVING count(nvt.product_id) > 1)`, sppz_models.STOCK_TYPE_AU))
		}
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("ivt.id ASC")
	}

	sb.Limit(offset, limit)

	inventory_list := make([]InventoryViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&inventory_list)
	core.HandleError(err)

	var ivt_count struct {
		Count int64 `xorm:"ivt_count"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&ivt_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_list,
		Total: ivt_count.Count,
	})
}

type InventoryDetailForm struct {
	Inventory      kcgl_models.InventoryInfo    `json:"inventory"`
	StockList      []sppz_models.StockInfo      `json:"stock_list"`
	StorageLocList []sppz_models.StorageLocInfo `json:"storage_loc_list"`
	ProductSku     string                       `json:"product_sku"`
	IsNew          bool                         `json:"is_new"`
}

func (f *InventoryDetailForm) IsCreate() bool {
	return f.Inventory.Id == 0
}

func (f *InventoryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.Inventory.StockId).Message("仓库不能为空！")

	validation.Required(f.Inventory.StorageLocId).Message("库位不能为空！")

	validation.Required(f.Inventory.ProductId).Message("商品不能为空！")

	validation.Required(f.Inventory.Qty).Message("数量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpInventory) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	var inventory_id int64
	c.Params.Bind(&inventory_id, "id")

	var stock_list = make([]sppz_models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	var detail_form InventoryDetailForm
	detail_form.StockList = stock_list

	var inventory kcgl_models.InventoryInfo
	if inventory_id != 0 {
		has, err := db_session.Id(inventory_id).Get(&inventory)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品库存不存在！"})
		}

		var storage_loc_list = make([]sppz_models.StorageLocInfo, 0)
		//err = db_session.Where(fmt.Sprintf("stock_id = %d", inventory.StockId)).Find(&storage_loc_list)
		err = db_session.Id(inventory.StorageLocId).Find(&storage_loc_list) // 修改不是新增时不能改仓位，所以取一个就够了，那列表实在太大了
		core.HandleError(err)

		detail_form.StorageLocList = storage_loc_list

		var product kcgl_models.ProductInfo
		_, err = db_session.Id(inventory.ProductId).Get(&product)
		core.HandleError(err)

		detail_form.ProductSku = product.Sku
	}

	detail_form.Inventory = inventory

	detail_form.IsNew = detail_form.IsCreate()

	return c.RenderJson(detail_form)
}

func (c ErpInventory) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	var detail_form InventoryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	inventory := detail_form.Inventory
	var product_sku = detail_form.ProductSku

	var product_id = product_service.Get_product_id_by_sku(product_sku, db_session)
	if product_id == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "保存失败，商品SKU不存在！",
		})
	}

	inventory.ProductId = product_id

	// 数据库中 stock_id, storage_loc_id, product_id 三者的组合应该是唯一的
	var inventory_in_db kcgl_models.InventoryInfo
	_, err = db_session.
	Where("stock_id = ? AND storage_loc_id = ? AND product_id = ?", inventory.StockId, inventory.StorageLocId, inventory.ProductId).
		Get(&inventory_in_db)
	core.HandleError(err)

	var qty_before = inventory_in_db.Qty

	var stock sppz_models.StockInfo
	_, err = db_session.Id(inventory.StockId).Get(&stock)
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	var affected int64
	if inventory_in_db.Id == 0 { //创建数据库记录
		if stock.StockTypeId == sppz_models.STOCK_TYPE_AU {
			inventory_count, err := db_session.Where("stock_id = ? AND storage_loc_id = ?", inventory.StockId, inventory.StorageLocId).Count(new(kcgl_models.InventoryInfo))
			core.HandleError(err)
			if inventory_count != 0 {
				return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，仓库仓位中SKU库存已存在，一个仓位只能一个SKU！"})
			}

			inventory_count, err = db_session.Where("stock_id = ? AND product_id = ?", inventory.StockId, inventory.ProductId).Count(new(kcgl_models.InventoryInfo))
			core.HandleError(err)
			if inventory_count != 0 {
				return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，仓库中SKU的库存已存在，一个SKU只能一个仓位！"})
			}
		}

		affected, err = db_session.Insert(&inventory)
		core.HandleError(err)

		inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("手工新增库存， 工号：%s", login_user.LoginName), db_session)
	} else { // 更新数据库记录
		inventory.Id = inventory_in_db.Id

		/* 不知道为什么对有些SKU就是不能通过这个检查，而且数据库里面也没发现有重复的，先禁用这个检查再说，因为更新时，已不能更改SKU和仓位，所以也没关系的。
		if stock.StockTypeId == models.STOCK_TYPE_AU {
			inventory_count, err := db_session.Where("id <> ? AND stock_id = ? AND storage_loc_id = ?", inventory.Id, inventory.StockId, inventory.StorageLocId).Count(new(models.InventoryInfo))
			core.HandleError(err)
			if inventory_count != 0 {
				return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，仓库仓位中SKU库存已存在，一个仓位只能一个SKU！" })
			}

			inventory_count, err = db_session.Where("id <> ? AND stock_id = ? AND product_id = ?", inventory.Id, inventory.StockId, inventory.ProductId).Count(new(models.InventoryInfo))
			core.HandleError(err)
			if inventory_count != 0 {
				return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，仓库中SKU的库存已存在，一个SKU只能一个仓位！" })
			}
		}
		*/

		affected, err = db_session.Id(inventory.Id).Update(&inventory)
		core.HandleError(err)

		var qty_after = inventory.Qty

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}

		inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("手工编辑库存 %d => %d，工号：%s", qty_before, qty_after, login_user.LoginName), db_session)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpInventory) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	inventory_id_list := make([]int64, 0)
	c.Params.Bind(&inventory_id_list, "id_list")

	var inventory_list = make([]kcgl_models.InventoryInfo, 0)
	err := db_session.In("id", inventory_id_list).Find(&inventory_list)
	core.HandleError(err)

	affected, err := db_session.In("id", inventory_id_list).Delete(new(kcgl_models.InventoryInfo))
	core.HandleError(err)

	var login_user = core.GetLoginUser(c.Session)

	for _, inventory := range inventory_list {
		inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("手工删除库存 %d => 0, 工号：%s", inventory.Qty, login_user.LoginName), db_session)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpInventory) ReCalPending() revel.Result {

	inventory_service.Calculate_inventory_pending_qty()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "更新待配货数量成功",
	})
}

func (c ErpInventory) UpdateAuInventoryCostPrice() revel.Result {

	go inventory_service.Update_au_stock_cost()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "启动澳仓成本更新后台任务成功！",
	})
}

func (c ErpInventory) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

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

	//序号 商品编号	仓库	仓位	库存

	type inventory_import_t struct {
		OrdNum         string  `json:"ord_num"`          // 序号
		ProductSku     string  `json:"product_sku"`      // 商品编号
		StockName      string  `json:"stock_name"`       // 仓库名称
		StorageLocName string  `json:"storage_loc_name"` // 仓位名称
		Qty            int64   `json:"qty"`              // 库存
		CostPrice      float64 `json:"cost_price"`       // 成本价

		// 下面这些字段是从数据库里面查询后保存在这里面的，这样过程中，就不用二次查询了
		stock_id       int64
		storage_loc_id int64
		product_id     int64
		ErrorMessage   string `json:"error_message"`
	}

	inventory_import_list := make([]inventory_import_t, 0)

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

			var inventory_import inventory_import_t

			col_index := 0

			var data_col_count = 6

			if len(row.Cells) < data_col_count-1 { // 6 列数据， 最后一列可能为空
				break
			}

			// 序号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			inventory_import.OrdNum = cell_val
			col_index++

			// 商品编号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.ProductSku = cell_val
			col_index++

			// 仓库名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StockName = cell_val
			col_index++

			// 仓位名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StorageLocName = cell_val
			col_index++

			// 库存
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			qty, _ := strconv.ParseInt(cell_val, 10, 64)
			inventory_import.Qty = qty
			col_index++

			if len(row.Cells) > data_col_count-1 { // 最后一列有数据
				// 成本价
				// cell = row.Cells[col_index]
				cell = sheet.Cell(row_index, col_index)
				cell_val, err = cell.String()
				core.HandleError(err)
				cell_val = strings.TrimSpace(cell_val)
				cost_price, _ := strconv.ParseFloat(cell_val, 64)
				inventory_import.CostPrice = cost_price
				col_index++
			}

			inventory_import_list = append(inventory_import_list, inventory_import)
		}
	}

	if len(inventory_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品库存列表为空！",
		})
	}

	var stock_id_list = stock_pick_crew_service.Get_login_user_stock_id_list(db_session, c.Session)

	var is_stock_id_can_view = func(stock_id int64) bool {
		for _, stock_id_in_list := range stock_id_list {
			if stock_id_in_list == stock_id {
				return true
			}
		}
		return false
	}

	// 用于检查 仓库和SKU 的重复
	var stock_product_map = make(map[string]string)

	var inventory_import_error_list = make([]inventory_import_t, 0)
	// 先检查一遍是不是有错误，
	for index, inventory_import := range inventory_import_list {

		// 检查 商品 sku
		var product_id int64 = product_service.Get_product_id_by_sku(inventory_import.ProductSku, db_session)

		if product_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的商品SKU %s 不存在！", inventory_import.OrdNum, inventory_import.ProductSku)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		}

		// 检查 仓库名称
		stock_id, type_id := stock_service.Get_stock_id_type_by_name(inventory_import.StockName, db_session)
		if type_id != sppz_models.STOCK_TYPE_CN {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "只能导入中国仓的成本！",
			})
		}

		var storage_loc_id int64

		if stock_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StockName)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		} else { //仓库名称存在
			if is_stock_id_can_view(stock_id) == false {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库 %s 你没有权限查看和修改！", inventory_import.OrdNum, inventory_import.StockName)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			var storage_loc_name = strings.TrimSpace(inventory_import.StorageLocName)
			if storage_loc_name != "" {
				// 如果提供了仓位， 检查仓位名称
				storage_loc_id = storage_loc_service.Get_storage_loc_id_by_name(storage_loc_name, stock_id, db_session)
				if storage_loc_id == 0 {
					inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓位名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StorageLocName)
					inventory_import_error_list = append(inventory_import_error_list, inventory_import)
				}
			}
		}

		// 检查导入的SKU是否有重复的，是 仓库和SKU 的组合的重复检查
		if stock_id != 0 && product_id != 0 {
			var map_key = fmt.Sprintf("%d_%d", stock_id, product_id)
			pre_order_num, found := stock_product_map[map_key]
			if found {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库（%s）和商品（%s）的组合和之前(序号%s)的有重复！", inventory_import.OrdNum, inventory_import.StockName, inventory_import.ProductSku, pre_order_num)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			stock_product_map[map_key] = inventory_import.OrdNum
		}

		inventory_import_list[index].product_id = product_id
		inventory_import_list[index].stock_id = stock_id
		inventory_import_list[index].storage_loc_id = storage_loc_id
	}

	if len(inventory_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": inventory_import_error_list,
		})
	}

	var login_user = core.GetLoginUser(c.Session)

	for _, inventory_import := range inventory_import_list {

		var stock_id = inventory_import.stock_id
		var storage_loc_id = inventory_import.storage_loc_id
		var product_id = inventory_import.product_id

		var inventory kcgl_models.InventoryInfo
		has, err := db_session.Where(fmt.Sprintf("stock_id = %d AND storage_loc_id = %d", stock_id, storage_loc_id)).Get(&inventory)
		core.HandleError(err)

		inventory.StockId = stock_id
		inventory.StorageLocId = storage_loc_id
		inventory.ProductId = product_id
		// inventory.Qty += inventory_import.Qty // 如果是插入的话， 开始时 inventory.Qty
		var qty_before = inventory.Qty
		inventory.Qty = inventory_import.Qty // 这里采用覆盖原来库存的方式
		var qty_after = inventory.Qty
		inventory.CostPrice = inventory_import.CostPrice

		if has == false {
			// 插入
			_, err = db_session.Insert(&inventory)
			core.HandleError(err)

			inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("导入新增库存, 工号：%s", login_user.LoginName), db_session)
		} else {
			// 更新
			_, err := db_session.Id(inventory.Id).Cols("product_id", "qty", "cost_price").Update(&inventory)
			core.HandleError(err)

			inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("导入更新库存 %d => %d, 工号：%s", qty_before, qty_after, login_user.LoginName), db_session)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpInventory) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY, c.Session)

	db_session := c.DbSession

	xlsx_file := c.get_export_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "库存列表" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpInventory) get_export_excel(db_session *xorm.Session) *xlsx.File {
	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	/*
	SELECT ivt.*, st.*, loc.*, prd.*
	FROM erp_inventory ivt
		INNER JOIN erp_stock st ON ivt.stock_id = st.id
		LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	var sb = core.NewSqlBuilder()
	sb.Select("ivt.*, st.*, loc.*, prd.*")
	sb.Count("count(ivt.id) AS ivt_count")
	sb.From(`erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id`)

	//query extra filter here
	var stock_id_str_list = stock_pick_crew_service.Get_login_user_stock_id_str_list(db_session, c.Session)
	var filter = fmt.Sprintf("st.id IN (%s)", strings.Join(stock_id_str_list, ","))

	if len(id_list) == 0 {
		if stock_id_query != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("st.id = %d", stock_id_query))
		}

		if product_code_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_code_query))
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("ivt.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	sb.Where(filter)
	sb.OrderBy("ivt.id")

	inventory_list := make([]InventoryViewItem, 0)
	err := db_session.Sql(sb.SqlList()).Find(&inventory_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("库存列表")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = horizontal_align
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	row_height := 0.6

	// 序号	商品编号	仓库	仓位	库存	成本价（RMB）
	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "序号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 25)
		set_head_cell_style(cell)
		cell.Value = "商品编号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 22)
		set_head_cell_style(cell)
		cell.Value = "仓库"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "仓位"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "库存"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "成本价"
	}

	for index, inventory_view := range inventory_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Product.Sku

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Stock.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.StorageLoc.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(inventory_view.Inventory.Qty)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(inventory_view.Inventory.CostPrice)
	}

	return xlsx_file
}

func (c ErpInventory) ImportProductStorageLoc() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_SKU_LOC, c.Session)

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

	//序号 商品编号	仓库	仓位	库存

	type inventory_import_t struct {
		OrdNum         string  `json:"ord_num"`          // 序号
		ProductSku     string  `json:"product_sku"`      // 商品编号
		StockName      string  `json:"stock_name"`       // 仓库名称
		StorageLocName string  `json:"storage_loc_name"` // 仓位名称
		Qty            int64   `json:"qty"`              // 库存
		CostPrice      float64 `json:"cost_price"`       // 成本价

		// 下面这些字段是从数据库里面查询后保存在这里面的，这样过程中，就不用二次查询了
		stock_id       int64
		stock_type_id  int64
		storage_loc_id int64
		product_id     int64
		ErrorMessage   string `json:"error_message"`
	}

	inventory_import_list := make([]inventory_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var inventory_import inventory_import_t

			col_index := 0

			// 序号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			inventory_import.OrdNum = cell_val
			col_index++

			// 商品编号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.ProductSku = cell_val
			col_index++

			// 仓库名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StockName = cell_val
			col_index++

			// 仓位名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StorageLocName = cell_val
			col_index++

			inventory_import_list = append(inventory_import_list, inventory_import)
		}
	}

	if len(inventory_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品库存列表为空！",
		})
	}

	var stock_id_list = stock_pick_crew_service.Get_login_user_stock_id_list(db_session, c.Session)

	var is_stock_id_can_view = func(stock_id int64) bool {
		for _, stock_id_in_list := range stock_id_list {
			if stock_id_in_list == stock_id {
				return true
			}
		}
		return false
	}

	// 用于检查 仓库和SKU 的重复
	var stock_product_map = make(map[string]string)

	var au_stock_id_map = make(map[int64]int64)
	var au_stock_import_list = make([]inventory_import_t, 0)
	var inventory_import_error_list = make([]inventory_import_t, 0)
	// 先检查一遍是不是有错误，
	for index, inventory_import := range inventory_import_list {

		// 检查 商品 sku
		var product_id int64 = product_service.Get_product_id_by_sku(inventory_import.ProductSku, db_session)

		if product_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的商品SKU %s 不存在！", inventory_import.OrdNum, inventory_import.ProductSku)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		}

		// 检查 仓库名称
		stock_id, stock_type_id := stock_service.Get_stock_id_type_by_name(inventory_import.StockName, db_session)
		var storage_loc_id int64

		if stock_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StockName)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		} else { //仓库名称存在
			if is_stock_id_can_view(stock_id) == false {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库 %s 你没有权限查看和修改！", inventory_import.OrdNum, inventory_import.StockName)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			var storage_loc_name = strings.TrimSpace(inventory_import.StorageLocName)
			if storage_loc_name != "" {
				// 如果提供了仓位， 检查仓位名称
				storage_loc_id = storage_loc_service.Get_storage_loc_id_by_name(storage_loc_name, stock_id, db_session)
				if storage_loc_id == 0 {
					inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓位名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StorageLocName)
					inventory_import_error_list = append(inventory_import_error_list, inventory_import)
				}
			}
		}

		// 检查导入的SKU是否有重复的，是 仓库和SKU 的组合的重复检查
		if stock_id != 0 && product_id != 0 {
			var map_key = fmt.Sprintf("%d_%d", stock_id, product_id)
			pre_order_num, found := stock_product_map[map_key]
			if found {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库（%s）和商品（%s）的组合和之前(序号%s)的有重复！", inventory_import.OrdNum, inventory_import.StockName, inventory_import.ProductSku, pre_order_num)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			stock_product_map[map_key] = inventory_import.OrdNum
		}

		inventory_import_list[index].product_id = product_id
		inventory_import_list[index].stock_id = stock_id
		inventory_import_list[index].stock_type_id = stock_type_id
		inventory_import_list[index].storage_loc_id = storage_loc_id

		if stock_type_id == sppz_models.STOCK_TYPE_AU {
			au_stock_id_map[stock_id] = stock_id
			au_stock_import_list = append(au_stock_import_list, inventory_import)
		}
	}

	if len(inventory_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": inventory_import_error_list,
		})
	}

	/*
	var au_stock_id_list = make([]int64, 0)
	for stock_id := range au_stock_id_map {
		au_stock_id_list = append(au_stock_id_list, stock_id)
	}

	if len(au_stock_id_list) > 0 {
		// 澳仓的做 一个仓位一个SKU的检查
		var au_stock_id_str_list = core.Int_list_to_str_list(au_stock_id_list)

//DROP TEMPORARY TABLE IF EXISTS tmp_stock_sku_storage;

		var sql_drop_tmp_table = `DROP TEMPORARY TABLE IF EXISTS tmp_stock_sku_storage;`
		_, err := db_session.Exec(sql_drop_tmp_table)
		core.HandleError(err)

//CREATE TEMPORARY TABLE tmp_stock_sku_storage (
//    stock_id       BIGINT NOT NULL,
//    product_id     BIGINT NOT NULL,
//    storage_loc_id BIGINT NOT NULL
//);

		var sql_create_tmp_table = `CREATE TEMPORARY TABLE tmp_stock_sku_storage (
    stock_id       BIGINT NOT NULL,
    product_id     BIGINT NOT NULL,
    storage_loc_id BIGINT NOT NULL
);`
		_, err = db_session.Exec(sql_create_tmp_table)
		core.HandleError(err)

		var sql_insert_record_in_db = `INSERT INTO tmp_stock_sku_storage (stock_id, product_id, storage_loc_id)
    SELECT nvt.stock_id, nvt.product_id, nvt.storage_loc_id
    FROM erp_inventory nvt
    WHERE nvt.stock_id IN (%s);`
		sql_insert_record_in_db = fmt.Sprintf(sql_insert_record_in_db, strings.Join(au_stock_id_str_list, ","))
		_, err = db_session.Exec(sql_insert_record_in_db)
		core.HandleError(err)

		// # do insert and update by stock_id and product_id

		// SELECT count(*) AS inventory_count FROM tmp_stock_sku_storage t WHERE t.stock_id = 0 AND t.product_id = 0;
		for _, inventory_import := range au_stock_import_list {
			var sql_select_count = `SELECT count(*) AS inventory_count FROM tmp_stock_sku_storage t WHERE t.stock_id = %d AND t.product_id = %d;`
			sql_select_count = fmt.Sprintf(sql_select_count, inventory_import.stock_id, inventory_import.storage_loc_id)
			var inventory_count struct {
				Count int64 `xorm:"'inventory_count'"`
			}
			_, err = db_session.Sql(sql_select_count).Get(&inventory_count)
			core.HandleError(err)

			if inventory_count.Count == 0 {
				var sql_insert = `INSERT INTO tmp_stock_sku_storage (stock_id, product_id, storage_loc_id) VALUES (%d, %d, %d);`
				sql_insert = fmt.Sprintf(sql_insert, inventory_import.stock_id, inventory_import.product_id, inventory_import.storage_loc_id)
				_, err = db_session.Exec(sql_insert)
				core.HandleError(err)
			} else {
				var sql_update = `UPDATE tmp_stock_sku_storage SET storage_loc_id = %d WHERE stock_id = %d AND product_id = %d;`
				sql_update = fmt.Sprintf(sql_update, inventory_import.storage_loc_id, inventory_import.stock_id, inventory_import.product_id)
				_, err = db_session.Exec(sql_update)
				core.HandleError(err)
			}
		}

		var sql_get_list = `SELECT s.name AS stock_name, sl.name AS storage_loc_name
FROM
    (
        SELECT t.stock_id, t.storage_loc_id FROM tmp_stock_sku_storage t
        GROUP BY t.stock_id, t.storage_loc_id
        HAVING count(*) > 1 # 如果大于一个记录，每个记录的product_id都不一样，所以就是一个仓位放了两个SKU
    ) ss
    INNER JOIN erp_stock s ON s.id = ss.stock_id
    INNER JOIN erp_storage_loc sl ON sl.id = ss.storage_loc_id;`

		type stock_storage_t struct {
			StockName      string `xorm:"''stock_name"`
			StorageLocName string `xorm:"''storage_loc_name"`
		}
		var stock_storage_list = make([]stock_storage_t, 0)
		err = db_session.Sql(sql_get_list).Find(&stock_storage_list)
		core.HandleError(err)

		if len(stock_storage_list) > 0 {
			var inventory_import_error_list = make([]inventory_import_t, 0)
			for _, stock_storage := range stock_storage_list {
				var inventory_import inventory_import_t
				inventory_import.StockName = stock_storage.StockName
				inventory_import.StorageLocName = stock_storage.StorageLocName

				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			return c.RenderJson(map[string]interface{}{
				"success":         false,
				"message":         "数据导入出现错误，详情查看页面下方表格！",
				"error_data_list": inventory_import_error_list,
			})
		}
	}
	*/

	for _, inventory_import := range inventory_import_list {

		var stock_id = inventory_import.stock_id
		var stock_type_id = inventory_import.stock_type_id
		var storage_loc_id = inventory_import.storage_loc_id
		var product_id = inventory_import.product_id

		if stock_type_id == sppz_models.STOCK_TYPE_AU {
			// 一个SKU一个仓位
			var inventory kcgl_models.InventoryInfo
			has, err := db_session.Where(fmt.Sprintf("stock_id = %d AND product_id = %d", stock_id, product_id)).Get(&inventory) //这里只能保证一个SKU一个仓位
			core.HandleError(err)

			inventory.StockId = stock_id
			inventory.StorageLocId = storage_loc_id
			inventory.ProductId = product_id

			if has == false {
				// 插入
				_, err = db_session.Insert(&inventory)
				core.HandleError(err)
			} else {
				// 更新
				_, err := db_session.Id(inventory.Id).Cols("storage_loc_id").Update(&inventory)
				core.HandleError(err)
			}
		} else {
			// 中国仓和在途仓，一个仓位可以多个SKU，就是 stock_id, product_id 组合好后，一个仓库一个SKU，storage_loc_id随便填
			var inventory kcgl_models.InventoryInfo
			has, err := db_session.Where(fmt.Sprintf("stock_id = %d AND product_id = %d", stock_id, product_id)).Get(&inventory)
			core.HandleError(err)

			inventory.StockId = stock_id
			inventory.StorageLocId = storage_loc_id
			inventory.ProductId = product_id

			if has == false {
				// 插入
				_, err = db_session.Insert(&inventory)
				core.HandleError(err)
			} else {
				// 更新
				_, err := db_session.Id(inventory.Id).Cols("storage_loc_id").Update(&inventory)
				core.HandleError(err)
			}
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpInventory) ExportProductStorageLoc() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_SKU_LOC, c.Session)

	db_session := c.DbSession

	xlsx_file := c.get_export_product_storage_loc_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "SKU仓位列表" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpInventory) get_export_product_storage_loc_excel(db_session *xorm.Session) *xlsx.File {
	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	/*
	SELECT ivt.*, st.*, loc.*, prd.*
	FROM erp_inventory ivt
		INNER JOIN erp_stock st ON ivt.stock_id = st.id
		LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	var sb = core.NewSqlBuilder()
	sb.Select("ivt.*, st.*, loc.*, prd.*")
	sb.Count("count(ivt.id) AS ivt_count")
	sb.From(`erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id`)

	//query extra filter here
	var stock_id_str_list = stock_pick_crew_service.Get_login_user_stock_id_str_list(db_session, c.Session)
	var filter = fmt.Sprintf("st.id IN (%s)", strings.Join(stock_id_str_list, ","))

	if len(id_list) == 0 {
		if stock_id_query != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("st.id = %d", stock_id_query))
		}

		if product_code_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_code_query))
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("ivt.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	sb.Where(filter)
	sb.OrderBy("ivt.id")

	inventory_list := make([]InventoryViewItem, 0)
	err := db_session.Sql(sb.SqlList()).Find(&inventory_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("库存列表")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = horizontal_align
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	row_height := 0.6

	// 序号	商品编号	仓库	仓位	库存	成本价（RMB）
	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "序号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 25)
		set_head_cell_style(cell)
		cell.Value = "商品编号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 22)
		set_head_cell_style(cell)
		cell.Value = "仓库"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "仓位"
	}

	for index, inventory_view := range inventory_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Product.Sku

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Stock.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.StorageLoc.Name
	}

	return xlsx_file
}

func (c ErpInventory) ImportInventoryQuantity() revel.Result {
	// pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY_OVER, c.Session)
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY, c.Session)

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

	//序号 商品编号	仓库	仓位	库存

	type inventory_import_t struct {
		OrdNum         string  `json:"ord_num"`          // 序号
		ProductSku     string  `json:"product_sku"`      // 商品编号
		StockName      string  `json:"stock_name"`       // 仓库名称
		StorageLocName string  `json:"storage_loc_name"` // 仓位名称
		Qty            int64   `json:"qty"`              // 库存
		CostPrice      float64 `json:"cost_price"`       // 成本价

		// 下面这些字段是从数据库里面查询后保存在这里面的，这样过程中，就不用二次查询了
		stock_id       int64
		storage_loc_id int64
		product_id     int64
		ErrorMessage   string `json:"error_message"`
	}

	inventory_import_list := make([]inventory_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var inventory_import inventory_import_t

			col_index := 0

			// 序号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			inventory_import.OrdNum = cell_val
			col_index++

			// 商品编号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.ProductSku = cell_val
			col_index++

			// 仓库名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StockName = cell_val
			col_index++

			/*
			// 仓位名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StorageLocName = cell_val
			col_index++
			*/

			// 库存
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			qty, _ := strconv.ParseInt(cell_val, 10, 64)
			inventory_import.Qty = qty
			col_index++

			inventory_import_list = append(inventory_import_list, inventory_import)
		}
	}

	if len(inventory_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品库存列表为空！",
		})
	}

	var stock_id_list = stock_pick_crew_service.Get_login_user_stock_id_list(db_session, c.Session)

	var is_stock_id_can_view = func(stock_id int64) bool {
		for _, stock_id_in_list := range stock_id_list {
			if stock_id_in_list == stock_id {
				return true
			}
		}
		return false
	}

	// 用于检查 仓库和SKU 的重复
	var stock_product_map = make(map[string]string)

	var inventory_import_error_list = make([]inventory_import_t, 0)
	// 先检查一遍是不是有错误，
	for index, inventory_import := range inventory_import_list {

		// 检查 商品 sku
		var product_id int64 = product_service.Get_product_id_by_sku(inventory_import.ProductSku, db_session)

		if product_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的商品SKU %s 不存在！", inventory_import.OrdNum, inventory_import.ProductSku)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		}

		// 检查 仓库名称
		stock_id, _ := stock_service.Get_stock_id_type_by_name(inventory_import.StockName, db_session)
		var storage_loc_id int64

		if stock_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StockName)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		} else { //仓库名称存在
			if is_stock_id_can_view(stock_id) == false {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库 %s 你没有权限查看和修改！", inventory_import.OrdNum, inventory_import.StockName)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			/*
			var storage_loc_name = strings.TrimSpace(inventory_import.StorageLocName)
			if storage_loc_name != "" {
				// 如果提供了仓位， 检查仓位名称
				storage_loc_id = storage_loc_service.Get_storage_loc_id_by_name(storage_loc_name, stock_id, db_session)
				if storage_loc_id == 0 {
					inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓位名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StorageLocName)
					inventory_import_error_list = append(inventory_import_error_list, inventory_import)
				}
			}
			*/
		}

		// 检查导入的SKU是否有重复的，是 仓库和SKU 的组合的重复检查
		if stock_id != 0 && product_id != 0 {
			var map_key = fmt.Sprintf("%d_%d", stock_id, product_id)
			pre_order_num, found := stock_product_map[map_key]
			if found {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库（%s）和商品（%s）的组合和之前(序号%s)的有重复！", inventory_import.OrdNum, inventory_import.StockName, inventory_import.ProductSku, pre_order_num)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			stock_product_map[map_key] = inventory_import.OrdNum
		}

		inventory_import_list[index].product_id = product_id
		inventory_import_list[index].stock_id = stock_id
		inventory_import_list[index].storage_loc_id = storage_loc_id
	}

	if len(inventory_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": inventory_import_error_list,
		})
	}

	revel.TRACE.Println("inventory_import_list", len(inventory_import_list))

	var login_user = core.GetLoginUser(c.Session)

	for _, inventory_import := range inventory_import_list {

		var stock_id = inventory_import.stock_id
		//var storage_loc_id = inventory_import.storage_loc_id
		var product_id = inventory_import.product_id

		var inventory kcgl_models.InventoryInfo
		has, err := db_session.Where(fmt.Sprintf("stock_id = %d AND product_id = %d", stock_id, product_id)).Get(&inventory)
		core.HandleError(err)

		inventory.StockId = stock_id
		//inventory.StorageLocId = storage_loc_id
		inventory.ProductId = product_id
		var qty_before = inventory.Qty
		inventory.Qty = inventory_import.Qty
		var qty_after = inventory.Qty

		if has == false {
			// 插入
			_, err = db_session.Insert(&inventory)

			inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("导入新增库存, 工号：%s", login_user.LoginName), db_session)
			core.HandleError(err)
		} else {
			// 更新
			_, err := db_session.Id(inventory.Id).Cols("qty").Update(&inventory)
			core.HandleError(err)

			inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("导入更新库存 %d => %d，工号：%s", qty_before, qty_after, login_user.LoginName), db_session)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpInventory) ExportInventoryQuantity() revel.Result {
	//pass_service.Check_login_user_pass_or([]int64{pass_res.KCGL_INVENTORY_IMPORT_QTY_OVER, pass_res.KCGL_INVENTORY_IMPORT_QTY_DIFF}, c.Session)
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_EXPORT_QTY, c.Session)

	db_session := c.DbSession

	xlsx_file := c.get_export_inventory_quantity_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "库存数量列表" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpInventory) get_export_inventory_quantity_excel(db_session *xorm.Session) *xlsx.File {
	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	/*
	SELECT ivt.*, st.*, loc.*, prd.*
	FROM erp_inventory ivt
		INNER JOIN erp_stock st ON ivt.stock_id = st.id
		LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	var sb = core.NewSqlBuilder()
	sb.Select("ivt.*, st.*, loc.*, prd.*")
	sb.Count("count(ivt.id) AS ivt_count")
	sb.From(`erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id`)

	//query extra filter here
	var stock_id_str_list = stock_pick_crew_service.Get_login_user_stock_id_str_list(db_session, c.Session)
	var filter = fmt.Sprintf("st.id IN (%s)", strings.Join(stock_id_str_list, ","))

	if len(id_list) == 0 {
		if stock_id_query != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("st.id = %d", stock_id_query))
		}

		if product_code_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_code_query))
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("ivt.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	sb.Where(filter)
	sb.OrderBy("ivt.id")

	inventory_list := make([]InventoryViewItem, 0)
	err := db_session.Sql(sb.SqlList()).Find(&inventory_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("库存列表")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = horizontal_align
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	row_height := 0.6

	// 序号	商品编号	仓库	仓位	库存	成本价（RMB）
	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "序号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 25)
		set_head_cell_style(cell)
		cell.Value = "商品编号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 22)
		set_head_cell_style(cell)
		cell.Value = "仓库"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "仓位"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "库存"
	}

	for index, inventory_view := range inventory_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Product.Sku

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Stock.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.StorageLoc.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(inventory_view.Inventory.Qty)
	}

	return xlsx_file
}

func (c ErpInventory) ImportInventoryPrice() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_SKU_PRICE, c.Session)

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

	//序号 商品编号	仓库	仓位	库存

	type inventory_import_t struct {
		OrdNum         string  `json:"ord_num"`          // 序号
		ProductSku     string  `json:"product_sku"`      // 商品编号
		StockName      string  `json:"stock_name"`       // 仓库名称
		StorageLocName string  `json:"storage_loc_name"` // 仓位名称
		Qty            int64   `json:"qty"`              // 库存
		CostPrice      float64 `json:"cost_price"`       // 成本价

		// 下面这些字段是从数据库里面查询后保存在这里面的，这样过程中，就不用二次查询了
		stock_id       int64
		storage_loc_id int64
		product_id     int64
		ErrorMessage   string `json:"error_message"`
	}

	inventory_import_list := make([]inventory_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var inventory_import inventory_import_t

			col_index := 0

			// 序号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			inventory_import.OrdNum = cell_val
			col_index++

			// 商品编号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.ProductSku = cell_val
			col_index++

			// 仓库名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StockName = cell_val
			col_index++

			/*
			// 仓位名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StorageLocName = cell_val
			col_index++
			*/

			// 成本价
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			cost_price, _ := strconv.ParseFloat(cell_val, 64)
			inventory_import.CostPrice = cost_price
			col_index++

			inventory_import_list = append(inventory_import_list, inventory_import)
		}
	}

	if len(inventory_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品库存列表为空！",
		})
	}

	var stock_id_list = stock_pick_crew_service.Get_login_user_stock_id_list(db_session, c.Session)

	var is_stock_id_can_view = func(stock_id int64) bool {
		for _, stock_id_in_list := range stock_id_list {
			if stock_id_in_list == stock_id {
				return true
			}
		}
		return false
	}

	// 用于检查 仓库和SKU 的重复
	var stock_product_map = make(map[string]string)

	var inventory_import_error_list = make([]inventory_import_t, 0)
	// 先检查一遍是不是有错误，
	for index, inventory_import := range inventory_import_list {

		// 检查 商品 sku
		var product_id int64 = product_service.Get_product_id_by_sku(inventory_import.ProductSku, db_session)

		if product_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的商品SKU %s 不存在！", inventory_import.OrdNum, inventory_import.ProductSku)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		}

		// 检查 仓库名称
		stock_id, type_id := stock_service.Get_stock_id_type_by_name(inventory_import.StockName, db_session)
		if type_id != sppz_models.STOCK_TYPE_CN {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "只能导入中国仓的成本！",
			})
		}

		var storage_loc_id int64

		if stock_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StockName)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		} else { //仓库名称存在
			if is_stock_id_can_view(stock_id) == false {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库 %s 你没有权限查看和修改！", inventory_import.OrdNum, inventory_import.StockName)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			/*
			var storage_loc_name = strings.TrimSpace(inventory_import.StorageLocName)
			if storage_loc_name != "" {
				// 如果提供了仓位， 检查仓位名称
				storage_loc_id = storage_loc_service.Get_storage_loc_id_by_name(storage_loc_name, stock_id, db_session)
				if storage_loc_id == 0 {
					inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓位名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StorageLocName)
					inventory_import_error_list = append(inventory_import_error_list, inventory_import)
				}
			}
			*/
		}

		// 检查导入的SKU是否有重复的，是 仓库和SKU 的组合的重复检查
		if stock_id != 0 && product_id != 0 {
			var map_key = fmt.Sprintf("%d_%d", stock_id, product_id)
			pre_order_num, found := stock_product_map[map_key]
			if found {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库（%s）和商品（%s）的组合和之前(序号%s)的有重复！", inventory_import.OrdNum, inventory_import.StockName, inventory_import.ProductSku, pre_order_num)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			stock_product_map[map_key] = inventory_import.OrdNum
		}

		inventory_import_list[index].product_id = product_id
		inventory_import_list[index].stock_id = stock_id
		inventory_import_list[index].storage_loc_id = storage_loc_id
	}

	if len(inventory_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": inventory_import_error_list,
		})
	}

	for _, inventory_import := range inventory_import_list {

		var stock_id = inventory_import.stock_id
		//var storage_loc_id = inventory_import.storage_loc_id
		var product_id = inventory_import.product_id

		var inventory kcgl_models.InventoryInfo
		has, err := db_session.Where(fmt.Sprintf("stock_id = %d AND product_id = %d", stock_id, product_id)).Get(&inventory)
		core.HandleError(err)

		inventory.StockId = stock_id
		// inventory.StorageLocId = storage_loc_id
		inventory.ProductId = product_id
		inventory.CostPrice = inventory_import.CostPrice

		if has == false {
			// 插入
			_, err = db_session.Insert(&inventory)
			core.HandleError(err)
		} else {
			// 更新
			_, err := db_session.Id(inventory.Id).Cols("cost_price").Update(&inventory)
			core.HandleError(err)
		}
	}

	// 导入中国仓价格之后，马上更新澳仓价格
	go inventory_service.Update_au_stock_cost()

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpInventory) ExportInventoryPrice() revel.Result {
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_SKU_PRICE, c.Session)

	db_session := c.DbSession

	xlsx_file := c.get_export_inventory_price_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "库存价格列表" + ".xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpInventory) get_export_inventory_price_excel(db_session *xorm.Session) *xlsx.File {
	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var product_code_query string
	c.Params.Bind(&product_code_query, "product_code_query")
	product_code_query = strings.TrimSpace(product_code_query)

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	/*
	SELECT ivt.*, st.*, loc.*, prd.*
	FROM erp_inventory ivt
		INNER JOIN erp_stock st ON ivt.stock_id = st.id
		LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		INNER JOIN erp_product prd ON ivt.product_id = prd.id
	*/

	var sb = core.NewSqlBuilder()
	sb.Select("ivt.*, st.*, loc.*, prd.*")
	sb.Count("count(ivt.id) AS ivt_count")
	sb.From(`erp_inventory ivt
		    INNER JOIN erp_stock st ON ivt.stock_id = st.id
		    LEFT JOIN erp_storage_loc loc ON ivt.storage_loc_id = loc.id
		    INNER JOIN erp_product prd ON ivt.product_id = prd.id`)

	//query extra filter here
	var stock_id_str_list = stock_pick_crew_service.Get_login_user_stock_id_str_list(db_session, c.Session)
	var filter = fmt.Sprintf("st.id IN (%s)", strings.Join(stock_id_str_list, ","))

	if len(id_list) == 0 {
		if stock_id_query != 0 {
			filter = core.SqlAnd(filter, fmt.Sprintf("st.id = %d", stock_id_query))
		}

		if product_code_query != "" {
			filter = core.SqlAnd(filter, fmt.Sprintf("prd.sku LIKE '%%%s%%'", product_code_query))
		}
	} else {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("ivt.id IN (%s)", strings.Join(id_str_list, ",")))
	}

	sb.Where(filter)
	sb.OrderBy("ivt.id")

	inventory_list := make([]InventoryViewItem, 0)
	err := db_session.Sql(sb.SqlList()).Find(&inventory_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("库存列表")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = horizontal_align
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	row_height := 0.6

	// 序号	商品编号	仓库	仓位	库存	成本价（RMB）
	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0
		cell := header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 8)
		set_head_cell_style(cell)
		cell.Value = "序号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 25)
		set_head_cell_style(cell)
		cell.Value = "商品编号"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 22)
		set_head_cell_style(cell)
		cell.Value = "仓库"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 18)
		set_head_cell_style(cell)
		cell.Value = "仓位"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "成本价"
	}

	for index, inventory_view := range inventory_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Product.Sku

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.Stock.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = inventory_view.StorageLoc.Name

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(inventory_view.Inventory.CostPrice)
	}

	return xlsx_file
}

func (c ErpInventory) ImportInventoryQuantityDiff() revel.Result {
	// pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY_DIFF, c.Session)
	pass_service.Check_login_user_pass(pass_res.KCGL_INVENTORY_IMPORT_QTY, c.Session)

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

	//序号 商品编号	仓库	仓位	库存

	type inventory_import_t struct {
		OrdNum         string  `json:"ord_num"`          // 序号
		ProductSku     string  `json:"product_sku"`      // 商品编号
		StockName      string  `json:"stock_name"`       // 仓库名称
		StorageLocName string  `json:"storage_loc_name"` // 仓位名称
		Qty            int64   `json:"qty"`              // 库存
		CostPrice      float64 `json:"cost_price"`       // 成本价

		// 下面这些字段是从数据库里面查询后保存在这里面的，这样过程中，就不用二次查询了
		stock_id       int64
		storage_loc_id int64
		product_id     int64
		ErrorMessage   string `json:"error_message"`
	}

	inventory_import_list := make([]inventory_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var inventory_import inventory_import_t

			col_index := 0

			// 序号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			inventory_import.OrdNum = cell_val
			col_index++

			// 商品编号
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.ProductSku = cell_val
			col_index++

			// 仓库名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StockName = cell_val
			col_index++

			/*
			// 仓位名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			inventory_import.StorageLocName = cell_val
			col_index++
			*/

			// 库存
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			qty, _ := strconv.ParseInt(cell_val, 10, 64)
			inventory_import.Qty = qty
			col_index++

			inventory_import_list = append(inventory_import_list, inventory_import)
		}
	}

	if len(inventory_import_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品库存列表为空！",
		})
	}

	var stock_id_list = stock_pick_crew_service.Get_login_user_stock_id_list(db_session, c.Session)

	var is_stock_id_can_view = func(stock_id int64) bool {
		for _, stock_id_in_list := range stock_id_list {
			if stock_id_in_list == stock_id {
				return true
			}
		}
		return false
	}

	// 用于检查 仓库和SKU 的重复
	var stock_product_map = make(map[string]string)

	var inventory_import_error_list = make([]inventory_import_t, 0)
	// 先检查一遍是不是有错误，
	for index, inventory_import := range inventory_import_list {

		// 检查 商品 sku
		var product_id int64 = product_service.Get_product_id_by_sku(inventory_import.ProductSku, db_session)

		if product_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的商品SKU %s 不存在！", inventory_import.OrdNum, inventory_import.ProductSku)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		}

		// 检查 仓库名称
		stock_id, _ := stock_service.Get_stock_id_type_by_name(inventory_import.StockName, db_session)
		var storage_loc_id int64

		if stock_id == 0 {
			inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StockName)
			inventory_import_error_list = append(inventory_import_error_list, inventory_import)
		} else { //仓库名称存在
			if is_stock_id_can_view(stock_id) == false {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库 %s 你没有权限查看和修改！", inventory_import.OrdNum, inventory_import.StockName)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			/*
			var storage_loc_name = strings.TrimSpace(inventory_import.StorageLocName)
			if storage_loc_name != "" {
				// 如果提供了仓位， 检查仓位名称
				storage_loc_id = storage_loc_service.Get_storage_loc_id_by_name(storage_loc_name, stock_id, db_session)
				if storage_loc_id == 0 {
					inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓位名称 %s 不存在！", inventory_import.OrdNum, inventory_import.StorageLocName)
					inventory_import_error_list = append(inventory_import_error_list, inventory_import)
				}
			}
			*/
		}

		// 检查导入的SKU是否有重复的，是 仓库和SKU 的组合的重复检查
		if stock_id != 0 && product_id != 0 {
			var map_key = fmt.Sprintf("%d_%d", stock_id, product_id)
			pre_order_num, found := stock_product_map[map_key]
			if found {
				inventory_import.ErrorMessage = fmt.Sprintf("序号 %s 的仓库（%s）和商品（%s）的组合和之前(序号%s)的有重复！", inventory_import.OrdNum, inventory_import.StockName, inventory_import.ProductSku, pre_order_num)
				inventory_import_error_list = append(inventory_import_error_list, inventory_import)
			}

			stock_product_map[map_key] = inventory_import.OrdNum
		}

		inventory_import_list[index].product_id = product_id
		inventory_import_list[index].stock_id = stock_id
		inventory_import_list[index].storage_loc_id = storage_loc_id
	}

	if len(inventory_import_error_list) != 0 {
		return c.RenderJson(map[string]interface{}{
			"success":         false,
			"message":         "数据导入出现错误，详情查看页面下方表格！",
			"error_data_list": inventory_import_error_list,
		})
	}

	revel.TRACE.Println("inventory_import_list", len(inventory_import_list))

	var login_user = core.GetLoginUser(c.Session)

	for _, inventory_import := range inventory_import_list {

		var stock_id = inventory_import.stock_id
		//var storage_loc_id = inventory_import.storage_loc_id
		var product_id = inventory_import.product_id

		var inventory kcgl_models.InventoryInfo
		has, err := db_session.Where(fmt.Sprintf("stock_id = %d AND product_id = %d", stock_id, product_id)).Get(&inventory)
		core.HandleError(err)

		inventory.StockId = stock_id
		//inventory.StorageLocId = storage_loc_id
		inventory.ProductId = product_id
		var qty_before = inventory.Qty
		inventory.Qty = inventory.Qty + inventory_import.Qty
		var qty_after = inventory.Qty

		if has == false {
			// 插入
			_, err = db_session.Insert(&inventory)
			core.HandleError(err)

			inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("导入差异新增， 工号：%s", login_user.LoginName), db_session)
		} else {
			// 更新
			_, err := db_session.Id(inventory.Id).Cols("qty").Update(&inventory)
			core.HandleError(err)

			inventory_history_service.Insert_inventory_history(inventory, fmt.Sprintf("导入差异更新 %d + %d => %d，工号：%s", qty_before, inventory_import.Qty, qty_after, login_user.LoginName), db_session)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}
