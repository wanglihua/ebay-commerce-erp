package purchase_bill_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/cggl_models"
	"matrix/core"
	"fmt"
)

func SetEntryStatusAndUpdateBillStatus(entry_id int64, status_id int64, db_session *xorm.Session) {
	// 获取明细的实体
	var purchase_bill_entry cggl_models.PurchaseBillEntryInfo
	_, err := db_session.Id(entry_id).Get(&purchase_bill_entry)
	core.HandleError(err)

	// 更新明细的状态
	purchase_bill_entry.StatusId = status_id
	_, err = db_session.Id(entry_id).Cols("status_id").Update(purchase_bill_entry)
	core.HandleError(err)

	// 查看同一个采购单的明细是不是都是当前设置的这个状态了
	var bill_id = purchase_bill_entry.BillId
	count, err := db_session.Where(fmt.Sprintf("bill_id = %d AND status_id <> %d", bill_id, status_id)).Count(new(cggl_models.PurchaseBillEntryInfo))
	core.HandleError(err)

	// 如果都是当前设置的这个状态了， 就更新采购单的状态
	if count == 0 {
		var purchase_bill cggl_models.PurchaseBillInfo
		_, err = db_session.Id(bill_id).Get(&purchase_bill)
		core.HandleError(err)

		purchase_bill.StatusId = status_id
		_, err = db_session.Id(purchase_bill.Id).Cols("status_id").Update(&purchase_bill)
		core.HandleError(err)
	}
}
