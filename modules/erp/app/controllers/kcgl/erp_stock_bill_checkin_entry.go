package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
)

type ErpStockBillCheckinEntry struct {
	*revel.Controller
}

func (c ErpStockBillCheckinEntry) Index(bill_id int64) revel.Result {
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

	c.RenderArgs["stock_bill"] = stock_bill
	c.RenderArgs["stock_name"] = stock.Name
	return c.RenderTemplate("erp/kcgl/stock_bill_checkin_entry/stock_bill_checkin_entry_index.html")
}

type ErpStockBillCheckinEntryViewItem struct {
	models.StockBillEntryInfo `xorm:"extends" json:"sbe"`
	models.StockBillInfo      `xorm:"extends" json:"sb"`
	models.StorageLocInfo     `xorm:"extends" json:"sl"`
	models.ProductInfo        `xorm:"extends" json:"p"`
}

func (c ErpStockBillCheckinEntry) ListData(bill_id int64) revel.Result {
	db_session := c.DbSession

	filter, order, _, _ := core.GetGridRequestParam(c.Request)

	/*
		SELECT sbe.*, sb.*, sl.*, p.*
		FROM erp_stock_bill_entry sbe
		    INNER JOIN erp_stock_bill sb ON sbe.bill_id = sb.id
		    INNER JOIN erp_storage_loc sl ON sbe.storage_loc_id = sl.id
		    INNER JOIN erp_product p ON sbe.product_id = p.id
	 */
	sb := core.NewSqlBuilder()
	sb.Select(`sbe.*, sb.*, sl.*, p.*`)
	sb.Count("count(*) AS entry_count")
	sb.From(`erp_stock_bill_entry sbe
		    INNER JOIN erp_stock_bill sb ON sbe.bill_id = sb.id
		    INNER JOIN erp_storage_loc sl ON sbe.storage_loc_id = sl.id
		    INNER JOIN erp_product p ON sbe.product_id = p.id`)
	sb.Where(core.SqlAnd(filter, fmt.Sprintf("bill_id = %d", bill_id)))

	//query extra filter here

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("sbe.id ASC")
	}

	stock_bill_entry_list := make([]ErpStockBillCheckinEntryViewItem, 0)
	err := db_session.Sql(sb.SqlList()).Find(&stock_bill_entry_list)
	core.HandleError(err)

	var entry_count struct {
		Count int64 `xorm:"'entry_count'" json:"entry_count"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&entry_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  stock_bill_entry_list,
		Total: entry_count.Count,
	})
}

type StockBillEntryCheckinDetailForm struct {
	Id             int64            `xorm:"'id'" json:"id"`
	BillId         int64            `xorm:"'bill_id'" json:"bill_id"`
	StorageLocName string           `xorm:"'storage_loc_name'" json:"storage_loc_name"`
	ProductSku     string           `xorm:"'product_sku'" json:"product_sku"`
	Qty            int64            `xorm:"'qty'" json:"qty"`
	Price          core.NullFloat   `xorm:"'price'" json:"price"`
}

func (f *StockBillEntryCheckinDetailForm) IsCreate() bool {
	return f.Id == 0
}

func (f *StockBillEntryCheckinDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.BillId).Message("入库单不能为空！")
	validation.Required(f.StorageLocName).Message("库位不能为空！")
	validation.Required(f.ProductSku).Message("商品SKU不能为空！")
	validation.Required(f.Qty).Message("数量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpStockBillCheckinEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var stock_bill_entry_id int64
	c.Params.Bind(&stock_bill_entry_id, "id")

	var detail_form StockBillEntryCheckinDetailForm
	if stock_bill_entry_id != 0 {
		/*
			SELECT sbe.id AS id, sb.id AS bill_id, sl.name AS storage_loc_name, p.sku AS product_sku, sbe.qty AS qty, sbe.price AS price
			FROM erp_stock_bill_entry sbe
			    INNER JOIN erp_stock_bill sb ON sbe.bill_id = sb.id
			    INNER JOIN erp_storage_loc sl ON sbe.storage_loc_id = sl.id
			    INNER JOIN erp_product p ON sbe.product_id = p.id
			WHERE sbe.id = 1
		*/
		sql := `
SELECT sbe.id AS id, sb.id AS bill_id, sl.name AS storage_loc_name, p.sku AS product_sku, sbe.qty AS qty, sbe.price AS price
FROM erp_stock_bill_entry sbe
    INNER JOIN erp_stock_bill sb ON sbe.bill_id = sb.id
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

func (c ErpStockBillCheckinEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form StockBillEntryCheckinDetailForm
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
	stock_bill_entry.BillId = detail_form.BillId
	stock_bill_entry.StorageLocId = storage_loc_id.StorageLocId.Int64
	stock_bill_entry.ProductId = product_id.ProductId.Int64
	stock_bill_entry.PreQty = product_storage_loc_qty
	stock_bill_entry.Qty = detail_form.Qty
	stock_bill_entry.Price = detail_form.Price

	var affected int64
	if detail_form.IsCreate() {
		if detail_form.Price.Valid == true && detail_form.Price.Float64 != 0 {
			stock_bill_entry.Amount = core.NewNullFloat(detail_form.Price.Float64*float64(detail_form.Qty), true)
		} else {
			stock_bill_entry.Price = core.NewNullFloat(0, false)
			stock_bill_entry.Amount = core.NewNullFloat(0, false)
		}
		affected, err = db_session.Insert(&stock_bill_entry)
		core.HandleError(err)
	} else {
		stock_bill_entry.Id = detail_form.Id
		cols := models.StockBillEntryCols
		affected, err = db_session.Id(stock_bill_entry.Id).Cols(cols.StorageLocId, cols.ProductId, cols.Qty, cols.Price).Update(&stock_bill_entry)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpStockBillCheckinEntry) Delete() revel.Result {
	db_session := c.DbSession

	stock_bill_entry_id_list := make([]int64, 0)
	c.Params.Bind(&stock_bill_entry_id_list, "id_list")

	affected, err := db_session.In("id", stock_bill_entry_id_list).Delete(new(models.StockBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
