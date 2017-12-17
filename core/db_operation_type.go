package core

// 数据库操作类型

type Db_operation_type int64

var Db_operation_types = struct {
	Insert Db_operation_type
	Update Db_operation_type
	Delete Db_operation_type
}{
	Insert: 1,
	Update: 2,
	Delete: 3,
}
