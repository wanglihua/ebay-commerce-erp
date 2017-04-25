package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpStockPickCrew struct {
	*revel.Controller
}

func (c ErpStockPickCrew) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/stock_pick_crew/stock_pick_crew_index.html")
}

func (c ErpStockPickCrew) PageData() revel.Result {
	var db_session = c.DbSession

	var stock_list = make([]models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"stock_id_query": 0,
		"stock_list":     stock_list,
	})
}

func (c ErpStockPickCrew) ListData() revel.Result {
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

	stock_pick_crew_list := make([]models.StockPickCrewInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&stock_pick_crew_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.StockPickCrewInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  stock_pick_crew_list,
		Total: count,
	})
}

type StockPickCrewDetailForm struct {
	StockPickCrew models.StockPickCrewInfo  `json:"stock_pick_crew"`
	StockList     []models.StockInfo        `json:"stock_list"`
}

func (f *StockPickCrewDetailForm) IsCreate() bool {
	return f.StockPickCrew.Id == 0
}

func (f *StockPickCrewDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.StockPickCrew.StockId).Message("仓库不能为空！")
	validation.Required(f.StockPickCrew.UserId).Message("用户不能为空！")
	validation.Required(f.StockPickCrew.IsManager).Message("是否管理员不能为空！")

	return validation.HasErrors() == false
}

func (c ErpStockPickCrew) DetailData() revel.Result {
	db_session := c.DbSession

	var stock_pick_crew_id int64
	c.Params.Bind(&stock_pick_crew_id, "id")

	var stock_pick_crew models.StockPickCrewInfo
	if stock_pick_crew_id != 0 {
		has, err := db_session.Id(stock_pick_crew_id).Get(&stock_pick_crew)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的仓库拣货人员不存在！"})
		}
	}

	var stock_list = make([]models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	var detail_form StockPickCrewDetailForm
	detail_form.StockPickCrew = stock_pick_crew
	detail_form.StockList = stock_list

	return c.RenderJson(detail_form)
}

func (c ErpStockPickCrew) Save() revel.Result {
	db_session := c.DbSession

	var detail_form StockPickCrewDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	stock_pick_crew := detail_form.StockPickCrew

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&stock_pick_crew)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(stock_pick_crew.Id).Update(&stock_pick_crew)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpStockPickCrew) Delete() revel.Result {
	db_session := c.DbSession

	stock_pick_crew_id_list := make([]int64, 0)
	c.Params.Bind(&stock_pick_crew_id_list, "id_list")

	affected, err := db_session.In("id", stock_pick_crew_id_list).Delete(new(models.StockPickCrewInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
