package product_sale_setting_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/xsgl_models"
)

func Get_product_sale_setting(db_session *xorm.Session) xsgl_models.ProductSaleSettingInfo {
	var product_sale_setting xsgl_models.ProductSaleSettingInfo
	sql := "SELECT * FROM erp_product_sale_setting pss ORDER BY pss.update_time DESC LIMIT 0, 1"
	_, err := db_session.Sql(sql).Get(&product_sale_setting)
	core.HandleError(err)

	return product_sale_setting
}
