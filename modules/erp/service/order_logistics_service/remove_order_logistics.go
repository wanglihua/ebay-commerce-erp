package order_logistics_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/models"
	"matrix/core"
)

//这个方法在一些Hold之后会被调用，所以不能在这时去掉订单的Hold信息
func remove_order_logistics(order models.OrderInfo, db_session *xorm.Session) {

	order.StockId = 0
	order.AddressId = 0
	order.CompanyId = 0
	order.ModeId = 0

	_, err := db_session.Id(order.Id).Cols("stock_id", "address_id", "company_id", "mode_id").Update(&order)
	core.HandleError(err)

	_, err = db_session.Where(fmt.Sprintf("order_id = %d", order.Id)).Delete(new(models.OrderLogisticsInfo))
	core.HandleError(err)
}
