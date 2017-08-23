package product_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/service/product_log_service"
	"matrix/modules/erp/models/kcgl_models"
)

func Get_product_id_by_sku(sku string, db_session *xorm.Session) int64 {
	sql := `SELECT p.id AS product_id FROM erp_product p WHERE p.sku = '%s'`
	sql = fmt.Sprintf(sql, sku)

	var product_id struct {
		ProductId int64 `xorm:"'product_id'"`
	}
	_, err := db_session.Sql(sql).Get(&product_id)
	core.HandleError(err)

	return product_id.ProductId
}

func Delete_product_list(product_id_list []int64, db_session *xorm.Session, web_session revel.Session) {

	//todo: InventoryCheckEntryInfo ProductId

	for _, product_id := range product_id_list {
		var product kcgl_models.ProductInfo
		has, err := db_session.Id(product_id).Get(&product)
		core.HandleError(err)

		if has {
			product_log_service.Insert_product_log(nil, &product, core.Db_operation_types.Delete, db_session, web_session)
		}
	}

	_, err := db_session.In("id", product_id_list).Delete(new(kcgl_models.ProductInfo))
	core.HandleError(err)
}
