package exchange_rate_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/xtpz_models"
)

type Exchange_rate_view_item_t struct {
	Currency     xtpz_models.CurrencyInfo        `xorm:"extends" json:"c"`
	ExchangeRate xtpz_models.ExchangeRateInfo    `xorm:"extends" json:"er"`
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

// 获取澳币的汇率
func Get_aud_exchange_rate(db_session *xorm.Session) float64 {
	var aud_rate struct {
		Rate float64 `xorm:"'rate'"`
	}

	sql := `SELECT er.rate
FROM erp_exchange_rate er
    INNER JOIN erp_currency c ON c.id = er.currency_id
WHERE c.code = 'AUD'`

	_, err := db_session.Sql(sql).Get(&aud_rate)
	core.HandleError(err)

	if aud_rate.Rate == 0 {
		return 1
	}

	return aud_rate.Rate
}
