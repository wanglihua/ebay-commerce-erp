package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpLogisticsFeeRate struct {
	*revel.Controller
}

func (c ErpLogisticsFeeRate) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_fee_rate/logistics_fee_rate_index.html")
}

func (c ErpLogisticsFeeRate) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type logistics_fee_rate_view_item_t struct {
		models.LogisticsFeeRateInfo  `xorm:"extends" json:"fr"`
		models.LogisticsModeInfo     `xorm:"extends" json:"lm"`
		models.LogisticsCompanyInfo  `xorm:"extends" json:"lc"`
	}

	/*
		SELECT fr.*, lm.*, lc.*
		FROM erp_logistics_fee_rate fr
		    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
		    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
	 */

	sb := core.NewSqlBuilder()
	sb.Select("fr.*, lm.*, lc.*")
	sb.Count("count(fr.id) AS fee_rate_count")
	sb.From(`erp_logistics_fee_rate fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id`)
	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("fr.id ASC")
	}

	sb.Limit(offset, limit)

	logistics_fee_rate_list := make([]logistics_fee_rate_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_fee_rate_list)
	core.HandleError(err)

	var fee_rate_count struct {
		Count int64 `xorm:"'fee_rate_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&fee_rate_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_fee_rate_list,
		Total: fee_rate_count.Count,
	})
}

type LogisticsFeeRateDetailForm struct {
	LogisticsFeeRate     models.LogisticsFeeRateInfo   `json:"logistics_fee_rate"`
	LogisticsCompanyList []models.LogisticsCompanyInfo `json:"logistics_company_list"`
	LogisticsModeList    []models.LogisticsModeInfo    `json:"logistics_mode_list"`
	CompanyId            int64                         `json:"company_id"`
}

func (f *LogisticsFeeRateDetailForm) IsCreate() bool {
	return f.LogisticsFeeRate.Id == 0
}

func (f *LogisticsFeeRateDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsFeeRate.ModeId).Message("物流公司不能为空！")

	validation.Required(f.LogisticsFeeRate.Postcode).Message("邮编不能为空！")
	if f.LogisticsFeeRate.Postcode != "" {
		validation.MinSize(f.LogisticsFeeRate.Postcode, 2).Message("邮编长度不能小于2！")
	}
	if f.LogisticsFeeRate.Postcode != "" {
		validation.MaxSize(f.LogisticsFeeRate.Postcode, 100).Message("邮编长度不能大于100！")
	}

	validation.Required(f.LogisticsFeeRate.Qbj).Message("起步价不能为空！")

	validation.Required(f.LogisticsFeeRate.Qbjzl).Message("起步价重量不能为空！")

	validation.Required(f.LogisticsFeeRate.Zdsf).Message("最低收费不能为空！")

	validation.Required(f.LogisticsFeeRate.Xzjt).Message("续重阶梯不能为空！")

	validation.Required(f.LogisticsFeeRate.Xzjtsf).Message("续重阶梯收费不能为空！")

	validation.Required(f.LogisticsFeeRate.Pydqjsyf).Message("偏远地区加收运费不能为空！")

	if f.LogisticsFeeRate.Desc.Valid == true {
		validation.MaxSize(f.LogisticsFeeRate.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsFeeRate) DetailData() revel.Result {
	db_session := c.DbSession

	var logistics_fee_rate_id int64
	c.Params.Bind(&logistics_fee_rate_id, "id")

	var logistics_fee_rate models.LogisticsFeeRateInfo
	if logistics_fee_rate_id != 0 {
		has, err := db_session.Id(logistics_fee_rate_id).Get(&logistics_fee_rate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流费率不存在！"})
		}
	}

	var detail_form LogisticsFeeRateDetailForm
	detail_form.LogisticsFeeRate = logistics_fee_rate

	var logistics_company_list = make([]models.LogisticsCompanyInfo, 0)
	err := db_session.Find(&logistics_company_list)
	core.HandleError(err)
	detail_form.LogisticsCompanyList = logistics_company_list

	if logistics_fee_rate_id == 0 {
		if len(logistics_company_list) != 0 {
			var selected_company_id = logistics_company_list[0].Id;
			var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
			err = db_session.Where(fmt.Sprintf("company_id = %d", selected_company_id)).Find(&logistics_mode_list)
			core.HandleError(err)

			detail_form.CompanyId = selected_company_id
			detail_form.LogisticsModeList = logistics_mode_list
		}
	} else {
		//get company id from database
		var mode_id = logistics_fee_rate.ModeId
		var logistics_mode models.LogisticsModeInfo
		_, err = db_session.Id(mode_id).Get(&logistics_mode)
		core.HandleError(err)
		var company_id = logistics_mode.CompanyId

		//get mode list from database
		var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
		err = db_session.Where(fmt.Sprintf("company_id = %d", company_id)).Find(&logistics_mode_list)
		core.HandleError(err)

		detail_form.CompanyId = company_id
		detail_form.LogisticsModeList = logistics_mode_list
	}

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsFeeRate) Save() revel.Result {
	db_session := c.DbSession

	var detail_form LogisticsFeeRateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	revel.TRACE.Println("postcode", detail_form.LogisticsFeeRate.Postcode)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_fee_rate := detail_form.LogisticsFeeRate

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&logistics_fee_rate)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(logistics_fee_rate.Id).Update(&logistics_fee_rate)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsFeeRate) Delete() revel.Result {
	db_session := c.DbSession

	logistics_fee_rate_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_fee_rate_id_list, "id_list")

	affected, err := db_session.In("id", logistics_fee_rate_id_list).Delete(new(models.LogisticsFeeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpLogisticsFeeRate) Import() revel.Result {
	return c.RenderEmpty()
}

func (c ErpLogisticsFeeRate) GetLogisticsModeListByCompanyId() revel.Result {
	db_session := c.DbSession

	var company_id int64
	c.Params.Bind(&company_id, "company_id")

	var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
	err := db_session.Where(fmt.Sprintf("company_id = %d", company_id)).Find(&logistics_mode_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"mode_list": logistics_mode_list,
	})
}
