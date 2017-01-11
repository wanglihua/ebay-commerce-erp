package controllers

import (
	"github.com/revel/revel"
)

type ErpDiaryIncomeReport struct {
	*revel.Controller
}

func (c ErpDiaryIncomeReport) Index() revel.Result {
	return c.RenderTemplate("erp/sppz/brand/brand_index.html")
}
