package product_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Get_product_cost_price(product_id int64, db_session *xorm.Session) float64 {
	// SELECT p.cost_price FROM erp_product p WHERE p.id = 1
	sql := `SELECT p.cost_price FROM erp_product p WHERE p.id = %d`
	sql = fmt.Sprintf(sql, product_id)

	var product_cost_price struct {
		Price float64 `xorm:"'cost_price'"`
	}

	_, err := db_session.Sql(sql).Get(&product_cost_price)
	core.HandleError(err)

	return product_cost_price.Price
}
