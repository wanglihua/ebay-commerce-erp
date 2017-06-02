package product_cate_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

type ProductCateSelectItem struct {
	CateId   int64  `xorm:"'cate_id'" json:"cate_id"`
	CateName string `xorm:"'cate_name'" json:"cate_name"`
}

func GetProductCateSelectItemList(db_session *xorm.Session) []ProductCateSelectItem {
	var product_cate_select_item_list = make([]ProductCateSelectItem, 0)

	sql := `
SELECT c.id AS cate_id, concat_ws('-', l1_p.name, l2_p.name, c.name) AS cate_name
FROM erp_product_cate c
    LEFT JOIN erp_product_cate l1_p ON l1_p.id = c.l1_parent
    LEFT JOIN erp_product_cate l2_p ON l2_p.id = c.l2_parent
ORDER BY concat_ws('-', l1_p.name, l2_p.name, c.name) ASC;`

	err := db_session.Sql(sql).Find(&product_cate_select_item_list)
	core.HandleError(err)

	return product_cate_select_item_list
}
