package shop_account_log_service

import (
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
)

// 插入店铺账户日志
func Insert_shop_account_log(operation_type core.Db_operation_type, db_session *xorm.Session, web_session revel.Session) {
	var shop_log models.ShopAccountLogInfo

	login_user := core.GetLoginUser(web_session)
	shop_log.UserId = login_user.UserId
	shop_log.LoginName = login_user.LoginName

	if operation_type == core.Db_operation_types.Delete {
		shop_log.Operation = "删除数据"
	} else if operation_type == core.Db_operation_types.Insert {
		shop_log.Operation = "插入数据"
	} else if operation_type == core.Db_operation_types.Update {
		shop_log.Operation = "更新数据"
	}

	_, err := db_session.Insert(&shop_log)
	core.HandleError(err)
}

// 插入店铺账户API的日志
func Insert_shop_account_api_log(operation_type core.Db_operation_type, db_session *xorm.Session, web_session revel.Session) {
	var shop_log models.ShopAccountLogInfo

	login_user := core.GetLoginUser(web_session)
	shop_log.UserId = login_user.UserId
	shop_log.LoginName = login_user.LoginName

	if operation_type == core.Db_operation_types.Delete {
		shop_log.Operation = "删除API配置数据"
	} else if operation_type == core.Db_operation_types.Insert {
		shop_log.Operation = "插入API配置数据"
	} else if operation_type == core.Db_operation_types.Update {
		shop_log.Operation = "更新API配置数据"
	}

	_, err := db_session.Insert(&shop_log)
	core.HandleError(err)
}
