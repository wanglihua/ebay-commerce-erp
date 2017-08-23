package pick_fee_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/ddpz_models"
)

// 获取最新的拣货人工费配置信息，如果数据库中不存在，就返回零值
func Get_pick_fee_setting(db_session *xorm.Session) ddpz_models.PickFeeInfo {
	var pick_fee ddpz_models.PickFeeInfo
	sql := "SELECT * FROM erp_pick_fee pf ORDER BY pf.update_time DESC LIMIT 0, 1"
	_, err := db_session.Sql(sql).Get(&pick_fee)
	core.HandleError(err)

	return pick_fee
}

//　根据　拣货人工费用的设置，　获取货物实际重量的拣货费用
func Get_product_pick_fee(actual_weight float64, pick_fee_setting ddpz_models.PickFeeInfo) float64 {
	if actual_weight < 5 { // 5Kg以下
		return pick_fee_setting.Fee1
	} else if actual_weight < 10 { // 5.1KG-10KG
		return pick_fee_setting.Fee2
	} else if actual_weight < 25 { // 10.1KG-25KG
		return pick_fee_setting.Fee3
	} else if actual_weight < 50 { // 25.1KG-50KG
		return pick_fee_setting.Fee4
	} else { // 50.1KG以上
		return pick_fee_setting.Fee5
	}
}
