package postal_code_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ddlc_models"
	"fmt"
	"strings"
	"unicode/utf8"
	"matrix/core"
)

func Get_orders_postal_code_list(order_list []ddlc_models.OrderInfo, db_session *xorm.Session) []string {

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
			postal_code_list = append(postal_code_list, fmt.Sprintf("'%s'", string(([]rune(postal_code_result.PostalCode))[1:])))
		}
	}

	return postal_code_list
}
