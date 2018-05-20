package logistics_fee_rate_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/ddlc_models"
	"strings"
	"matrix/modules/erp/service/postal_code_service"
)

func Get_logistics_fee_rate2_enabled_and_orders_related(order_list []ddlc_models.OrderInfo, db_session *xorm.Session) []ddpz_models.LogisticsFeeRate2Info {

	var postal_code_list = postal_code_service.Get_orders_postal_code_list(order_list, db_session)

	/*
	SELECT fr.*
FROM erp_logistics_fee_rate2 fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = 5000 AND fr.postcode IN (?)
	 */
	sql_fr := `SELECT fr.*
FROM erp_logistics_fee_rate2 fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = %d AND fr.postcode IN (%s)`
	sql_fr = fmt.Sprintf(sql_fr, models.LOGISTICS_COMPANY_ENABLED, strings.Join(postal_code_list, ","))
	var logistics_fee_rate_all_list = make([]ddpz_models.LogisticsFeeRate2Info, 0)
	err := db_session.Sql(sql_fr).Find(&logistics_fee_rate_all_list)
	core.HandleError(err)

	return logistics_fee_rate_all_list
}
