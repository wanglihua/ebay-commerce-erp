package product_cate_service

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
)

type ProductCateSelectItem struct {
	CateId   int64  `xorm:"'cate_id'" json:"cate_id"`
	CateName string `xorm:"'cate_name'" json:"cate_name"`
}

func GetProductCateSelectItemList(db_session *xorm.Session) []ProductCateSelectItem {
	var product_cate_select_item_list = make([]ProductCateSelectItem, 0)

	sql := `
SELECT c.id AS cate_id, concat_ws('-', p.name, c.name) AS cate_name
FROM erp_product_cate c
    LEFT JOIN erp_product_cate p ON c.parent_id = p.id
ORDER BY concat_ws('-', p.name, c.name) ASC;`

	err := db_session.Sql(sql).Find(&product_cate_select_item_list)
	core.HandleError(err)

	return product_cate_select_item_list
}
