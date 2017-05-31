package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/sea_container_cost_service"
)

type ErpSeaContainerCost struct {
	*revel.Controller
}

func (c ErpSeaContainerCost) Index() revel.Result {
	// pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_SETTING, c.Session)

	return c.RenderTemplate("erp/sppz/sea_container_cost/sea_container_cost_index.html")
}

type SeaContainerCostDetailForm struct {
	SeaContainerCost models.SeaContainerCostInfo `json:"sea_container_cost"`
}

func (f *SeaContainerCostDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.SeaContainerCost.Volume).Message("集装箱容积不能为空！")
	validation.Required(f.SeaContainerCost.Cost).Message("集装箱海运成本不能为空！")

	return validation.HasErrors() == false
}

func (c ErpSeaContainerCost) DetailData() revel.Result {
	// pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_SETTING, c.Session)

	db_session := c.DbSession

	var sea_container_cost = sea_container_cost_service.Get_sea_container_cost(db_session)

	var detail_form SeaContainerCostDetailForm
	detail_form.SeaContainerCost = sea_container_cost

	return c.RenderJson(detail_form)
}

func (c ErpSeaContainerCost) Save() revel.Result {
	// pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_SETTING, c.Session)

	db_session := c.DbSession

	var detail_form SeaContainerCostDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	sea_container_cost := &(detail_form.SeaContainerCost) // 取地址，这样下面的对 logistics_setting 的修改都会反映到 detail_form 里面去的，返回的 json 里面也会有反映

	// 先插入再删除，这样从表里面取数据的时候，总能取到的

	// 插入当前的配置信息
	sea_container_cost.Id = 0
	_, err = db_session.Insert(sea_container_cost) // 因为前面取的是地址，所以这里就不用取地址了的
	core.HandleError(err)

	// 删除原来表里面的信息记录
	sql_delete := "DELETE FROM erp_sea_container_cost WHERE id <> %d"
	sql_delete = fmt.Sprintf(sql_delete, sea_container_cost.Id)
	_, err = db_session.Exec(sql_delete)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"detail":  detail_form, // detail_form 里面的数据有过优化和更改的（对 AddressSpecialWord 的优化），所以就再返回一下
		"message": "数据保存成功!",
	})
}
