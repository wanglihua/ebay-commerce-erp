package pick_fee_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

// 获取最新的拣货人工费配置信息，如果数据库中不存在，就返回零值
func Get_pick_fee(db_session *xorm.Session) models.PickFeeInfo {
	var pick_fee models.PickFeeInfo
	sql := "SELECT * FROM erp_pick_fee pf ORDER BY pf.update_time DESC LIMIT 0, 1"
	_, err := db_session.Sql(sql).Get(&pick_fee)
	core.HandleError(err)

	return pick_fee
}
