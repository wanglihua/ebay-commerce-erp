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
