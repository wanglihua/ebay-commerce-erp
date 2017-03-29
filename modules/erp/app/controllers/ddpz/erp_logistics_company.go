package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpLogisticsCompany struct {
	*revel.Controller
}

func (c ErpLogisticsCompany) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_company/logistics_company_index.html")
}

func (c ErpLogisticsCompany) ListData() revel.Result {
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

	logistics_company_list := make([]models.LogisticsCompanyInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&logistics_company_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.LogisticsCompanyInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_company_list,
		Total: count,
	})
}

type LogisticsCompanyDetailForm struct {
	LogisticsCompany models.LogisticsCompanyInfo `json:"logistics_company"`
}

func (f *LogisticsCompanyDetailForm) IsCreate() bool {
	return f.LogisticsCompany.Id == 0
}

func (f *LogisticsCompanyDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsCompany.Name).Message("名称不能为空！")
	if f.LogisticsCompany.Name != "" {
		validation.MinSize(f.LogisticsCompany.Name, 2).Message("名称长度不能小于2！")
	}
	if f.LogisticsCompany.Name != "" {
		validation.MaxSize(f.LogisticsCompany.Name, 100).Message("名称长度不能大于100！")
	}

	if f.LogisticsCompany.Desc.Valid == true {
		validation.MaxSize(f.LogisticsCompany.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsCompany) DetailData() revel.Result {
	db_session := c.DbSession

	var logistics_company_id int64
	c.Params.Bind(&logistics_company_id, "id")

	var logistics_company models.LogisticsCompanyInfo
	if logistics_company_id != 0 {
		has, err := db_session.Id(logistics_company_id).Get(&logistics_company)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流公司不存在！"})
		}
	}

	var detail_form LogisticsCompanyDetailForm
	detail_form.LogisticsCompany = logistics_company

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsCompany) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsCompanyDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_company := detail_form.LogisticsCompany

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where("name = ?", logistics_company.Name).Count(&logistics_company)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&logistics_company)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", logistics_company.Id, logistics_company.Name).Count(&logistics_company)
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Id(logistics_company.Id).Update(&logistics_company)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsCompany) Delete() revel.Result {
	db_session := c.DbSession

	logistics_company_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_company_id_list, "id_list")

	var mode_count struct {
		Count int64 `xorm:"'mode_count'"`
	}
	sql := "SELECT count(lm.id) AS mode_count FROM erp_logistics_mode lm WHERE lm.company_id IN (%s)"
	company_id_str_list := core.Int_list_to_str_list(logistics_company_id_list)
	sql = fmt.Sprintf(sql, strings.Join(company_id_str_list, ","))
	_, err := db_session.Sql(sql).Get(&mode_count)
	core.HandleError(err)

	if mode_count.Count != 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "数据在别处被引用，不能删除！",
		})
	}

	affected, err := db_session.In("id", logistics_company_id_list).Delete(new(models.LogisticsCompanyInfo))
	core.HandleError(err)

	//todo: 删除时要考虑别的关联

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
