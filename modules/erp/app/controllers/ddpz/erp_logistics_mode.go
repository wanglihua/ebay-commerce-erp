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
	"matrix/modules/erp/models/ddpz_models"
	"matrix/modules/erp/models/sppz_models"
)

type ErpLogisticsMode struct {
	*revel.Controller
}

func (c ErpLogisticsMode) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_MODE, c.Session)

	return c.RenderTemplate("erp/ddpz/logistics_mode/logistics_mode_index.html")
}

func (c ErpLogisticsMode) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_MODE, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type logistics_mode_view_item_t struct {
		ddpz_models.LogisticsModeInfo          `xorm:"extends" json:"m"`
		ddpz_models.LogisticsCompanyInfo       `xorm:"extends" json:"c"`
		sppz_models.StockInfo                  `xorm:"extends" json:"sk"`
		models.StatusInfo                 `xorm:"extends" json:"st"`
	}

	/*
		SELECT m.*, c.*, sk.*, st.*
		FROM erp_logistics_mode m
		    INNER JOIN erp_logistics_company c ON c.id = m.company_id
		    INNER JOIN erp_stock sk ON sk.id = m.stock_id
		    INNER JOIN erp_status st ON st.id = m.status_id
	*/

	sb := core.NewSqlBuilder()
	sb.Select("m.*, c.*, sk.*, st.*")
	sb.Count("count(m.id) as mode_count")
	sb.From(`erp_logistics_mode m
		    INNER JOIN erp_logistics_company c ON c.id = m.company_id
		    INNER JOIN erp_stock sk ON sk.id = m.stock_id
		    INNER JOIN erp_status st ON st.id = m.status_id`)

	sb.Where(filter)

	//query extra filter here

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("m.id ASC")
	}

	sb.Limit(offset, limit)

	logistics_mode_list := make([]logistics_mode_view_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&logistics_mode_list)
	core.HandleError(err)

	type mode_count_t struct {
		Count int64 `xorm:"'mode_count'"`
	}

	var mode_count mode_count_t
	_, err = db_session.Sql(sb.SqlCount()).Get(&mode_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  logistics_mode_list,
		Total: mode_count.Count,
	})
}

type LogisticsModeDetailForm struct {
	LogisticsMode         ddpz_models.LogisticsModeInfo       `json:"logistics_mode"`
	LogisticsCompanyList  []ddpz_models.LogisticsCompanyInfo  `json:"logistics_company_list"`
	StockList             []sppz_models.StockInfo             `json:"stock_list"`
	LogisticsModeTypeList []ddpz_models.LogisticsModeTypeInfo `json:"logistics_mode_type_list"`
	StatusList            []models.StatusInfo            `json:"status_list"`
}

func (f *LogisticsModeDetailForm) IsCreate() bool {
	return f.LogisticsMode.Id == 0
}

func (f *LogisticsModeDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LogisticsMode.Name).Message("名称不能为空！")
	if f.LogisticsMode.Name != "" {
		validation.MinSize(f.LogisticsMode.Name, 2).Message("名称长度不能小于2！")
	}
	if f.LogisticsMode.Name != "" {
		validation.MaxSize(f.LogisticsMode.Name, 100).Message("名称长度不能大于100！")
	}

	if f.LogisticsMode.Desc.Valid == true {
		validation.MaxSize(f.LogisticsMode.Desc.String, 500).Message("描述长度不能大于500！")
	}

	return validation.HasErrors() == false
}

