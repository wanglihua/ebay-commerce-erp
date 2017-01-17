package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"strconv"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
)

type ErpPaypalAccount struct {
	*revel.Controller
}

func (c ErpPaypalAccount) Index() revel.Result {
	return c.RenderTemplate("erp/xtpz/paypal_account/paypal_account_index.html")
}

type ErpPaypalAccountViewItem struct {
	models.PaypalAccountInfo `xorm:"extends" json:"a"`
	models.StatusInfo        `xorm:"extends" json:"s"`
}

func (c ErpPaypalAccount) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT a.*, s.* FROM erp_paypal_account a
		    INNER JOIN erp_status s ON a.status_id = s.id
	*/

	query := db_session.
		Select("a.*, s.*").
		Table("erp_paypal_account").Alias("a").
		Join("INNER", []string{"erp_status", "s"}, "a.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("a.id")
	}

	paypal_account_list := make([]ErpPaypalAccountViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&paypal_account_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpPaypalAccountViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  paypal_account_list,
		Total: count,
	})
}

type PaypalAccountDetailForm struct {
	PaypalAccount models.PaypalAccountInfo `json:"paypal_account"`
	StatusList    []models.StatusInfo      `json:"status_list"`
}

func (f *PaypalAccountDetailForm) IsCreate() bool {
	return f.PaypalAccount.Id == 0
}

func (f *PaypalAccountDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.PaypalAccount.PaypalAccount).Message("Paypal账号不能为空！")
	if f.PaypalAccount.PaypalAccount != "" {
		validation.MinSize(f.PaypalAccount.PaypalAccount, 2).Message("Paypal账号长度不能小于2！")
	}
	if f.PaypalAccount.PaypalAccount != "" {
		validation.MaxSize(f.PaypalAccount.PaypalAccount, 100).Message("Paypal账号长度不能大于100！")
	}

	validation.Required(f.PaypalAccount.ApiAccount).Message("API账号不能为空！")
	if f.PaypalAccount.ApiAccount != "" {
		validation.MinSize(f.PaypalAccount.ApiAccount, 2).Message("API账号长度不能小于2！")
	}
	if f.PaypalAccount.ApiAccount != "" {
		validation.MaxSize(f.PaypalAccount.ApiAccount, 100).Message("API账号长度不能大于100！")
	}

	validation.Required(f.PaypalAccount.ApiToken).Message("API令牌不能为空！")
	if f.PaypalAccount.ApiToken != "" {
		validation.MinSize(f.PaypalAccount.ApiToken, 2).Message("API令牌长度不能小于2！")
	}
	if f.PaypalAccount.ApiToken != "" {
		validation.MaxSize(f.PaypalAccount.ApiToken, 100).Message("API令牌长度不能大于100！")
	}

	if f.PaypalAccount.Desc.Valid == true {
		validation.MaxSize(f.PaypalAccount.Desc.String, 500).Message("描述长度不能大于500！")
	}

	validation.Required(f.PaypalAccount.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPaypalAccount) DetailData() revel.Result {
	db_session := c.DbSession

	var paypal_account_id int64
	c.Params.Bind(&paypal_account_id, "id")

	var paypal_account models.PaypalAccountInfo
	if paypal_account_id != 0 {
		has, err := db_session.Id(paypal_account_id).Get(&paypal_account)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的Paypal账号不存在！"})
		}
	}

	status_list := status_service.GetStatusList(db_session, models.STATUS_TYPE_PAYPAL_ACCOUNT)

	var detail_form PaypalAccountDetailForm
	detail_form.PaypalAccount = paypal_account
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpPaypalAccount) Save() revel.Result {
	db_session := c.DbSession

	var detail_form PaypalAccountDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	paypal_account := detail_form.PaypalAccount

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&paypal_account)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(paypal_account.Id).Update(&paypal_account)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPaypalAccount) Delete() revel.Result {
	db_session := c.DbSession

	paypal_account_id_list := make([]int64, 0)
	c.Params.Bind(&paypal_account_id_list, "id_list")

	affected, err := db_session.In("id", paypal_account_id_list).Delete(new(models.PaypalAccountInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
