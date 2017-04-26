package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
)

type ErpProductCate struct {
	*revel.Controller
}

func (c ErpProductCate) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CATE, c.Session)

	return c.RenderTemplate("erp/sppz/product_cate/product_cate_index.html")
}

type ErpProductCateViewItem struct {
	models.ProductCateInfo         `xorm:"extends" json:"c"`
	Parent models.ProductCateInfo  `xorm:"extends" json:"p"`
	models.StatusInfo              `xorm:"extends" json:"s"`
}

func (c ErpProductCate) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CATE, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT c.*, p.*, s.*
		FROM erp_product_cate c
		    INNER JOIN erp_product_cate p ON c.parent_id = p.id
		    INNER JOIN erp_status s ON c.status_id = s.id
	*/

	query := db_session.
	Select("c.*, p.*, s.*").
		Table("erp_product_cate").Alias("c").
		Join("LEFT", []string{"erp_product_cate", "p"}, "c.parent_id = p.id").
		Join("INNER", []string{"erp_status", "s"}, "c.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("c.id")
	}

	product_cate_list := make([]ErpProductCateViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_cate_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpProductCateViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_cate_list,
		Total: count,
	})
}

type ProductCateDetailForm struct {
	ProductCate    models.ProductCateInfo   `json:"product_cate"`
	ParentCateList []models.ProductCateInfo `json:"parent_cate_list"`
	StatusList     []models.StatusInfo      `json:"status_list"`
}

func (f *ProductCateDetailForm) IsCreate() bool {
	return f.ProductCate.Id == 0
}

func (f *ProductCateDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ProductCate.Name).Message("名称不能为空！")
	if f.ProductCate.Name != "" {
		validation.MinSize(f.ProductCate.Name, 2).Message("名称长度不能小于2！")
	}
	if f.ProductCate.Name != "" {
		validation.MaxSize(f.ProductCate.Name, 100).Message("名称长度不能大于100！")
	}

	if f.ProductCate.Desc.Valid == true {
		validation.MaxSize(f.ProductCate.Desc.String, 500).Message("描述长度不能大于500！")
	}

	validation.Required(f.ProductCate.StatusId).Message("状态不能为空！")

	return validation.HasErrors() == false
}

func (c ErpProductCate) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CATE, c.Session)

	db_session := c.DbSession

	var product_cate_id int64
	c.Params.Bind(&product_cate_id, "id")

	var product_cate models.ProductCateInfo
	if product_cate_id != 0 {
		has, err := db_session.Id(product_cate_id).Get(&product_cate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的商品目录不存在！"})
		}
	}

	var parent_cate_list []models.ProductCateInfo
	if product_cate_id == 0 {
		// 新的商品类别的默认状态
		product_cate.StatusId = models.PRODUCT_CATE_STATUS_ENABLED
		err := db_session.Where(fmt.Sprintf("%[1]s is null or %[1]s = 0", "parent_id")).Find(&parent_cate_list)
		core.HandleError(err)
	} else {
		// 编辑是的父类别不应该和自己类别是一样的
		err := db_session.
		Where(fmt.Sprintf("(%[1]s is null or %[1]s = 0) and (%[2]s <> %[3]d)", "parent_id", "id", product_cate_id)).
			Find(&parent_cate_list)
		core.HandleError(err)
	}

	status_list := status_service.Get_status_list(db_session, models.STATUS_TYPE_PRODUCT_CATE)

	var detail_form ProductCateDetailForm
	detail_form.ProductCate = product_cate
	detail_form.StatusList = status_list
	detail_form.ParentCateList = parent_cate_list

	return c.RenderJson(detail_form)
}

func (c ErpProductCate) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CATE, c.Session)

	db_session := c.DbSession

	var detail_form ProductCateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	product_cate := detail_form.ProductCate

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where("name = ?", product_cate.Name).Count(new(models.ProductCateInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		product_cate.FullName = c.calculate_product_cate_full_name(product_cate, db_session)

		affected, err = db_session.Insert(&product_cate)
		core.HandleError(err)

		// 插入操作日志
		c.insert_product_cate_log(&product_cate, "插入数据", c.DbSession, c.Session)
	} else {
		name_count, err := db_session.Where("id <> ? and name = ?", product_cate.Id, product_cate.Name).Count(new(models.ProductCateInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！" })
		}

		product_cate.FullName = c.calculate_product_cate_full_name(product_cate, db_session)

		affected, err = db_session.Id(product_cate.Id).Update(&product_cate)
		core.HandleError(err)

		// 插入操作日志
		c.insert_product_cate_log(&product_cate, "编辑数据", c.DbSession, c.Session)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpProductCate) calculate_product_cate_full_name(product_cate models.ProductCateInfo, db_session *xorm.Session) string {
	if product_cate.ParentId.Valid == false || product_cate.ParentId.Int64 == 0 {
		return product_cate.Name
	}

	parent_cate_id := product_cate.ParentId.Int64
	var parent_cate models.ProductCateInfo
	_, err := db_session.Id(parent_cate_id).Get(&parent_cate)
	core.HandleError(err)

	return fmt.Sprintf("%s-%s", parent_cate.Name, product_cate.Name)
}

func (c ErpProductCate) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_CATE, c.Session)

	db_session := c.DbSession

	product_cate_id_list := make([]int64, 0)
	c.Params.Bind(&product_cate_id_list, "id_list")

	// 插入删除日志
	for _, product_cate_id := range product_cate_id_list {
		var product_cate models.ProductCateInfo
		has, err := db_session.Id(product_cate_id).Get(&product_cate)
		core.HandleError(err)

		if has {
			c.insert_product_cate_log(&product_cate, "删除数据", c.DbSession, c.Session)
		}
	}

	affected, err := db_session.In("id", product_cate_id_list).Delete(new(models.ProductCateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpProductCate) insert_product_cate_log(product_cate *models.ProductCateInfo, log_desc string, db_session *xorm.Session, web_session revel.Session) {
	var product_cate_log models.ProductCateLogInfo
	product_cate_log.Name = product_cate.FullName
	product_cate_log.Desc = log_desc
	product_cate_log.LoginName = core.GetLoginUser(c.Session).LoginName

	_, err := db_session.Insert(&product_cate_log)
	core.HandleError(err)
}
