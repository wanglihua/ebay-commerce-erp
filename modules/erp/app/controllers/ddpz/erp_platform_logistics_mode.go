
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpPlatformLogisticsMode struct {
	*revel.Controller
}

func (c ErpPlatformLogisticsMode) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/platform_logistics_mode/platform_logistics_mode_index.html")
}

func (c ErpPlatformLogisticsMode) ListData() revel.Result {
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

	platform_logistics_mode_list := make([]models.PlatformLogisticsModeInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&platform_logistics_mode_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.PlatformLogisticsModeInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  platform_logistics_mode_list,
		Total: count,
	})
}

type PlatformLogisticsModeDetailForm struct {
	PlatformLogisticsMode models.PlatformLogisticsModeInfo `json:"platform_logistics_mode"`
}

func (f *PlatformLogisticsModeDetailForm) IsCreate() bool {
	return f.PlatformLogisticsMode.Id == 0
}

func (f *PlatformLogisticsModeDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.PlatformLogisticsMode.PlatformId).Message("电商平台不能为空！")

	validation.Required(f.PlatformLogisticsMode.Name).Message("名称不能为空！")
	if f.PlatformLogisticsMode.Name != "" {
		validation.MinSize(f.PlatformLogisticsMode.Name, 2).Message("名称长度不能小于2！")
	}
	if f.PlatformLogisticsMode.Name != "" {
		validation.MaxSize(f.PlatformLogisticsMode.Name, 100).Message("名称长度不能大于100！")
	}

	if f.PlatformLogisticsMode.Desc.Valid == true {
		validation.MaxSize(f.PlatformLogisticsMode.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpPlatformLogisticsMode) DetailData() revel.Result {
	db_session := c.DbSession

	var platform_logistics_mode_id int64
	c.Params.Bind(&platform_logistics_mode_id, "id")

	var platform_logistics_mode models.PlatformLogisticsModeInfo
	if platform_logistics_mode_id != 0 {
		has, err := db_session.Id(platform_logistics_mode_id).Get(&platform_logistics_mode)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的电商平台物流方式不存在！"})
		}
	}

	var detail_form PlatformLogisticsModeDetailForm
	detail_form.PlatformLogisticsMode = platform_logistics_mode

	return c.RenderJson(detail_form)
}

func (c ErpPlatformLogisticsMode) Save() revel.Result {
	db_session := c.DbSession

	var detail_form PlatformLogisticsModeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	platform_logistics_mode := detail_form.PlatformLogisticsMode

	var affected int64
	if detail_form.IsCreate() { 
		name_count, err := db_session.Where("name = ?", platform_logistics_mode.Name).Count(&platform_logistics_mode)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&platform_logistics_mode)
		core.HandleError(err)
	} else { 
		name_count, err := db_session.Where("id <> ? and name = ?", platform_logistics_mode.Id, platform_logistics_mode.Name).Count(&platform_logistics_mode)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Id(platform_logistics_mode.Id).Update(&platform_logistics_mode)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPlatformLogisticsMode) Delete() revel.Result {
	db_session := c.DbSession

	platform_logistics_mode_id_list := make([]int64, 0)
	c.Params.Bind(&platform_logistics_mode_id_list, "id_list")

	affected, err := db_session.In("id", platform_logistics_mode_id_list).Delete(new(models.PlatformLogisticsModeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

