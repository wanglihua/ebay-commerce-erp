package controllers

import (
	"github.com/revel/revel"
	"matrix/modules/erp/service/region_service"
)

type ErpRegion struct {
	*revel.Controller
}

func (c ErpRegion) GetProvinceList() revel.Result {
	var db_session = c.DbSession

	var province_list = region_service.Get_province_list(db_session)

	return c.RenderJson(map[string]interface{}{
		"success":       true,
		"province_list": province_list,
	})
}

func (c ErpRegion) GetCityList() revel.Result {
	var db_session = c.DbSession

	var province_id int64
	c.Params.Bind(&province_id, "province")

	var city_list = region_service.Get_city_list(province_id, db_session)

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
