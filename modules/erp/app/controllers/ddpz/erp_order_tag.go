package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpOrderTag struct {
	*revel.Controller
}

func (c ErpOrderTag) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/order_tag/order_tag_index.html")
}

func (c ErpOrderTag) ListData() revel.Result {
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

	order_tag_list := make([]models.OrderTagInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&order_tag_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.OrderTagInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  order_tag_list,
		Total: count,
	})
}

type OrderTagDetailForm struct {
	OrderTag models.OrderTagInfo `json:"order_tag"`
}

func (f *OrderTagDetailForm) IsCreate() bool {
	return f.OrderTag.Id == 0
}

func (f *OrderTagDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.OrderTag.Name).Message("名称不能为空！")
	if f.OrderTag.Name != "" {
		validation.MinSize(f.OrderTag.Name, 2).Message("名称长度不能小于2！")
	}
	if f.OrderTag.Name != "" {
		validation.MaxSize(f.OrderTag.Name, 100).Message("名称长度不能大于100！")
	}

	validation.Required(f.OrderTag.BgColor).Message("背景色不能为空！")
	if f.OrderTag.BgColor != "" {
		validation.MaxSize(f.OrderTag.BgColor, 500).Message("背景色长度不能大于500！")
	}

	validation.Required(f.OrderTag.Desc).Message("备注不能为空！")
	if f.OrderTag.Desc != "" {
		validation.MaxSize(f.OrderTag.Desc, 500).Message("备注长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpOrderTag) DetailData() revel.Result {
	db_session := c.DbSession

	var order_tag_id int64
	c.Params.Bind(&order_tag_id, "id")

	var order_tag models.OrderTagInfo
	if order_tag_id != 0 {
		has, err := db_session.Id(order_tag_id).Get(&order_tag)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的订单标签不存在！"})
		}
	}

	var detail_form OrderTagDetailForm
	detail_form.OrderTag = order_tag

	return c.RenderJson(detail_form)
}

func (c ErpOrderTag) Save() revel.Result {
	db_session := c.DbSession

	var detail_form OrderTagDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	order_tag := detail_form.OrderTag

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where("name = ?", order_tag.Name).Count(new(models.OrderTagInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&order_tag)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", order_tag.Id, order_tag.Name).Count(new(models.OrderTagInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Id(order_tag.Id).Update(&order_tag)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpOrderTag) Delete() revel.Result {
	db_session := c.DbSession

	order_tag_id_list := make([]int64, 0)
	c.Params.Bind(&order_tag_id_list, "id_list")

	affected, err := db_session.In("id", order_tag_id_list).Delete(new(models.OrderTagInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
