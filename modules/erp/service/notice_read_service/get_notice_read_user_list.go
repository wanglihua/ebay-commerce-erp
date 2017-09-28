package notice_read_service

import (
	"matrix/core"
	"github.com/go-xorm/xorm"
	"fmt"
)

type NoticeReadUser struct {
	UserName string    `xorm:"'user_name'"`
	ReadTime core.Time `xorm:"'read_time'"`
}

/*
SELECT u.nick_name AS user_name, nr.create_time AS read_time
FROM erp_notice n
    INNER JOIN erp_notice_read nr ON n.id = nr.notice_id
    INNER JOIN auth_user u ON u.id = nr.user_id
WHERE n.id = 2
ORDER BY n.create_time ASC;

 */

func GetNoticeReadUserList(notice_id int64, db_session *xorm.Session) []NoticeReadUser {

	var sql = `
SELECT u.nick_name AS user_name, nr.create_time AS read_time
FROM erp_notice n
    INNER JOIN erp_notice_read nr ON n.id = nr.notice_id
    INNER JOIN auth_user u ON u.id = nr.user_id
WHERE n.id = %d
ORDER BY n.create_time ASC;
	`
	sql = fmt.Sprintf(sql, notice_id)

	var notice_read_user_list = make([]NoticeReadUser, 0)
	core.HandleError(db_session.Sql(sql).Find(&notice_read_user_list))

	return notice_read_user_list
}
