package order_type_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func Get_order_type_id_name_map(type_id_list []int64, db_session *xorm.Session) map[int64]string {
	var type_id_name_map = make(map[int64]string)

	if len(type_id_list) == 0 {
		return type_id_name_map
	}

	var order_type_list = make([]models.OrderTypeInfo, 0)
	err := db_session.In("id", type_id_list).Find(&order_type_list)
	core.HandleError(err)


	for _, order_type := range order_type_list {
		type_id_name_map[order_type.Id] = order_type.Name
	}

	return type_id_name_map
}
