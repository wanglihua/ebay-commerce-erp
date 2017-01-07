package controllers

import (
	"github.com/revel/revel"
)

type ErpProductQualityFeedback struct {
	*revel.Controller
}

func (c ErpProductQualityFeedback) Index() revel.Result {
	return c.RenderTemplate("erp/kcgl/product_quality_feedback/product_quality_feedback_index.html")
}
