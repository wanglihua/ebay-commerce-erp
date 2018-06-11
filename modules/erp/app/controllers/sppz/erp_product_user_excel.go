package controllers

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/tealeg/xlsx"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"strings"
	"matrix/modules/erp/models/kcgl_models"
	"bytes"
	"time"
	auth_models "matrix/modules/auth/models"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/excel_service"
)

type ErpProductUserExcel struct {
	*revel.Controller
}

func (c ErpProductUserExcel) Import() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_USER, c.Session)

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

	type product_user_t struct {
		sku                       string
		sale_user_login_name      string
		service_user_login_name   string
		dev_user_login_name       string
		replenish_user_login_name string
		product                   kcgl_models.ProductInfo
		sale_user_id              int64
		service_user_id           int64
		dev_user_id               int64
		replenish_user_id         int64
	}

	product_user_list := make([]product_user_t, 0)

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	// 会导入所有 sheets
	sheet_count := len(xlsx_file.Sheets)
	//revel.TRACE.Printf("sheet_count: %d\n", sheet_count)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		//revel.TRACE.Printf("max_row: %d\n", max_row)
		//略过标题行
		for row_index := 1; row_index < max_row; row_index ++ {
			// row := sheet.Rows[row_index]

			var cell *xlsx.Cell
			var cell_val string
			var err error

			var product_user product_user_t

			col_index := 1 //从第二列开始导入，略去第一列的编号

			//SKU
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			if cell_val == "" {
				break
			}
			product_user.sku = cell_val
			col_index++

			// sale_user_login_name
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_user.sale_user_login_name = cell_val
			col_index++

			// service_user_login_name
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_user.service_user_login_name = cell_val
			col_index++

			// dev_user_login_name
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_user.dev_user_login_name = cell_val
			col_index++

			// replenish_user_login_name
			cell = sheet.Cell(row_index, col_index)
			cell_val, err = cell.String()
			core.HandleError(err)
			cell_val = strings.TrimSpace(cell_val)
			product_user.replenish_user_login_name = cell_val
			col_index++

			product_user_list = append(product_user_list, product_user)
		}
	}

	if len(product_user_list) == 0 {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "Excel文件中导入商品负责人列表为空！",
		})
	}

	for index, product_user := range product_user_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_user.sku)).Get(&product)
		core.HandleError(err)

		if has == false {
			panic(fmt.Sprintf("商品SKU %s 数据库中不存在！", product_user.sku))
		}

		product_user_list[index].product = product

		// 销售人员
		if product_user.sale_user_login_name == "" {
			product_user_list[index].sale_user_id = 0
		} else {
			var sale_user auth_models.UserInfo
			has, err = db_session.Where("login_name = ?", product_user.sale_user_login_name).Get(&sale_user)
			core.HandleError(err)

			if has == false {
				panic(fmt.Sprintf("指定的销售人员工号数据库中不存在！"))
			} else {
				product_user_list[index].sale_user_id = sale_user.Id
			}
		}

		// 客服人员
		if product_user.service_user_login_name == "" {
			product_user_list[index].service_user_id = 0
		} else {
			var service_user auth_models.UserInfo
			has, err = db_session.Where("login_name = ?", product_user.service_user_login_name).Get(&service_user)
			core.HandleError(err)

			if has == false {
				panic(fmt.Sprintf("指定的客服人员工号数据库中不存在！"))
			} else {
				product_user_list[index].service_user_id = service_user.Id
			}
		}

		// 开发人员
		if product_user.dev_user_login_name == "" {
			product_user_list[index].dev_user_id = 0
		} else {
			var dev_user auth_models.UserInfo
			has, err = db_session.Where("login_name = ?", product_user.dev_user_login_name).Get(&dev_user)
			core.HandleError(err)

			if has == false {
				panic(fmt.Sprintf("指定的开发人员工号数据库中不存在！"))
			} else {
				product_user_list[index].dev_user_id = dev_user.Id
			}
		}

		// 补货人员
		if product_user.replenish_user_login_name == "" {
			product_user_list[index].replenish_user_id = 0
		} else {
			var replenish_user auth_models.UserInfo
			has, err = db_session.Where("login_name = ?", product_user.replenish_user_login_name).Get(&replenish_user)
			core.HandleError(err)

			if has == false {
				panic(fmt.Sprintf("指定的开发人员工号数据库中不存在！"))
			} else {
				product_user_list[index].dev_user_id = replenish_user.Id
			}
		}
	}

	for _, product_user := range product_user_list {
		var product = product_user.product
		product.SaleUserId = product_user.sale_user_id
		product.ServiceUserId = product_user.service_user_id
		product.DevUserId = product_user.dev_user_id
		product.ReplenishUserId = product_user.replenish_user_id

		update_col_list := []string{"sale_user_id", "service_user_id", "dev_user_id", "replenish_user_id"}
		_, err = db_session.Id(product.Id).Cols(update_col_list...).Update(&product)
		core.HandleError(err)
	}

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "数据导入成功！",
	})
}

