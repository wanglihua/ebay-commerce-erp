package sea_container_cost_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/sppz_models"
)

func Get_sea_container_cost(db_session *xorm.Session) sppz_models.SeaContainerCostInfo {
	var sea_container_cost sppz_models.SeaContainerCostInfo
	sql := "SELECT * FROM erp_sea_container_cost ls ORDER BY ls.update_time DESC LIMIT 0, 1"
	_, err := db_session.Sql(sql).Get(&sea_container_cost)
	core.HandleError(err)

	return sea_container_cost
}
