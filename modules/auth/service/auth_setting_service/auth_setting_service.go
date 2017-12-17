package auth_setting_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/auth/models"
)

func Get_auth_setting(db_session *xorm.Session) models.SettingInfo {
	var sql = `SELECT s.* FROM auth_setting s ORDER BY s.update_time DESC LIMIT 0, 1`

	var auth_setting models.SettingInfo
	_, err := db_session.Sql(sql).Get(&auth_setting)
	core.HandleError(err)

	return auth_setting
}
