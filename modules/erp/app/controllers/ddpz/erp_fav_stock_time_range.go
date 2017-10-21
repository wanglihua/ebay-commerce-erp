package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/ddpz_models"
)

type ErpFavStockTimeRange struct {
	*revel.Controller
}

func (c ErpFavStockTimeRange) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/fav_stock_time_range/fav_stock_time_range_index.html")
}

func (c ErpFavStockTimeRange) ListData() revel.Result {
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

	fav_stock_time_range_list := make([]ddpz_models.FavStockTimeRangeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&fav_stock_time_range_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ddpz_models.FavStockTimeRangeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  fav_stock_time_range_list,
		Total: count,
	})
}

type FavStockTimeRangeDetailForm struct {
	FavStockTimeRange ddpz_models.FavStockTimeRangeInfo `json:"fav_stock_time_range"`
}

func (f *FavStockTimeRangeDetailForm) IsCreate() bool {
	return f.FavStockTimeRange.Id == 0
}

func (f *FavStockTimeRangeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.FavStockTimeRange.Weekday).Message("星期不能为空！")

	validation.Required(f.FavStockTimeRange.BeginHour).Message("开始小时不能为空！")

	validation.Required(f.FavStockTimeRange.BeginMinute).Message("开始分钟不能为空！")

	validation.Required(f.FavStockTimeRange.EndHour).Message("结束小时不能为空！")

	validation.Required(f.FavStockTimeRange.EndMinute).Message("结束分钟不能为空！")

	return validation.HasErrors() == false
}

func (c ErpFavStockTimeRange) DetailData() revel.Result {
	db_session := c.DbSession

	var fav_stock_time_range_id int64
	c.Params.Bind(&fav_stock_time_range_id, "id")

	var fav_stock_time_range ddpz_models.FavStockTimeRangeInfo
	if fav_stock_time_range_id != 0 {
		has, err := db_session.Id(fav_stock_time_range_id).Get(&fav_stock_time_range)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的优先仓库时间段不存在！"})
		}
	}

	var detail_form FavStockTimeRangeDetailForm
	detail_form.FavStockTimeRange = fav_stock_time_range

	return c.RenderJson(detail_form)
}

func (c ErpFavStockTimeRange) Save() revel.Result {
	db_session := c.DbSession

	var detail_form FavStockTimeRangeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	fav_stock_time_range := detail_form.FavStockTimeRange

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&fav_stock_time_range)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(fav_stock_time_range.Id).Update(&fav_stock_time_range)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpFavStockTimeRange) Delete() revel.Result {
	db_session := c.DbSession

	fav_stock_time_range_id_list := make([]int64, 0)
	c.Params.Bind(&fav_stock_time_range_id_list, "id_list")

	affected, err := db_session.In("id", fav_stock_time_range_id_list).Delete(new(ddpz_models.FavStockTimeRangeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
