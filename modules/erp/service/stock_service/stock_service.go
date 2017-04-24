package stock_service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Get_stock_id_by_name(stock_name string, db_session *xorm.Session) int64 {

	var stock_id struct {
		StockId int64 `xorm:"'stock_id'"`
	}

	/*
	SELECT s.id AS stock_id FROM erp_stock s WHERE s.name = 'abc'
	 */

	var sql_get_stock_id = `SELECT s.id AS stock_id FROM erp_stock s WHERE s.name = '%s'`
	sql_get_stock_id = fmt.Sprintf(sql_get_stock_id, stock_name)
	_, err := db_session.Sql(sql_get_stock_id).Get(&stock_id)
	core.HandleError(err)

	return stock_id.StockId
}
