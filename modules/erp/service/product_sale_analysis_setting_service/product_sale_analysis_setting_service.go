package product_sale_analysis_setting_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models/spbb_models"
)

// 获取最新的物流配置信息，如果数据库中不存在，就返回零值
func Get_product_sale_analysis_setting(db_session *xorm.Session) spbb_models.ProductSaleAnalysisSettingInfo {
	var product_sale_analysis_setting spbb_models.ProductSaleAnalysisSettingInfo
	sql := "SELECT * FROM erp_product_sale_analysis_setting s ORDER BY s.update_time DESC LIMIT 0, 1"
	_, err := db_session.Sql(sql).Get(&product_sale_analysis_setting)
	core.HandleError(err)

	return product_sale_analysis_setting
}
