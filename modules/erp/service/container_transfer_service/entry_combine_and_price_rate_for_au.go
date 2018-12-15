package container_transfer_service

import (
	"matrix/modules/erp/models/hgbg_models"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/kcgl_models"
	"math"
	"strings"
)

func EntryCombineAndPriceRateForAu(p_transfer *hgbg_models.ContainerTransferInfo, entry_list []ContainerTransferEntryItem, db_session *xorm.Session) ([]ContainerTransferEntryItem) {
	var find_entry_in_list = func(au_decl_name string, entry_list []ContainerTransferEntryItem) (index_ret int, entry_ret ContainerTransferEntryItem) {
		for index, item := range entry_list {
			// if strings.TrimSpace(item.Product.AuDeclName) == strings.TrimSpace(au_decl_name) {
			if strings.TrimSpace(item.ContainerTransferEntry.AuHsEnName) == strings.TrimSpace(au_decl_name) {
				return index, item
			}
		}

		return -1, ContainerTransferEntryItem{}
	}

	var combined_entry_list = make([]ContainerTransferEntryItem, 0)

	for _, item := range entry_list {
		found_index, combined_entry := find_entry_in_list(item.ContainerTransferEntry.AuHsEnName, combined_entry_list)

		// weight volume 的计算
		var product kcgl_models.ProductInfo
		_, err := db_session.Id(item.ContainerTransferEntry.ProductId).Get(&product)
		core.HandleError(err)

		var net_weight_sum = product.Weight * float64(item.ContainerTransferEntry.Qty)
		var vol_sum = (product.Length * product.Width * product.Height / 1000000) * float64(item.ContainerTransferEntry.Qty)

		// 澳洲清关的都是合并的，不管开不开票。
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
		if math.Abs(p_transfer.PriceUpRate) > core.FLOAT_MIN_DIFF { // 如果上调率不是零
			if combined_entry.ContainerTransferEntry.AuHsDutyRate == 0 {
				combined_entry_list[index].ContainerTransferEntry.Price = combined_entry.ContainerTransferEntry.Price * (1 + p_transfer.PriceUpRate) // 不征税商品上调
			} else {
				combined_entry_list[index].ContainerTransferEntry.Price = combined_entry.ContainerTransferEntry.Price * p_transfer.PriceDownRate // 征税商品打折
			}
		} else {
			// 如果上调率是零
			combined_entry_list[index].ContainerTransferEntry.Price = combined_entry.ContainerTransferEntry.Price
		}
	}

	return combined_entry_list
}
