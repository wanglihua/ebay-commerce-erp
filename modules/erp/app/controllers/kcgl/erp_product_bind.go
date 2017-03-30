package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
)

type ErpProductBind struct {
	*revel.Controller
}

func (c ErpProductBind) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_bind/product_bind_index.html")
}

type ErpProductBindViewItem struct {
	models.ProductBindInfo     `xorm:"extends" json:"b"`
	models.PackageMaterialInfo `xorm:"extends" json:"m"`
	models.StatusInfo          `xorm:"extends" json:"s"`
}

func (c ErpProductBind) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT b.*, m.*, s.*
		FROM erp_product_bind b
		    INNER JOIN erp_package_material m ON b.pkg_mtrl_id = m.id
		    INNER JOIN erp_status s ON b.status_id = s.id
	*/

	query := db_session.
	Select("b.*, m.*, s.*").
		Table("erp_product_bind").Alias("b").
		Join("INNER", []string{"erp_package_material", "m"}, "b.pkg_mtrl_id = m.id").
		Join("INNER", []string{"erp_status", "s"}, "b.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_bind_list := make([]ErpProductBindViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_bind_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpProductBindViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_bind_list,
		Total: count,
	})
}

type ProductBindDetailForm struct {
	ProductBind         models.ProductBindInfo       `json:"product_bind"`
	PackageMaterialList []models.PackageMaterialInfo `json:"package_material_list"`
	StatusList          []models.StatusInfo          `json:"status_list"`
}

func (f *ProductBindDetailForm) IsCreate() bool {
	return f.ProductBind.Id == 0
}

func (f *ProductBindDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductBind.Code).Message("商品编号不能为空！")
	if f.ProductBind.Code != "" {
		validation.MinSize(f.ProductBind.Code, 2).Message("商品编号长度不能小于2！")
	}
	if f.ProductBind.Code != "" {
		validation.MaxSize(f.ProductBind.Code, 100).Message("商品编号长度不能大于100！")
	}

	validation.Required(f.ProductBind.Name).Message("名称不能为空！")
	if f.ProductBind.Name != "" {
		validation.MinSize(f.ProductBind.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ProductBind.Name != "" {
		validation.MaxSize(f.ProductBind.Name, 200).Message("名称长度不能大于200！")
	}

	validation.Required(f.ProductBind.BindList).Message("捆绑商品列表不能为空！")
	if f.ProductBind.BindList != "" {
		validation.MinSize(f.ProductBind.BindList, 1).Message("捆绑商品列表长度不能小于1！")
	}
	if f.ProductBind.BindList != "" {
		validation.MaxSize(f.ProductBind.BindList, 500).Message("捆绑商品列表长度不能大于500！")
	}

	validation.Required(f.ProductBind.PkgMtrlId).Message("包材不能为空！")

	validation.Required(f.ProductBind.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductBind) DetailData() revel.Result {
	db_session := c.DbSession

	var product_bind_id int64
	c.Params.Bind(&product_bind_id, "id")

	var product_bind models.ProductBindInfo
	if product_bind_id != 0 {
		has, err := db_session.Id(product_bind_id).Get(&product_bind)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的绑定商品不存在！"})
		}
	} else {
		product_bind.StatusId = models.PRODUCT_BIND_STATUS_NORMAL
	}

	var package_material_list []models.PackageMaterialInfo
	err := db_session.Find(&package_material_list)
	core.HandleError(err)

	status_list := status_service.GetStatusList(db_session, models.STATUS_TYPE_PRODUCT_BIND)

	var detail_form ProductBindDetailForm
	detail_form.ProductBind = product_bind
	detail_form.PackageMaterialList = package_material_list
	detail_form.StatusList = status_list

	return c.RenderJson(detail_form)
}

func (c ErpProductBind) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ProductBindDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_bind := detail_form.ProductBind

	valid, err := check_product_bind_valid(product_bind.BindList, db_session)
	if valid == false {
		return c.RenderJson(map[string]interface{}{
			"success":  false,
			"messsage": err.Error(),
		})
	}

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&product_bind)
		core.HandleError(err)
	} else {
		update_cols := []string{"code", "name", "bind_list", "status_id"}
		affected, err = db_session.Id(product_bind.Id).Cols(update_cols...).Update(&product_bind)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductBind) Delete() revel.Result {
	db_session := c.DbSession

	product_bind_id_list := make([]int64, 0)
	c.Params.Bind(&product_bind_id_list, "id_list")

	affected, err := db_session.In("id", product_bind_id_list).Delete(new(models.ProductBindInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

// 暂时先写在这里
// bind_produt_list: bind_product_expression
func check_product_bind_valid(bind_product_list string, db_session *xorm.Session) (result bool, err error) {
	if strings.Index(bind_product_list, "，") != -1 { //中文的逗号，看看有没有输入错误
		return false, errors.New("多个商品时不能用中文的逗号（，）分割。")
	}

	product_sku_with_number_list := strings.Split(bind_product_list, ",")
	for _, product_sku_with_number := range product_sku_with_number_list {
		product_sku_number_slice := strings.Split(product_sku_with_number, "*")
		if len(product_sku_number_slice) != 2 {
			return false, errors.New("请用 SKU*数量 表示商品的数量，单个数量就用 *1 的方式！")
		}

		product_sku := product_sku_number_slice[0];
		// 检查SKU的正确性
		sql := `SELECT count(p.id) AS product_count FROM erp_product p WHERE p.sku = '%s'`
		sql = fmt.Sprintf(sql, product_sku)
		var product_count struct {
			Count int64 `xorm:"'product_count'"`
		}
		_, err := db_session.Sql(sql).Get(&product_count)
		core.HandleError(err)

		if product_count.Count == 0 {
			return false, errors.New("请用 SKU*数量 表示商品的数量，商品SKU不存在！")
		}

		product_number_str := product_sku_number_slice[1];
		// 检查数量的格式
		_, err = strconv.ParseInt(product_number_str, 10, 64)
		if err != nil {
			return false, errors.New("请用 SKU*数量 表示商品的数量，数量应为正整数！")
		}
	}

	return true, nil
}
