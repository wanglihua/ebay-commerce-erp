package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/xtpz_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpEmailAddress struct {
	*revel.Controller
}

func (c ErpEmailAddress) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EMAIL_ADDRESS, c.Session)

	var db_session = c.DbSession

	var shop_list = make([]xtpz_models.ShopAccountInfo, 0)
	err := db_session.Find(&shop_list)
	core.HandleError(err)

	c.RenderArgs["shop_list"] = shop_list

	return c.RenderTemplate("erp/xtpz/email_address/email_address_index.html")
}

func (c ErpEmailAddress) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EMAIL_ADDRESS, c.Session)

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

	email_address_list := make([]xtpz_models.EmailAddressInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&email_address_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(xtpz_models.EmailAddressInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  email_address_list,
		Total: count,
	})
}

type EmailAddressDetailForm struct {
	EmailAddress xtpz_models.EmailAddressInfo `json:"email_address"`
}

func (f *EmailAddressDetailForm) IsCreate() bool {
	return f.EmailAddress.Id == 0
}

func (f *EmailAddressDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.EmailAddress.ShopId).Message("店铺不能为空！")

	validation.Required(f.EmailAddress.Email).Message("邮件地址不能为空！")

	validation.Required(f.EmailAddress.Account).Message("账号不能为空！")

	validation.Required(f.EmailAddress.Password).Message("密码不能为空！")

	validation.Required(f.EmailAddress.SmtpServer).Message("发件服务器不能为空！")

	validation.Required(f.EmailAddress.SmtpPort).Message("发件服务器端口不能为空！")

	return validation.HasErrors() == false
}

func (c ErpEmailAddress) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EMAIL_ADDRESS, c.Session)

	db_session := c.DbSession

	var email_address_id int64
	c.Params.Bind(&email_address_id, "id")

	var email_address xtpz_models.EmailAddressInfo
	if email_address_id != 0 {
		has, err := db_session.Id(email_address_id).Get(&email_address)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的电子邮件地址不存在！"})
		}
	} else {
		email_address.SmtpPort = "25"
	}

	var detail_form EmailAddressDetailForm
	detail_form.EmailAddress = email_address

	return c.RenderJson(detail_form)
}

func (c ErpEmailAddress) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EMAIL_ADDRESS, c.Session)

	db_session := c.DbSession

	var detail_form EmailAddressDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	email_address := detail_form.EmailAddress

	var shop_id = email_address.ShopId
	var shop_account xtpz_models.ShopAccountInfo
	_, err = db_session.Id(shop_id).Get(&shop_account)
	core.HandleError(err)

	email_address.ShopName = shop_account.Name

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&email_address)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(email_address.Id).Update(&email_address)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpEmailAddress) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EMAIL_ADDRESS, c.Session)

	db_session := c.DbSession

	email_address_id_list := make([]int64, 0)
	c.Params.Bind(&email_address_id_list, "id_list")

	affected, err := db_session.In("id", email_address_id_list).Delete(new(xtpz_models.EmailAddressInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
