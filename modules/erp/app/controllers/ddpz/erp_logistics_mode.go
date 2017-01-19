package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
)

type ErpLogisticsMode struct {
	*revel.Controller
}

func (c ErpLogisticsMode) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_mode/logistics_mode_index.html")
}

type ErpLogisticsModeViewItem struct {
	models.LogisticsModeInfo               `xorm:"extends" json:"m"`
	models.PlatformLogisticsModeInfo       `xorm:"extends" json:"pm"`
}

func (c ErpLogisticsMode) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT m.*, pm.*
		FROM erp_logistics_mode m
		    INNER JOIN erp_platform_logistics_mode pm ON m.ebay_mode_id = pm.id
	*/

	query := db_session.
	Select("m.*, pm.*").
		Table("erp_logistics_mode").Alias("m").
		Join("INNER", []string{"erp_platform_logistics_mode", "pm"}, "m.ebay_mode_id = pm.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("m.id")
	}

	logistics_mode_list := make([]ErpLogisticsModeViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&logistics_mode_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpLogisticsModeViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_mode_list,
		Total: count,
	})
}

type LogisticsModeDetailForm struct {
	LogisticsMode         models.LogisticsModeInfo            `json:"logistics_mode"`
	EbayLogisticsModeList []models.PlatformLogisticsModeInfo  `json:"ebay_logistics_mode_list"`
}

func (f *LogisticsModeDetailForm) IsCreate() bool {
	return f.LogisticsMode.Id == 0
}

func (f *LogisticsModeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsMode.EbayModeId).Message("EBay物流方式不能为空！")

	validation.Required(f.LogisticsMode.Name).Message("名称不能为空！")
	if f.LogisticsMode.Name != "" {
		validation.MinSize(f.LogisticsMode.Name, 2).Message("名称长度不能小于2！")
	}
	if f.LogisticsMode.Name != "" {
		validation.MaxSize(f.LogisticsMode.Name, 100).Message("名称长度不能大于100！")
	}

	if f.LogisticsMode.Desc.Valid == true {
		validation.MaxSize(f.LogisticsMode.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsMode) DetailData() revel.Result {
	db_session := c.DbSession

	var logistics_mode_id int64
	c.Params.Bind(&logistics_mode_id, "id")

	var logistics_mode models.LogisticsModeInfo
	if logistics_mode_id != 0 {
		has, err := db_session.Id(logistics_mode_id).Get(&logistics_mode)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流方式不存在！"})
		}
	}

	sql := `
SELECT pm.*
FROM erp_platform_logistics_mode pm
WHERE pm.platform_id = %d
ORDER BY pm.id ASC;
	`
	sql = fmt.Sprintf(sql, models.PLATFORM_ID_EBAY) // models.PLATFORM_ID_EBAY == 100
	var ebay_logistics_mode_list = make([]models.PlatformLogisticsModeInfo, 0)
	err := db_session.Sql(sql).Find(&ebay_logistics_mode_list)
	core.HandleError(err)

	var detail_form LogisticsModeDetailForm
	detail_form.LogisticsMode = logistics_mode
	detail_form.EbayLogisticsModeList = ebay_logistics_mode_list

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsMode) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsModeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_mode := detail_form.LogisticsMode

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where("name = ?", logistics_mode.Name).Count(&logistics_mode)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&logistics_mode)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", logistics_mode.Id, logistics_mode.Name).Count(&logistics_mode)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Id(logistics_mode.Id).Update(&logistics_mode)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsMode) Delete() revel.Result {
	db_session := c.DbSession

	logistics_mode_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_mode_id_list, "id_list")

	affected, err := db_session.In("id", logistics_mode_id_list).Delete(new(models.LogisticsModeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
