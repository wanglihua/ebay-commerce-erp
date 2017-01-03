
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpProductLog struct {
	*revel.Controller
}

func (c ErpProductLog) Index() revel.Result {
	return c.RenderTemplate("erp/product_log/product_log_index.html")
}

func (c ErpProductLog) ListData() revel.Result {
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

	product_log_list := make([]models.ProductLogInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_log_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_log_list,
		Total: count,
	})
}

type ProductLogDetailForm struct {
	ProductLog models.ProductLogInfo `json:"product_log"`
}

func (f *ProductLogDetailForm) IsCreate() bool {
	return f.ProductLog.Id == 0
}

func (f *ProductLogDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.ProductLog.ProductId).Message("商品Id不能为空！")

	validation.Required(f.ProductLog.ProductCode).Message("商品编号不能为空！")

	validation.Required(f.ProductLog.UserId).Message("用户Id不能为空！")

	validation.Required(f.ProductLog.LoginName).Message("用户工号不能为空！")

	validation.Required(f.ProductLog.Operation).Message("操作不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductLog) DetailData() revel.Result {
	db_session := c.DbSession

	var product_log_id int64
	c.Params.Bind(&product_log_id, "id")

	var product_log models.ProductLogInfo
	if product_log_id != 0 {
		has, err := db_session.Id(product_log_id).Get(&product_log)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品修改日志不存在！"})
		}
	}

	var detail_form ProductLogDetailForm
	detail_form.ProductLog = product_log

	return c.RenderJson(detail_form)
}

func (c ErpProductLog) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductLogDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_log := detail_form.ProductLog

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&product_log)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(product_log.Id).Update(&product_log)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductLog) Delete() revel.Result {
	db_session := c.DbSession

	product_log_id_list := make([]int64, 0)
	c.Params.Bind(&product_log_id_list, "id_list")

	affected, err := db_session.In("id", product_log_id_list).Delete(new(models.ProductLogInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

