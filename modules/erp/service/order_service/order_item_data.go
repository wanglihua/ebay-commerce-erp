package order_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Order_item_view_item_t struct {
	OrderItem   models.OrderItemInfo       `xorm:"extends" json:"oi"`
	Order       models.OrderInfo           `xorm:"extends" json:"o"`
	ShopAccount models.ShopAccountInfo     `xorm:"extends" json:"sa"`
	Platform    models.ShopAccountInfo     `xorm:"extends" json:"pf"`
	Product     models.ProductInfo         `xorm:"extends" json:"p"`
	Status      models.StatusInfo          `xorm:"extends" json:"s"`
}

/*
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
*/

func get_order_item_data_list_sql_builder() *core.SqlBuilder {
	sb := core.NewSqlBuilder()
	sb.Select("oi.*, o.*, sa.*, pf.*, p.*, s.*")
	sb.From(`erp_order_item oi
    INNER JOIN erp_order o ON o.id = oi.order_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
    INNER JOIN erp_ecommerce_platform pf ON pf.id = sa.platform_id
    INNER JOIN erp_product p ON p.id = oi.item_product_id
    INNER JOIN erp_status s ON s.id = o.status_id`)

	sb.OrderBy("oi.id DESC")

	return sb
}

func Get_order_item_data_list_by_product_id(db_session *xorm.Session, product_id int64) []Order_item_view_item_t {
	sb := get_order_item_data_list_sql_builder()
	sb.Where(fmt.Sprintf("oi.item_product_id = %d", product_id))

	order_item_list := make([]Order_item_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&order_item_list)
	core.HandleError(err)

	return order_item_list
}

func Get_order_item_data_list_by_product_id_top_n(db_session *xorm.Session, product_id int64, top_n int64) []Order_item_view_item_t {
	sb := get_order_item_data_list_sql_builder()
	sb.Where(fmt.Sprintf("oi.item_product_id = %d", product_id))
	sb.Limit(0, top_n)

	order_item_list := make([]Order_item_view_item_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&order_item_list)
	core.HandleError(err)

	return order_item_list
}
