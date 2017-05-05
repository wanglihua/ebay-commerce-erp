package controllers

import (
	"fmt"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/region_service"
)

type ErpRegion struct {
	*revel.Controller
}

type region_t struct {
	Id   int64      `xorm:"'id'" json:"id"`
	Name string     `xorm:"'name'" json:"name"`
}

func (c ErpRegion) GetProvinceList() revel.Result {
	var db_session = c.DbSession

	/*
	SELECT r.id AS id, r.province AS name FROM erp_region r WHERE r.parent = 1 ORDER BY r.id ASC;
	 */

	var sql_province_list = `SELECT r.id AS id, r.province AS name FROM erp_region r WHERE r.parent = %d ORDER BY r.id ASC`
	sql_province_list = fmt.Sprintf(sql_province_list, models.REGION_ROOT_ID)

	var province_list = make([]region_t, 0)
	err := db_session.Sql(sql_province_list).Find(&province_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success":       true,
		"province_list": province_list,
	})
}

func (c ErpRegion) GetCityList() revel.Result {
	var db_session = c.DbSession

	var province_id int64
	c.Params.Bind(&province_id, "province")

	/*
	SELECT r.id AS id, concat(r.city, r.district) AS name FROM erp_region r WHERE r.parent = 2 ORDER BY r.id ASC;
	 */

	var sql_city_list = `SELECT r.id AS id, concat(r.city, r.district) AS name FROM erp_region r WHERE r.parent = %d ORDER BY r.id ASC`
	sql_city_list = fmt.Sprintf(sql_city_list, province_id)

	var city_list = make([]region_t, 0)
	err := db_session.Sql(sql_city_list).Find(&city_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"success":   true,
		"city_list": city_list,
	})
}

func (c ErpRegion) GetRegionFullName() revel.Result {
	var db_session = c.DbSession

	var region_id int64
	c.Params.Bind(&region_id, "region_id")

	var region_full_name = region_service.Get_region_full_name(region_id, db_session)

	return c.RenderJson(map[string]interface{}{
		"success":          true,
		"region_full_name": region_full_name,
	})
}
