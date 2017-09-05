package logistics_fee_rate_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/ddlc_models"
	"unicode/utf8"
	"strings"
)

// 改用下面的更高效的方式
func Get_logistics_fee_rate_enabled(db_session *xorm.Session) []ddpz_models.LogisticsFeeRateInfo {
	/*
	SELECT fr.*
FROM erp_logistics_fee_rate fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = 5000
	 */
	sql_fr := `SELECT fr.*
FROM erp_logistics_fee_rate fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = %d`
	sql_fr = fmt.Sprintf(sql_fr, models.LOGISTICS_COMPANY_ENABLED)
	var logistics_fee_rate_all_list = make([]ddpz_models.LogisticsFeeRateInfo, 0)
	err := db_session.Sql(sql_fr).Find(&logistics_fee_rate_all_list)
	core.HandleError(err)

	return logistics_fee_rate_all_list
}

func Get_logistics_fee_rate_enabled_and_orders_related(order_list []ddlc_models.OrderInfo, db_session *xorm.Session) []ddpz_models.LogisticsFeeRateInfo {
	/*
SELECT DISTINCT osa.postal_code
FROM erp_order_shipping_address osa
WHERE osa.order_id IN (88415, 88416, 88417)
	 */

	var order_id_str_list = make([]string, 0)
	for _, order := range order_list {
		order_id_str_list = append(order_id_str_list, fmt.Sprintf("%d", order.Id))
	}

	type postal_code_result struct {
		PostalCode string `xorm:"'postal_code'"`
	}

	var postal_code_result_list = make([]postal_code_result, 0)
	var sql_postal_code_list = `SELECT DISTINCT osa.postal_code
FROM erp_order_shipping_address osa
WHERE osa.order_id IN (%s)`
	sql_postal_code_list = fmt.Sprintf(sql_postal_code_list, strings.Join(order_id_str_list, ","))
	core.HandleError(db_session.Sql(sql_postal_code_list).Find(&postal_code_result_list))

	var postal_code_list = make([]string, 0)
	for _, postal_code_result := range postal_code_result_list {
		postal_code_list = append(postal_code_list, fmt.Sprintf("'%s'", postal_code_result.PostalCode))

		// 三位的补个 0， 放入列表
		if utf8.RuneCountInString(postal_code_result.PostalCode) == 3 {
			postal_code_list = append(postal_code_list, fmt.Sprintf("'%s'", "0"+postal_code_result.PostalCode))
		}

		// 四位的以 0 开头的， 去掉 0， 放入列表
		if utf8.RuneCountInString(postal_code_result.PostalCode) == 4 && strings.HasPrefix(postal_code_result.PostalCode, "0") {
			postal_code_list = append(postal_code_list, fmt.Sprintf("'%s'", strings.TrimPrefix(postal_code_result.PostalCode, "0")))
		}
	}

	/*
	SELECT fr.*
FROM erp_logistics_fee_rate fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = 5000 AND fr.postcode IN (?)
	 */
	sql_fr := `SELECT fr.*
FROM erp_logistics_fee_rate fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = %d AND fr.postcode IN (%s)`
	sql_fr = fmt.Sprintf(sql_fr, models.LOGISTICS_COMPANY_ENABLED, strings.Join(postal_code_list, ","))
	var logistics_fee_rate_all_list = make([]ddpz_models.LogisticsFeeRateInfo, 0)
	err := db_session.Sql(sql_fr).Find(&logistics_fee_rate_all_list)
	core.HandleError(err)

	return logistics_fee_rate_all_list
}
