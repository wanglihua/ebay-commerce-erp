package container_transfer_service

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func EntryCombineAndPriceRateForAu(p_transfer *models.ContainerTransferInfo, entry_list []ContainerTransferEntryItem, db_session *xorm.Session) ([]ContainerTransferEntryItem) {
	var find_entry_in_list = func(au_decl_name string, entry_list []ContainerTransferEntryItem) (index_ret int, entry_ret ContainerTransferEntryItem) {
		for index, item := range entry_list {
			if item.Product.AuDeclName == au_decl_name {
				return index, item
			}
		}

		return -1, ContainerTransferEntryItem{}
	}

	var combined_entry_list = make([]ContainerTransferEntryItem, 0)

	for _, item := range entry_list {
		found_index, combined_entry := find_entry_in_list(item.Product.AuDeclName, combined_entry_list)

		// weight volume 的计算
		var product models.ProductInfo
		_, err := db_session.Id(item.ContainerTransferEntry.ProductId).Get(&product)
		core.HandleError(err)

		var net_weight_sum = product.Weight * float64(item.ContainerTransferEntry.Qty)
		var vol_sum = (product.Length * product.Width * product.Height / 1000000) * float64(item.ContainerTransferEntry.Qty)

		if found_index < 0 {
			item.ContainerTransferEntry.WeightSum = net_weight_sum
			item.ContainerTransferEntry.VolumeSum = vol_sum
			combined_entry_list = append(combined_entry_list, item)
		} else {
			// do combine
			combined_entry_list[found_index].ContainerTransferEntry.Price = (combined_entry.ContainerTransferEntry.Price*float64(combined_entry.ContainerTransferEntry.Qty) + item.ContainerTransferEntry.Price*float64(item.ContainerTransferEntry.Qty)) / float64(combined_entry.ContainerTransferEntry.Qty+item.ContainerTransferEntry.Qty)
			combined_entry_list[found_index].ContainerTransferEntry.Qty = combined_entry.ContainerTransferEntry.Qty + item.ContainerTransferEntry.Qty
			combined_entry_list[found_index].ContainerTransferEntry.BoxQty = combined_entry.ContainerTransferEntry.BoxQty + item.ContainerTransferEntry.BoxQty
			combined_entry_list[found_index].ContainerTransferEntry.BoxQty = combined_entry.ContainerTransferEntry.BoxQty + item.ContainerTransferEntry.BoxQty
			combined_entry_list[found_index].ContainerTransferEntry.WeightSum = combined_entry.ContainerTransferEntry.WeightSum + net_weight_sum
			combined_entry_list[found_index].ContainerTransferEntry.VolumeSum = combined_entry.ContainerTransferEntry.VolumeSum + vol_sum
		}
	}

	// do price rate
	for index, combined_entry := range combined_entry_list {
		if combined_entry.ContainerTransferEntry.AuHsDutyRate == 0 {
			combined_entry_list[index].ContainerTransferEntry.Price = combined_entry.ContainerTransferEntry.Price * (1 + p_transfer.PriceUpRate) // 不征税商品上调
		} else {
			combined_entry_list[index].ContainerTransferEntry.Price = combined_entry.ContainerTransferEntry.Price * p_transfer.PriceDownRate // 征税商品打折
		}
	}

	return combined_entry_list
}
