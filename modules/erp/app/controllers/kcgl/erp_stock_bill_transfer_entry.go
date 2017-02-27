package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpStockBillTransferEntry struct {
	*revel.Controller
}

func (c ErpStockBillTransferEntry) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/stock_bill_transfer_entry/stock_bill_transfer_entry_index.html")
}

func (c ErpStockBillTransferEntry) ListData() revel.Result {
	db_session := c.DbSession

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	var stock_bill models.StockBillInfo
	if bill_id != 0 {
		has, err := db_session.Id(bill_id).Get(&stock_bill)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的调拨单不存在！"})
		}
	}

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

	return c.RenderJson(map[string]interface{}{
		"bill_code": func() string {
			if stock_bill.Code == "" {
				return "请先单击调拨单行"
			}
			return stock_bill.Code
		}(),
		"list_data": core.GridResult{
			Data:  stock_bill_entry_list,
			Total: count,
		},
	})
}

type StockBillEntryTransferDetailForm struct {
	StockBill      models.StockBillInfo        `json:"stock_bill"`
	StockBillEntry models.StockBillEntryInfo   `json:"bill_entry"`
	ProductName    string                      `json:"product_name"`
	//StockName      string                      `json:"stock_name"`
	//StorageLocList []models.StorageLocInfo     `json:"storage_loc_list"`
}

func (f *StockBillEntryTransferDetailForm) IsCreate() bool {
	return f.StockBill.Id == 0
}

func (f *StockBillEntryTransferDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.StockBillEntry.BillId).Message("盘点批次不能为空！")
	//validation.Required(f.StockBillEntry.StorageLocId).Message("库位不能为空！")
	validation.Required(f.StockBillEntry.ProductId).Message("商品不能为空！")
	validation.Required(f.StockBillEntry.Qty).Message("盘点库存不能为空！")

	return validation.HasErrors() == false
}

func (c ErpStockBillTransferEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	if bill_id == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "请先选择盘点单！"})
	}

	var stock_bill models.StockBillInfo
	has, err := db_session.Id(bill_id).Get(&stock_bill)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的盘点单不存在！"})
	}

	/*
	var stock models.StockInfo
	_, err = db_session.Id(stock_bill.StockId).Get(&stock)

	storage_loc_list := make([]models.StorageLocInfo, 0)
	err = db_session.Where(fmt.Sprintf("stock_id = %d", stock_bill.StockId)).OrderBy("name asc").Find(&storage_loc_list)
	core.HandleError(err)
	*/

	var bill_entry_id int64
	c.Params.Bind(&bill_entry_id, "id")

	var stock_bill_entry models.StockBillEntryInfo
	var product_name string
	if bill_entry_id != 0 {
		has, err := db_session.Id(bill_entry_id).Get(&stock_bill_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的库存盘点明细不存在！"})
		}

		var product models.ProductInfo
		product_id := stock_bill_entry.ProductId
		_, err = db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		product_name = product.Name
	}

	var detail_form StockBillEntryTransferDetailForm
	detail_form.StockBill = stock_bill
	detail_form.StockBillEntry = stock_bill_entry
	detail_form.ProductName = product_name
	//detail_form.StockName = stock.Name
	//detail_form.StorageLocList = storage_loc_list

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

	stock_bill_entry := detail_form.StockBillEntry
	stock_bill_entry.BillId = detail_form.StockBill.Id

	bill_id := detail_form.StockBill.Id
	var stock_bill models.StockBillInfo
	_, err = db_session.Id(bill_id).Get(&stock_bill)
	core.HandleError(err)

	//需要判断某些状态下，信息已不能修改

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&stock_bill_entry)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(stock_bill_entry.Id).Update(&stock_bill_entry)
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

func (c ErpStockBillTransferEntry) ProductSelectListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type ErpProductViewItem struct {
		models.ProductInfo          `xorm:"extends" json:"prd"`
		models.ProductCateInfo      `xorm:"extends" json:"cate"`
		models.SupplierInfo         `xorm:"extends" json:"sup"`
		models.PackageMaterialInfo  `xorm:"extends" json:"pm"`
		models.StatusInfo           `xorm:"extends" json:"st"`
	}

	/*
SELECT prd.*, cate.*, sup.*, pm.*, st.*
FROM erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id
	*/

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("prd.*, cate.*, sup.*, pm.*, st.*")
	sql_builder.Count("count(prd.id) as product_count")
	sql_builder.From(`erp_product prd
    LEFT JOIN erp_product_cate cate ON prd.cate_id = cate.id
    LEFT JOIN erp_supplier sup ON prd.supplier_id = sup.id
    LEFT JOIN erp_package_material pm ON prd.pkg_mtrl_id = pm.id
    INNER JOIN erp_status st ON prd.status_id = st.id`)

	sql_builder.Where(filter)

	//query extra filter here

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("prd.id ASC")
	}

	sql_builder.Limit(offset, limit)

	product_list := make([]ErpProductViewItem, 0, limit)
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
