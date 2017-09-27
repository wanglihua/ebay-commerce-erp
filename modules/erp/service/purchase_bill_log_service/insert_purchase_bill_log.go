package purchase_bill_log_service

import (
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models/cggl_models"
)

func Insert_purchase_bill_log(purchase_bill cggl_models.PurchaseBillInfo, operation string, db_session *xorm.Session, web_session revel.Session) {
	var purchase_bill_log cggl_models.PurchaseBillLogInfo

	purchase_bill_log.BillId = purchase_bill.Id
	purchase_bill_log.BillCode = purchase_bill.Code

	var login_user = core.GetLoginUser(web_session)
	purchase_bill_log.UserId = login_user.UserId
	purchase_bill_log.LoginName = login_user.LoginName

	purchase_bill_log.Operation = operation

	_, err := db_session.Insert(&purchase_bill_log)
	core.HandleError(err)
}
