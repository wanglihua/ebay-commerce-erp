package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"fmt"
	"matrix/modules/erp/service/container_transfer_service"
	"matrix/modules/erp/service/exchange_rate_service"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"matrix/modules/erp/service/container_transfer_log_service"
	"matrix/modules/erp/service/product_service"
)

type ErpContainerTransferEntry struct {
	*revel.Controller
}

func (c ErpContainerTransferEntry) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	return c.RenderTemplate("erp/hgbg/container_transfer_entry/container_transfer_entry_index.html")
}

func (c ErpContainerTransferEntry) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type cte_item_t struct {
		ContainerTransferEntry models.ContainerTransferEntryInfo `xorm:"extends" json:"cte"`
		ContainerTransfer      models.ContainerTransferInfo      `xorm:"extends" json:"ct"`
		Product                models.ProductInfo                `xorm:"extends" json:"p"`
	}

	/*

SELECT cte.*, ct.*, p.*
FROM erp_container_transfer_entry cte
    INNER JOIN erp_container_transfer ct ON cte.transfer_id = ct.id
    INNER JOIN erp_product p ON cte.product_id = p.id

	 */

	var sb = core.NewSqlBuilder()
	sb.Select("cte.*, ct.*, p.*")
	sb.Count("count(*) as entry_count")
	sb.From(`erp_container_transfer_entry cte
    INNER JOIN erp_container_transfer ct ON cte.transfer_id = ct.id
    INNER JOIN erp_product p ON cte.product_id = p.id`)

	filter = core.SqlAnd(filter, fmt.Sprintf("cte.transfer_id = %d", transfer_id))

	sb.Where(filter)

	sb.Limit(offset, limit)

	if order != "" {
		sb.OrderBy(order)
	} else {
		sb.OrderBy("id asc")
	}
	container_transfer_entry_list := make([]cte_item_t, 0, limit)
	err := db_session.Sql(sb.SqlList()).Find(&container_transfer_entry_list)
	core.HandleError(err)

	var entry_count struct {
		EntryCount int64 `xorm:"'entry_count'"`
	}
	_, err = db_session.Sql(sb.SqlCount()).Get(&entry_count)
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  container_transfer_entry_list,
		Total: entry_count.EntryCount,
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
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

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
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	var detail_form ContainerTransferEntryDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	var ct models.ContainerTransferInfo
	_, err = db_session.Id(detail_form.ContainerTransferEntry.TransferId).Get(&ct)
	core.HandleError(err)

	if container_transfer_can_modify(ct) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前状态不能修改！"})
	}

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	container_transfer_entry := detail_form.ContainerTransferEntry

	var product_sku = container_transfer_entry.ProductSku

	var product models.ProductInfo
	has, err := db_session.Where(fmt.Sprintf("sku = '%s'", product_sku)).Get(&product)
	core.HandleError(err)

	if has == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "商品SKU不存在！"})
	}

	var au_exchange_rate = exchange_rate_service.Get_aud_exchange_rate(db_session)
	//container_transfer_entry.Price = inventory_service.GetProductCnStockCostPrice(product.Id, db_session) / au_exchange_rate
	container_transfer_entry.Price = product_service.GetProductPrice(product.Id, db_session) / au_exchange_rate

	if product.AuHsCode == "" {
		return c.RenderJson(core.JsonResult{Success: false, Message: "商品澳洲海关编码未设置！"})
	} else {
		var hsCodeAu models.HsCodeAuDutyInfo
		has, err := db_session.Where("hs_code = ?", product.AuHsCode).Get(&hsCodeAu)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "商品澳洲海关编码已设置，但该海关编码不存在！"})
		}

		container_transfer_entry.AuHsId = hsCodeAu.Id
		container_transfer_entry.AuHsCode = hsCodeAu.HsCode
		container_transfer_entry.AuHsDutyRate = hsCodeAu.DutyRate
		container_transfer_entry.AuHsCnName = hsCodeAu.CnName
		container_transfer_entry.AuHsEnName = hsCodeAu.EnName
	}

	if product.CnHsCode == "" {
		return c.RenderJson(core.JsonResult{Success: false, Message: "商品中国海关编码未设置！"})
	} else {
		var hsCodeCn models.HsCodeCnReportElemInfo
		has, err := db_session.Where("hs_code = ?", product.CnHsCode).Get(&hsCodeCn)
		core.HandleError(err)

		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "商品中国海关编码已设置，但该海关编码不存在！"})
		}

		container_transfer_entry.CnHsId = hsCodeCn.Id
		container_transfer_entry.CnHsCode = hsCodeCn.HsCode
		container_transfer_entry.CnHsCnName = hsCodeCn.CnName
		container_transfer_entry.CnHsEnName = hsCodeCn.EnName
	}

	container_transfer_entry.ProductId = product.Id

	container_transfer_entry.BatchCode = ct.BatchCode

	var affected int64
	if detail_form.IsCreate() {
		affected, err = db_session.Insert(&container_transfer_entry)
		core.HandleError(err)

		container_transfer_log_service.Insert_container_transfer_log(ct, "新增调拨明细", db_session, c.Session)
	} else {
		//todo: 哪些字段需要更新
		affected, err = db_session.Id(container_transfer_entry.Id).Update(&container_transfer_entry)
		core.HandleError(err)

		container_transfer_log_service.Insert_container_transfer_log(ct, "修改调拨明细", db_session, c.Session)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	container_transfer_service.CalculateAndSaveWeightVolumeAmount(container_transfer_entry.TransferId, db_session)

	if container_transfer_entry.Price == 0 {
		return c.RenderJson(core.JsonResult{
			Success: false,
			Message: "数据保存成功，但商品采购价为零！修改商品采购价后，请重新保存本条数据！",
		})
	} else {
		return c.RenderJson(core.JsonResult{
			Success: true,
			Message: strconv.FormatInt(affected, 10) + "条数据保存成功!",
		})
	}
}

func (c ErpContainerTransferEntry) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

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

	var ct models.ContainerTransferInfo
	_, err = db_session.Id(containerTransferEntry.TransferId).Get(&ct)
	core.HandleError(err)

	if container_transfer_can_modify(ct) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前状态不能修改！"})
	}

	container_transfer_log_service.Insert_container_transfer_log(ct, "删除调拨明细", db_session, c.Session)

	affected, err := db_session.In("id", container_transfer_entry_id_list).Delete(new(models.ContainerTransferEntryInfo))
	core.HandleError(err)

	container_transfer_service.CalculateAndSaveWeightVolumeAmount(containerTransferEntry.TransferId, db_session)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}
