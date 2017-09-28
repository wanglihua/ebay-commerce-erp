package purchase_bill_service

import (
	"github.com/revel/revel"
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/models"
)

// 判断采购单当前能不能修改删除，新增导入修改删除明细等
func CanPurchaseBillModify(purchase_bill cggl_models.PurchaseBillInfo, web_session revel.Session) bool {
	/*
		models.PURCHASE_BILL_STATUS_YCJ,
		models.PURCHASE_BILL_STATUS_YTJ,
		models.PURCHASE_BILL_STATUS_YSP,
		models.PURCHASE_BILL_STATUS_CGZ,
		models.PURCHASE_BILL_STATUS_YRK,
	 */

	if purchase_bill.StatusId == models.PURCHASE_BILL_STATUS_YCJ {
		return true
	} else if purchase_bill.StatusId == models.PURCHASE_BILL_STATUS_YTJ || purchase_bill.StatusId == models.PURCHASE_BILL_STATUS_YSP {
		var user_has_approve_pass = pass_service.Is_login_user_has_pass(pass_res.CGGL_PURCHASE_APPROVE, web_session)
		if user_has_approve_pass == true {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
