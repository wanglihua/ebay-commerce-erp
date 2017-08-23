package envelope_fee_rate_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/ddpz_models"
)

// 获得信封物流费率列表，按费用从小到大排序
func Get_envelope_fee_rate_list(db_session *xorm.Session) []ddpz_models.EnvelopeFeeRateInfo {
	var envelope_fee_rate_list = make([]ddpz_models.EnvelopeFeeRateInfo, 0)
	err := db_session.OrderBy("fee ASC").Find(&envelope_fee_rate_list)
	core.HandleError(err)
	return envelope_fee_rate_list
}
