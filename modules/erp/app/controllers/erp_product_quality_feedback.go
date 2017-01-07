package controllers

import (
	"github.com/revel/revel"
)

type ErpProductQualityFeedback struct {
	*revel.Controller
}

func (c ErpProductQualityFeedback) Index() revel.Result {
	return c.RenderTemplate("erp/product_quality_feedback/product_quality_feedback_index.html")
}
