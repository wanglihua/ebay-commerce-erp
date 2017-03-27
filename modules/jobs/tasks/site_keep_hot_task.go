package tasks

import (
	"net/http"
	"strconv"
	"github.com/revel/revel"
)

type SiteKeepHotTask struct {
}

const Site_keep_hot_task_name = "网站预热"

func (t SiteKeepHotTask) Name() string {
	return Site_keep_hot_task_name
}

func (t SiteKeepHotTask) Run() {
	http.Get("http://localhost:" + strconv.Itoa(revel.HttpPort) + "/hotload")
}
