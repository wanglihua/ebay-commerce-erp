package controllers

import (
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"time"
)

type ErpRefundReport struct {
	*revel.Controller
}

const REPORT_TYPE_DAY int64 = 1

const REPORT_TYPE_MONTH int64 = 2

func (c ErpRefundReport) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	// var db_session = c.DbSession

	c.RenderArgs["report_type_day"] = REPORT_TYPE_DAY
	c.RenderArgs["report_type_month"] = REPORT_TYPE_MONTH

	return c.RenderTemplate("erp/ppgl/refund_report/refund_report_index.html")
}

func (c ErpRefundReport) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	return c.RenderJson(core.GridResult{
		Data:  nil,
		Total: 0,
	})
}

func (c ErpRefundReport) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.PPGL_REFUND_REPORT, c.Session)

	// var db_session = c.DbSession

	/*

	xlsx_file := c.get_export_excel_content(db_session, order_refund_list)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("订单退款导出_%s.xlsx", time.Now().Format("2016-01-02 15:04"))

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
	*/

	return c.RenderBinary(nil, "", revel.Attachment, time.Now())
}
