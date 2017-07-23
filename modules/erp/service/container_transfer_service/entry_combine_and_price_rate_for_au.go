package container_transfer_service

import (
	"matrix/modules/erp/models"
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func EntryCombineAndPriceRateForAu(p_transfer *models.ContainerTransferInfo, entry_list []models.ContainerTransferEntryInfo, db_session *xorm.Session) ([]models.ContainerTransferEntryInfo) {
	var find_entry_in_list = func(au_hs_code string, entry_list []models.ContainerTransferEntryInfo) (index_ret int, entry_ret models.ContainerTransferEntryInfo) {
		for index, entry := range entry_list {
			if entry.AuHsCode == au_hs_code {
				return index, entry
			}
		}

		return -1, models.ContainerTransferEntryInfo{}
	}

	var combined_entry_list = make([]models.ContainerTransferEntryInfo, 0)

	for _, entry := range entry_list {
		found_index, combined_entry := find_entry_in_list(entry.AuHsCode, combined_entry_list)

		// weight volume 的计算
		var product models.ProductInfo
		_, err := db_session.Id(entry.ProductId).Get(&product)
		core.HandleError(err)

		var net_weight = product.Weight * float64(entry.Qty)
		var vol = (product.Length * product.Width * product.Height / 1000000) * float64(entry.Qty)

		if found_index < 0 {
			entry.WeightSum = net_weight * float64(entry.Qty)
			entry.VolumeSum = vol * float64(entry.Qty)
			combined_entry_list = append(combined_entry_list, entry)
		} else {
			// do combine
			combined_entry_list[found_index].Price = (combined_entry.Price*float64(combined_entry.Qty) + entry.Price*float64(entry.Qty)) / float64(combined_entry.Qty+entry.Qty)
			combined_entry_list[found_index].Qty = combined_entry.Qty + entry.Qty
			combined_entry_list[found_index].BoxQty = combined_entry.BoxQty + entry.BoxQty
			combined_entry_list[found_index].BoxQty = combined_entry.BoxQty + entry.BoxQty
			combined_entry_list[found_index].WeightSum = combined_entry.WeightSum + net_weight*float64(entry.Qty)
			combined_entry_list[found_index].VolumeSum = combined_entry.VolumeSum + vol*float64(entry.Qty)
		}
	}

	// do price rate
	for index, combined_entry := range combined_entry_list {
		if combined_entry.AuHsDutyRate == 0 {
			combined_entry_list[index].Price = combined_entry.Price * (1 + p_transfer.PriceUpRate) // 不征税商品上调
		} else {
			combined_entry_list[index].Price = combined_entry.Price * p_transfer.PriceDownRate // 征税商品打折
		}
	}

	return combined_entry_list
}
