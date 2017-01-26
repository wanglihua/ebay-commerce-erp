package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"html/template"
)

type ErpStockBillTransferEntry struct {
	*revel.Controller
}

func (c ErpStockBillTransferEntry) Index(bill_id int64) revel.Result {
	db_session := c.DbSession

	var stock_bill models.StockBillInfo
	has, err := db_session.Id(bill_id).Get(&stock_bill)
	core.HandleError(err)

	if has == false {
		panic("指定的批次号不存在！")
	}

	var stock models.StockInfo
	_, err = db_session.Id(stock_bill.StockId).Get(&stock)
	core.HandleError(err)

	breadcrumb_nav := template.HTML("")
	bill_type := ""
	if stock_bill.TypeId == models.STOCK_BILL_TYPE_CHECKIN {
		bill_type = "入库单"
		breadcrumb_nav = template.HTML("<li class='active' data-menu='kcgl-stockbill-checkin-menu'><a href='/erp/stock/bill/checkin'>手工入库</a></li>")
	} else if stock_bill.TypeId == models.STOCK_BILL_TYPE_CHECKOUT {
		bill_type = "出库单"
		breadcrumb_nav = template.HTML("<li class='active' data-menu='kcgl-stockbill-checkout-menu'><a href='/erp/stock/bill/checkout'>手工出库</a></li>")
	} else if stock_bill.TypeId == models.STOCK_BILL_TYPE_TRANSFER {
		bill_type = "调拨单"
		breadcrumb_nav = template.HTML("<li class='active' data-menu='kcgl-stockbill-transfer-menu'><a href='/erp/stock/bill/transfer'>分仓调拨</a></li>")
	}

	c.RenderArgs["breadcrumb_nav"] = breadcrumb_nav
	c.RenderArgs["bill_type"] = bill_type
	c.RenderArgs["stock_bill"] = stock_bill
	c.RenderArgs["stock_name"] = stock.Name
	return c.RenderTemplate("erp/kcgl/stock_bill_transfer_entry/stock_bill_transfer_entry_index.html")
}

func (c ErpStockBillTransferEntry) ListData(bill_id int64) revel.Result {
	db_session := c.DbSession

	filter, order, _, _ := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here
	query = query.Where("bill_id = ?", bill_id)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	stock_bill_entry_list := make([]models.StockBillEntryInfo, 0)
	err := data_query.Find(&stock_bill_entry_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.StockBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  stock_bill_entry_list,
		Total: count,
	})
}

type StockBillEntryTransferDetailForm struct {
	Id             int64   `xorm:"'id'" json:"id"`
	BillId         int64   `xorm:"'bill_id'" json:"bill_id"`
	StorageLocName string  `xorm:"'storage_loc_name'" json:"storage_loc_name"`
	ProductSku     string  `xorm:"'product_sku'" json:"product_sku"`
	Qty            int64   `xorm:"'qty'" json:"qty"`
	Price          float64 `xorm:"'price'" json:"price"`
}

func (f *StockBillEntryTransferDetailForm) IsCreate() bool {
	return f.Id == 0
}

func (f *StockBillEntryTransferDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.StorageLocName).Message("库位不能为空！")
	validation.Required(f.ProductSku).Message("商品SKU不能为空！")
	validation.Required(f.Qty).Message("数量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpStockBillTransferEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var stock_bill_entry_id int64
	c.Params.Bind(&stock_bill_entry_id, "id")

	var detail_form StockBillEntryTransferDetailForm
	if stock_bill_entry_id != 0 {
		/*
			SELECT sbe.id AS id, sl.name AS storage_loc_name, p.sku AS product_sku, sbe.qty AS qty, sbe.price AS price
			FROM erp_stock_bill_entry sbe
			    INNER JOIN erp_storage_loc sl ON sbe.storage_loc_id = sl.id
			    INNER JOIN erp_product p ON sbe.product_id = p.id
			WHERE sbe.id = 1
		*/
		sql := `
SELECT sbe.id AS id, sl.name AS storage_loc_name, p.sku AS product_sku, sbe.qty AS qty, sbe.price AS price
FROM erp_stock_bill_entry sbe
    INNER JOIN erp_storage_loc sl ON sbe.storage_loc_id = sl.id
    INNER JOIN erp_product p ON sbe.product_id = p.id
WHERE sbe.id = ?
`
		has, err := db_session.Sql(sql, stock_bill_entry_id).Get(&detail_form)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的明细信息不存在！"})
		}
	}

	return c.RenderJson(detail_form)
}

func (c ErpStockBillTransferEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form StockBillEntryTransferDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	// 获取仓位 id
	sql := `
SELECT id AS storage_loc_id
FROM erp_storage_loc
WHERE name = ?;
`
	var storage_loc_id struct {
		StorageLocId core.NullInt `xorm:"'storage_loc_id'"`
	}

	has, err := db_session.Sql(sql, detail_form.StorageLocName).Get(&storage_loc_id)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "输入的仓位不存在！"})
	}

	// 获取商品 id
	sql = `
SELECT id AS product_id
FROM erp_product
WHERE sku = ?;
`
	var product_id struct {
		ProductId core.NullInt `xorm:"'product_id'"`
	}

	has, err = db_session.Sql(sql, detail_form.ProductSku).Get(&product_id)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "输入的商品SKU不存在！"})
	}

	sql = `
SELECT v.qty AS qty
FROM erp_inventory v
WHERE v.storage_loc_id = ? AND v.product_id = ?
`
	var qty struct {
		Qty core.NullInt `xorm:"'qty'"`
	}

	_, err = db_session.Sql(sql, storage_loc_id.StorageLocId.Int64, product_id.ProductId.Int64).Get(&qty)
	core.HandleError(err)

	var product_storage_loc_qty int64 = 0
	if qty.Qty.Valid == true {
		product_storage_loc_qty = qty.Qty.Int64
	}

	var stock_bill_entry models.StockBillEntryInfo

	var affected int64
	if detail_form.IsCreate() {
		stock_bill_entry.Qty = product_storage_loc_qty
		affected, err = db_session.Insert(&stock_bill_entry)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(stock_bill_entry.Id).Cols("").Update(&stock_bill_entry)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpStockBillTransferEntry) Delete() revel.Result {
	db_session := c.DbSession

	stock_bill_entry_id_list := make([]int64, 0)
	c.Params.Bind(&stock_bill_entry_id_list, "id_list")

	affected, err := db_session.In("id", stock_bill_entry_id_list).Delete(new(models.StockBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
