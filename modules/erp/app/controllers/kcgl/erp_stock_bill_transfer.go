package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/models"
	"time"
	"fmt"
)

type ErpStockBillTransfer struct {
	*revel.Controller
}

func (c ErpStockBillTransfer) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/stock_bill_transfer/stock_bill_transfer_index.html")
}

type ErpStockBillTransferViewItem struct {
	models.StockBillInfo `xorm:"extends" json:"sb"`
	SrcStock             models.StockInfo `xorm:"extends" json:"sc"`
	DestStock            models.StockInfo `xorm:"extends" json:"dc"`
	auth_models.UserInfo `xorm:"extends" json:"u"`
}

func (c ErpStockBillTransfer) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT sb.*, sc.*, dc.*, u.*
		FROM erp_stock_bill sb
		    LEFT JOIN erp_stock sc ON sb.src_stock_id = sc.id
		    INNER JOIN erp_stock dc ON sb.stock_id = dc.id
		    INNER JOIN auth_user u ON sb.user_id = u.id
		WHERE sb.type_id = 300
	*/

	query := db_session.
		Select("sb.*, sc.*, dc.*, u.*").
		Table("erp_stock_bill").Alias("sb").
		Join("LEFT", []string{"erp_stock", "sc"}, "sb.src_stock_id = sc.id").
		Join("INNER", []string{"erp_stock", "dc"}, "sb.stock_id = dc.id").
		Join("INNER", []string{"auth_user", "u"}, "sb.user_id = u.id").
		Where("sb.type_id = ?", models.STOCK_BILL_TYPE_TRANSFER).
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("sb.id")
	}

	stock_bill_list := make([]ErpStockBillTransferViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&stock_bill_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpStockBillTransferViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  stock_bill_list,
		Total: count,
	})
}

type StockBillTransferDetailForm struct {
	StockBill models.StockBillInfo `json:"stock_bill"`
	StockList []models.StockInfo   `json:"stock_list"`
}

func (f *StockBillTransferDetailForm) IsCreate() bool {
	return f.StockBill.Id == 0
}

func (f *StockBillTransferDetailForm) Valid(validation *revel.Validation) bool {

	if f.StockBill.SrcStockId.Valid == false || f.StockBill.SrcStockId.Int64 == 0 {
		validation.Required(f.StockBill.SrcStockId).Message("调出仓库不能为空！")
	}

	validation.Required(f.StockBill.StockId).Message("调入仓库不能为空！")

	if f.StockBill.Desc.Valid == true {
		validation.MinSize(f.StockBill.Desc.String, 2).Message("备注长度不能小于2！")
	}
	if f.StockBill.Desc.Valid == true {
		validation.MaxSize(f.StockBill.Desc.String, 400).Message("备注长度不能大于400！")
	}

	validation.Required(f.StockBill.Date).Message("日期不能为空！")

	validation.Required(f.StockBill.UserId).Message("操作人不能为空！")

	return validation.HasErrors() == false
}

func (c ErpStockBillTransfer) DetailData() revel.Result {
	db_session := c.DbSession

	var stock_bill_id int64
	c.Params.Bind(&stock_bill_id, "id")

	var stock_bill models.StockBillInfo
	if stock_bill_id != 0 {
		has, err := db_session.Id(stock_bill_id).Get(&stock_bill)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的出入库单不存在！"})
		}
	} else {
		stock_bill.Date = core.Date(time.Now())
	}

	var stock_list []models.StockInfo
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	var detail_form StockBillTransferDetailForm
	detail_form.StockBill = stock_bill
	detail_form.StockList = stock_list

	return c.RenderJson(detail_form)
}

func (c ErpStockBillTransfer) Save() revel.Result {
	db_session := c.DbSession

	var detail_form StockBillTransferDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	stock_bill := detail_form.StockBill

	var affected int64
	if detail_form.IsCreate() {
		loginuser := core.GetLoginUser(c.Session)
		stock_bill.UserId = loginuser.UserId
		stock_bill.TypeId = models.STOCK_BILL_TYPE_TRANSFER
		affected, err = db_session.Insert(&stock_bill)
		core.HandleError(err)
	} else {
		sb_cols := models.StockBillCols
		affected, err = db_session.Id(stock_bill.Id).Cols(sb_cols.Code, sb_cols.SrcStockId, sb_cols.StockId, sb_cols.Desc, sb_cols.Date).Update(&stock_bill)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpStockBillTransfer) Delete() revel.Result {
	db_session := c.DbSession

	stock_bill_id_list := make([]int64, 0)
	c.Params.Bind(&stock_bill_id_list, "id_list")

	affected, err := db_session.In("id", stock_bill_id_list).Delete(new(models.StockBillInfo))
	core.HandleError(err)

	//同时删除关联明细
	/*
		DELETE FROM erp_stock_bill_entry WHERE bill_id IN (123, 456, 789);
	 */

	stock_bill_id_str_list := make([]string, 0)
	for _, stock_bill_id := range stock_bill_id_list {
		stock_bill_id_str_list = append(stock_bill_id_str_list, fmt.Sprint(stock_bill_id))
	}

	sql := `DELETE FROM erp_stock_bill_entry WHERE bill_id IN (%s);`
	sql = fmt.Sprintf(sql, stock_bill_id_str_list)
	_, err = db_session.Exec(sql)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
