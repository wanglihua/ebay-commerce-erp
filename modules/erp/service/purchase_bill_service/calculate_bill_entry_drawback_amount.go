package purchase_bill_service

import (
	"matrix/modules/erp/models/cggl_models"
	"matrix/modules/erp/models/kcgl_models"
)

func CalculateBillEntryDrawbackAmount(bill_entry cggl_models.PurchaseBillEntryInfo) float64 {
	if bill_entry.NeedInvoice <= 0 {
		return 0
	}

	// return ((bill_entry.Price * float64(bill_entry.Qty)) / 1.17) * bill_entry.DrawbackPoint
	// 退税额，只需计算单个商品的就可以了，不用考虑数量的
	return (bill_entry.Price / 1.17) * bill_entry.DrawbackPoint / 100
}

func CalculateProductDrawbackAmount(product kcgl_models.ProductInfo, purchase_price float64) float64 {
	if product.NeedInvoice <= 0 {
		return 0
	}

	// return ((bill_entry.Price * float64(bill_entry.Qty)) / 1.17) * bill_entry.DrawbackPoint
	// 退税额，只需计算单个商品的就可以了，不用考虑数量的
	return (purchase_price / 1.17) * product.DrawbackPoint / 100
}
