package storage_loc_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Get_storage_loc_id_by_name(storage_loc_name string, stock_id int64, db_session *xorm.Session) int64 {

	var storage_loc_id struct {
		StorageLocId int64 `xorm:"'storage_loc_id'"`
	}

	/*
	SELECT sl.id AS storage_loc_id FROM erp_storage_loc sl WHERE sl.name = 'abc' AND sl.stock_id = 0
	 */

	var sql_get_storage_loc_id = `SELECT sl.id AS storage_loc_id FROM erp_storage_loc sl WHERE sl.name = '%s' AND sl.stock_id = %d`
	sql_get_storage_loc_id = fmt.Sprintf(sql_get_storage_loc_id, storage_loc_name, stock_id)
	_, err := db_session.Sql(sql_get_storage_loc_id).Get(&storage_loc_id)
	core.HandleError(err)

	return storage_loc_id.StorageLocId
}
