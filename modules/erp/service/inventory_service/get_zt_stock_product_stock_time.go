package inventory_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/modules/erp/models"
	"matrix/modules/erp/models/sppz_models"
	"matrix/core"
	"time"
)

/*
	CONTAINER_TRANSFER_STATUS_NEW    = 5400 // 集装箱调拨 已创建
	CONTAINER_TRANSFER_STATUS_DEPART = 5500 // 集装箱调拨 已出发
	CONTAINER_TRANSFER_STATUS_STOCK  = 5600 // 集装箱调拨 已到仓

# 获取某个商品的在途仓库的库存的到港时间

SELECT ct.id AS transfer_id, ct.batch_code AS batch_code, ct.stock_time AS stock_time
FROM erp_container_transfer ct
WHERE ct.status_id = 5500 AND ct.id IN (
    SELECT cte.transfer_id
    FROM erp_container_transfer_entry cte
    WHERE cte.product_id = 123
) AND ct.dest_port_id = (
    SELECT s.dest_port_id
    FROM erp_stock s
    WHERE s.stock_type_id = 3 AND s.id = 12
)
ORDER BY ct.stock_time ASC
LIMIT 1;

*/

type ZtProductInfo struct {
	TransferId int64     `xorm:"'transfer_id'"`
	BatchCode  string    `xorm:"'batch_code'"`
	StockTime  time.Time `xorm:"'stock_time'"`
}

func Get_zt_stock_product_stock_time(stock_id int64, product_id int64, db_session *xorm.Session) ZtProductInfo {
	var sql = `SELECT ct.id AS transfer_id, ct.batch_code AS batch_code, ct.stock_time AS stock_time
FROM erp_container_transfer ct
WHERE ct.status_id = %d AND ct.id IN (
    SELECT cte.transfer_id
    FROM erp_container_transfer_entry cte
    WHERE cte.product_id = %d
) AND ct.dest_port_id = (
    SELECT s.dest_port_id
    FROM erp_stock s
    WHERE s.stock_type_id = %d AND s.id = %d
)
ORDER BY ct.stock_time ASC
LIMIT 1;`
	sql = fmt.Sprintf(sql, models.CONTAINER_TRANSFER_STATUS_DEPART, product_id, sppz_models.STOCK_TYPE_ZT, stock_id)

	var zt_product_info ZtProductInfo
	_, err := db_session.Sql(sql).Get(&zt_product_info)
	core.HandleError(err)

	return zt_product_info
}
