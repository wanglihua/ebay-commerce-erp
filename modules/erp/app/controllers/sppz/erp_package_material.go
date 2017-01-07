package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
)

type ErpPackageMaterial struct {
	*revel.Controller
}

func (c ErpPackageMaterial) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/package_material/package_material_index.html")
}

type ErpPackageMaterialViewItem struct {
	models.PackageMaterialInfo      `xorm:"extends" json:"p"`
	models.StatusInfo               `xorm:"extends" json:"s"`
}

func (c ErpPackageMaterial) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT p.*, s.* FROM erp_package_material p
		    INNER JOIN erp_status s ON p.status_id = s.id
	 */

	query := db_session.
	Select("p.*, s.*").
		Table("erp_package_material").Alias("p").
		Join("INNER", []string{"erp_status", "s"}, "p.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("p.id")
	}

	package_material_list := make([]ErpPackageMaterialViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&package_material_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpPackageMaterialViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  package_material_list,
		Total: count,
	})
}

type PackageMaterialDetailForm struct {
	PackageMaterial models.PackageMaterialInfo  `json:"package_material"`
	StatusList      []models.StatusInfo         `json:"status_list"`
}

func (f *PackageMaterialDetailForm) IsCreate() bool {
	return f.PackageMaterial.Id == 0
}

func (f *PackageMaterialDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.PackageMaterial.Name).Message("名称不能为空！")
	if f.PackageMaterial.Name != "" {
		validation.MinSize(f.PackageMaterial.Name, 2).Message("名称长度不能小于2！")
	}
	if f.PackageMaterial.Name != "" {
		validation.MaxSize(f.PackageMaterial.Name, 100).Message("名称长度不能大于100！")
	}

	if f.PackageMaterial.Desc.Valid == true {
		validation.MaxSize(f.PackageMaterial.Desc.String, 500).Message("描述长度不能大于500！")
	}

	validation.Required(f.PackageMaterial.Level).Message("级别不能为空！")

	validation.Required(f.PackageMaterial.Length).Message("长(cm)不能为空！")

	validation.Required(f.PackageMaterial.Width).Message("宽(cm)不能为空！")

	validation.Required(f.PackageMaterial.Height).Message("高(cm)不能为空！")

	validation.Required(f.PackageMaterial.Weight).Message("重量(克)不能为空！")

	validation.Required(f.PackageMaterial.Price).Message("价格(元)不能为空！")

	validation.Required(f.PackageMaterial.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpPackageMaterial) DetailData() revel.Result {
	db_session := c.DbSession

	var package_material_id int64
	c.Params.Bind(&package_material_id, "id")

	var package_material models.PackageMaterialInfo
	if package_material_id != 0 {
		has, err := db_session.Id(package_material_id).Get(&package_material)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的包材不存在！"})
		}
	}

	status_list := status_service.GetStatusList(db_session, models.STATUS_TYPE_PACKAGE)

	var detail_form PackageMaterialDetailForm
	detail_form.PackageMaterial = package_material
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpPackageMaterial) Save() revel.Result {
	db_session := c.DbSession

	var detail_form PackageMaterialDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	package_material := detail_form.PackageMaterial

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&package_material)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(package_material.Id).Update(&package_material)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpPackageMaterial) Delete() revel.Result {
	db_session := c.DbSession

	package_material_id_list := make([]int64, 0)
	c.Params.Bind(&package_material_id_list, "id_list")

	affected, err := db_session.In("id", package_material_id_list).Delete(new(models.PackageMaterialInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
