package notice_read_service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"matrix/core"
)

func GetUserUnreadNoticeCount(user_id int64, db_session *xorm.Session) int64 {
	var sql = `SELECT count(n.id) AS count
FROM erp_notice n
WHERE n.id NOT IN (
    SELECT nr.notice_id
    FROM erp_notice_read nr
    WHERE nr.user_id = %d)`

	sql = fmt.Sprintf(sql, user_id)

	var count core.Count
	_, err := db_session.Sql(sql).Get(&count)
	core.HandleError(err)

	return count.Count
}