func (c ErpLogisticsMode) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_MODE, c.Session)

	db_session := c.DbSession

	var logistics_mode_id int64
	c.Params.Bind(&logistics_mode_id, "id")

	var logistics_mode ddpz_models.LogisticsModeInfo
	if logistics_mode_id != 0 {
		has, err := db_session.Id(logistics_mode_id).Get(&logistics_mode)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的物流方式不存在！"})
		}
	} else {
		logistics_mode.Tjzysjzqx = 1
		logistics_mode.StatusId = models.LOGISTICS_MODE_ENABLED
		logistics_mode.TypeId = ddpz_models.LOGISTICS_MODE_TYPE_PACKAGE
	}

	var logistics_company_list = make([]ddpz_models.LogisticsCompanyInfo, 0)
	err := db_session.Find(&logistics_company_list)
	core.HandleError(err)

	var stock_list = make([]sppz_models.StockInfo, 0)
	err = db_session.Where(fmt.Sprintf("stock_type_id = %d", sppz_models.STOCK_TYPE_AU)).Find(&stock_list)
	core.HandleError(err)

	var status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_LOGISTICS_MODE)).Find(&status_list)
	core.HandleError(err)

	var logistics_mode_type_list = make([]ddpz_models.LogisticsModeTypeInfo, 0)
	err = db_session.OrderBy("id ASC").Find(&logistics_mode_type_list)
	core.HandleError(err)

	var detail_form LogisticsModeDetailForm
	detail_form.LogisticsMode = logistics_mode
	detail_form.LogisticsCompanyList = logistics_company_list
	detail_form.StockList = stock_list
	detail_form.StatusList = status_list
	detail_form.LogisticsModeTypeList = logistics_mode_type_list

	return c.RenderJson(detail_form)
}

func (c ErpLogisticsMode) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_MODE, c.Session)

	db_session := c.DbSession

	var detail_form LogisticsModeDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	logistics_mode := detail_form.LogisticsMode

	var affected int64
	if detail_form.IsCreate() {
		name_count, err := db_session.Where(fmt.Sprintf("company_id = %d AND name = '%s'", logistics_mode.CompanyId, logistics_mode.Name)).Count(new(ddpz_models.LogisticsModeInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
		}

		affected, err = db_session.Insert(&logistics_mode)
		core.HandleError(err)
	} else {
		name_count, err := db_session.Where(fmt.Sprintf("id <> %d AND company_id = %d AND name = '%s'", logistics_mode.Id, logistics_mode.CompanyId, logistics_mode.Name)).Count(new(ddpz_models.LogisticsModeInfo))
		core.HandleError(err)
		if name_count != 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "保存失败，名称已存在！"})
		}

		affected, err = db_session.Id(logistics_mode.Id).Update(&logistics_mode)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLogisticsMode) DeleteConfirm() revel.Result {
	var db_session = c.DbSession

	logistics_mode_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_mode_id_list, "id_list")

	if len(logistics_mode_id_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "请选择要删除的物流方式！",
		})
	}

	var mode_id_str_list = core.Int_list_to_str_list(logistics_mode_id_list)

	// 查看物流方式是否已经在匹配物流方式的订单里面使用

	var mode_use_count struct {
		Count int64 `xorm:"'ol_count'"`
	}
	/*
SELECT count(ol.id) AS ol_count
FROM erp_order_logistics ol
WHERE ol.mode_id IN (1, 2, 3)
	 */
	var sql_mode_use_count = `SELECT count(ol.id) AS ol_count
FROM erp_order_logistics ol
WHERE ol.mode_id IN (%s)`
	sql_mode_use_count = fmt.Sprintf(sql_mode_use_count, strings.Join(mode_id_str_list, ","))

	_, err := db_session.Sql(sql_mode_use_count).Get(&mode_use_count)
	core.HandleError(err)

	if mode_use_count.Count != 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "物流方式当前有匹配成功的订单，如果删除，这些订单的物流信息将不全，是否继续删除？",
		})
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "无需确认，直接删除！",
	})
}

func (c ErpLogisticsMode) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_MODE, c.Session)

	db_session := c.DbSession

	logistics_mode_id_list := make([]int64, 0)
	c.Params.Bind(&logistics_mode_id_list, "id_list")

	/*
	// 如果数据在使用的就不用被删除
	var fee_rate_count struct {
		Count int64 `xorm:"'fee_rate_count'"`
	}
	sql := "SELECT count(fr.id) AS fee_rate_count FROM erp_logistics_fee_rate fr WHERE fr.mode_id IN (%s)"
	mode_id_str_list := core.Int_list_to_str_list(logistics_mode_id_list)
	sql = fmt.Sprintf(sql, strings.Join(mode_id_str_list, ","))
	_, err := db_session.Sql(sql).Get(&fee_rate_count)
	core.HandleError(err)
	if fee_rate_count.Count != 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "数据在别处被引用，不能删除！",
		})
	}
	*/

	// 一并删除物流费率
	_, err := db_session.In("mode_id", logistics_mode_id_list).Delete(new(ddpz_models.LogisticsFeeRateInfo))
	core.HandleError(err)

	_, err = db_session.In("id", logistics_mode_id_list).Delete(new(ddpz_models.LogisticsModeInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "数据删除成功!"})
}

