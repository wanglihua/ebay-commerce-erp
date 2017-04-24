package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/setting"
)

type ErpLogisticsCompany struct {
	*revel.Controller
}

func (c ErpLogisticsCompany) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_COMPANY, c.Session)

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	return c.RenderTemplate("erp/ddpz/logistics_company/logistics_company_index.html")
}

func (c ErpLogisticsCompany) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_COMPANY, c.Session)

	db_session := c.DbSession

	type logitics_company_view_item_t struct {
		models.LogisticsCompanyInfo  `xorm:"extends" json:"lc"`
		models.StatusInfo            `xorm:"extends" json:"s"`
	}

	/*
SELECT lc.*, s.*
FROM erp_logistics_company lc
    INNER JOIN erp_status s ON s.id = lc.status_id
	 */

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	sb := core.NewSqlBuilder()
	sb.Select("lc.*, s.*")
	sb.Count("count(lc.id) as logistics_company_count")
	sb.From(`erp_logistics_company lc
    INNER JOIN erp_status s ON s.id = lc.status_id`)
	sb.Where(filter)

	//query extra filter here

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("lc.id ASC")
	}

	sb.Limit(offset, limit)

	logistics_company_list := make([]logitics_company_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_company_list)
	core.HandleError(err)

	var logistics_company_count struct {
		Count int64 `xorm:"'logistics_company_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&logistics_company_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_company_list,
		Total: logistics_company_count.Count,
	})
}

type LogisticsCompanyDetailForm struct {
	LogisticsCompany models.LogisticsCompanyInfo `json:"logistics_company"`
	StatusList       []models.StatusInfo         `json:"status_list"`
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
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_COMPANY, c.Session)

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
	} else {
		logistics_company.StatusId = models.LOGISTICS_COMPANY_ENABLED
	}

	var status_list = make([]models.StatusInfo, 0)
	err := db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_LOGISTICS_COMPANY)).Find(&status_list)
	core.HandleError(err)

	var detail_form LogisticsCompanyDetailForm
	detail_form.LogisticsCompany = logistics_company
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsCompany) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_COMPANY, c.Session)

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
		name_count, err := db_session.Where("name = ?", logistics_company.Name).Count(new(models.LogisticsCompanyInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		affected, err = db_session.Insert(&logistics_company)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", logistics_company.Id, logistics_company.Name).Count(new(models.LogisticsCompanyInfo))
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
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_COMPANY, c.Session)

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

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
