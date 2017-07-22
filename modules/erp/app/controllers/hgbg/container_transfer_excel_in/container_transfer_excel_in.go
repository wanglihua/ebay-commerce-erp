package container_transfer_excel_in

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_in/elife_pty"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_in/giant_star"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_in/ozo_estore"
	"github.com/Luxurioust/excelize"
)

type ContainerTransferExcelInExporter interface {
	GetExcelFileName(p_transfer *models.ContainerTransferInfo) string
	GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File
}

func GetExporter(company_code string) ContainerTransferExcelInExporter {
	if company_code == "elife" {
		return elife_pty.Exporter{}
	}

	if company_code == "giant" {
		return giant_star.Exporter{}
	}

	if company_code == "ozo" {
		return ozo_estore.Exporter{}
	}

	panic("进口清关资料的公司编码不正确！")

	return nil
}

func EntryCombineAndPriceRate(p_transfer *models.ContainerTransferInfo, entry_list []models.ContainerTransferEntryInfo) []models.ContainerTransferEntryInfo {
	var find_entry_in_list = func(au_hs_code string, entry_list []models.ContainerTransferEntryInfo) models.ContainerTransferEntryInfo {
		for _, entry := range entry_list {
			if entry.AuHsCode == au_hs_code {
				return entry
			}
		}

		return models.ContainerTransferEntryInfo{}
	}

	var combined_entry_list = make([]models.ContainerTransferEntryInfo, 0)

	for index, entry := range entry_list {
		var combined_entry = find_entry_in_list(entry.AuHsCode, combined_entry_list)

		if combined_entry.Id == 0 {
			combined_entry_list = append(combined_entry_list, entry)
		} else {
			// do combine
			entry_list[index].Qty = combined_entry.Qty + entry.Qty
			entry_list[index].BoxQty = combined_entry.BoxQty + entry.BoxQty
			entry_list[index].Price = (combined_entry.Price*float64(combined_entry.Qty) + entry.Price*float64(entry.Qty)) / float64(combined_entry.Qty+entry.Qty)
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
