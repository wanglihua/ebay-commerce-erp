package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"matrix/core"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpStockPickCrew struct {
	*revel.Controller
}

func (c ErpStockPickCrew) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STOCK_PICK_CREW, c.Session)

	return c.RenderTemplate("erp/sppz/stock_pick_crew/stock_pick_crew_index.html")
}

func (c ErpStockPickCrew) PageData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STOCK_PICK_CREW, c.Session)

	var db_session = c.DbSession

	var stock_list = make([]models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"stock_id_query":   0,
		"stock_list":       stock_list,
		"login_name_query": "",
		"nick_name_query":  "",
		"is_manager_query": 0,
	})
}

func (c ErpStockPickCrew) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STOCK_PICK_CREW, c.Session)

	db_session := c.DbSession

	var stock_id_query int64
	c.Params.Bind(&stock_id_query, "stock_id_query")

	var login_name_query string
	c.Params.Bind(&login_name_query, "login_name_query")
	login_name_query = strings.TrimSpace(login_name_query)

	var nick_name_query string
	c.Params.Bind(&nick_name_query, "nick_name_query")
	nick_name_query = strings.TrimSpace(nick_name_query)

	var is_manager_query int64
	c.Params.Bind(&is_manager_query, "is_manager_query")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type stock_pick_crew_view_item_t struct {
		models.StockPickCrewInfo  `xorm:"extends" json:"spc"`
		models.StockInfo          `xorm:"extends" json:"s"`
		auth_models.UserInfo      `xorm:"extends" json:"u"`
	}

	/*
SELECT spc.*, s.*, u.*
FROM erp_stock_pick_crew spc
    INNER JOIN erp_stock s ON s.id = spc.stock_id
    INNER JOIN auth_user u ON u.id = spc.user_id
WHERE spc.stock_id = 0 AND u.login_name LIKE '%abc%' AND u.nick_name LIKE '%abc%' AND spc.is_manager = -1
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("spc.*, s.*, u.*")
	sb.Count("count(spc.id) AS stock_pick_crew_count")
	sb.From(`erp_stock_pick_crew spc
    INNER JOIN erp_stock s ON s.id = spc.stock_id
    INNER JOIN auth_user u ON u.id = spc.user_id`)

	if stock_id_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("spc.stock_id = %d", stock_id_query))
	}

	if login_name_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("u.login_name LIKE '%%%s%%'", login_name_query))
	}

	if nick_name_query != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("u.nick_name LIKE '%%%s%%'", nick_name_query))
	}

	if is_manager_query != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("spc.is_manager = %d", is_manager_query))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("spc.id ASC")
	}

	sb.Limit(offset, limit)

	stock_pick_crew_list := make([]stock_pick_crew_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&stock_pick_crew_list)
	core.HandleError(err)

	var stock_pick_view_count struct {
		Count int64 `xorm:"'stock_pick_crew_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&stock_pick_view_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  stock_pick_crew_list,
		Total: stock_pick_view_count.Count,
	})
}

type StockPickCrewDetailForm struct {
	StockPickCrew models.StockPickCrewInfo  `json:"stock_pick_crew"`
	StockList     []models.StockInfo        `json:"stock_list"`
	UserLoginName string                    `json:"user_login_name"`
}

func (f *StockPickCrewDetailForm) IsCreate() bool {
	return f.StockPickCrew.Id == 0
}

func (f *StockPickCrewDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.StockPickCrew.StockId).Message("仓库不能为空！")

	if f.StockPickCrew.StockId == 0 {
		validation.Error("请选择仓库！");
	}

	validation.Required(f.StockPickCrew.UserId).Message("用户不能为空！")
	validation.Required(f.StockPickCrew.IsManager).Message("是否管理员不能为空！")

	return validation.HasErrors() == false
}

func (c ErpStockPickCrew) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STOCK_PICK_CREW, c.Session)

	db_session := c.DbSession

	var stock_pick_crew_id int64
	c.Params.Bind(&stock_pick_crew_id, "id")

	var detail_form StockPickCrewDetailForm

	var stock_pick_crew models.StockPickCrewInfo
	if stock_pick_crew_id != 0 {
		has, err := db_session.Id(stock_pick_crew_id).Get(&stock_pick_crew)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的仓库拣货人员不存在！"})
		}

		var user_id = stock_pick_crew.UserId
		var user auth_models.UserInfo
		_, err = db_session.Id(user_id).Get(&user)
		core.HandleError(err)

		detail_form.UserLoginName = user.LoginName

	} else {
		stock_pick_crew.IsManager = -1 // 否，不是管理员
	}

	var stock_list = make([]models.StockInfo, 0)
	err := db_session.Find(&stock_list)
	core.HandleError(err)

	detail_form.StockPickCrew = stock_pick_crew
	detail_form.StockList = stock_list

	return c.RenderJson(detail_form)
}

func (c ErpStockPickCrew) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STOCK_PICK_CREW, c.Session)

	db_session := c.DbSession

	var detail_form StockPickCrewDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	stock_pick_crew := detail_form.StockPickCrew

	//todo: 重复数据的判断， stock_id, user_id
	var stock_pick_crew_count struct {
		Count int64 `xorm:"'count'"`
	}

	var affected int64
	if detail_form.IsCreate() {
		/*
		SELECT count(id) AS count FROM erp_stock_pick_crew WHERE stock_id = 0 AND user_id = 0
		 */

		var sql_count = `SELECT count(id) AS count FROM erp_stock_pick_crew WHERE stock_id = %d AND user_id = %d`
		sql_count = fmt.Sprintf(sql_count, stock_pick_crew.StockId, stock_pick_crew.UserId)
		_, err := db_session.Sql(sql_count).Get(&stock_pick_crew_count)
		core.HandleError(err)

		if stock_pick_crew_count.Count != 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "该仓库已存在该拣货人员！",
			})
		}

		affected, err = db_session.Insert(&stock_pick_crew)
		core.HandleError(err)
	} else {
		/*
		SELECT count(id) AS count FROM erp_stock_pick_crew WHERE id <> 0 AND stock_id = 0 AND user_id = 0
		 */

		var sql_count = `SELECT count(id) AS count FROM erp_stock_pick_crew WHERE id <> %d AND stock_id = %d AND user_id = %d`
		sql_count = fmt.Sprintf(sql_count, stock_pick_crew.Id, stock_pick_crew.StockId, stock_pick_crew.UserId)
		_, err := db_session.Sql(sql_count).Get(&stock_pick_crew_count)
		core.HandleError(err)

		if stock_pick_crew_count.Count != 0 {
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "该仓库已存在该拣货人员！",
			})
		}

		affected, err = db_session.Id(stock_pick_crew.Id).Update(&stock_pick_crew)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpStockPickCrew) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STOCK_PICK_CREW, c.Session)

	db_session := c.DbSession

	stock_pick_crew_id_list := make([]int64, 0)
	c.Params.Bind(&stock_pick_crew_id_list, "id_list")

	affected, err := db_session.In("id", stock_pick_crew_id_list).Delete(new(models.StockPickCrewInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpStockPickCrew) SelectUserListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_STOCK_PICK_CREW, c.Session)

	db_session := c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type user_view_item_t struct {
		User auth_models.UserInfo `xorm:"extends" json:"u"`
	}

	/*
SELECT u.*
FROM auth_user u
WHERE u.id NOT IN (SELECT spc.user_id FROM erp_stock_pick_crew spc WHERE spc.stock_id = 0)
	 */

	var sb = core.NewSqlBuilder()
	sb.Select("u.*")
	sb.Count("count(u.id) AS user_count")
	sb.From("auth_user u")

	filter = core.SqlAnd(filter, fmt.Sprintf("u.id IN (SELECT gu.user_id FROM auth_group_user gu WHERE gu.group_id = %d)", auth_models.GROUP_PICK_CREW))

	if stock_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("u.id NOT IN (SELECT spc.user_id FROM erp_stock_pick_crew spc WHERE spc.stock_id = %d)", stock_id))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("u.id ASC")
	}

	sb.Limit(offset, limit)

	userList := make([]user_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&userList)
	core.HandleError(err)

	var user_count struct {
		Count int64 `xorm:"'user_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&user_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  userList,
		Total: user_count.Count,
	})
}
