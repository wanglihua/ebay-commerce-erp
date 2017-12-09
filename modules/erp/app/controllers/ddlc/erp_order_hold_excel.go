package controllers

import (
	"github.com/revel/revel"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"bytes"
	"time"
	"matrix/core"
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

	return nil
}
