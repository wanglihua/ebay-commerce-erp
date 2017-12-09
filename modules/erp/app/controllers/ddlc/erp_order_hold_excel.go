package controllers

import (
	"github.com/revel/revel"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"bytes"
	"time"
	"matrix/core"
	"fmt"
)

type ErpOrderHoldExcel struct {
	*revel.Controller
}

// 最佳仓缺货 Hold 导出
func (c ErpOrderHoldExcel) ExportBestStockOos() revel.Result {
	var db_session = c.DbSession

	var xlsx_file = c.get_export_best_stock_oos_excel(db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err := xlsx_file.Write(&buffer)
	core.HandleError(err)

	var reader = bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口
	var file_name = "最佳仓缺货SKU列表" + ".xlsx"       //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}

func (c ErpOrderHoldExcel) get_export_best_stock_oos_excel(db_session *xorm.Session) *xlsx.File {

	type order_hold_export_item_t struct {
		OrderCode        string            `xorm:"'order_code'"`
		OrderPaidTime    core.NullLongTime `xorm:"'order_paid_time'"`
		OrderShippedTime core.NullLongTime `xorm:"'order_shipped_time'"`
		StockName        string            `xorm:"'stock_name'"`
		ProductSku       string            `xorm:"'product_sku'"`
		ProductName      string            `xorm:"'product_name'"`
		ProductQty       int64             `xorm:"'product_qty'"`
		Price            float64           `xorm:"'price'"`
		StatusName       string            `xorm:"'status_name'"`
	}

	/*
SELECT
    o.code                         AS order_code,
    o.paid_time                    AS order_paid_time,
    o.shipped_time                 AS order_shipped_time,
    st.name                        AS stock_name,
    p.sku                          AS product_sku,
    p.cn_name                      AS product_name,
    oi.quantity                    AS product_qty,
    round(oi.transaction_price, 2) AS price,
    s.name                         AS status_name
FROM
    erp_order_item oi
    INNER JOIN erp_order o ON o.id = oi.order_id
    INNER JOIN erp_status s ON s.id = o.status_id
    INNER JOIN erp_product p ON p.id = oi.item_product_id
    LEFT JOIN erp_order_logistics ol ON o.id = ol.order_id
    LEFT JOIN erp_stock st ON st.id = ol.stock_id
WHERE o.hold_type_id = 10060 # ORDER_HOLD_TYPE_OOS_BEST_STOCK 最佳仓库无货
ORDER BY o.paid_time ASC;
	 */

	var sql = `SELECT
    o.code                         AS order_code,
    o.paid_time                    AS order_paid_time,
    o.shipped_time                 AS order_shipped_time,
    st.name                        AS stock_name,
    p.sku                          AS product_sku,
    p.cn_name                      AS product_name,
    oi.quantity                    AS product_qty,
    round(oi.transaction_price, 2) AS price,
    s.name                         AS status_name
FROM
    erp_order_item oi
    INNER JOIN erp_order o ON o.id = oi.order_id
    INNER JOIN erp_status s ON s.id = o.status_id
    INNER JOIN erp_product p ON p.id = oi.item_product_id
    LEFT JOIN erp_order_logistics ol ON o.id = ol.order_id
    LEFT JOIN erp_stock st ON st.id = ol.stock_id
WHERE o.hold_type_id = 10060 # ORDER_HOLD_TYPE_OOS_BEST_STOCK 最佳仓库无货
ORDER BY o.paid_time ASC;`

	var item_list = make([]order_hold_export_item_t, 0)
	err := db_session.Sql(sql).Find(&item_list)
	core.HandleError(err)

	//创建工作薄
	xlsx_file := xlsx.NewFile()

	sheet, err := xlsx_file.AddSheet("订单编码订单代码")
	core.HandleError(err)

	//统一设置标题单元格格式的函数
	set_head_cell_style := func(cell *xlsx.Cell) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = true
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyBorder = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = "center"
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	set_text_cell_style := func(cell *xlsx.Cell, horizontal_align string) {
		style := xlsx.NewStyle()

		font := * xlsx.NewFont(9, "微软雅黑")
		font.Bold = false
		font.Italic = false
		font.Underline = false
		style.Font = font
		style.ApplyFont = true

		fill := *xlsx.NewFill("solid", "00FFFFFF", "00000000") //CCFFCC
		style.Fill = fill
		style.ApplyFill = true

		border := *xlsx.NewBorder("thin", "thin", "thin", "thin")
		style.Border = border
		style.ApplyBorder = true

		alignment := *new(xlsx.Alignment)
		alignment.Horizontal = horizontal_align
		alignment.Vertical = "center"
		alignment.WrapText = true
		style.Alignment = alignment
		style.ApplyAlignment = true

		cell.SetStyle(style)
	}

	row_height := 0.6

	{
		header_row := sheet.AddRow()
		header_row.SetHeightCM(row_height)

		col_index := 0

		var cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 20)
		set_head_cell_style(cell)
		cell.Value = "订单编码"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 14)
		set_head_cell_style(cell)
		cell.Value = "支付时间"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 14)
		set_head_cell_style(cell)
		cell.Value = "发货时间"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 12)
		set_head_cell_style(cell)
		cell.Value = "仓库"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 15)
		set_head_cell_style(cell)
		cell.Value = "商品SKU"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 45)
		set_head_cell_style(cell)
		cell.Value = "商品名称"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 6)
		set_head_cell_style(cell)
		cell.Value = "数量"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "价格"

		col_index++
		cell = header_row.AddCell()
		sheet.SetColWidth(col_index, col_index, 10)
		set_head_cell_style(cell)
		cell.Value = "状态"
	}

	for _, item := range item_list {
		row := sheet.AddRow()
		row.SetHeightCM(row_height)

		// 订单编码
		var cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = item.OrderCode

		// 支付时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if item.OrderPaidTime.Valid {
			cell.Value = item.OrderPaidTime.Time.Format("2006-01-02 15:04")
		} else {
			cell.Value = ""
		}

		// 发货时间
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		if item.OrderShippedTime.Valid {
			cell.Value = item.OrderShippedTime.Time.Format("2006-01-02 15:04")
		} else {
			cell.Value = ""
		}

		// 仓库
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = item.StockName

		// 商品SKU
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = item.ProductSku

		// 商品名称
		cell = row.AddCell()
		set_text_cell_style(cell, "left")
		cell.Value = item.ProductName

		// 数量
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprint(item.ProductQty)

		// 价格
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = fmt.Sprintf("%0.2f", item.Price)

		// 状态
		cell = row.AddCell()
		set_text_cell_style(cell, "center")
		cell.Value = item.StatusName
	}

	return xlsx_file
}
