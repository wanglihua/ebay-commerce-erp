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
	"matrix/modules/erp/models/ddpz_models"
	"bytes"
	"time"
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/excel_service"
)

type ErpLogisticsFeeRate2 struct {
	*revel.Controller
}

func (c ErpLogisticsFeeRate2) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	return c.RenderTemplate("erp/ddpz/logistics_fee_rate2/logistics_fee_rate2_index.html")
}

type logistics_fee_rate2_view_item_t struct {
	LogisticsFeeRate ddpz_models.LogisticsFeeRate2Info `xorm:"extends" json:"fr"`
	LogisticsMode    ddpz_models.LogisticsModeInfo     `xorm:"extends" json:"lm"`
	LogisticsCompany ddpz_models.LogisticsCompanyInfo  `xorm:"extends" json:"lc"`
}

func (c ErpLogisticsFeeRate2) ListData() revel.Result {
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

	/*
		SELECT fr.*, lm.*, lc.*
		FROM erp_logistics_fee_rate2 fr
		    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
		    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
		WHERE lc.name LIKE '%abc%' AND lm.name LIKE '%bcd%' AND fr.postcode LIKE '%123%'
	 */

	sb := core.NewSqlBuilder()
	sb.Select("fr.*, lm.*, lc.*")
	sb.Count("count(fr.id) AS fee_rate_count")
	sb.From(`erp_logistics_fee_rate2 fr
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

	logistics_fee_rate_list := make([]logistics_fee_rate2_view_item_t, 0, limit)
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

type LogisticsFeeRate2DetailForm struct {
	LogisticsFeeRate     ddpz_models.LogisticsFeeRate2Info  `json:"logistics_fee_rate"`
	LogisticsCompanyList []ddpz_models.LogisticsCompanyInfo `json:"logistics_company_list"`
	LogisticsModeList    []ddpz_models.LogisticsModeInfo    `json:"logistics_mode_list"`
	CompanyId            int64                              `json:"company_id"`
}

func (f *LogisticsFeeRate2DetailForm) IsCreate() bool {
	return f.LogisticsFeeRate.Id == 0
}

func (f *LogisticsFeeRate2DetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsFeeRate.ModeId).Message("物流公司不能为空！")

	validation.Required(f.LogisticsFeeRate.Postcode).Message("邮编不能为空！")
	if f.LogisticsFeeRate.Postcode != "" {
		validation.MinSize(f.LogisticsFeeRate.Postcode, 2).Message("邮编长度不能小于2！")
	}
	if f.LogisticsFeeRate.Postcode != "" {
		validation.MaxSize(f.LogisticsFeeRate.Postcode, 100).Message("邮编长度不能大于100！")
	}

	if f.LogisticsFeeRate.Description != "" {
		validation.MaxSize(f.LogisticsFeeRate.Description, 500).Message("描述长度不能大于500！")
	}

	if f.LogisticsFeeRate.Fee01 < 0 ||
		f.LogisticsFeeRate.Fee02 < 0 ||
		f.LogisticsFeeRate.Fee03 < 0 ||
		f.LogisticsFeeRate.Fee04 < 0 ||
		f.LogisticsFeeRate.Fee05 < 0 ||
		f.LogisticsFeeRate.Fee06 < 0 ||
		f.LogisticsFeeRate.Fee07 < 0 ||
		f.LogisticsFeeRate.Fee08 < 0 ||
		f.LogisticsFeeRate.Fee09 < 0 ||
		f.LogisticsFeeRate.Fee10 < 0 ||
		f.LogisticsFeeRate.Fee11 < 0 ||
		f.LogisticsFeeRate.Fee12 < 0 {
		validation.Error("收费不能小于零！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsFeeRate2) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	var logistics_fee_rate_id int64
	c.Params.Bind(&logistics_fee_rate_id, "id")

	var logistics_fee_rate ddpz_models.LogisticsFeeRate2Info
	if logistics_fee_rate_id != 0 {
		has, err := db_session.Id(logistics_fee_rate_id).Get(&logistics_fee_rate)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流费率不存在！"})
		}
	}

	var detail_form LogisticsFeeRate2DetailForm
	detail_form.LogisticsFeeRate = logistics_fee_rate

	var logistics_company_list = make([]ddpz_models.LogisticsCompanyInfo, 0)
	err := db_session.Find(&logistics_company_list)
	core.HandleError(err)
	detail_form.LogisticsCompanyList = logistics_company_list

	if logistics_fee_rate_id == 0 {
		if len(logistics_company_list) != 0 {
			var selected_company_id = logistics_company_list[0].Id
			var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
			err = db_session.Where(fmt.Sprintf("company_id = %d", selected_company_id)).Find(&logistics_mode_list)
			core.HandleError(err)

			detail_form.CompanyId = selected_company_id
			detail_form.LogisticsModeList = logistics_mode_list
		}
	} else {
		//get company id from database
		var mode_id = logistics_fee_rate.ModeId
		var logistics_mode ddpz_models.LogisticsModeInfo
		_, err = db_session.Id(mode_id).Get(&logistics_mode)
		core.HandleError(err)
		var company_id = logistics_mode.CompanyId

		//get mode list from database
		var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
		err = db_session.Where(fmt.Sprintf("company_id = %d", company_id)).Find(&logistics_mode_list)
		core.HandleError(err)

		detail_form.CompanyId = company_id
		detail_form.LogisticsModeList = logistics_mode_list
	}

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsFeeRate2) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	var detail_form LogisticsFeeRate2DetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	revel.TRACE.Println("postcode", detail_form.LogisticsFeeRate.Postcode)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_fee_rate := detail_form.LogisticsFeeRate

	var logistics_fee_rate_in_db ddpz_models.LogisticsFeeRate2Info

	var affected int64
	if detail_form.IsCreate() {
		has, err := db_session.Where("mode_id = ? AND postcode = ?", logistics_fee_rate.ModeId, logistics_fee_rate.Postcode).Limit(1).Get(&logistics_fee_rate_in_db)
		core.HandleError(err)

		if has == true {
			var logistics_mode ddpz_models.LogisticsModeInfo
			_, err := db_session.Id(logistics_fee_rate.ModeId).Get(&logistics_mode)
			core.HandleError(err)

			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("物流方式 %s, 邮编为 %s 的费率已存在！", logistics_mode.Name, logistics_fee_rate.Postcode),
			})
		}

		affected, err = db_session.Insert(&logistics_fee_rate)
		core.HandleError(err)
	} else {
		has, err := db_session.Where("id <> ? AND mode_id = ? AND postcode = ?", logistics_fee_rate.Id, logistics_fee_rate.ModeId, logistics_fee_rate.Postcode).Limit(1).Get(&logistics_fee_rate_in_db)
		core.HandleError(err)

		if has == true {
			var logistics_mode ddpz_models.LogisticsModeInfo
			_, err := db_session.Id(logistics_fee_rate.ModeId).Get(&logistics_mode)
			core.HandleError(err)

			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("物流方式 %s, 邮编为 %s 的费率已存在！", logistics_mode.Name, logistics_fee_rate.Postcode),
			})
		}

		affected, err = db_session.Id(logistics_fee_rate.Id).Update(&logistics_fee_rate)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsFeeRate2) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	logistics_fee_rate_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_fee_rate_id_list, "id_list")

	affected, err := db_session.In("id", logistics_fee_rate_id_list).Delete(new(ddpz_models.LogisticsFeeRate2Info))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpLogisticsFeeRate2) BatchDelete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	var logistics_company_name string
	c.Params.Bind(&logistics_company_name, "logistics_company")
	logistics_company_name = strings.TrimSpace(logistics_company_name)

	var logistics_mode_name string
	c.Params.Bind(&logistics_mode_name, "logistics_mode")
	logistics_mode_name = strings.TrimSpace(logistics_mode_name)

	if logistics_company_name == "" || logistics_mode_name == "" {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请指定 物流公司 和 物流方式！",
		})
	}

	var db_session = c.DbSession

	var sql_get_logistics_mode = `SELECT m.*
FROM erp_logistics_mode m
    INNER JOIN erp_logistics_company c ON c.id = m.company_id
WHERE c.name = '%s' AND m.name = '%s'
LIMIT 1;`

	sql_get_logistics_mode = fmt.Sprintf(sql_get_logistics_mode, logistics_company_name, logistics_mode_name)

	var logistics_mode ddpz_models.LogisticsModeInfo
	has, err := db_session.Sql(sql_get_logistics_mode).Get(&logistics_mode)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "指定的 物流公司 物流方式 不存在！",
		})
	}

	_, err = db_session.Where("mode_id = ?", logistics_mode.Id).Delete(new(ddpz_models.LogisticsFeeRate2Info))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: fmt.Sprintf("物流方式 %s 所有费率数据删除成功!", logistics_mode.Name)})
}

// 这个的权限检查就不用加上了，别的多个地方会调用
func (c ErpLogisticsFeeRate2) GetLogisticsModeListByCompanyId() revel.Result {
	db_session := c.DbSession

	var company_id int64
	c.Params.Bind(&company_id, "company_id")

	var logistics_mode_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.Where(fmt.Sprintf("company_id = %d", company_id)).Find(&logistics_mode_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"mode_list": logistics_mode_list,
	})
}

func (c ErpLogisticsFeeRate2) Import() revel.Result {
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
		LogisticsFeeRate ddpz_models.LogisticsFeeRate2Info
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
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			logistics_fee_rate.CompanyName = cell_val
			col_index++

			// 物流方式名称
			//cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.ModeName = cell_val
			col_index++

			// 邮编
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Postcode = cell_val
			col_index++

			// 物流公司 物流方式 邮编 0-500g 500g-1kg 1kg-2kg 2kg-3kg 3kg-4kg 4kg-5kg 5kg-7kg 7kg-10kg 10kg-15kg 15kg-22kg 22kg-30kg 30kg-50kg
			// 是否为偏远地区 偏远地区加收费用

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee01, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee02, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee03, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee04, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee05, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee06, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee07, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee08, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee09, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee10, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee11, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_fee_rate.LogisticsFeeRate.Fee12, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 是否为偏远地区
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
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
				// cell = row.Cells[col_index]
				cell = sheet.Cell(row_index, col_index)
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
		CompanyName string `xorm:"'company_name'"`
		ModeId      int64  `xorm:"'mode_id'"`
		ModeName    string `xorm:"'mode_name'"`
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

		if logistics_fee_rate.Fee01 < 0 ||
			logistics_fee_rate.Fee02 < 0 ||
			logistics_fee_rate.Fee03 < 0 ||
			logistics_fee_rate.Fee04 < 0 ||
			logistics_fee_rate.Fee05 < 0 ||
			logistics_fee_rate.Fee06 < 0 ||
			logistics_fee_rate.Fee07 < 0 ||
			logistics_fee_rate.Fee08 < 0 ||
			logistics_fee_rate.Fee09 < 0 ||
			logistics_fee_rate.Fee10 < 0 ||
			logistics_fee_rate.Fee11 < 0 ||
			logistics_fee_rate.Fee12 < 0 {
			panic("收费不能小于零！")
		}

		logistics_fee_rate.ModeId = mode_id

		var logistics_fee_rate_in_db ddpz_models.LogisticsFeeRate2Info
		has, err := db_session.Where("mode_id = ? AND postcode = ?", logistics_fee_rate.ModeId, logistics_fee_rate.Postcode).Limit(1).Get(&logistics_fee_rate_in_db)
		core.HandleError(err)

		if has == false {
			_, err = db_session.Insert(&logistics_fee_rate)
			core.HandleError(err)
		} else {
			_, err = db_session.Where("mode_id = ? AND postcode = ?", logistics_fee_rate.ModeId, logistics_fee_rate.Postcode).Update(&logistics_fee_rate)
			core.HandleError(err)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpLogisticsFeeRate2) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_FEE_RATE, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("fr.id IN (%s)", strings.Join(id_str_list, ",")))
	}

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

	/*
		SELECT fr.*, lm.*, lc.*
		FROM erp_logistics_fee_rate2 fr
		    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
		    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id
		WHERE lc.name LIKE '%abc%' AND lm.name LIKE '%bcd%' AND fr.postcode LIKE '%123%'
	 */

	sb := core.NewSqlBuilder()
	sb.Select("fr.*, lm.*, lc.*")
	sb.Count("count(fr.id) AS fee_rate_count")
	sb.From(`erp_logistics_fee_rate2 fr
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
	sb.OrderBy("fr.id ASC")

	logistics_fee_rate_list := make([]logistics_fee_rate2_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_fee_rate_list)
	core.HandleError(err)

	xlsx_file := c.export_fee_rate_list_to_excel(db_session, logistics_fee_rate_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "物流费率2导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpLogisticsFeeRate2) export_fee_rate_list_to_excel(db_session *xorm.Session, logistics_fee_rate_list []logistics_fee_rate2_view_item_t) *xlsx.File {

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		excel_service.Set_head_cell_style(cell)
	}

	//统一设置内容单元格格式的函数
	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		excel_service.Set_text_cell_style(cell, horizontal_align)
	}

	head_row_height := 0.8
	common_row_height := 0.6

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	//添加工作表
	sheet, err := xlsx_file.AddSheet("物流费率")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	var product_status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT)).Find(&product_status_list)
	core.HandleError(err)

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "物流公司"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 30)
		set_head_cell_style(cell)
		cell.Value = "物流方式"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "邮编"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "0-500g"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "500g-1kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "1kg-2kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "2kg-3kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "3kg-4kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "4kg-5kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "5kg-7kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "7kg-10kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "10kg-15kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "15kg-22kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "22kg-30kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "30kg-50kg"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "是否为偏远地区"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "偏远地区加收运费(澳币)"
		col_index++
	}

	for _, fee_rate := range logistics_fee_rate_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 物流公司
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fee_rate.LogisticsCompany.Name

		// 物流方式
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fee_rate.LogisticsMode.Name

		// 邮编
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fee_rate.LogisticsFeeRate.Postcode

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee01)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee02)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee03)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee04)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee05)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee06)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee07)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee08)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee09)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee10)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee11)

		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Fee12)

		// 是否为偏远地区
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		if fee_rate.LogisticsFeeRate.Sfpydq == 1 {
			cell.Value = "是"
		} else {
			cell.Value = "否"
		}

		// 偏远地区加收运费(澳币)
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = fmt.Sprintf("%0.2f", fee_rate.LogisticsFeeRate.Pydqjsyf)
	}

	return xlsx_file
}
