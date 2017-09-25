package notice_read_service

import (
	"github.com/go-xorm/xorm"
	"matrix/modules/erp/models/xtgl_models"
	"matrix/core"
)

// 标记用户已读通知
func MarkNoticeUserRead(user_id int64, notice_id int64, db_session *xorm.Session) {
	var notice_read xtgl_models.NoticeReadInfo
	has, err := db_session.Where("user_id = ? AND notice_id = ?", user_id, notice_id).Get(&notice_read)
	core.HandleError(err)

	if has == false {
		notice_read.UserId = user_id
		notice_read.NoticeId = notice_id

		_, err := db_session.Insert(&notice_read)
		core.HandleError(err)
	}
}
