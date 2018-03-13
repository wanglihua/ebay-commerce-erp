package order_combine_service

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/ddlc_models"
	"matrix/modules/erp/models/ddpz_models"
	"time"
	"strings"
	"matrix/modules/erp/models"
)

var order_code_seq int64 = 1

func insert_combine_order_from_first_order(first_order_id int64, order_id_str_list []string, db_session *xorm.Session) int64 {
	// 获取被合并订单中的第一个订单
	var first_order ddlc_models.OrderInfo
	_, err := db_session.Id(first_order_id).Get(&first_order)
	core.HandleError(err)

	// 根据第一个订单复制生成合并的订单
	var combine_order ddlc_models.OrderInfo = first_order
	combine_order.Id = 0

	combine_order.Code = time.Now().Format("200601021504") + fmt.Sprintf("%02d", order_code_seq)
	order_code_seq = (order_code_seq + 1) % 100 // 加一个序列号

	combine_order.TypeId = ddpz_models.ORDER_TYPE_HBDD // 合并订单
	combine_order.StatusId = models.ORDER_STATUS_YZF   // 只有 已支付的，待配货的 订单才会被合并，合并后的订单的状态设置为 已支付

	combine_order.TagId = 0
	combine_order.StockId = 0
	combine_order.CompanyId = 0
	combine_order.ModeId = 0
	combine_order.EbayOrderId = ""

	_, err = db_session.Insert(&combine_order)
	core.HandleError(err)

	var combine_order_id = combine_order.Id

	// 插入合并订单和被合并订单的关联记录
	var combined_order_list = make([]ddlc_models.OrderInfo, 0)
	err = db_session.Where(fmt.Sprintf("id IN (%s)", strings.Join(order_id_str_list, ","))).Find(&combined_order_list)
	core.HandleError(err)

	for _, combined_order := range combined_order_list {

		var combine_order_info ddlc_models.CombineOrderInfo

		combine_order_info.OrderId = combine_order_id
		combine_order_info.OrderCode = combine_order.Code
		combine_order_info.SubOrderId = combined_order.Id
		combine_order_info.SubOrderCode = combined_order.Code
		combine_order_info.CreateTime = core.Time(time.Now())

		_, err := db_session.Insert(&combine_order_info)
		core.HandleError(err)
	}

	return combine_order_id
}
