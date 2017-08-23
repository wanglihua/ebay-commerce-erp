package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/service/pick_fee_service"
	"matrix/modules/erp/models/ddpz_models"
)

type ErpPickFee struct {
	*revel.Controller
}

func (c ErpPickFee) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_PICK_FEE, c.Session)

	return c.RenderTemplate("erp/ddpz/pick_fee/pick_fee_index.html")
}

type PickFeeDetailForm struct {
	PickFee ddpz_models.PickFeeInfo `json:"pick_fee"`
}

func (f *PickFeeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.PickFee.Fee1).Message("5Kg以下不能为空！")

	validation.Required(f.PickFee.Fee2).Message("5.1KG-10KG不能为空！")

	validation.Required(f.PickFee.Fee3).Message("10.1KG-25KG不能为空！")

	validation.Required(f.PickFee.Fee4).Message("25.1KG-50KG不能为空！")

	validation.Required(f.PickFee.Fee5).Message("50.1KG以上不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPickFee) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_PICK_FEE, c.Session)

	db_session := c.DbSession
	var pick_fee = pick_fee_service.Get_pick_fee_setting(db_session)

	var detail_form PickFeeDetailForm
	detail_form.PickFee = pick_fee

	return c.RenderJson(detail_form)
}

func (c ErpPickFee) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_PICK_FEE, c.Session)

	db_session := c.DbSession

	var detail_form PickFeeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	pick_fee := detail_form.PickFee

	// 先插入再删除，这样从表里面取数据的时候，总能取到的

	//插入当前的配置信息
	pick_fee.Id = 0
	_, err = db_session.Insert(&pick_fee)
	core.HandleError(err)

	//删除原来表里面旧的信息记录
	sql_delete := "DELETE FROM erp_pick_fee WHERE id <> %d"
	sql_delete = fmt.Sprintf(sql_delete, pick_fee.Id)
	_, err = db_session.Exec(sql_delete)
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据保存成功!"})

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据保存成功！",
	})
}
