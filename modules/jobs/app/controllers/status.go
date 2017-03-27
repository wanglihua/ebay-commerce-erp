package controllers

import (
	"github.com/revel/cron"
	"github.com/revel/revel"
	"matrix/modules/jobs/app/jobs"
	"matrix/modules/jobs/tasks"
)

type Jobs struct {
	*revel.Controller
}

func (c Jobs) Status() revel.Result {
	/* 远程也能查看
	remoteAddress := c.Request.RemoteAddr
	if revel.Config.BoolDefault("jobs.acceptproxyaddress", false) {
		if proxiedAddress, isProxied := c.Request.Header["X-Forwarded-For"]; isProxied {
			remoteAddress = proxiedAddress[0]
		}
	}
	if !strings.HasPrefix(remoteAddress, "127.0.0.1") && !strings.HasPrefix(remoteAddress, "::1") &&
		!strings.HasPrefix(remoteAddress, "[127.0.0.1]") && !strings.HasPrefix(remoteAddress, "[::1]") {
		return c.Forbidden("%s is not local", remoteAddress)
	}
	*/

	entries := jobs.MainCron.Entries()

	filtered_entries := make([]*cron.Entry, 0)
	for _, job_entry := range entries {
		job := cast_job(job_entry.Job)
		if job.Name != tasks.Site_keep_hot_task_name {
			filtered_entries = append(filtered_entries, job_entry)
		}
	}

	c.RenderArgs["entries"] = filtered_entries
	return c.RenderTemplate("jobs/status.html")
}

func init() {
	revel.TemplateFuncs["cast_job"] = cast_job
}

func cast_job(job cron.Job) *jobs.Job {
	return job.(*jobs.Job)
}
