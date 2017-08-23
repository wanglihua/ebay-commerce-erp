package container_transfer_excel_out

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/hgbg_models"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_out/ozo"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_out/third_party"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_out/wzt"
	"github.com/Luxurioust/excelize"
)

type ContainerTransferExcelOutExporter interface {
	GetExcelFileName(p_transfer *hgbg_models.ContainerTransferInfo) string
	GetExcelFile(p_transfer *hgbg_models.ContainerTransferInfo, db_session *xorm.Session) *excelize.File
}

func GetExporter(company_code string) ContainerTransferExcelOutExporter {
	if company_code == "ozo" {
		return ozo.Exporter{}
	}

	if company_code == "third" {
		return third_party.Exporter{}
	}

	if company_code == "wzt" {
		return wzt.Exporter{}
	}

	panic("出口报关资料的公司编码不正确！")

	return nil
}
