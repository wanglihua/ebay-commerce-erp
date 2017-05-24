package tasks

import (
	"matrix/app"
	"matrix/core/lic"
)

type SystemCheckAppLicTask struct {
}

const System_check_app_lic_task_name = "check_al"

func (t SystemCheckAppLicTask) Name() string {
	return System_check_app_lic_task_name
}

func (t SystemCheckAppLicTask) Run() {
	lic.App_lic_date = lic.GetAppLicDate(app.AppName)
}
