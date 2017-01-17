package controllers

import (
	"github.com/revel/revel"
)

type ErpProductFeedReport struct {
	*revel.Controller
}

func (c ErpProductFeedReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/product_feed_report/product_feed_report_index.html")
}
