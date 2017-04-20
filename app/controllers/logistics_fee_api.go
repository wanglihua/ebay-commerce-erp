package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/revel/revel"
)

type LogisticsFeeApi struct {
	*revel.Controller
}

type RenderJavascriptResult struct {
	Script string
}

func (r RenderJavascriptResult) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "text/javascript; charset=utf-8")
	resp.Out.Write([]byte(r.Script))
}

func (c LogisticsFeeApi) Calc() revel.Result {
	var item_num string
	c.Params.Bind(&item_num, "item_num")
	item_num = strings.TrimSpace(item_num)

	var zip_code string
	c.Params.Bind(&zip_code, "zip_code")
	zip_code = strings.TrimSpace(zip_code)

	//todo: 匹配的逻辑
	var fee float64 = rand.Float64()*5 + 5

	var callback_code = `if(log_fee_callback) { log_fee_callback(%f); }`
	callback_code = fmt.Sprintf(callback_code, fee)

	return &RenderJavascriptResult{Script: callback_code}
}

func (c LogisticsFeeApi) Test() revel.Result {
	return c.RenderTemplate("logistics_fee_api/logistics_fee_api_test.html")
}
