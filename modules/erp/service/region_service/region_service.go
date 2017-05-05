package region_service

import (
	"fmt"
	"math"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func Get_region_full_name(region_id int64, db_session *xorm.Session) string {

	type region_name_parent_id_t struct {
		Name     string   `xorm:"'name'"`
		ParentId int64    `xorm:"'parent_id'"`
	}
	/*
	SELECT concat(r.province, r.city, r.district) AS name, r.parent AS parent_id FROM erp_region r WHERE r.id = 3
	 */

	var sql_template = `SELECT concat(r.province, r.city, r.district) AS name, r.parent AS parent_id FROM erp_region r WHERE r.id = %d`

	var region_full_name = ""

	var parent_id int64 = math.MaxInt64
	for parent_id > models.REGION_ROOT_ID {
		var region_name_parent_id region_name_parent_id_t
		var sql = fmt.Sprintf(sql_template, region_id)
		_, err := db_session.Sql(sql).Get(&region_name_parent_id)
		core.HandleError(err)

		region_full_name = region_name_parent_id.Name + region_full_name
		parent_id = region_name_parent_id.ParentId
	}

	return region_full_name
}
