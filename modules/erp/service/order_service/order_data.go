package order_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

type Order_view_item_t struct {
	models.OrderInfo              `xorm:"extends" json:"o"`
	models.OrderTypeInfo          `xorm:"extends" json:"ot"`
	models.OrderHoldTypeInfo      `xorm:"extends" json:"ht"`
	models.ECommercePlatformInfo  `xorm:"extends" json:"pf"`
	models.ShopAccountInfo        `xorm:"extends" json:"sa"`
	models.StatusInfo             `xorm:"extends" json:"s"`
}

func Get_order_data_list(db_session *xorm.Session, filter string, order string, offset int64, limit int64) ([]Order_view_item_t, int64) {

	/*
SELECT o.*, ot.*, ht.*, pf.*, sa.*, s.*
FROM erp_order o
    INNER JOIN erp_order_type ot ON ot.id = o.type_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
    inner join erp_ecommerce_platform pf on pf.id = sa.platform_id
    INNER JOIN erp_status s ON s.id = o.status_id
	 */

	sql_builder := core.NewSqlBuilder()
	sql_builder.Select("o.*, ot.*, ht.*, pf.*, sa.*, s.*")
	sql_builder.Count("count(o.id) as order_count")
	sql_builder.From(`
erp_order o
    INNER JOIN erp_order_type ot ON ot.id = o.type_id
    LEFT JOIN erp_order_hold_type ht ON ht.id = o.hold_type_id
    INNER JOIN erp_shop_account sa ON sa.id = o.shop_id
    inner join erp_ecommerce_platform pf on pf.id = sa.platform_id
    INNER JOIN erp_status s ON s.id = o.status_id
	`)
	sql_builder.Where(filter)

	if order == "" {
		sql_builder.OrderBy("o.id DESC")
	} else {
		sql_builder.OrderBy(order)
	}

	sql_builder.Limit(offset, limit)

	order_list := make([]Order_view_item_t, 0)
	err := db_session.Sql(sql_builder.SqlList()).Find(&order_list)
	core.HandleError(err)

	var order_count struct {
		Count int64 `xorm:"'order_count'"`
	}
	_, err = db_session.Sql(sql_builder.SqlCount()).Get(&order_count)
	core.HandleError(err)

	return order_list, order_count.Count
}
