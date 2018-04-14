package purchase_bill_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/cggl_models"
	"matrix/core"
)

func SetEntriesStatus(bill_id_list []int64, old_status_id int64, status_id int64, db_session *xorm.Session) {

	/*
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
	*/

	var purchase_bill_entry_list = make([]cggl_models.PurchaseBillEntryInfo, 0)
	core.HandleError(db_session.In("bill_id", bill_id_list).Find(&purchase_bill_entry_list))

	for _, purchase_bill_entry := range purchase_bill_entry_list {
		if purchase_bill_entry.StatusId == old_status_id {
			purchase_bill_entry.StatusId = status_id
			_, err := db_session.Id(purchase_bill_entry.Id).Cols("status_id").Update(&purchase_bill_entry)
			core.HandleError(err)
		}
	}
}
