package purchase_bill_service

import (
	"matrix/modules/erp/models/cggl_models"
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"matrix/core"
	"math"
)

// qty 需要进行分配的数量， 在保存的时候是 entry.CheckinQty, 在提建议的时候是 entry.Qty
func Get_new_checkin_stock_qty(entry cggl_models.PurchaseBillEntryInfo, qty int64, db_session *xorm.Session) (mel_qty int64, bne_qty int64) {

	var mel_qty_in_entry = entry.MelQty
	var bne_qty_in_entry = entry.BneQty

	var intel_repl cggl_models.IntelReplInfo
	_, err := db_session.Where("product_id = ? AND status_id = ?", entry.ProductId, models.INTEL_REPL_STATUS_XSC).Get(&intel_repl)
	core.HandleError(err)

	// 当查询的智能补货不存在时，这里的两个值都会是零，也不影响接下来的逻辑
	var mel_qty_in_ir = intel_repl.MelQty
	var bne_qty_in_ir = intel_repl.BneQty

	var total_qty = mel_qty_in_entry + bne_qty_in_entry + mel_qty_in_ir + bne_qty_in_ir

	mel_qty = int64(math.Floor((float64(mel_qty_in_entry+mel_qty_in_ir)/float64(total_qty))*float64(qty) + 0.5))
	bne_qty = int64(qty - mel_qty)

	return mel_qty, bne_qty
}
