package product_bind_log_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/status_service"
)

// 插入商品绑定日志
func Insert_product_bind_log(product_bind_from_ui *models.ProductBindInfo /*更新操作时使用*/ , product_bind_in_db *models.ProductBindInfo, operation_type core.Db_operation_type, db_session *xorm.Session, web_session revel.Session) {
	var product_bind_log models.ProductBindLogInfo

	product_bind_log.ProductBindId = product_bind_in_db.Id
	product_bind_log.ProductBindCode = product_bind_in_db.Code

	login_user := core.GetLoginUser(web_session)
	product_bind_log.UserId = login_user.UserId
	product_bind_log.LoginName = login_user.LoginName

	if operation_type == core.Db_operation_types.Delete {
		product_bind_log.Operation = "删除数据"
	} else if operation_type == core.Db_operation_types.Insert {
		//获得详细的插入数据
		product_bind_log.Operation = "插入数据"
	} else if operation_type == core.Db_operation_types.Update {
		//获得详细的字段修改数据
		product_bind_log.Operation = "更新数据。"

		// todo 商品日志里面也可以这样判断
		var operation_list = make([]string, 0)
		if product_bind_in_db.Code != product_bind_from_ui.Code {
			operation_list = append(operation_list, fmt.Sprintf("编号: %s => %s", product_bind_in_db.Code, product_bind_from_ui.Code))
		}
		if product_bind_in_db.Name != product_bind_from_ui.Name {
			operation_list = append(operation_list, fmt.Sprintf("名称: %s => %s", product_bind_in_db.Name, product_bind_from_ui.Name))
		}
		if product_bind_in_db.BindList != product_bind_from_ui.BindList {
			operation_list = append(operation_list, fmt.Sprintf("捆绑商品列表: %s => %s", product_bind_in_db.BindList, product_bind_from_ui.BindList))
		}
		if product_bind_in_db.Desc != product_bind_from_ui.Desc {
			operation_list = append(operation_list, fmt.Sprintf("描述: %s => %s", product_bind_in_db.Desc, product_bind_from_ui.Desc))
		}
		if product_bind_in_db.StatusId != product_bind_from_ui.StatusId {
			var status_name_in_db = status_service.Get_status_name(product_bind_in_db.StatusId, db_session)
			var status_name_from_ui = status_service.Get_status_name(product_bind_from_ui.StatusId, db_session)
			operation_list = append(operation_list, fmt.Sprintf("状态： %s => %s", status_name_in_db, status_name_from_ui))
		}

		product_bind_log.Operation += strings.Join(operation_list, ", ")
	}

	_, err := db_session.Insert(&product_bind_log)
	core.HandleError(err)
}
