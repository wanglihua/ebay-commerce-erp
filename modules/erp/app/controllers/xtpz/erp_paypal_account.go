package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
	"matrix/modules/erp/models/xtpz_models"
)

type ErpPaypalAccount struct {
	*revel.Controller
}

func (c ErpPaypalAccount) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_PAYPAL_ACCOUNT, c.Session)

	return c.RenderTemplate("erp/xtpz/paypal_account/paypal_account_index.html")
}

type ErpPaypalAccountViewItem struct {
	xtpz_models.PaypalAccountInfo `xorm:"extends" json:"a"`
	models.StatusInfo             `xorm:"extends" json:"s"`
}

func (c ErpPaypalAccount) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_PAYPAL_ACCOUNT, c.Session)

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
	PaypalAccount xtpz_models.PaypalAccountInfo `json:"paypal_account"`
	StatusList    []models.StatusInfo           `json:"status_list"`
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
		validation.MaxSize(f.PaypalAccount.PaypalAccount, 200).Message("Paypal账号长度不能大于200！")
	}

	validation.Required(f.PaypalAccount.ClientId).Message("ClientId不能为空！")
	if f.PaypalAccount.ClientId != "" {
		validation.MinSize(f.PaypalAccount.ClientId, 2).Message("ClientId长度不能小于2！")
	}
	if f.PaypalAccount.ClientId != "" {
		validation.MaxSize(f.PaypalAccount.ClientId, 200).Message("ClientId长度不能大于200！")
	}

	validation.Required(f.PaypalAccount.Secret).Message("Secret不能为空！")
	if f.PaypalAccount.Secret != "" {
		validation.MinSize(f.PaypalAccount.Secret, 2).Message("Secret长度不能小于2！")
	}
	if f.PaypalAccount.Secret != "" {
		validation.MaxSize(f.PaypalAccount.Secret, 200).Message("Secret长度不能大于100！")
	}

	if f.PaypalAccount.Desc.Valid == true {
		validation.MaxSize(f.PaypalAccount.Desc.String, 500).Message("描述长度不能大于500！")
	}

	validation.Required(f.PaypalAccount.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPaypalAccount) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_PAYPAL_ACCOUNT, c.Session)

	db_session := c.DbSession

	var paypal_account_id int64
	c.Params.Bind(&paypal_account_id, "id")

	var paypal_account xtpz_models.PaypalAccountInfo
	if paypal_account_id != 0 {
		has, err := db_session.Id(paypal_account_id).Get(&paypal_account)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的Paypal账号不存在！"})
		}
	}

	status_list := status_service.Get_status_list(db_session, models.STATUS_TYPE_PAYPAL_ACCOUNT)

	var detail_form PaypalAccountDetailForm
	detail_form.PaypalAccount = paypal_account
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpPaypalAccount) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_PAYPAL_ACCOUNT, c.Session)

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
		paypal_account_count, err := db_session.Where(fmt.Sprintf("paypal_account = '%s'", paypal_account.PaypalAccount)).Count(new(xtpz_models.PaypalAccountInfo))
		if paypal_account_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，账号已存在！"})
		}

		affected, err = db_session.Insert(&paypal_account)
		core.HandleError(err)
	} else {
		paypal_account_count, err := db_session.Where(fmt.Sprintf("id <> %d AND paypal_account = '%s'", paypal_account.Id, paypal_account.PaypalAccount)).Count(new(xtpz_models.PaypalAccountInfo))
		if paypal_account_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，账号已存在！"})
		}

		affected, err = db_session.Id(paypal_account.Id).Update(&paypal_account)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPaypalAccount) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_PAYPAL_ACCOUNT, c.Session)

	db_session := c.DbSession

	paypal_account_id_list := make([]int64, 0)
	c.Params.Bind(&paypal_account_id_list, "id_list")

	affected, err := db_session.In("id", paypal_account_id_list).Delete(new(xtpz_models.PaypalAccountInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
