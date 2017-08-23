package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
	"matrix/modules/erp/models/xtpz_models"
)

type ErpExchangeRate struct {
	*revel.Controller
}

func (c ErpExchangeRate) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EXCHANGE_RATE, c.Session)

	return c.RenderTemplate("erp/xtpz/exchange_rate/exchange_rate_index.html")
}

type ErpExchangeRateViewItem struct {
	xtpz_models.ExchangeRateInfo `xorm:"extends" json:"r"`
	xtpz_models.CurrencyInfo     `xorm:"extends" json:"c"`
	models.StatusInfo       `xorm:"extends" json:"s"`
}

func (c ErpExchangeRate) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EXCHANGE_RATE, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
SELECT r.*, c.*, s.*
FROM erp_exchange_rate r
    INNER JOIN erp_currency c ON r.currency_id = c.id
    INNER JOIN erp_status s ON r.status_id = s.id
	*/

	sql_builder := core.NewSqlBuilder()

	sql_builder.Select("r.*, c.*, s.*")
	sql_builder.Count("count(r.id) as rate_count")
	sql_builder.From(`erp_exchange_rate r
    INNER JOIN erp_currency c ON r.currency_id = c.id
    INNER JOIN erp_status s ON r.status_id = s.id`)
	sql_builder.Where(filter)

	if order != "" {
		sql_builder.OrderBy(order)
	} else {
		sql_builder.OrderBy("r.id ASC")
	}

	sql_builder.Limit(offset, limit)

	exchange_rate_list := make([]ErpExchangeRateViewItem, 0)
	err := db_session.Sql(sql_builder.SqlList()).Find(&exchange_rate_list)
	core.HandleError(err)

	var rate_count struct {
		Count int64 `xorm:"'rate_count'"`
	}

	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&rate_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  exchange_rate_list,
		Total: rate_count.Count,
	})
}

type ExchangeRateDetailForm struct {
	ExchangeRate xtpz_models.ExchangeRateInfo `json:"exchange_rate"`
	StatusList   []models.StatusInfo     `json:"status_list"`
	CurrencyList []xtpz_models.CurrencyInfo   `json:"currency_list"`
}

func (f *ExchangeRateDetailForm) IsCreate() bool {
	return f.ExchangeRate.Id == 0
}

