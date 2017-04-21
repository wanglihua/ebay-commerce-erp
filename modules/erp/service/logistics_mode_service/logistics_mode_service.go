package logistics_mode_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func Get_logistics_mode_enabled(db_session *xorm.Session) []models.LogisticsModeInfo {
	/*
SELECT lm.*
FROM erp_logistics_mode lm
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = 5000 AND lm.status_id = 5200
	 */
	sql_lm := `SELECT lm.*
FROM erp_logistics_mode lm
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.status_id = %d AND lm.status_id = %d`
	sql_lm = fmt.Sprintf(sql_lm, models.LOGISTICS_COMPANY_ENABLED, models.LOGISTICS_MODE_ENABLED)
	var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
	err := db_session.Sql(sql_lm).Find(&logistics_mode_list)
	core.HandleError(err)

	return logistics_mode_list
}
