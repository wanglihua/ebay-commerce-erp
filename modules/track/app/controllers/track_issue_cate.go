
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/track/models"
)

type TrackIssueCate struct {
	*revel.Controller
}

func (c TrackIssueCate) Index() revel.Result {
	return c.RenderTemplate("track/issue_cate/issue_cate_index.html")
}

func (c TrackIssueCate) ListData() revel.Result {
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

	issue_cate_list := make([]models.IssueCateInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&issue_cate_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.IssueCateInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  issue_cate_list,
		Total: count,
	})
}

type IssueCateDetailForm struct {
	IssueCate models.IssueCateInfo `json:"issue_cate"`
}

func (f *IssueCateDetailForm) IsCreate() bool {
	return f.IssueCate.Id == 0
}

func (f *IssueCateDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.IssueCate.Name).Message("名称不能为空！")
	if f.IssueCate.Name != "" {
		validation.MinSize(f.IssueCate.Name, 2).Message("名称长度不能小于2！")
	}
	if f.IssueCate.Name != "" {
		validation.MaxSize(f.IssueCate.Name, 100).Message("名称长度不能大于100！")
	}

	validation.Required(f.IssueCate.Desc).Message("描述不能为空！")
	if f.IssueCate.Desc != "" {
		validation.MinSize(f.IssueCate.Desc, 2).Message("描述长度不能小于2！")
	}
	if f.IssueCate.Desc != "" {
		validation.MaxSize(f.IssueCate.Desc, 4000).Message("描述长度不能大于4000！")
	}

	return validation.HasErrors() == false
}

func (c TrackIssueCate) DetailData() revel.Result {
	db_session := c.DbSession

	var issue_cate_id int64
	c.Params.Bind(&issue_cate_id, "id")

	var issue_cate models.IssueCateInfo
	if issue_cate_id != 0 {
		has, err := db_session.Id(issue_cate_id).Get(&issue_cate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的缺陷类别不存在！"})
		}
	}

	var detail_form IssueCateDetailForm
	detail_form.IssueCate = issue_cate

	return c.RenderJson(detail_form)
}

func (c TrackIssueCate) Save() revel.Result {
	db_session := c.DbSession

	var detail_form IssueCateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	issue_cate := detail_form.IssueCate

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&issue_cate)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(issue_cate.Id).Update(&issue_cate)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c TrackIssueCate) Delete() revel.Result {
	db_session := c.DbSession

	issue_cate_id_list := make([]int64, 0)
	c.Params.Bind(&issue_cate_id_list, "id_list")

	affected, err := db_session.In("id", issue_cate_id_list).Delete(new(models.IssueCateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

