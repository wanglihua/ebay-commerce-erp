package order_import_batch_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
	"fmt"
)

func Get_order_import_batch_code_list(status_id int64, db_session *xorm.Session) []int64 {
	// 获得导入的批次列表
	type batch_code_t struct {
		BatchCode int64 `xorm:"'batch_code'"`
	}

	// SELECT DISTINCT batch_code FROM erp_order WHERE batch_code IS NOT NULL AND batch_code <> 0 AND status_id = 0 ORDER BY batch_code DESC;

	var sb = core.NewSqlBuilder()
	sb.Select("DISTINCT batch_code")
	sb.From("erp_order")

	var filter = "batch_code IS NOT NULL AND batch_code <> 0"
	if status_id != 0 {
		filter = core.SqlAnd(filter, fmt.Sprintf("status_id = %d", status_id))
	}

	sb.Where(filter)
	sb.OrderBy("batch_code DESC")

	var batch_code_t_list = make([]batch_code_t, 0)
	err := db_session.Sql(sb.SqlList()).Find(&batch_code_t_list)
	core.HandleError(err)

	var batch_code_list = make([]int64, 0)
	for _, batch_code := range batch_code_t_list {
		batch_code_list = append(batch_code_list, batch_code.BatchCode)
	}

	return batch_code_list
}
