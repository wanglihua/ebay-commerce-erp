package purchase_bill_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/cggl_models"
	"matrix/core"
	"matrix/modules/erp/models/kcgl_models"
	"fmt"
)

func Get_propose_purchase_price(product_id int64, db_session *xorm.Session) float64 {
	// 获取采购历史
	var sql_get_latest_purchase_price = `SELECT pbe.*
FROM erp_purchase_bill_entry pbe
WHERE pbe.product_id = %d
ORDER BY pbe.create_time DESC
LIMIT 1;
`
	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	sql_get_latest_purchase_price = fmt.Sprintf(sql_get_latest_purchase_price, product_id)
	has, err := db_session.Sql(sql_get_latest_purchase_price).Get(&purchase_bill_entry)
	core.HandleError(err)

	if has == true {
		return purchase_bill_entry.Price
	} else {
		// 获取商品采购价
		var product kcgl_models.ProductInfo
		_, err := db_session.Id(product_id).Get(&product)
		core.HandleError(err)
		return product.CostPrice
	}
}
