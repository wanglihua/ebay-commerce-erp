package tasks

import (
	"net/http"
	"strconv"
	"github.com/revel/revel"
)

type SiteKeepHotTask struct {
}

func (t SiteKeepHotTask) Name() string {
	return "网站预热"
}

func (t SiteKeepHotTask) Run() {
	http.Get("http://localhost:" + strconv.Itoa(revel.HttpPort) + "/hotload")
}
