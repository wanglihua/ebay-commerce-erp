package container_transfer_log_service

import (
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/modules/erp/models/hgbg_models"
	"matrix/core"
)

func Insert_container_transfer_log(container_transfer hgbg_models.ContainerTransferInfo, operation string, db_session *xorm.Session, web_session revel.Session) {
	var container_transfer_log hgbg_models.ContainerTransferLogInfo

	container_transfer_log.TransferId = container_transfer.Id
	container_transfer_log.BatchCode = container_transfer.BatchCode

	var login_user = core.GetLoginUser(web_session)
	container_transfer_log.UserId = login_user.UserId
	container_transfer_log.LoginName = login_user.LoginName

	container_transfer_log.Operation = operation

	_, err := db_session.Insert(&container_transfer_log)
	core.HandleError(err)
}
