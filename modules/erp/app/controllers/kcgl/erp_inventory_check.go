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
)

type ErpInventoryCheck struct {
	*revel.Controller
}

func (c ErpInventoryCheck) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/inventory_check/inventory_check_index.html")
}

type ErpInventoryCheckViewItem struct {
	models.InventoryCheckInfo `xorm:"extends" json:"ic"`
	models.StockInfo          `xorm:"extends" json:"st"`
	CheckUser   auth_models.UserInfo `xorm:"extends" json:"cu"`
	ApproveUser auth_models.UserInfo `xorm:"extends" json:"au"`
	models.StatusInfo         `xorm:"extends" json:"ss"`
}

func (c ErpInventoryCheck) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT ic.*, st.*, cu.*, au.*, ss.*
		FROM erp_inventory_check ic.
		    INNER JOIN erp_stock st ON ic.stock_id = st.id
		    INNER JOIN auth_user cu ON ic.check_user_id = cu.id
		    LEFT JOIN auth_user au ON ic.approve_user_id = au.id
		    INNER JOIN erp_status ss ON ic.status_id = ss.id
	*/

	query := db_session.
	Select("ic.*, st.*, cu.*, au.*, ss.*").
		Table("erp_inventory_check").Alias("ic").
		Join("INNER", []string{"erp_stock", "st"}, "ic.stock_id = st.id").
		Join("INNER", []string{"auth_user", "cu"}, "ic.check_user_id = cu.id").
		Join("LEFT", []string{"auth_user", "au"}, "ic.approve_user_id = au.id").
		Join("INNER", []string{"erp_status", "ss"}, "ic.status_id = ss.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("ic.id")
	}

	inventory_check_list := make([]ErpInventoryCheckViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&inventory_check_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpInventoryCheckViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  inventory_check_list,
		Total: count,
	})
}

type InventoryCheckDetailForm struct {
	InventoryCheck models.InventoryCheckInfo `json:"inventory_check"`
	StockList      []models.StockInfo        `json:"stock_list"`
}

func (f *InventoryCheckDetailForm) IsCreate() bool {
	return f.InventoryCheck.Id == 0
}

func (f *InventoryCheckDetailForm) Valid(validation *revel.Validation) bool {

	validation.Required(f.InventoryCheck.StockId).Message("仓库不能为空！")

	if f.InventoryCheck.Desc.Valid == true {
		validation.MinSize(f.InventoryCheck.Desc.String, 2).Message("备注长度不能小于2！")
	}
	if f.InventoryCheck.Desc.Valid == true {
		validation.MaxSize(f.InventoryCheck.Desc.String, 400).Message("备注长度不能大于400！")
	}

	validation.Required(f.InventoryCheck.CheckDate).Message("盘点日期不能为空！")

	validation.Required(f.InventoryCheck.CheckUserId).Message("盘点人不能为空！")

	validation.Required(f.InventoryCheck.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpInventoryCheck) DetailData() revel.Result {
	db_session := c.DbSession

	var inventory_check_id int64
	c.Params.Bind(&inventory_check_id, "id")

	var inventory_check models.InventoryCheckInfo
	if inventory_check_id != 0 {
		has, err := db_session.Id(inventory_check_id).Get(&inventory_check)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的库存盘点不存在！"})
		}
	} else {
		inventory_check.Code = bill_service.GetNextBillNo(db_session)
		inventory_check.CheckDate = core.Date(time.Now())
	}

	var stock_list []models.StockInfo
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	var detail_form InventoryCheckDetailForm
	detail_form.InventoryCheck = inventory_check
	detail_form.StockList = stock_list

	return c.RenderJson(detail_form)
}

func (c ErpInventoryCheck) Save() revel.Result {
	db_session := c.DbSession

	var detail_form InventoryCheckDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	inventory_check := detail_form.InventoryCheck

	var affected int64
	if detail_form.IsCreate() {
		loginuser := core.GetLoginUser(c.Session)
		inventory_check.CheckUserId = loginuser.UserId
		inventory_check.StatusId = models.INVENTORY_CHECK_STATUS_YCJ
		affected, err = db_session.Insert(&inventory_check)
		core.HandleError(err)
	} else {
		if inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_YCJ && inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_WTG {
			return c.RenderJson(core.JsonResult{Success: false, Message: "当前盘点单不是已创建或未通过状态，信息不能修改！"})
		}

		ic_cols := models.InventoryCheckCols
		affected, err = db_session.Id(inventory_check.Id).Cols(ic_cols.Code, ic_cols.StockId, ic_cols.Desc, ic_cols.CheckDate).Update(&inventory_check)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpInventoryCheck) Delete() revel.Result {
	db_session := c.DbSession

	inventory_check_id_list := make([]int64, 0)
	c.Params.Bind(&inventory_check_id_list, "id_list")

	_, err := db_session.In("check_id", inventory_check_id_list).Delete(new(models.InventoryCheckEntryInfo))
	core.HandleError(err)

	_, err = db_session.In("id", inventory_check_id_list).Delete(new(models.InventoryCheckInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据删除成功!"})
}

func (c ErpInventoryCheck) Submit() revel.Result {
	db_session := c.DbSession

	var check_id int64
	c.Params.Bind(&check_id, "check_id")

	var inventory_check models.InventoryCheckInfo
	has, err := db_session.Id(check_id).Get(&inventory_check)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的盘点单不存在!"})
	}

	if inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_YCJ && inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_WTG {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前盘点单状态不是已创建或未通过，不能提交!"})
	}

	inventory_check.StatusId = models.INVENTORY_CHECK_STATUS_SPZ

	_, err = db_session.Id(check_id).Cols("status_id").Update(&inventory_check)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据提交成功!"})
}

func (c ErpInventoryCheck) Back() revel.Result {
	db_session := c.DbSession

	var check_id int64
	c.Params.Bind(&check_id, "check_id")

	var inventory_check models.InventoryCheckInfo
	has, err := db_session.Id(check_id).Get(&inventory_check)
	core.HandleError(err)
	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "指定的盘点单不存在!"})
	}

	if inventory_check.StatusId != models.INVENTORY_CHECK_STATUS_SPZ {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前盘点单状态不是审批中，不能提交!"})
	}

	inventory_check.StatusId = models.INVENTORY_CHECK_STATUS_YCJ

	_, err = db_session.Id(check_id).Cols("status_id").Update(&inventory_check)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据回退成功!"})
}
