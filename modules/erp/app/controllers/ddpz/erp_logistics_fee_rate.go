package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpLogisticsFeeRate struct {
	*revel.Controller
}

func (c ErpLogisticsFeeRate) Index() revel.Result {
	return c.RenderTemplate("erp/ddpz/logistics_fee_rate/logistics_fee_rate_index.html")
}

func (c ErpLogisticsFeeRate) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

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
	 */

	sb := core.NewSqlBuilder()
	sb.Select("fr.*, lm.*, lc.*")
	sb.Count("count(fr.id) AS fee_rate_count")
	sb.From(`erp_logistics_fee_rate fr
    INNER JOIN erp_logistics_mode lm ON lm.id = fr.mode_id
    INNER JOIN erp_logistics_company lc ON lc.id = lm.company_id`)
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

	validation.Required(f.LogisticsFeeRate.Xzjtsf).Message("续重阶梯收费不能为空！")

	validation.Required(f.LogisticsFeeRate.Pydqjsyf).Message("偏远地区加收运费不能为空！")

	if f.LogisticsFeeRate.Desc.Valid == true {
		validation.MaxSize(f.LogisticsFeeRate.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsFeeRate) DetailData() revel.Result {
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
	db_session := c.DbSession

	logistics_fee_rate_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_fee_rate_id_list, "id_list")

	affected, err := db_session.In("id", logistics_fee_rate_id_list).Delete(new(models.LogisticsFeeRateInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpLogisticsFeeRate) Import() revel.Result {
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

	type logistics_mode_t struct {
		CompanyName   string
		LogisticsMode models.LogisticsModeInfo
	}

	logistics_mode_list := make([]logistics_mode_t, 0)

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

			var logistics_mode logistics_mode_t

			col_index := 0

			// 物流公司名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.CompanyName = cell_val
			col_index++

			// 物流方式名称
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Name = cell_val
			col_index++

			// 是否考虑 体积重
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Sfkltjz = 1
			} else {
				logistics_mode.LogisticsMode.Sfkltjz = 0
			}
			col_index++

			// 体积重公式 ( 长X宽X高/250) 250
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Tjzgs, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 体积重 系数
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			cell_val = strings.TrimSpace(cell_val)
			core.HandleError(err)
			logistics_mode.LogisticsMode.Tjzxs, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 体积重与 实际重取向 （取大 or 取小） 取大(1) or 取小(2)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			cell_val = strings.TrimSpace(cell_val)
			core.HandleError(err)
			if strings.TrimSpace(cell_val) == "取大" {
				logistics_mode.LogisticsMode.Tjzysjzqx = 1
			} else if strings.TrimSpace(cell_val) == "取小" {
				logistics_mode.LogisticsMode.Tjzysjzqx = 2
			}
			col_index++

			// 体积重与实际重 是否有比值要求
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Zlsfybzyq = 1
			} else {
				logistics_mode.LogisticsMode.Zlsfybzyq = 0
			}
			col_index++

			// 体积重/实际重 比值
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Tjzsjzbz, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 最长边是否 有限制
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Zcbsfyxz = 1
			} else {
				logistics_mode.LogisticsMode.Zcbsfyxz = 0
			}
			col_index++

			// 最长边限制 长度（cm）
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Zcbxzcd, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 重量是否 有限制
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Zlsfyxz = 1
			} else {
				logistics_mode.LogisticsMode.Zlsfyxz = 0
			}
			col_index++

			// 重量不能 超过（Kg)
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Zlbncg, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 是否允许 超长
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Sfyxcc = 1
			} else {
				logistics_mode.LogisticsMode.Sfyxcc = 0
			}
			col_index++

			// 超长加收多少 运费（澳币）
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Ccjsyf, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 是否 允许超重
			cell = row.Cells[col_index]
			core.HandleError(err)
			cell_val, err = cell.String()
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Sfyxcz = 1
			} else {
				logistics_mode.LogisticsMode.Sfyxcz = 0
			}
			col_index++

			// 超重加收多少 运费（澳币）
			cell = row.Cells[col_index]
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Czjsyf, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			logistics_mode_list = append(logistics_mode_list, logistics_mode)
		}
	}

	if len(logistics_mode_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品列表为空！",
		})
	}

	for _, logistics_mode_import := range logistics_mode_list {

		var logistics_mode = logistics_mode_import.LogisticsMode

		// 物流公司
		var logistics_company models.LogisticsCompanyInfo
		has, err := db_session.Where(fmt.Sprintf("name = '%s'", logistics_mode_import.CompanyName)).Get(&logistics_company)
		core.HandleError(err)

		if has == false {
			logistics_company.Name = logistics_mode_import.CompanyName
			logistics_company.Desc = core.NullStringFrom(logistics_mode_import.CompanyName)

			_, err := db_session.Insert(&logistics_company)
			core.HandleError(err)
		}

		logistics_mode.CompanyId = logistics_company.Id

		_, err = db_session.Insert(&logistics_mode)
		core.HandleError(err)

	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

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
