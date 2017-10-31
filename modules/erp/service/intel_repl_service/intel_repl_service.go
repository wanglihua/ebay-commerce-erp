package intel_repl_service

import (
	"github.com/go-xorm/xorm"
	"matrix/core"
)

func Calculate_intel_repl(db_session *xorm.Session) {
	_, err := db_session.Exec("CALL intel_repl();")
	core.HandleError(err)
}
