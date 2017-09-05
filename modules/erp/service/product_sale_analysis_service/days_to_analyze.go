package product_sale_analysis_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/service/product_sale_analysis_setting_service"
)

func GetDaysToAnalyze(db_session *xorm.Session) int {
	var product_sale_analysis_setting = product_sale_analysis_setting_service.Get_product_sale_analysis_setting(db_session)
	return int(product_sale_analysis_setting.DayAnalysisDays)
}
