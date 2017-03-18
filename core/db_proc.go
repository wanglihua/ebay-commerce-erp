package core

import (
	"fmt"
	"strconv"

	"github.com/go-xorm/xorm"
)

// 约定 proc_name 以版本号 _v0001 等结尾
func CreateProc(proc_name string, create_proc_sql string, db_session *xorm.Session) {
	proc_name_full_len := len(proc_name)
	proc_version_len := len("_v0001")

	// proc_name 长度不符合要求直接返回
	if proc_name_full_len <= proc_version_len {
		return
	}

	//判断该版本的存储过程在数据库中是否存在，如果存在就直接返回
	count_sql := `SELECT count(*) AS proc_count FROM mysql.proc WHERE db = SCHEMA() AND name = ?;`

	var proc_count struct {
		ProcCount int `xorm:"'proc_count'"`
	}
	_, err := db_session.Sql(count_sql, proc_name).Get(&proc_count)
	HandleError(err)

	if proc_count.ProcCount != 0 {
		return
	}

	//创建之前先删除以前的版本
	proc_name_without_version := string(([]rune(proc_name))[0:proc_name_full_len-proc_version_len]) // proc_name

	proc_version_str := ([]rune(proc_name))[proc_name_full_len-proc_version_len+2: proc_name_full_len] // 0001
	proc_version, err := strconv.Atoi(string(proc_version_str))
	HandleError(err)

	for version := 1; version < proc_version; version++ {
		proc_name_with_version := fmt.Sprintf("%s_v%04d", proc_name_without_version, version)
		drop_sql := fmt.Sprintf(`DROP PROCEDURE IF EXISTS %s;`, proc_name_with_version)
		_, err = db_session.Exec(drop_sql)
		HandleError(err)
	}

	//创建存储过程
	_, err = db_session.Exec(create_proc_sql)
	HandleError(err)
}
