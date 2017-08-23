package container_transfer_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/kcgl_models"
)

func EntryCombineForCn(item_list []ContainerTransferEntryItem, is_third bool, db_session *xorm.Session) []ContainerTransferEntryItem {
	var find_entry_in_list = func(cn_decl_name string, entry_list []ContainerTransferEntryItem) (index_ret int, entry_ret ContainerTransferEntryItem) {
		for index, entry := range entry_list {
			if entry.Product.CnDeclName == cn_decl_name {
				return index, entry
			}
		}

		return -1, ContainerTransferEntryItem{}
	}

	var not_combined_entry_list = make([]ContainerTransferEntryItem, 0)
	var combined_entry_list = make([]ContainerTransferEntryItem, 0)

	for _, item := range item_list {
		// weight volume 的计算
		var product kcgl_models.ProductInfo
		_, err := db_session.Id(item.ContainerTransferEntry.ProductId).Get(&product)
		core.HandleError(err)

		var net_weight_sum = product.Weight * float64(item.ContainerTransferEntry.Qty)
		var vol_sum = (product.Length * product.Width * product.Height / 1000000) * float64(item.ContainerTransferEntry.Qty)

		// 需要开票的不合并
		//if is_product_need_invoice(item.Product.Id, db_session) == true {
		if product.NeedInvoice > 0 {
			//if item.ContainerTransferEntry.NeedInvoice > 0 { // 如需要开票
			item.ContainerTransferEntry.WeightSum = net_weight_sum
			item.ContainerTransferEntry.VolumeSum = vol_sum
			not_combined_entry_list = append(not_combined_entry_list, item)
		} else { //不需要开票的合并
			found_index, combined_entry := find_entry_in_list(item.Product.CnDeclName, combined_entry_list)

			if found_index < 0 {
				item.ContainerTransferEntry.WeightSum = net_weight_sum
				item.ContainerTransferEntry.VolumeSum = vol_sum
				combined_entry_list = append(combined_entry_list, item)
			} else {
				// do combine
				combined_entry_list[found_index].ContainerTransferEntry.Price = (combined_entry.ContainerTransferEntry.Price*float64(combined_entry.ContainerTransferEntry.Qty) + item.ContainerTransferEntry.Price*float64(item.ContainerTransferEntry.Qty)) / float64(combined_entry.ContainerTransferEntry.Qty+item.ContainerTransferEntry.Qty)
				combined_entry_list[found_index].ContainerTransferEntry.Qty = combined_entry.ContainerTransferEntry.Qty + item.ContainerTransferEntry.Qty
				combined_entry_list[found_index].ContainerTransferEntry.BoxQty = combined_entry.ContainerTransferEntry.BoxQty + item.ContainerTransferEntry.BoxQty
				combined_entry_list[found_index].ContainerTransferEntry.WeightSum = combined_entry.ContainerTransferEntry.WeightSum + net_weight_sum
				combined_entry_list[found_index].ContainerTransferEntry.VolumeSum = combined_entry.ContainerTransferEntry.VolumeSum + vol_sum
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
