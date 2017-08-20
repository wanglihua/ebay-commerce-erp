package container_transfer_service

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func EntryCombineForCn(entry_list []models.ContainerTransferEntryInfo, is_third bool, db_session *xorm.Session) []models.ContainerTransferEntryInfo {
	var find_entry_in_list = func(au_hs_code string, entry_list []models.ContainerTransferEntryInfo) (index_ret int, entry_ret models.ContainerTransferEntryInfo) {
		for index, entry := range entry_list {
			if entry.AuHsCode == au_hs_code {
				return index, entry
			}
		}

		return -1, models.ContainerTransferEntryInfo{}
	}

	var not_combined_entry_list = make([]models.ContainerTransferEntryInfo, 0)
	var combined_entry_list = make([]models.ContainerTransferEntryInfo, 0)

	for _, entry := range entry_list {
		// weight volume 的计算
		var product models.ProductInfo
		_, err := db_session.Id(entry.ProductId).Get(&product)
		core.HandleError(err)

		var net_weight_sum = product.Weight * float64(entry.Qty)
		var vol_sum = (product.Length * product.Width * product.Height / 1000000) * float64(entry.Qty)

		// 需要开票的不合并
		//if is_product_need_invoice(entry.ProductId, db_session) == true {
		if entry.NeedInvoice > 0 { // 如需要开票
			entry.WeightSum = net_weight_sum
			entry.VolumeSum = vol_sum
			not_combined_entry_list = append(not_combined_entry_list, entry)
		} else { //不需要开票的合并
			found_index, combined_entry := find_entry_in_list(entry.AuHsCode, combined_entry_list)

			if found_index < 0 {
				entry.WeightSum = net_weight_sum
				entry.VolumeSum = vol_sum
				combined_entry_list = append(combined_entry_list, entry)
			} else {
				// do combine
				combined_entry_list[found_index].Price = (combined_entry.Price*float64(combined_entry.Qty) + entry.Price*float64(entry.Qty)) / float64(combined_entry.Qty+entry.Qty)
				combined_entry_list[found_index].Qty = combined_entry.Qty + entry.Qty
				combined_entry_list[found_index].BoxQty = combined_entry.BoxQty + entry.BoxQty
				combined_entry_list[found_index].WeightSum = combined_entry.WeightSum + net_weight_sum
				combined_entry_list[found_index].VolumeSum = combined_entry.VolumeSum + vol_sum
			}
		}
	}

	/*
	// 开票的和不开票的合并
	combined_entry_list = append(not_combined_entry_list, combined_entry_list...)

	return combined_entry_list
	*/

	// 不开票的导出到第三方, 就是合并的那部分
	if is_third {
		return combined_entry_list
	} else {
		return not_combined_entry_list
	}
}

/*
// 需要开票的不合并，不需要开票的合并
func is_product_need_invoice(product_id int64, db_session *xorm.Session) bool {
	var need_invoice struct {
		NeedInvoice int64 `xorm:"'need_invoice'"`
	}

	var sql = `SELECT
    p.need_invoice AS need_invoice
FROM erp_product p
WHERE p.id = %d`
	sql = fmt.Sprintf(sql, product_id)

	_, err := db_session.Sql(sql).Get(&need_invoice)
	core.HandleError(err)

	if need_invoice.NeedInvoice == 0 {
		return false
	} else {
		return true
	}
}
*/
