package shop_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

type Platform_shop_name_t struct {
	ShopId       int64         `xorm:"'shop_id'"`
	PlatformName string        `xorm:"'platform_name'"`
	ShopCode     string        `xorm:"'shop_code'"`
	ShopName     string        `xorm:"'shop_name'"`
}

func Get_shop_id_name_map(shop_id_list []int64, db_session *xorm.Session) map[int64]Platform_shop_name_t {
	var shop_id_name_map = make(map[int64]Platform_shop_name_t)

	if len(shop_id_list) == 0 {
		return shop_id_name_map
	}

	/*
SELECT sa.id AS shop_id, p.name AS platform_name, sa.code AS shop_code, sa.name AS shop_name
FROM erp_shop_account sa
    INNER JOIN erp_ecommerce_platform p ON sa.platform_id = p.id
WHERE sa.id IN (1, 2, 3)
	 */
	var sql = `SELECT sa.id AS shop_id, p.name AS platform_name, sa.code AS shop_code, sa.name AS shop_name
FROM erp_shop_account sa
    INNER JOIN erp_ecommerce_platform p ON sa.platform_id = p.id
WHERE sa.id IN (%s)`
	sql = fmt.Sprintf(sql, strings.Join(core.Int_list_to_str_list(shop_id_list), ","))

	var platform_shop_name_list = make([]Platform_shop_name_t, 0)
	err := db_session.Sql(sql).Find(&platform_shop_name_list)
	core.HandleError(err)

	for _, platform_shop_name := range platform_shop_name_list {
		shop_id_name_map[platform_shop_name.ShopId] = platform_shop_name
	}

	return shop_id_name_map
}
