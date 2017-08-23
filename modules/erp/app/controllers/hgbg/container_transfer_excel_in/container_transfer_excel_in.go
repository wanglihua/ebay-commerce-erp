package container_transfer_excel_in

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/hgbg_models"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_in/elife_pty"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_in/giant_star"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_in/ozo_estore"
	"github.com/Luxurioust/excelize"
)

type ContainerTransferExcelInExporter interface {
	GetExcelFileName(p_transfer *hgbg_models.ContainerTransferInfo) string
	GetExcelFile(p_transfer *hgbg_models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File
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
