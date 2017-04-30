package status_service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func Get_status_list(db_session *xorm.Session, status_type int64) []models.StatusInfo {
	var status_list = make([]models.StatusInfo, 0)
	err := db_session.Where(fmt.Sprintf("type = %d", status_type)).Asc("id").Find(&status_list)
	core.HandleError(err)

	return status_list
}

func Get_status_name(status_id int64, db_session *xorm.Session) string {
	var status_name struct {
		Name string `xorm:"'status_name'"`
	}

	/*
	SELECT s.name FROM erp_status s WHERE s.id = 0
	 */
	var sql = `SELECT s.name FROM erp_status s WHERE s.id = %d`
	sql = fmt.Sprintf(sql, status_id)

	_, err := db_session.Sql(sql).Get(&status_name)
	core.HandleError(err)

	return status_name.Name
}

func Get_status_id_by_name(status_name string, status_type int64, db_session *xorm.Session) int64 {

	var status_id struct {
		StatusId int64 `xorm:"'status_id'"`
	}
	/*
	SELECT s.id AS status_id FROM erp_status s WHERE s.type = 0 AND s.name = 'abc'
	 */

	var sql = `SELECT s.id AS status_id FROM erp_status s WHERE s.type = %d AND s.name = '%s'`
	sql = fmt.Sprint(sql, status_type, status_name)
	_, err := db_session.Sql(sql).Get(&status_id)
	core.HandleError(err)

	return status_id.StatusId
}

func Get_status_id_name_map(status_id_list []int64, db_session *xorm.Session) map[int64]string {
	var status_id_name_map = make(map[int64]string)

	if len(status_id_list) == 0 {
		return status_id_name_map
	}

	var status_list = make([]models.StatusInfo, 0)
	err := db_session.In("id", status_id_list).Find(&status_list)
	core.HandleError(err)

	for _, status := range status_list {
		status_id_name_map[status.Id] = status.Name
	}

	return status_id_name_map
}
