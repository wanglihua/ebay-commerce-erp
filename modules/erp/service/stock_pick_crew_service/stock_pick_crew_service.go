package stock_pick_crew_service

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"matrix/core"
)

// 用户是不是墨尔本仓库的拣货人员
func Is_pick_crew_has_mel_stock(user_id int64, db_session *xorm.Session) bool {
	var stock_pick_crew_count struct {
		Count int64 `xorm:"'stock_pick_crew_count'"`
	}
	/*
SELECT count(spc.id) AS stock_pick_crew_count
FROM erp_stock_pick_crew spc
WHERE spc.user_id = 0
      AND spc.stock_id = (SELECT s.id AS stock_id FROM erp_stock s WHERE s.short_name = 'MEL' AND s.name LIKE '%墨尔本%')
	 */
	var sql_is_mel_pick_screw = `
SELECT count(spc.id) AS stock_pick_crew_count
FROM erp_stock_pick_crew spc
WHERE spc.user_id = %d
      AND spc.stock_id = (SELECT s.id AS stock_id FROM erp_stock s WHERE s.short_name = 'MEL' AND s.name LIKE '%%墨尔本%%')`
	sql_is_mel_pick_screw = fmt.Sprintf(sql_is_mel_pick_screw, user_id)

	_, err := db_session.Sql(sql_is_mel_pick_screw).Get(&stock_pick_crew_count)
	core.HandleError(err)

	if stock_pick_crew_count.Count == 0 {
		return false
	} else {
		return true
	}
}
