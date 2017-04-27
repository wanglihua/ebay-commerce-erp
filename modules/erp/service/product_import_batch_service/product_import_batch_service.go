package product_import_batch_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func  Get_product_import_batch_code_list(db_session *xorm.Session) []int64 {
	// 获得导入的批次列表
	type batch_code_t struct {
		BatchCode int64 `xorm:"'batch_code'"`
	}
	var sql = `SELECT DISTINCT batch_code FROM erp_product_import_batch ORDER BY batch_code DESC`
	var batch_code_t_list = make([]batch_code_t, 0)
	err := db_session.Sql(sql).Find(&batch_code_t_list)
	core.HandleError(err)

	var batch_code_list = make([]int64, 0)
	for _, batch_code := range batch_code_t_list {
		batch_code_list = append(batch_code_list, batch_code.BatchCode)
	}

	return batch_code_list
}
