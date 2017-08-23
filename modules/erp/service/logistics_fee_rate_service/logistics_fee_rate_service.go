package logistics_fee_rate_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/ddpz_models"
)

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
