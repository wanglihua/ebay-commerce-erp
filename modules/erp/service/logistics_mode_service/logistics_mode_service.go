package logistics_mode_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/ddpz_models"
)

func Get_logistics_mode_enabled(db_session *xorm.Session) ([]ddpz_models.LogisticsModeInfo, error) {
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
	var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.Sql(sql_lm).Find(&logistics_mode_list)
	// core.HandleError(err)

	return logistics_mode_list, err
}

func Get_logistics_mode_id_name_map(mode_id_list []int64, db_session *xorm.Session) map[int64]string {
	var mode_id_name_map = make(map[int64]string)

	if len(mode_id_list) == 0 {
		return mode_id_name_map
	}

	var mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.In("id", mode_id_list).Find(&mode_list)
	core.HandleError(err)

	for _, mode := range mode_list {
		mode_id_name_map[mode.Id] = mode.Name
	}

	return mode_id_name_map
}

func Get_logistics_mode_map(mode_id_list []int64, db_session *xorm.Session) map[int64]ddpz_models.LogisticsModeInfo {
	var mode_map = make(map[int64]ddpz_models.LogisticsModeInfo)

	if len(mode_id_list) == 0 {
		return mode_map
	}

	var mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.In("id", mode_id_list).Find(&mode_list)
	core.HandleError(err)

	for _, mode := range mode_list {
		mode_map[mode.Id] = mode
	}

	return mode_map
}
