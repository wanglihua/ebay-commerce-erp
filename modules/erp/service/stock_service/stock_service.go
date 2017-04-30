package stock_service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
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

func Get_stock_id_name_map(stock_id_list []int64, db_session *xorm.Session) map[int64]string {
	var stock_id_name_map = make(map[int64]string)

	if len(stock_id_list) == 0 {
		return stock_id_name_map
	}

	var stock_list = make([]models.StockInfo, 0)
	err := db_session.In("id", stock_id_list).Find(&stock_list)
	core.HandleError(err)

	for _, stock := range stock_list {
		stock_id_name_map[stock.Id] = stock.Name
	}

	return stock_id_name_map
}
