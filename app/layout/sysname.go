package layout

import (
	"github.com/go-xorm/xorm"
)

func GetSysName(dbsession *xorm.Session) string {
	return "ERP"
}