func (f *ExchangeRateDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ExchangeRate.CurrencyId).Message("货币不能为空！")
	validation.Required(f.ExchangeRate.Rate).Message("汇率不能为空！")
	validation.Required(f.ExchangeRate.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpExchangeRate) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EXCHANGE_RATE, c.Session)

	db_session := c.DbSession

	var exchange_rate_id int64
	c.Params.Bind(&exchange_rate_id, "id")

	var exchange_rate xtpz_models.ExchangeRateInfo
	if exchange_rate_id != 0 {
		has, err := db_session.Id(exchange_rate_id).Get(&exchange_rate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的货币汇率不存在！"})
		}
	} else {
		//如果是新增，就给个初始的汇率状态
		exchange_rate.StatusId = models.EXCHANGE_RATE_STATUS_ENABLED
	}

	status_list := status_service.Get_status_list(db_session, models.STATUS_TYPE_EXCHANGE_RATE)

	currency_list := make([]xtpz_models.CurrencyInfo, 0)
	err := db_session.Sql(fmt.Sprintf("SELECT * FROM erp_currency WHERE id <> %d", xtpz_models.CURRENCY_RMB)).Find(&currency_list)
	core.HandleError(err)

	var detail_form ExchangeRateDetailForm
	detail_form.ExchangeRate = exchange_rate
	detail_form.CurrencyList = currency_list
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpExchangeRate) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EXCHANGE_RATE, c.Session)

	db_session := c.DbSession

	var detail_form ExchangeRateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	exchange_rate_from_ui := detail_form.ExchangeRate

	var affected int64
	if detail_form.IsCreate() {
		currency_rate_count, err := db_session.Where("currency_id = ?", exchange_rate_from_ui.CurrencyId).Count(new(xtpz_models.ExchangeRateInfo))
		core.HandleError(err)
		if currency_rate_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，货币汇率已存在！"})
		}

		affected, err = db_session.Insert(&exchange_rate_from_ui)
		core.HandleError(err)

		var currency xtpz_models.CurrencyInfo
		_, err = db_session.Id(exchange_rate_from_ui.CurrencyId).Get(&currency)
		core.HandleError(err)

		var exchange_rate_log xtpz_models.ExchangeRateLogInfo
		exchange_rate_log.Name = fmt.Sprintf("%s(%s)", currency.Name, currency.Code)
		exchange_rate_log.Desc = fmt.Sprintf("新增数据 %.4f", exchange_rate_from_ui.Rate)
		exchange_rate_log.LoginName = core.GetLoginUser(c.Session).LoginName

		_, err = db_session.Insert(&exchange_rate_log)
		core.HandleError(err)

	} else {
		currency_rate_count, err := db_session.Where("id <> ? and currency_id = ?", exchange_rate_from_ui.Id, exchange_rate_from_ui.CurrencyId).Count(new(xtpz_models.ExchangeRateInfo))
		core.HandleError(err)
		if currency_rate_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，货币汇率已存在！"})
		}

		var exchange_rate_in_db xtpz_models.ExchangeRateInfo
		_, err = db_session.Id(exchange_rate_from_ui.Id).Get(&exchange_rate_in_db)
		core.HandleError(err)

		affected, err = db_session.Id(exchange_rate_from_ui.Id).Update(&exchange_rate_from_ui)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}

		var currency xtpz_models.CurrencyInfo
		_, err = db_session.Id(exchange_rate_from_ui.CurrencyId).Get(&currency)
		core.HandleError(err)

		var exchange_rate_log xtpz_models.ExchangeRateLogInfo
		exchange_rate_log.Name = fmt.Sprintf("%s(%s)", currency.Name, currency.Code)
		exchange_rate_log.Desc = fmt.Sprintf("编辑数据 %.4f=>%.4f", exchange_rate_in_db.Rate, exchange_rate_from_ui.Rate)
		exchange_rate_log.LoginName = core.GetLoginUser(c.Session).LoginName

		_, err = db_session.Insert(&exchange_rate_log)
		core.HandleError(err)
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpExchangeRate) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_EXCHANGE_RATE, c.Session)

	db_session := c.DbSession

	exchange_rate_id_list := make([]int64, 0)
	c.Params.Bind(&exchange_rate_id_list, "id_list")

	for _, exchange_rate_id := range exchange_rate_id_list {
		var exchange_rate xtpz_models.ExchangeRateInfo
		has, err := db_session.Id(exchange_rate_id).Get(&exchange_rate)
		core.HandleError(err)

		if has {
			var exchange_rate_log xtpz_models.ExchangeRateLogInfo
			currency := c.get_exchange_rate_currency(exchange_rate_id, db_session)
			exchange_rate_log.Name = fmt.Sprintf("%s(%s)", currency.Name, currency.Code)
			exchange_rate_log.Desc = "删除数据"
			exchange_rate_log.LoginName = core.GetLoginUser(c.Session).LoginName

			_, err := db_session.Insert(&exchange_rate_log)
			core.HandleError(err)
		}
	}

	affected, err := db_session.In("id", exchange_rate_id_list).Delete(new(xtpz_models.ExchangeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpExchangeRate) get_exchange_rate_currency(exchange_rate_id int64, db_session *xorm.Session) xtpz_models.CurrencyInfo {
	sql := `
SELECT c.*
FROM erp_exchange_rate r
    INNER JOIN erp_currency c ON r.currency_id = c.id
WHERE r.id = %d
`
	sql = fmt.Sprintf(sql, exchange_rate_id)

	var currency xtpz_models.CurrencyInfo
	_, err := db_session.Sql(sql).Get(&currency)
	core.HandleError(err)

	return currency
}
