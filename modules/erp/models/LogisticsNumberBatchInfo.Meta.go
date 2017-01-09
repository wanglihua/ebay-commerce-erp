package models

var LogisticsNumberBatchCols = struct {
	Id string
	ModeId string
	Prefix string
	Suffix string
	BeginNum string
	EndNum string
	CurNum string
	Remind string
	StatusId string
	CreateTime string
	UpdateTime string
	Version string
}{
	Id: "id",
	ModeId: "mode_id",
	Prefix: "prefix",
	Suffix: "suffix",
	BeginNum: "begin_num",
	EndNum: "end_num",
	CurNum: "cur_num",
	Remind: "remind",
	StatusId: "status_id",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version: "version",
}
