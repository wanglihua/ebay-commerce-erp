package controllers

import (
	"github.com/revel/cron"
	"github.com/revel/revel"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/jobs/app/jobs"
	"matrix/modules/jobs/tasks"
)

type ErpScheduleTask struct {
	*revel.Controller
}

func (c ErpScheduleTask) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SCHEDULE_TASK, c.Session)

	return c.RenderTemplate("erp/xtpz/schedule_task/schedule_task_index.html")
}

func (c ErpScheduleTask) TaskData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SCHEDULE_TASK, c.Session)

	cast_job := func(job cron.Job) *jobs.Job {
		return job.(*jobs.Job)
	}

	type schedule_task_view_item_t struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		PrevTime string `json:"prev_time"`
		NextTime string `json:"next_time"`
	}

	entries := jobs.MainCron.Entries()

	schedule_task_list := make([]schedule_task_view_item_t, 0)
	for _, job_entry := range entries {
		job := cast_job(job_entry.Job)
		if job.Name != tasks.Site_keep_hot_task_name {
			var schedule_task schedule_task_view_item_t
			schedule_task.Name = job.Name
			schedule_task.Status = job.Status()

			if job_entry.Prev.IsZero() == false {
				schedule_task.PrevTime = job_entry.Prev.Format("2006-01-02 15:04:05")
			} else {
				schedule_task.PrevTime = ""
			}

			if job_entry.Next.IsZero() == false {
				schedule_task.NextTime = job_entry.Next.Format("2006-01-02 15:04:05")
			} else {
				schedule_task.NextTime = ""
			}

			schedule_task_list = append(schedule_task_list, schedule_task)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"task_list": schedule_task_list,
	})
}

func (c ErpScheduleTask) HistoryData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_SCHEDULE_TASK, c.Session)

	return c.RenderEmpty()
}
