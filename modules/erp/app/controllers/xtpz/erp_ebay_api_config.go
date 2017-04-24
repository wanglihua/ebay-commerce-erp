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
	"matrix/modules/erp/service/shop_account_log_service"
)

type ErpEBayAPIConfig struct {
	*revel.Controller
}

type EBayAPIConfigDetailForm struct {
	EBayAPIConfig models.EBayAPIConfigInfo `json:"api_config"`
	ShopName      string                   `json:"shop_name"`
}

func (f *EBayAPIConfigDetailForm) IsCreate() bool {
	return f.EBayAPIConfig.Id == 0
}

func (f *EBayAPIConfigDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.EBayAPIConfig.ShopId).Message("店铺不能为空！")
	validation.Required(f.EBayAPIConfig.AppId).Message("AppID不能为空！")
	validation.Required(f.EBayAPIConfig.CertId).Message("CertID不能为空！")
	validation.Required(f.EBayAPIConfig.DevId).Message("DevID不能为空！")
	validation.Required(f.EBayAPIConfig.Token).Message("Token不能为空！")

	return validation.HasErrors() == false
}

func (c ErpEBayAPIConfig) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	var shop_id int64
	c.Params.Bind(&shop_id, "shop_id")

	var ebay_api_config models.EBayAPIConfigInfo
	if shop_id == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择ebay店铺进行api参数配置",
		})
	}

	has, err := db_session.Where("shop_id = ?", shop_id).Get(&ebay_api_config)
	core.HandleError(err)
	if has == false {
		ebay_api_config.ShopId = shop_id //如果店铺配置不存在，也返回店铺配置的shop_id
	}

	var shop_name struct {
		Name string `xorm:"'shop_name'"`
	}
	/*
	SELECT concat(ep.name, ' ', sa.name) AS shop_name
	FROM erp_shop_account sa
		INNER JOIN erp_ecommerce_platform ep ON sa.platform_id = ep.id
	WHERE sa.id = 0
	 */
	sql := fmt.Sprintf(`SELECT concat(ep.name, ' ', sa.name) AS shop_name
FROM erp_shop_account sa
    INNER JOIN erp_ecommerce_platform ep ON sa.platform_id = ep.id
WHERE sa.id = %d`, shop_id)
	_, err = db_session.Sql(sql).Get(&shop_name)
	core.HandleError(err)

	var detail_form EBayAPIConfigDetailForm
	detail_form.EBayAPIConfig = ebay_api_config
	detail_form.ShopName = shop_name.Name

	return c.RenderJson(detail_form)
}

func (c ErpEBayAPIConfig) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	var detail_form EBayAPIConfigDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	ebay_api_config := detail_form.EBayAPIConfig

	var shop_account models.ShopAccountInfo
	_, err = db_session.Id(ebay_api_config.ShopId).Get(&shop_account)
	core.HandleError(err)

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&ebay_api_config)
		core.HandleError(err)

		shop_account_log_service.Insert_shop_account_api_log(&shop_account, core.Db_operation_types.Insert, db_session, c.Session)
	} else {
		affected, err = db_session.Id(ebay_api_config.Id).Update(&ebay_api_config)
		core.HandleError(err)

		shop_account_log_service.Insert_shop_account_api_log(&shop_account, core.Db_operation_types.Update, db_session, c.Session)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpEBayAPIConfig) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SHOP_ACCOUNT, c.Session)

	db_session := c.DbSession

	var config_id int64
	c.Params.Bind(&config_id, "id")

	_, err := db_session.Id(config_id).Delete(new(models.EBayAPIConfigInfo))
	core.HandleError(err)

	var ebay_api_config models.EBayAPIConfigInfo
	_, err = db_session.Id(config_id).Get(&ebay_api_config)
	core.HandleError(err)

	var shop_account models.ShopAccountInfo
	_, err = db_session.Id(ebay_api_config.ShopId).Get(&shop_account)
	core.HandleError(err)

	shop_account_log_service.Insert_shop_account_api_log(&shop_account, core.Db_operation_types.Delete, db_session, c.Session)

	return c.RenderJson(core.JsonResult{Success: true, Message: "API配置数据删除成功!"})
}
