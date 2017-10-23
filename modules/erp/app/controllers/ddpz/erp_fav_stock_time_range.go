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

	for index, time_range := range fav_stock_time_range_list {
		if time_range.Weekday == 1 {
			fav_stock_time_range_list[index].WeekdayName = "星期一"
		}

		if time_range.Weekday == 2 {
			fav_stock_time_range_list[index].WeekdayName = "星期二"
		}

		if time_range.Weekday == 3 {
			fav_stock_time_range_list[index].WeekdayName = "星期三"
		}

		if time_range.Weekday == 4 {
			fav_stock_time_range_list[index].WeekdayName = "星期四"
		}

		if time_range.Weekday == 5 {
			fav_stock_time_range_list[index].WeekdayName = "星期五"
		}

		if time_range.Weekday == 6 {
			fav_stock_time_range_list[index].WeekdayName = "星期六"
		}

		if time_range.Weekday == 7 {
			fav_stock_time_range_list[index].WeekdayName = "星期日"
		}
	}

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
	} else {
		fav_stock_time_range.Weekday = 1
	}

	var detail_form FavStockTimeRangeDetailForm
	detail_form.FavStockTimeRange = fav_stock_time_range

	type value_name_item struct {
		Value int64  `json:"value"`
		Name  string `json:"name"`
	}

	var weekday_list = make([]value_name_item, 0)
	weekday_list = append(weekday_list, value_name_item{Value: 1, Name: "星期一"})
	weekday_list = append(weekday_list, value_name_item{Value: 2, Name: "星期二"})
	weekday_list = append(weekday_list, value_name_item{Value: 3, Name: "星期三"})
	weekday_list = append(weekday_list, value_name_item{Value: 4, Name: "星期四"})
	weekday_list = append(weekday_list, value_name_item{Value: 5, Name: "星期五"})
	weekday_list = append(weekday_list, value_name_item{Value: 6, Name: "星期六"})
	weekday_list = append(weekday_list, value_name_item{Value: 7, Name: "星期日"})

	var hour_list = make([]int64, 0)
	for i := 0; i < 24; i++ {
		hour_list = append(hour_list, int64(i))
	}

	var minute_list = make([]int64, 0)
	for i := 0; i < 60; i++ {
		minute_list = append(minute_list, int64(i))
	}

	return c.RenderJson(map[string]interface{}{
		"detail":       detail_form,
		"weekday_list": weekday_list,
		"hour_list":    hour_list,
		"minute_list":  minute_list,
	})
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

	var begin = fav_stock_time_range.BeginHour*60 + fav_stock_time_range.BeginMinute
	var end = fav_stock_time_range.EndHour*60 + fav_stock_time_range.EndMinute
	if begin >= end {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "时间段开始应该早于时间段结束！",
		})
	}

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