func (c ErpLogisticsMode) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.DDPZ_LOGISTICS_MODE, c.Session)

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
		StockName     string
		LogisticsMode ddpz_models.LogisticsModeInfo
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

			//todo 这里也有可能会出现问题的，如果最后几列用户没有填怎么办
			if len(row.Cells) < 18 {
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
			logistics_mode.CompanyName = cell_val
			col_index++

			// 仓库名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.StockName = cell_val
			col_index++

			// 物流方式名称
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Name = cell_val
			col_index++

			// 是否考虑 体积重
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
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
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Tjzgs, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 体积重 系数
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			cell_val = strings.TrimSpace(cell_val)
			core.HandleError(err)
			logistics_mode.LogisticsMode.Tjzxs, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 体积重与 实际重取向 （取大 or 取小） 取大(1) or 取小(2)
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
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
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
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
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Tjzsjzbz, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 最长边是否 有限制
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
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
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Zcbxzcd, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 重量是否 有限制
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
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
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Zlbncg, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 是否允许 超长
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Sfyxcc = 1
			} else {
				logistics_mode.LogisticsMode.Sfyxcc = 0
			}
			col_index++

			// 超长收费阶梯（cm）
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Ccsfjt, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 超长加收多少 运费（澳币）
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			logistics_mode.LogisticsMode.Ccjsyf, _ = strconv.ParseFloat(cell_val, 64)
			col_index++

			// 是否 允许超重
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
			core.HandleError(err)
			cell_val, err = cell.String()
			if strings.TrimSpace(cell_val) == "是" {
				logistics_mode.LogisticsMode.Sfyxcz = 1
			} else {
				logistics_mode.LogisticsMode.Sfyxcz = 0
			}
			col_index++

			// 超重加收多少 运费（澳币）
			// cell = row.Cells[col_index]
			cell = sheet.Cell(row_index, col_index)
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
			"message": "Excel文件中导入物流方式列表为空！",
		})
	}

	for _, logistics_mode_import := range logistics_mode_list {

		var logistics_mode = logistics_mode_import.LogisticsMode

		// 物流公司
		var logistics_company ddpz_models.LogisticsCompanyInfo
		has, err := db_session.Where(fmt.Sprintf("name = '%s'", logistics_mode_import.CompanyName)).Get(&logistics_company)
		core.HandleError(err)

		if has == false {
			logistics_company.Name = logistics_mode_import.CompanyName
			logistics_company.Desc = core.NullStringFrom(logistics_mode_import.CompanyName)

			_, err := db_session.Insert(&logistics_company)
			core.HandleError(err)
		}

		logistics_mode.CompanyId = logistics_company.Id

		// 仓库
		var stock sppz_models.StockInfo
		has, err = db_session.Where(fmt.Sprintf("name = '%s'", logistics_mode_import.StockName)).Get(&stock)
		core.HandleError(err)
		if has == false {
			/*
			return c.RenderJson(map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("仓库名称 %s 不存在！", logistics_mode_import.StockName),
			})
			*/

			panic(fmt.Sprintf("仓库名称 %s 不存在！", logistics_mode_import.StockName))
		}

		logistics_mode.StockId = stock.Id
		logistics_mode.StatusId = models.LOGISTICS_MODE_ENABLED

		_, err = db_session.Insert(&logistics_mode)
		core.HandleError(err)

	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpLogisticsMode) GetListByStock() revel.Result {
	var db_session = c.DbSession

	var stock_id int64
	c.Params.Bind(&stock_id, "stock_id")

	var model_list = make([]ddpz_models.LogisticsModeInfo, 0)
	err := db_session.Where(fmt.Sprintf("stock_id = %d", stock_id)).Find(&model_list)
	core.HandleError(err)

	return c.RenderJson(map[string]interface{}{
		"mode_list": model_list,
	})
}
