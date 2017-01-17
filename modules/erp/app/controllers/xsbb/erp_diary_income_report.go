package controllers

import (
	"github.com/revel/revel"
)

type ErpDiaryIncomeReport struct {
	*revel.Controller
}

func (c ErpDiaryIncomeReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/diary_income_report/diary_income_report_index.html")
}
