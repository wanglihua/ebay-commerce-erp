package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
)

type ErpLogisticsFeeRate struct {
	*revel.Controller
}

func (c ErpLogisticsFeeRate) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	return c.RenderTemplate("erp/ddpz/logistics_fee_rate/logistics_fee_rate_index.html")
}

func (c ErpLogisticsFeeRate) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		aoData.push({'name': 'logistics_company_query', 'value': $('#logistics_company_query').val()});
		aoData.push({'name': 'logistics_mode_query', 'value': $('#logistics_mode_query').val()});
		aoData.push({'name': 'postcode_query', 'value': $('#postcode_query').val()});
	 */

	var logistics_company_query string
	var logistics_mode_query string
	var postcode_query string

	c.Params.Bind(&logistics_company_query, "logistics_company_query")
	c.Params.Bind(&logistics_mode_query, "logistics_mode_query")
	c.Params.Bind(&postcode_query, "postcode_query")

	type logistics_fee_rate_view_item_t struct {
		models.LogisticsFeeRateInfo  `xorm:"extends" json:"fr"`
		models.LogisticsModeInfo     `xorm:"extends" json:"lm"`
		models.LogisticsCompanyInfo  `xorm:"extends" json:"lc"`
	}

	/*
		SELECT fr.*, lm.*, lc.*
		FROM erp_logistics_fee_rate fr
		    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
		    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
		WHERE lc.name LIKE '%abc%' AND lm.name LIKE '%bcd%' AND fr.postcode LIKE '%123%'
	 */

	sb := core.NewSqlBuilder()
	sb.Select("fr.*, lm.*, lc.*")
	sb.Count("count(fr.id) AS fee_rate_count")
	sb.From(`erp_logistics_fee_rate fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id`)

	if strings.TrimSpace(logistics_company_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("lc.name LIKE '%%%s%%'", logistics_company_query))
	}

	if strings.TrimSpace(logistics_mode_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("lm.name LIKE '%%%s%%'", logistics_mode_query))
	}

	if strings.TrimSpace(postcode_query) != "" {
		filter = core.SqlAnd(filter, fmt.Sprintf("fr.postcode LIKE '%%%s%%'", postcode_query))
	}

	sb.Where(filter)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("fr.id ASC")
	}

	sb.Limit(offset, limit)

	logistics_fee_rate_list := make([]logistics_fee_rate_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_fee_rate_list)
	core.HandleError(err)

	var fee_rate_count struct {
		Count int64 `xorm:"'fee_rate_count'"`
	}

	_, err = db_session.Sql(sb.SqlCount()).Get(&fee_rate_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_fee_rate_list,
		Total: fee_rate_count.Count,
	})
}

type LogisticsFeeRateDetailForm struct {
	LogisticsFeeRate     models.LogisticsFeeRateInfo   `json:"logistics_fee_rate"`
	LogisticsCompanyList []models.LogisticsCompanyInfo `json:"logistics_company_list"`
	LogisticsModeList    []models.LogisticsModeInfo    `json:"logistics_mode_list"`
	CompanyId            int64                         `json:"company_id"`
}

func (f *LogisticsFeeRateDetailForm) IsCreate() bool {
	return f.LogisticsFeeRate.Id == 0
}

