package shop_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

type Shop_select_item_t struct {
	Id   int64     `xorm:"'shop_id'" json:"id"`
	Name string    `xorm:"'shop_name'" json:"name"`
}

func Get_shop_select_item_list(db_session *xorm.Session) []Shop_select_item_t {

	/*
SELECT sa.id AS shop_id, concat(ep.name, ' ', sa.name) AS shop_name
FROM erp_shop_account sa
    INNER JOIN erp_ecommerce_platform ep ON ep.id = sa.platform_id
	 */

	sql := `
SELECT sa.id AS shop_id, concat(ep.name, ' ', sa.name) AS shop_name
FROM erp_shop_account sa
    INNER JOIN erp_ecommerce_platform ep ON ep.id = sa.platform_id`

	shop_list := make([]Shop_select_item_t, 0)
	err := db_session.Sql(sql).Find(&shop_list)
	core.HandleError(err)

	return shop_list
}
