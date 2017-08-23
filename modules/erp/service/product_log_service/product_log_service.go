package product_log_service

import (
	"fmt"
	"strings"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/erp/models/kcgl_models"
)

// 插入商品日志
func Insert_product_log(product_from_ui *kcgl_models.ProductInfo /*更新操作时使用*/ , product_in_db *kcgl_models.ProductInfo, operation_type core.Db_operation_type, db_session *xorm.Session, web_session revel.Session) {
	var product_log kcgl_models.ProductLogInfo

	product_log.ProductId = product_in_db.Id
	product_log.ProductCode = product_in_db.Sku

	login_user := core.GetLoginUser(web_session)
	product_log.UserId = login_user.UserId
	product_log.LoginName = login_user.LoginName

	if operation_type == core.Db_operation_types.Delete {
		product_log.Operation = "删除数据"
	} else if operation_type == core.Db_operation_types.Insert {
		//获得详细的插入数据
		product_log.Operation = "插入数据"
	} else if operation_type == core.Db_operation_types.Update {
		//获得详细的字段修改数据
		product_log.Operation = "更新数据。"
		var operation_list = make([]string, 0)
		if product_in_db.Sku != product_from_ui.Sku {
			operation_list = append(operation_list, fmt.Sprintf("SKU: %s => %s", product_in_db.Sku, product_from_ui.Sku))
		}
		if product_in_db.Name != product_from_ui.Name {
			operation_list = append(operation_list, fmt.Sprintf("名称: %s => %s", product_in_db.Name, product_from_ui.Name))
		}
		if product_in_db.Desc != product_from_ui.Desc {
			operation_list = append(operation_list, fmt.Sprintf("描述: %s => %s", product_in_db.Desc.String, product_from_ui.Desc.String))
		}
		product_log.Operation += strings.Join(operation_list, ", ")

		//todo: 可以继续补充其它详细信息
	}

	_, err := db_session.Insert(&product_log)
	core.HandleError(err)
}
