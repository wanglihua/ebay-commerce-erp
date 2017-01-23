package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpStockBillEntry struct {
	*revel.Controller
}

func (c ErpStockBillEntry) Index(bill_id int64) revel.Result {
	c.RenderArgs["bill_id"] = bill_id
	return c.RenderTemplate("erp/kcgl/stock_bill_entry/stock_bill_entry_index.html")
}

func (c ErpStockBillEntry) ListData(bill_id int64) revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	stock_bill_entry_list := make([]models.StockBillEntryInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&stock_bill_entry_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.StockBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  stock_bill_entry_list,
		Total: count,
	})
}

type StockBillEntryDetailForm struct {
	StockBillEntry models.StockBillEntryInfo `json:"bill_entry"`
}

func (f *StockBillEntryDetailForm) IsCreate() bool {
	return f.StockBillEntry.Id == 0
}

func (f *StockBillEntryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.StockBillEntry.BillId).Message("批次不能为空！")

	validation.Required(f.StockBillEntry.StorageLocId).Message("仓位不能为空！")

	validation.Required(f.StockBillEntry.ProductId).Message("商品不能为空！")

	validation.Required(f.StockBillEntry.PreQty).Message("当前库存不能为空！")

	validation.Required(f.StockBillEntry.Qty).Message("出入库数量不能为空！")

	return validation.HasErrors() == false
}

func (c ErpStockBillEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var stock_bill_entry_id int64
	c.Params.Bind(&stock_bill_entry_id, "id")

	var stock_bill_entry models.StockBillEntryInfo
	if stock_bill_entry_id != 0 {
		has, err := db_session.Id(stock_bill_entry_id).Get(&stock_bill_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的出入库单明细不存在！"})
		}
	}

	var detail_form StockBillEntryDetailForm
	detail_form.StockBillEntry = stock_bill_entry

	return c.RenderJson(detail_form)
}

func (c ErpStockBillEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form StockBillEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	stock_bill_entry := detail_form.StockBillEntry

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

func (c ErpStockBillEntry) Delete() revel.Result {
	db_session := c.DbSession

	stock_bill_entry_id_list := make([]int64, 0)
	c.Params.Bind(&stock_bill_entry_id_list, "id_list")

	affected, err := db_session.In("id", stock_bill_entry_id_list).Delete(new(models.StockBillEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
