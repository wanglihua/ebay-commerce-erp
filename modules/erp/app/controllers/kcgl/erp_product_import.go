package controllers

import (
	"github.com/revel/revel"

	"fmt"
	"matrix/core"
	"matrix/modules/erp/models"
	"github.com/tealeg/xlsx"
)

type ErpProductImport struct {
	*revel.Controller
}

func (c ErpProductImport) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_import/product_import.html")
}

func (c ErpProductImport) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	product_import_list := make([]models.ProductImportInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&product_import_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ProductImportInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  product_import_list,
		Total: count,
	})
}

func (c ErpProductImport) Import() revel.Result {
	excel_file_bytes := make([]byte, 0)
	c.Params.Bind(&excel_file_bytes, "excel_file")

	c.Validation.Required(excel_file_bytes).Message("指定文件内容未成功上传！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	xlsx_file, err := xlsx.OpenBinary(excel_file_bytes)
	core.HandleError(err)

	sheet_count := len(xlsx_file.Sheets)
	revel.TRACE.Printf("sheet_count: %d\n", sheet_count)
	for sheet_index := 0; sheet_index < sheet_count; sheet_index ++ {
		sheet := xlsx_file.Sheets[sheet_index]
		max_row := sheet.MaxRow
		revel.TRACE.Printf("max_row: %d\n", max_row)
		for row_index := 0; row_index < max_row; row_index ++ {
			row := sheet.Rows[row_index]
			cell := row.Cells[0]
			cell_val, _ := cell.String()
			revel.TRACE.Println(cell_val)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"success":  true,
		"message":  "文件上传成功！",
	})

	return c.RenderJson(core.JsonResult{Success: true, Message: "Excel数据上传成功!"})
}

func (c ErpProductImport) Confirm() revel.Result {
	return c.RenderJson(core.JsonResult{Success: true, Message: "商品数据导入成功!"})
}

func (c ErpProductImport) Clear() revel.Result {
	db_session := c.DbSession

	login_user := core.GetLoginUser(c.Session)
	login_user_id := login_user.UserId

	_, err := db_session.
	Where(fmt.Sprintf("%s = ?", models.ProductImportCols.UserId), login_user_id).
		Delete(new(models.ProductImportInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: "临时导入数据清空成功!"})
}
