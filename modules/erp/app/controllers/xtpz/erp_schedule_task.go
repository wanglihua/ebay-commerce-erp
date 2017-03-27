package controllers

import (
	"github.com/revel/cron"
	"github.com/revel/revel"
	"matrix/modules/jobs/app/jobs"
	"matrix/modules/jobs/tasks"
)

type ErpScheduleTask struct {
	*revel.Controller
}

func (c ErpScheduleTask) Index() revel.Result {
	return c.RenderTemplate("erp/xtpz/schedule_task/schedule_task_index.html")
}

func (c ErpScheduleTask) TaskData() revel.Result {
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
			schedule_task.PrevTime = job_entry.Prev.Format("2006-01-02 15:04:05")
			schedule_task.NextTime = job_entry.Next.Format("2006-01-02 15:04:05")

			schedule_task_list = append(schedule_task_list, schedule_task)
		}
	}

	return c.RenderJson(map[string]interface{}{
		"task_list": schedule_task_list,
	})
}

func (c ErpScheduleTask) HistoryData() revel.Result {
	return c.RenderEmpty()
}
