package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpIntelRepl struct {
	*revel.Controller
}

func (c ErpIntelRepl) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	return c.RenderTemplate("erp/cggl/intel_repl/intel_repl_index.html")
}

func (c ErpIntelRepl) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

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

	intel_repl_list := make([]cggl_models.IntelReplInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&intel_repl_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(cggl_models.IntelReplInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  intel_repl_list,
		Total: count,
	})
}

type IntelReplDetailForm struct {
	IntelRepl cggl_models.IntelReplInfo `json:"intel_repl"`
}

func (f *IntelReplDetailForm) IsCreate() bool {
	return f.IntelRepl.Id == 0
}

func (f *IntelReplDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.IntelRepl.ProductId).Message("商品不能为空！")

	validation.Required(f.IntelRepl.StockId).Message("仓库不能为空！")

	validation.Required(f.IntelRepl.Qty).Message("数量不能为空！")

	validation.Required(f.IntelRepl.Profit).Message("毛利率不能为空！")

	validation.Required(f.IntelRepl.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpIntelRepl) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	var intel_repl_id int64
	c.Params.Bind(&intel_repl_id, "id")

	var intel_repl cggl_models.IntelReplInfo
	if intel_repl_id != 0 {
		has, err := db_session.Id(intel_repl_id).Get(&intel_repl)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的智能补货不存在！"})
		}
	}

	var detail_form IntelReplDetailForm
	detail_form.IntelRepl = intel_repl

	return c.RenderJson(detail_form)
}

func (c ErpIntelRepl) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	var detail_form IntelReplDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	intel_repl := detail_form.IntelRepl

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&intel_repl)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(intel_repl.Id).Update(&intel_repl)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpIntelRepl) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.CGGL_INTEL_REPL, c.Session)

	db_session := c.DbSession

	intel_repl_id_list := make([]int64, 0)
	c.Params.Bind(&intel_repl_id_list, "id_list")

	affected, err := db_session.In("id", intel_repl_id_list).Delete(new(cggl_models.IntelReplInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
