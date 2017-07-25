package order_hold_type_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"matrix/modules/erp/models"
)

func Get_order_hold_type_id_name_map( db_session *xorm.Session) map[int64]string {
	var hold_type_id_name_map = make(map[int64]string)


	var hold_type_list = make([]models.OrderHoldTypeInfo, 0)
	err := db_session.Find(&hold_type_list)
	core.HandleError(err)

	for _, hold_type := range hold_type_list {
		hold_type_id_name_map[hold_type.Id] = hold_type.Name
	}

	return hold_type_id_name_map
}
