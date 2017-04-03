package order_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/service/logistics_setting_service"
)

// 计算订单的物流方式，
func Calculate_order_logistics(order_id int64, db_session *xorm.Session) {
	logistics_setting := logistics_setting_service.Get_logistics_setting(db_session)
	_ = logistics_setting


}
