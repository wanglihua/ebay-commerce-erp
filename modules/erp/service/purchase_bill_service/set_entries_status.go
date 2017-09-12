package purchase_bill_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"strings"
	"matrix/core"
)

func SetEntriesStatus(bill_id_list []int64, status_id int64, db_session *xorm.Session) {
	/*
UPDATE erp_purchase_bill_entry pbe
SET pbe.status_id = 4400
WHERE pbe.bill_id IN (1, 2, 3)
	 */

	var bill_id_str_list = make([]string, 0)
	for _, bill_id := range bill_id_list {
		bill_id_str_list = append(bill_id_str_list, fmt.Sprint(bill_id))
	}

	var sql = `UPDATE erp_purchase_bill_entry pbe
SET pbe.status_id = %d
WHERE pbe.bill_id IN (%s)`

	sql = fmt.Sprintf(sql, status_id, strings.Join(bill_id_str_list, ","))

	_, err := db_session.Exec(sql)
	core.HandleError(err)
}
