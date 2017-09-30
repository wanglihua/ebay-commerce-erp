package purchase_bill_service

import (
	"matrix/modules/erp/models/cggl_models"
)

func CalculateBillEntryDrawbackAmount(bill_entry cggl_models.PurchaseBillEntryInfo) float64 {
	if bill_entry.NeedInvoice <= 0 {
		return 0
	}

	return (bill_entry.Price / 1.17) * bill_entry.DrawbackPoint
}
