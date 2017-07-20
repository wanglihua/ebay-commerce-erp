package container_transfer_excel_out

import (
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	"matrix/modules/erp/models"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_out/ozo"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_out/third_party"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_out/wzt"
)

type ContainerTransferExcelOutExporter interface {
	GetExcelFileName(p_transfer *models.ContainerTransferInfo) string
	GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *xlsx.File
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
