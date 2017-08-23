package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models/hgbg_models"
	"bytes"
	"time"
	"matrix/modules/erp/app/controllers/hgbg/container_transfer_excel_in"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

type ErpContainerTransferExcelIn struct {
	*revel.Controller
}

func (c ErpContainerTransferExcelIn) Export() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var company_code string // elife, giant, ozo
	c.Params.Bind(&company_code, "company_code")

	var containerTransfer hgbg_models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&containerTransfer)
	core.HandleError(err)

	var excel_exporter = container_transfer_excel_in.GetExporter(company_code)

	xlsx_file := excel_exporter.GetExcelFile(&containerTransfer, db_session)

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := excel_exporter.GetExcelFileName(&containerTransfer) // fmt.Sprintf("中国海关申报元素导出_%s.xlsx", containerTransfer.BatchCode) //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}
