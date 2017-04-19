package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
)

type ErpShopAccount struct {
	*revel.Controller
}

func (c ErpShopAccount) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	return c.RenderTemplate("erp/xtpz/shop_account/shop_account_index.html")
}

type ErpShopAccountViewItem struct {
	models.ShopAccountInfo        `xorm:"extends" json:"a"`
	models.ECommercePlatformInfo  `xorm:"extends" json:"p"`
	models.StatusInfo             `xorm:"extends" json:"s"`
}

func (c ErpShopAccount) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT a.*, p.*, s.*
FROM erp_shop_account a
    INNER JOIN erp_ecommerce_platform p ON a.platform_id = p.id
    INNER JOIN erp_status s ON a.status_id = s.id
	*/

	query := db_session.
	Select("a.*, p.*, s.*").
		Table("erp_shop_account").Alias("a").
		Join("INNER", []string{"erp_ecommerce_platform", "p"}, "a.platform_id = p.id").
		Join("INNER", []string{"erp_status", "s"}, "a.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("a.id")
	}

	shop_account_list := make([]ErpShopAccountViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&shop_account_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpShopAccountViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  shop_account_list,
		Total: count,
	})
}

type ShopAccountDetailForm struct {
	ShopAccount  models.ShopAccountInfo         `json:"shop_account"`
	StatusList   []models.StatusInfo            `json:"status_list"`
	PlatformList []models.ECommercePlatformInfo `json:"platform_list"`
}

func (f *ShopAccountDetailForm) IsCreate() bool {
	return f.ShopAccount.Id == 0
}

func (f *ShopAccountDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ShopAccount.Name).Message("账号不能为空！")
	if f.ShopAccount.Name != "" {
		validation.MinSize(f.ShopAccount.Name, 2).Message("账号长度不能小于2！")
	}
	if f.ShopAccount.Name != "" {
		validation.MaxSize(f.ShopAccount.Name, 100).Message("账号长度不能大于100！")
	}

	validation.Required(f.ShopAccount.PlatformId).Message("平台不能为空！")

	if f.ShopAccount.Desc.Valid == true {
		validation.MaxSize(f.ShopAccount.Desc.String, 500).Message("描述长度不能大于500！")
	}

	if f.ShopAccount.Manager.Valid == true {
		validation.MaxSize(f.ShopAccount.Manager.String, 30).Message("负责人长度不能大于30！")
	}

	validation.Required(f.ShopAccount.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpShopAccount) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	var shop_account_id int64
	c.Params.Bind(&shop_account_id, "id")

	var shop_account models.ShopAccountInfo
	if shop_account_id != 0 {
		has, err := db_session.Id(shop_account_id).Get(&shop_account)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的店铺账号不存在！"})
		}
	}

	if shop_account_id == 0 {
		shop_account.StatusId = models.SHOP_ACCOUNT_STATUS_ENABLED
	}

	status_list := status_service.GetStatusList(db_session, models.STATUS_TYPE_SHOP_ACCOUNT)

	platform_list := make([]models.ECommercePlatformInfo, 0)
	err := db_session.Find(&platform_list)
	core.HandleError(err)

	var detail_form ShopAccountDetailForm
	detail_form.ShopAccount = shop_account
	detail_form.StatusList = status_list
	detail_form.PlatformList = platform_list

	return c.RenderJson(detail_form)
}

func (c ErpShopAccount) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	var detail_form ShopAccountDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	shop_account := detail_form.ShopAccount

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&shop_account)
		core.HandleError(err)
	} else {
		update_col_list := []string{"name", "platform_id", "desc", "paypal_account", "manager", "status_id"}
		affected, err = db_session.Id(shop_account.Id).Cols(update_col_list...).Update(&shop_account)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpShopAccount) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	shop_account_id_list := make([]int64, 0)
	c.Params.Bind(&shop_account_id_list, "id_list")

	// 删除关联的 eBay API 配置
	_, err := db_session.In("shop_id", shop_account_id_list).Delete(new(models.EBayAPIConfigInfo))
	core.HandleError(err)

	// 删除店铺信息
	_, err = db_session.In("id", shop_account_id_list).Delete(new(models.ShopAccountInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "删除成功!"})
}

func (c ErpShopAccount) GetShopAccountInfo() revel.Result {
	// 这里面就不做权限检查了，别的地方调用的

	db_session := c.DbSession

	var shop_account_id int64
	c.Params.Bind(&shop_account_id, "shop_account_id")

	var shop_account models.ShopAccountInfo
	_, err := db_session.Id(shop_account_id).Get(&shop_account)
	core.HandleError(err)

	var platform models.ECommercePlatformInfo
	_, err = db_session.Id(shop_account.PlatformId).Get(&platform)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"shop_account": shop_account,
		"platform":     platform,
	})
}
