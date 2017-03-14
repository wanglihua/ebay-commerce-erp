package order_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Order_item_view_item_t struct {
	OrderItem models.OrderItemInfo `xorm:"extends" json:"oi"`
	Order     models.OrderInfo     `xorm:"extends" json:"o"`
	Product   models.ProductInfo   `xorm:"extends" json:"p"`
}

func Get_order_item_data_list_by_order_id(db_session *xorm.Session, order_id int64) []Order_item_view_item_t {
	sql := `
SELECT oi.*, o.*, p.*
FROM erp_order_item oi
    INNER JOIN erp_order o ON o.id = oi.order_id
    INNER JOIN erp_product p ON p.id = oi.item_product_id
WHERE oi.order_id = %d`
	sql = fmt.Sprintf(sql, order_id)

	order_item_list := make([]Order_item_view_item_t, 0)
	err := db_session.Sql(sql).Find(&order_item_list)
	core.HandleError(err)

	return order_item_list
}

func Get_order_item_data_list_by_product_id(db_session *xorm.Session, product_id int64) []Order_item_view_item_t {
	sql := `
SELECT oi.*, o.*, p.*
FROM erp_order_item oi
    INNER JOIN erp_order o ON o.id = oi.order_id
    INNER JOIN erp_product p ON p.id = oi.item_product_id
WHERE oi.item_product_id = %d`
	sql = fmt.Sprintf(sql, product_id)

	order_item_list := make([]Order_item_view_item_t, 0)
	err := db_session.Sql(sql).Find(&order_item_list)
	core.HandleError(err)

	return order_item_list
}