func (c ErpProductUserExcel) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.SPPZ_PRODUCT_USER, c.Session)

	db_session := c.DbSession

	var id_list = make([]int64, 0)
	c.Params.Bind(&id_list, "id_list")

	var filter = "1=1"

	if len(id_list) != 0 {
		var id_str_list = core.Int_list_to_str_list(id_list)
		filter = core.SqlAnd(filter, fmt.Sprintf("id IN (%s)", strings.Join(id_str_list, ",")))
	}

	product_list := make([]kcgl_models.ProductInfo, 0)
	err := db_session.Where(filter).Find(&product_list)
	core.HandleError(err)

	xlsx_file := c.export_product_list_to_excel(db_session, product_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := "商品负责人导出.xlsx" //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpProductUserExcel) export_product_list_to_excel(db_session *xorm.Session,
	product_list []kcgl_models.ProductInfo) *xlsx.File {

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
	sheet, err := xlsx_file.AddSheet("商品负责人列表")
	core.HandleError(err)

	//添加标题行
	row := sheet.AddRow()
	row.SetHeightCM(head_row_height)

	col_index := 0

	var product_status_list = make([]models.StatusInfo, 0)
	err = db_session.Where(fmt.Sprintf("type = %d", models.STATUS_TYPE_PRODUCT)).Find(&product_status_list)
	core.HandleError(err)

	var find_product_status_name = func(status_id int64) string {
		for _, status := range product_status_list {
			if status.Id == status_id {
				return status.Name
			}
		}

		return ""
	}

	var user_list = make([]auth_models.UserInfo, 0)
	core.HandleError(db_session.Find(&user_list))

	var find_user = func(user_id int64) auth_models.UserInfo {
		for _, user := range user_list {
			if user.Id == user_id {
				return user
			}
		}

		return auth_models.UserInfo{}
	}

	//序号	SKU	中文名称	产品长度（cm）	产品宽度（cm）	产品高度（cm）	实际重量(kg)	状态

	{
		cell := row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "序号"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "SKU"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 30)
		set_head_cell_style(cell)
		cell.Value = "中文名称"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "销售人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "销售人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "客服人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "客服人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "开发人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "开发人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "补货人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "补货人员"
		col_index++

		cell = row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
		col_index++
	}

	for index, product := range product_list {
		row = sheet.AddRow()
		row.SetHeightCM(common_row_height)

		// 序号
		cell := row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(index + 1)

		// SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = product.Sku

		// 中文名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = product.CnName

		// 销售人员

		var sale_user = find_user(product.SaleUserId)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = sale_user.LoginName

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = sale_user.NickName

		// 客服人员
		var service_user = find_user(product.ServiceUserId)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = service_user.LoginName

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = service_user.NickName

		// 开发人员
		var dev_user = find_user(product.DevUserId)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = dev_user.LoginName

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = dev_user.NickName

		// 补货人员
		var replenish_user = find_user(product.ReplenishUserId)

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = replenish_user.LoginName

		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = replenish_user.NickName

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "right")
		cell.Value = find_product_status_name(product.StatusId)
	}

	return xlsx_file
}
