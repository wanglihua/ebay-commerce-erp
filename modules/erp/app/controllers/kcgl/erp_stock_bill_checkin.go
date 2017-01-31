package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/models"
	"time"
	"matrix/modules/erp/service/bill_service"
	"fmt"
)

type ErpStockBillCheckin struct {
	*revel.Controller
}

func (c ErpStockBillCheckin) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/stock_bill_checkin/stock_bill_checkin_index.html")
}

type ErpStockBillCheckinViewItem struct {
	models.StockBillInfo `xorm:"extends" json:"sb"`
	models.StockInfo     `xorm:"extends" json:"s"`
	auth_models.UserInfo `xorm:"extends" json:"u"`
	models.StatusInfo    `xorm:"extends" json:"st"`
}

func (c ErpStockBillCheckin) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT sb.*, s.*, u.*, st.*
FROM erp_stock_bill sb
    INNER JOIN erp_stock s ON sb.stock_id = s.id
    INNER JOIN auth_user u ON sb.user_id = u.id
    INNER JOIN erp_status st ON sb.status_id = st.id
WHERE sb.type_id = 100;
	*/

	sb := core.NewSqlBuilder()
	sb.Select("sb.*, s.*, u.*, st.*")
	sb.Count("count(sb.id) AS bill_count")
	sb.From(`
		erp_stock_bill sb
		    INNER JOIN erp_stock s ON sb.stock_id = s.id
		    INNER JOIN auth_user u ON sb.user_id = u.id
		    INNER JOIN erp_status st ON sb.status_id = st.id
    `)

	sb.Where(core.SqlAnd(filter, fmt.Sprintf("sb.type_id = %d", models.STOCK_BILL_TYPE_CHECKIN)))

	//query extra filter here

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("sb.id ASC")
	}

	sb.Limit(offset, limit)

	stock_bill_list := make([]ErpStockBillCheckinViewItem, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&stock_bill_list)
	core.HandleError(err)

	var bill_count struct {
		Count int64 `xorm:"'bill_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&bill_count)

	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  stock_bill_list,
		Total: bill_count.Count,
	})
}

type StockBillCheckinDetailForm struct {
	StockBill models.StockBillInfo `json:"stock_bill"`
	StockList []models.StockInfo   `json:"stock_list"`
}

func (f *StockBillCheckinDetailForm) IsCreate() bool {
	return f.StockBill.Id == 0
}

func (f *StockBillCheckinDetailForm) Valid(validation *revel.Validation) bool {

	validation.Required(f.StockBill.StockId).Message("仓库不能为空！")

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

func (c ErpStockBillCheckin) DetailData() revel.Result {
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
		stock_bill.Code = bill_service.GetNextBillNo(db_session)
		stock_bill.Date = core.Date(time.Now())
	}

	var stock_list []models.StockInfo
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	var detail_form StockBillCheckinDetailForm
	detail_form.StockBill = stock_bill
	detail_form.StockList = stock_list

	return c.RenderJson(detail_form)
}

func (c ErpStockBillCheckin) Save() revel.Result {
	db_session := c.DbSession

	var detail_form StockBillCheckinDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	var stock_bill_in_db models.StockBillInfo
	_, err = db_session.Id(detail_form.StockBill.Id).Get(&stock_bill_in_db)
	core.HandleError(err)

	if stock_bill_in_db.StatusId == models.STOCK_BILL_CHECKIN_SUBMIT {
		return c.RenderJson(core.JsonResult{Success: false, Message: "入库单已提交，不能更改！"})
	}

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	stock_bill_in_ui := detail_form.StockBill

	var affected int64
	if detail_form.IsCreate() {
		loginuser := core.GetLoginUser(c.Session)
		stock_bill_in_ui.UserId = loginuser.UserId
		stock_bill_in_ui.TypeId = models.STOCK_BILL_TYPE_CHECKIN
		stock_bill_in_ui.StatusId = models.STOCK_BILL_CHECKIN_SAVE

		affected, err = db_session.Insert(&stock_bill_in_ui)
		core.HandleError(err)
	} else {
		sb_cols := models.StockBillCols
		affected, err = db_session.Id(stock_bill_in_ui.Id).Cols(sb_cols.Code, sb_cols.StockId, sb_cols.Desc, sb_cols.Date).Update(&stock_bill_in_ui)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"bill_id": stock_bill_in_ui.Id,
		"message": strconv.FormatInt(affected, 10) + "条数据保存成功!",
	})
}

func (c ErpStockBillCheckin) Submit() revel.Result {
	db_session := c.DbSession

	var bill_id int64
	c.Params.Bind(&bill_id, "bill_id")

	var stock_bill models.StockBillInfo
	has, err := db_session.Id(bill_id).Get(&stock_bill)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "入库单Id不存在!"})
	}

	if stock_bill.StatusId == models.STOCK_BILL_CHECKIN_SUBMIT {
		return c.RenderJson(core.JsonResult{Success: false, Message: "入库失败，该入库单已入库!"})
	}

	//todo: 检查明细是不是空的

	stock_bill.StatusId = models.STOCK_BILL_CHECKIN_SUBMIT

	cols := models.StockBillCols
	_, err = db_session.Id(bill_id).Cols(cols.StatusId).Update(&stock_bill)
	core.HandleError(err)

	// 更新商品库存
	sql_proc := `
CREATE PROCEDURE erp_stock_bill_checkin_to_inventory_v0001(p_bill_id BIGINT)
    BEGIN
        DECLARE v_stock_id BIGINT;
        DECLARE v_storage_loc_id BIGINT;
        DECLARE v_product_id BIGINT;
        DECLARE v_qty BIGINT;
        DECLARE v_inventory_id BIGINT;
        DECLARE v_done INT;

        DECLARE cur_stock_bill_entry CURSOR FOR
            SELECT storage_loc_id, product_id, qty
            FROM erp_stock_bill_entry
            WHERE bill_id = p_bill_id;
        DECLARE CONTINUE HANDLER FOR SQLSTATE '02000' SET v_done = 1;

        SELECT stock_id INTO v_stock_id FROM erp_stock_bill WHERE id = p_bill_id;

        OPEN cur_stock_bill_entry;
        REPEAT
            FETCH cur_stock_bill_entry INTO v_storage_loc_id, v_product_id, v_qty;

            SELECT id INTO v_inventory_id
            FROM erp_inventory
            WHERE stock_id = v_stock_id AND storage_loc_id = v_storage_loc_id AND product_id = v_product_id;

            IF v_inventory_id IS NULL THEN
                INSERT INTO erp_inventory (stock_id, storage_loc_id, product_id, qty, create_time, update_time)
                VALUES (v_stock_id, v_storage_loc_id, v_product_id, v_qty, now(), now());
            ELSE
                UPDATE erp_inventory SET qty = qty + v_qty, update_time = now(), version = version + 1 WHERE id = v_inventory_id;
            END IF;
        UNTIL v_done END REPEAT;
        CLOSE cur_stock_bill_entry;
    END;
`
	core.CreateProc("erp_stock_bill_checkin_to_inventory_v0001", sql_proc, db_session)

	db_session.Exec(fmt.Sprintf("CALL erp_stock_bill_checkin_to_inventory_v0001(%d);", bill_id))

	return c.RenderJson(core.JsonResult{Success: true, Message: "入库成功!"})
}

func (c ErpStockBillCheckin) Delete() revel.Result {
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