func (f *LogisticsFeeRateDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsFeeRate.ModeId).Message("物流公司不能为空！")

	validation.Required(f.LogisticsFeeRate.Postcode).Message("邮编不能为空！")
	if f.LogisticsFeeRate.Postcode != "" {
		validation.MinSize(f.LogisticsFeeRate.Postcode, 2).Message("邮编长度不能小于2！")
	}
	if f.LogisticsFeeRate.Postcode != "" {
		validation.MaxSize(f.LogisticsFeeRate.Postcode, 100).Message("邮编长度不能大于100！")
	}

	validation.Required(f.LogisticsFeeRate.Qbj).Message("起步价不能为空！")

	validation.Required(f.LogisticsFeeRate.Qbjzl).Message("起步价重量不能为空！")

	validation.Required(f.LogisticsFeeRate.Zdsf).Message("最低收费不能为空！")

	validation.Required(f.LogisticsFeeRate.Xzjt).Message("续重阶梯不能为空！")

	if f.LogisticsFeeRate.Xzjt <= 0 {
		validation.Error("续重阶梯应该大于零！")
	}

	validation.Required(f.LogisticsFeeRate.Xzjtsf).Message("续重阶梯收费不能为空！")

	validation.Required(f.LogisticsFeeRate.Pydqjsyf).Message("偏远地区加收运费不能为空！")

	if f.LogisticsFeeRate.Desc.Valid == true {
		validation.MaxSize(f.LogisticsFeeRate.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsFeeRate) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	var logistics_fee_rate_id int64
	c.Params.Bind(&logistics_fee_rate_id, "id")

	var logistics_fee_rate models.LogisticsFeeRateInfo
	if logistics_fee_rate_id != 0 {
		has, err := db_session.Id(logistics_fee_rate_id).Get(&logistics_fee_rate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流费率不存在！"})
		}
	}

	var detail_form LogisticsFeeRateDetailForm
	detail_form.LogisticsFeeRate = logistics_fee_rate

	var logistics_company_list = make([]models.LogisticsCompanyInfo, 0)
	err := db_session.Find(&logistics_company_list)
	core.HandleError(err)
	detail_form.LogisticsCompanyList = logistics_company_list

	if logistics_fee_rate_id == 0 {
		if len(logistics_company_list) != 0 {
			var selected_company_id = logistics_company_list[0].Id
			var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
			err = db_session.Where(fmt.Sprintf("company_id = %d", selected_company_id)).Find(&logistics_mode_list)
			core.HandleError(err)

			detail_form.CompanyId = selected_company_id
			detail_form.LogisticsModeList = logistics_mode_list
		}
	} else {
		//get company id from database
		var mode_id = logistics_fee_rate.ModeId
		var logistics_mode models.LogisticsModeInfo
		_, err = db_session.Id(mode_id).Get(&logistics_mode)
		core.HandleError(err)
		var company_id = logistics_mode.CompanyId

		//get mode list from database
		var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
		err = db_session.Where(fmt.Sprintf("company_id = %d", company_id)).Find(&logistics_mode_list)
		core.HandleError(err)

		detail_form.CompanyId = company_id
		detail_form.LogisticsModeList = logistics_mode_list
	}

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsFeeRate) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	var detail_form LogisticsFeeRateDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	revel.TRACE.Println("postcode", detail_form.LogisticsFeeRate.Postcode)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_fee_rate := detail_form.LogisticsFeeRate

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&logistics_fee_rate)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(logistics_fee_rate.Id).Update(&logistics_fee_rate)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsFeeRate) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	logistics_fee_rate_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_fee_rate_id_list, "id_list")

	affected, err := db_session.In("id", logistics_fee_rate_id_list).Delete(new(models.LogisticsFeeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpLogisticsFeeRate) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	type logistics_fee_rate_import_t struct {
		CompanyName      string
		ModeName         string
		LogisticsFeeRate models.LogisticsFeeRateInfo
	}

	logistics_fee_rate_list := make([]logistics_fee_rate_import_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var logistics_fee_rate logistics_fee_rate_import_t

			col_index := 0

			if len(row.Cells) < 9 { // 本来是有10列数据的，但最后一列是可以为空的，有时里面的数据列里面没数据，不算在 len(row.Cells) 里面
				break
			}

			// 物流公司名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			logistics_fee_rate.CompanyName = cell_val
			col_index++

			// 物流方式名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.ModeName = cell_val
			col_index++

			// 邮编
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Postcode = cell_val
			col_index++

			// 起步价(澳币)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Qbj, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 起步价重量(Kg)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Qbjzl, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 最低收费(澳币)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Zdsf, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 续重阶梯(Kg)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Xzjt, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 续重阶梯收费(澳币)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Xzjtsf, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 是否为偏远地区
			cell = row.Cells[col_index]
			core.HandleError(err)
			cell_val, err = cell.String()
			if strings.TrimSpace(cell_val) == "是" {
				logistics_fee_rate.LogisticsFeeRate.Sfpydq = 1
			} else {
				logistics_fee_rate.LogisticsFeeRate.Sfpydq = 0
			}
			col_index++

			// 偏远地区加收运费(澳币)
			if col_index <= len(row.Cells)-1 {
				cell = row.Cells[col_index]
				cell_val, err = cell.String()
				core.HandleError(err)
				cell_val = strings.TrimSpace(cell_val)
				logistics_fee_rate.LogisticsFeeRate.Pydqjsyf, _ = strconv.ParseFloat(cell_val, 64)
			}
			col_index++

			logistics_fee_rate_list = append(logistics_fee_rate_list, logistics_fee_rate)
		}
	}

	if len(logistics_fee_rate_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入费率列表为空！",
		})
	}

	/*

	加载所有物流公司和物流方式的名称的组合，方便接下来的验证，加快速度

SELECT
    DISTINCT
    lc.name AS company_name,
    lm.name AS mode_name
FROM erp_logistics_mode lm
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id

	 */

	type logistics_company_and_mode_t struct {
		CompanyName string  `xorm:"'company_name'"`
		ModeId      int64   `xorm:"'mode_id'"`
		ModeName    string  `xorm:"'mode_name'"`
	}

	sql := `SELECT
    DISTINCT
    lc.name AS company_name,
    lm.id   AS mode_id,
    lm.name AS mode_name
FROM erp_logistics_mode lm
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
`

	var logistics_company_and_mode_list = make([]logistics_company_and_mode_t, 0)
	err = db_session.Sql(sql).Find(&logistics_company_and_mode_list)
	core.HandleError(err)

	var is_logistics_company_and_mode_exist = func(company_name string, mode_name string) (has_return bool, mode_id_return int64) {
		for _, logistics_company_and_mode := range logistics_company_and_mode_list {
			if strings.TrimSpace(company_name) == strings.TrimSpace(logistics_company_and_mode.CompanyName) &&
				strings.TrimSpace(mode_name) == strings.TrimSpace(logistics_company_and_mode.ModeName) {
				return true, logistics_company_and_mode.ModeId
			}
		}
		return false, 0
	}

	for _, logistics_fee_rate_import := range logistics_fee_rate_list {

		var logistics_fee_rate = logistics_fee_rate_import.LogisticsFeeRate

		/* 使用了更好的方式
		// 物流方式
		//SELECT lm.id AS mode_id
		//FROM erp_logistics_mode lm
		//    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
		//WHERE lc.name = 'company_name' AND lm.name = 'mode_name'

		var mode_id struct {
			ModeId int64 `xorm:"'mode_id'"`
		}
		sql := `SELECT lm.id AS mode_id
FROM erp_logistics_mode lm
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
WHERE lc.name = '%s' AND lm.name = '%s'`
		sql = fmt.Sprintf(sql, logistics_fee_rate_import.CompanyName, logistics_fee_rate_import.ModeName)
		has, err := db_session.Sql(sql).Get(&mode_id)
		core.HandleError(err)

		*/

		has, mode_id := is_logistics_company_and_mode_exist(logistics_fee_rate_import.CompanyName, logistics_fee_rate_import.ModeName)
		if has == false {
			/*
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": "导入文件中的物流公司和物流方式不存在！",
			})
			*/

			panic(fmt.Sprintf("导入文件中的物流公司 %s 和物流方式 %s不存在", logistics_fee_rate_import.CompanyName, logistics_fee_rate_import.ModeName))
		}

		if logistics_fee_rate.Xzjt < 0 {
			panic(fmt.Sprintf("续重阶梯不能小于零！ %d", logistics_fee_rate.Xzjt))
		}

		logistics_fee_rate.ModeId = mode_id

		_, err = db_session.Insert(&logistics_fee_rate)
		core.HandleError(err)

	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

// 这个的权限检查就不用加上了，别的多个地方会调用
func (c ErpLogisticsFeeRate) GetLogisticsModeListByCompanyId() revel.Result {
	db_session := c.DbSession

	var company_id int64
	c.Params.Bind(&company_id, "company_id")

	var logistics_mode_list = make([]models.LogisticsModeInfo, 0)
	err := db_session.Where(fmt.Sprintf("company_id = %d", company_id)).Find(&logistics_mode_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"mode_list": logistics_mode_list,
	})
}
