package models

var ProductBindCols = struct {
	Id         string
	Code       string
	Name       string
	BindList   string
	PkgMtrlId  string
	StatusId   string
	CreateTime string
	UpdateTime string
	Version    string
}{
	Id:         "id",
	Code:       "code",
	Name:       "name",
	BindList:   "bind_list",
	PkgMtrlId:  "pkg_mtrl_id",
	StatusId:   "status_id",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Version:    "version",
}
