package exchange_rate_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Exchange_rate_view_item_t struct {
	Currency     models.CurrencyInfo        `xorm:"extends" json:"c"`
	ExchangeRate models.ExchangeRateInfo    `xorm:"extends" json:"er"`
}

func Get_exchange_rate_list(db_session *xorm.Session) []Exchange_rate_view_item_t {

	sql_get_exchange_rate_list := `SELECT er.*, c.*
FROM erp_exchange_rate er
    INNER JOIN erp_currency c ON c.id = er.currency_id`

	var exchange_rate_list = make([]Exchange_rate_view_item_t, 0)
	err := db_session.Sql(sql_get_exchange_rate_list).Find(&exchange_rate_list)
	core.HandleError(err)

	return exchange_rate_list
}
