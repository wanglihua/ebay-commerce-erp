package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"github.com/tealeg/xlsx"
	"bytes"
	"time"
)

type ErpContainerTransferExcelIn struct {
	*revel.Controller
}

func (c ErpContainerTransferExcelIn) Export() revel.Result {
	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var company_code string
	c.Params.Bind(&company_code, "company_code")

	var containerTransfer models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&containerTransfer)
	core.HandleError(err)

	xlsx_file := xlsx.NewFile()

	//先把文件内容写到buffer中，再从中读取
	var buffer bytes.Buffer
	err = xlsx_file.Write(&buffer)
	core.HandleError(err)

	reader := bytes.NewReader(buffer.Bytes()) //bytes.Reader 实现了 io.Reader 接口

	file_name := fmt.Sprintf("中国海关申报元素导出_%s.xlsx", containerTransfer.BatchCode) //"file.xlsx" //可以根据参数生成文件名

	return c.RenderBinary(reader, file_name, revel.Attachment, time.Now())
}
