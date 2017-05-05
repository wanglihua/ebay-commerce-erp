package region_service

import (
	"fmt"
	"math"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Region_t struct {
	Id   int64      `xorm:"'id'" json:"id"`
	Name string     `xorm:"'name'" json:"name"`
}

func Get_province_list(db_session *xorm.Session) []Region_t {
	/*
	SELECT r.id AS id, r.province AS name FROM erp_region r WHERE r.parent = 1 ORDER BY r.id ASC;
	 */

	var sql_province_list = `SELECT r.id AS id, r.province AS name FROM erp_region r WHERE r.parent = %d ORDER BY r.id ASC`
	sql_province_list = fmt.Sprintf(sql_province_list, models.REGION_ROOT_ID)

	var province_list = make([]Region_t, 0)
	err := db_session.Sql(sql_province_list).Find(&province_list)
	core.HandleError(err)

	return province_list
}

func Get_city_list(province_id int64, db_session *xorm.Session) []Region_t {
	/*
	SELECT r.id AS id, concat(r.city, r.district) AS name FROM erp_region r WHERE r.parent = 2 ORDER BY r.id ASC;
	 */

	var sql_city_list = `SELECT r.id AS id, concat(r.city, r.district) AS name FROM erp_region r WHERE r.parent = %d ORDER BY r.id ASC`
	sql_city_list = fmt.Sprintf(sql_city_list, province_id)

	var city_list = make([]Region_t, 0)
	err := db_session.Sql(sql_city_list).Find(&city_list)
	core.HandleError(err)

	return city_list
}

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

func Get_region_parent_id(region_id int64, db_session *xorm.Session) int64 {
	var parent_id struct {
		ParentId int64 `xorm:"'parent_id'"`
	}
	/*
	SELECT r.parent AS parent_id FROM erp_region r WHERE r.id = 3
	 */

	var sql = `SELECT r.parent AS parent_id FROM erp_region r WHERE r.id = %d`
	sql = fmt.Sprintf(sql, region_id)
	_, err := db_session.Sql(sql).Get(&parent_id)
	core.HandleError(err)

	return parent_id.ParentId
}
