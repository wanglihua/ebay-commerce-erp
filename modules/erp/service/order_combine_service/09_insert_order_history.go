package order_combine_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/service/order_history_service"
	"fmt"
	"strings"
	"matrix/core"
)

func insert_order_history(combine_order_id int64, order_id_list []int64, user core.LoginUser, db_session *xorm.Session) {

	// 订单历史
	var combine_order ddlc_models.OrderInfo
	_, err := db_session.Id(combine_order_id).Get(&combine_order)
	core.HandleError(err)

	var combined_order_list = make([]ddlc_models.OrderInfo, 0)
	core.HandleError(db_session.In("id", order_id_list).Find(&combined_order_list))

	var combined_order_code_list = make([]string, 0)
	for _, order := range combined_order_list {
		combined_order_code_list = append(combined_order_code_list, order.Code)
	}

	order_history_service.Insert_order_history(combine_order, user.UserId, fmt.Sprintf("订单合并,被和并订单编码: %s, 工号：%s", strings.Join(combined_order_code_list, ","), user.LoginName), db_session)

	for _, order := range combined_order_list {
		order_history_service.Insert_order_history(order, user.UserId, fmt.Sprintf("订单被合并而取消, 并入订单编码: %s， 工号：%s", combine_order.Code, user.LoginName), db_session)
	}
}
