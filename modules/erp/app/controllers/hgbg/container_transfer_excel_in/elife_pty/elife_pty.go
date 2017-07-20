package elife_pty

/*
type ContainerTransferExcelInExportInterface interface {
	GetExcelFileName(p_transfer *models.ContainerTransferInfo) string
	GetExcelFile(p_transfer *models.ContainerTransferInfo, db_session *xorm.Session) *xlsx.File
}

 */

type ElifePty struct {
}

func (c ElifePty) GetExcelFileName() string {
	return ""
}
