package purchase_bill_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func UpdatePurchaseBillAmount(purchase_bill_id int64, db_session *xorm.Session) {
	/*
UPDATE erp_purchase_bill
SET total_amount     = ifnull((SELECT sum(pbe.price * pbe.qty)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0),
    advance_amount   = ifnull((SELECT sum(pbe.pay_amount)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0),
    final_amount     = ifnull((SELECT sum(pbe.final_amount)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0),
    advance_expenses = ifnull((SELECT sum(pbe.freight_charge)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0)
WHERE erp_purchase_bill.id = 2;
	 */
	var sql = `
UPDATE erp_purchase_bill
SET total_amount     = ifnull((SELECT sum(pbe.price * pbe.qty)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0),
    advance_amount   = ifnull((SELECT sum(pbe.pay_amount)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0),
    final_amount     = ifnull((SELECT sum(pbe.final_amount)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0),
    advance_expenses = ifnull((SELECT sum(pbe.freight_charge)
                               FROM erp_purchase_bill_entry pbe
                               WHERE pbe.bill_id = erp_purchase_bill.id), 0)
WHERE erp_purchase_bill.id = %d;
`

	sql = fmt.Sprintf(sql, purchase_bill_id)

	_, err := db_session.Exec(sql)
	core.HandleError(err)
}
