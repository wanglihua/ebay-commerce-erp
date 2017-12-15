package controllers

import (
	"github.com/revel/revel"
)

type ErpDiaryIncomeAmountReport struct {
	*revel.Controller
}

func (c ErpDiaryIncomeAmountReport) Index() revel.Result {
	return c.RenderTemplate("erp/xsbb/diary_income_amount_report/diary_income_amount_report_index.html")
}
