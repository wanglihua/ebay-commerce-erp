package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"matrix/modules/erp/service/container_transfer_service"
	"matrix/modules/erp/service/inventory_service"
)

type ErpContainerTransferEntry struct {
	*revel.Controller
}

func (c ErpContainerTransferEntry) Index() revel.Result {
	return c.RenderTemplate("erp/hgbg/container_transfer_entry/container_transfer_entry_index.html")
}

func (c ErpContainerTransferEntry) ListData() revel.Result {
	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	filter = core.SqlAnd(filter, fmt.Sprintf("transfer_id = %d", transfer_id))

	query := db_session.Where(filter)

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	container_transfer_entry_list := make([]models.ContainerTransferEntryInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&container_transfer_entry_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ContainerTransferEntryInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  container_transfer_entry_list,
		Total: count,
	})
}

type ContainerTransferEntryDetailForm struct {
	ContainerTransferEntry models.ContainerTransferEntryInfo `json:"transfer_entry"`
}

func (f *ContainerTransferEntryDetailForm) IsCreate() bool {
	return f.ContainerTransferEntry.Id == 0
}

func (f *ContainerTransferEntryDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ContainerTransferEntry.TransferId).Message("集装箱调拨Id不能为空！")

	validation.Required(f.ContainerTransferEntry.ProductId).Message("商品Id不能为空！")

	validation.Required(f.ContainerTransferEntry.ProductSku).Message("SKU不能为空！")

	validation.Required(f.ContainerTransferEntry.Qty).Message("数量不能为空！")

	validation.Required(f.ContainerTransferEntry.BoxQty).Message("箱数不能为空！")

	validation.Required(f.ContainerTransferEntry.Price).Message("价格不能为空！")

	return validation.HasErrors() == false
}

func (c ErpContainerTransferEntry) DetailData() revel.Result {
	db_session := c.DbSession

	var container_transfer_entry_id int64
	c.Params.Bind(&container_transfer_entry_id, "id")

	var container_transfer_entry models.ContainerTransferEntryInfo
	if container_transfer_entry_id != 0 {
		has, err := db_session.Id(container_transfer_entry_id).Get(&container_transfer_entry)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的集装箱调拨明细不存在！"})
		}
	}

	var detail_form ContainerTransferEntryDetailForm
	detail_form.ContainerTransferEntry = container_transfer_entry

	return c.RenderJson(detail_form)
}

func (c ErpContainerTransferEntry) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ContainerTransferEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	container_transfer_entry := detail_form.ContainerTransferEntry

	var product_sku = container_transfer_entry.ProductSku

	var product models.ProductInfo
	has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_sku)).Get(&product)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "商品SKU不存在！" })
	}

	container_transfer_entry.Price = inventory_service.GetProductCnStockCostPrice(product.Id, db_session)

	if product.AuHsCode == "" {
		return c.RenderJson(core.JsonResult{Success: false, Message: "商品澳洲海关编码未设置！" })
	} else {
		var hsCodeAu models.HsCodeAuDutyInfo
		has, err := db_session.Where("hs_code = ?", product.AuHsCode).Get(&hsCodeAu)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "商品澳洲海关编码已设置，但该海关编码不存在！" })
		}

		container_transfer_entry.AuHsId = hsCodeAu.Id
		container_transfer_entry.AuHsCode = hsCodeAu.HsCode
		container_transfer_entry.AuHsDutyRate = hsCodeAu.DutyRate
		container_transfer_entry.AuHsCnName = hsCodeAu.CnName
		container_transfer_entry.AuHsEnName = hsCodeAu.EnName
	}

	if product.CnHsCode == "" {
		return c.RenderJson(core.JsonResult{Success: false, Message: "商品中国海关编码未设置！" })
	} else {
		var hsCodeCn models.HsCodeCnReportElemInfo
		has, err := db_session.Where("hs_code = ?", product.CnHsCode).Get(&hsCodeCn)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "商品中国海关编码已设置，但该海关编码不存在！" })
		}

		container_transfer_entry.CnHsId = hsCodeCn.Id
		container_transfer_entry.CnHsCode = hsCodeCn.HsCode
		container_transfer_entry.CnHsCnName = hsCodeCn.CnName
		container_transfer_entry.CnHsEnName = hsCodeCn.EnName
	}

	container_transfer_entry.ProductId = product.Id

	var containerTransfer models.ContainerTransferInfo
	_, err = db_session.Id(container_transfer_entry.TransferId).Get(&containerTransfer)
	core.HandleError(err)

	container_transfer_entry.BatchCode = containerTransfer.BatchCode

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&container_transfer_entry)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(container_transfer_entry.Id).Update(&container_transfer_entry)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	container_transfer_service.CalculateAndSaveWeightVolumeAmount(container_transfer_entry.TransferId, db_session)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpContainerTransferEntry) Delete() revel.Result {
	db_session := c.DbSession

	container_transfer_entry_id_list := make([]int64, 0)
	c.Params.Bind(&container_transfer_entry_id_list, "id_list")

	if len(container_transfer_entry_id_list) == 0 {
		return c.RenderJson(core.JsonResult{Success: false, Message: "请选择要删除的数据!"})
	}

	var entry_id = container_transfer_entry_id_list[0]
	var containerTransferEntry models.ContainerTransferEntryInfo
	_, err := db_session.Id(entry_id).Get(&containerTransferEntry)
	core.HandleError(err)

	affected, err := db_session.In("id", container_transfer_entry_id_list).Delete(new(models.ContainerTransferEntryInfo))
	core.HandleError(err)

	container_transfer_service.CalculateAndSaveWeightVolumeAmount(containerTransferEntry.TransferId, db_session)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
