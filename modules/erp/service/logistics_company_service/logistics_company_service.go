package logistics_company_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func Get_company_id_name_map(company_id_list []int64, db_session *xorm.Session) map[int64]string {
	var company_id_name_map = make(map[int64]string)

	if len(company_id_list) == 0 {
		return company_id_name_map
	}

	var company_list = make([]models.LogisticsCompanyInfo, 0)
	err := db_session.In("id", company_id_list).Find(&company_list)
	core.HandleError(err)

	for _, company := range company_list {
		company_id_name_map[company.Id] = company.Name
	}

	return company_id_name_map
}
