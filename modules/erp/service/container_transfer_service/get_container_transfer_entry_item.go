package container_transfer_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func GetContainerTransferEntryItem(transfer_id int64, db_session *xorm.Session) []ContainerTransferEntryItem {
	/*
SELECT cte.*, p.*
FROM erp_container_transfer_entry cte
    INNER JOIN erp_product p ON cte.product_id = p.id
WHERE cte.transfer_id = 0;
	 */

	var sql = `SELECT cte.*, p.*
FROM erp_container_transfer_entry cte
    INNER JOIN erp_product p ON cte.product_id = p.id
WHERE cte.transfer_id = %d;`
	sql = fmt.Sprintf(sql, transfer_id)

	var cte_item_list = make([]ContainerTransferEntryItem, 0)
	core.HandleError(db_session.Sql(sql).Find(&cte_item_list))

	return cte_item_list
}
