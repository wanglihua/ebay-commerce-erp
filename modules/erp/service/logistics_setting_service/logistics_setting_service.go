package logistics_setting_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

// 获取最新的物流配置信息，如果数据库中不存在，就返回零值
func Get_logistics_setting(db_session *xorm.Session) models.LogisticsSettingInfo {
	var logistics_setting models.LogisticsSettingInfo
	sql := "SELECT * FROM erp_logistics_setting ls ORDER BY ls.update_time DESC LIMIT 0, 1"
	_, err := db_session.Sql(sql).Get(&logistics_setting)
	core.HandleError(err)

	return logistics_setting
}
