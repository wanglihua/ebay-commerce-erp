package product_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func GetProductPrice(product_id int64, db_session *xorm.Session) float64 {
	var product_price struct {
		Price float64 `xorm:"'cost_price'"`
	}

	var sql = `SELECT p.cost_price
FROM erp_product p
WHERE p.id = %d;`
	sql = fmt.Sprintf(sql, product_id)

	_, err := db_session.Sql(sql).Get(&product_price)
	core.HandleError(err)

	return product_price.Price
}
