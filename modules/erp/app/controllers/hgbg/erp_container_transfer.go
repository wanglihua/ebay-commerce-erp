
package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpContainerTransfer struct {
	*revel.Controller
}

func (c ErpContainerTransfer) Index() revel.Result {
	return c.RenderTemplate("erp/hgbg/container_transfer/container_transfer_index.html")
}

func (c ErpContainerTransfer) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	container_transfer_list := make([]models.ContainerTransferInfo, 0, limit)
	err := data_query.Limit(limit, offset).Find(&container_transfer_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(models.ContainerTransferInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  container_transfer_list,
		Total: count,
	})
}

type ContainerTransferDetailForm struct {
	ContainerTransfer models.ContainerTransferInfo `json:"transfer"`
}

func (f *ContainerTransferDetailForm) IsCreate() bool {
	return f.ContainerTransfer.Id == 0
}

func (f *ContainerTransferDetailForm) Valid(validation *revel.Validation) bool { 
	validation.Required(f.ContainerTransfer.BatchCode).Message("批次号不能为空！")

	validation.Required(f.ContainerTransfer.DestStockId).Message("目的仓库不能为空！")

	validation.Required(f.ContainerTransfer.DestStockName).Message("目的仓库不能为空！")

	validation.Required(f.ContainerTransfer.ContainerCode).Message("集装箱号不能为空！")
	if f.ContainerTransfer.ContainerCode != "" {
		validation.MinSize(f.ContainerTransfer.ContainerCode, 2).Message("集装箱号长度不能小于2！")
	}
	if f.ContainerTransfer.ContainerCode != "" {
		validation.MaxSize(f.ContainerTransfer.ContainerCode, 100).Message("集装箱号长度不能大于100！")
	}

	validation.Required(f.ContainerTransfer.BlNo).Message("提单号不能为空！")
	if f.ContainerTransfer.BlNo != "" {
		validation.MinSize(f.ContainerTransfer.BlNo, 2).Message("提单号长度不能小于2！")
	}
	if f.ContainerTransfer.BlNo != "" {
		validation.MaxSize(f.ContainerTransfer.BlNo, 100).Message("提单号长度不能大于100！")
	}

	validation.Required(f.ContainerTransfer.SealNo).Message("铅封号不能为空！")
	if f.ContainerTransfer.SealNo != "" {
		validation.MinSize(f.ContainerTransfer.SealNo, 2).Message("铅封号长度不能小于2！")
	}
	if f.ContainerTransfer.SealNo != "" {
		validation.MaxSize(f.ContainerTransfer.SealNo, 100).Message("铅封号长度不能大于100！")
	}

	validation.Required(f.ContainerTransfer.TotalVolumn).Message("总体积不能为空！")

	validation.Required(f.ContainerTransfer.TotalWeight).Message("总重量不能为空！")

	validation.Required(f.ContainerTransfer.TotalAmount1).Message("总金额1不能为空！")

	validation.Required(f.ContainerTransfer.TotalAmount2).Message("总金额2不能为空！")

	validation.Required(f.ContainerTransfer.PriceUpRate).Message("非征品上调不能为空！")

	validation.Required(f.ContainerTransfer.PriceDownRate).Message("征税品下调不能为空！")

	validation.Required(f.ContainerTransfer.CostUpRate).Message("成本上调不能为空！")

	validation.Required(f.ContainerTransfer.PackingTime).Message("装箱时间不能为空！")

	validation.Required(f.ContainerTransfer.DepartureTime).Message("到港时间不能为空！")

	validation.Required(f.ContainerTransfer.StockTime).Message("到仓时间不能为空！")

	return validation.HasErrors() == false
}

func (c ErpContainerTransfer) DetailData() revel.Result {
	db_session := c.DbSession

	var container_transfer_id int64
	c.Params.Bind(&container_transfer_id, "id")

	var container_transfer models.ContainerTransferInfo
	if container_transfer_id != 0 {
		has, err := db_session.Id(container_transfer_id).Get(&container_transfer)
		core.HandleError(err)
		if has == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "指定的集装箱调拨不存在！"})
		}
	}

	var detail_form ContainerTransferDetailForm
	detail_form.ContainerTransfer = container_transfer

	return c.RenderJson(detail_form)
}

func (c ErpContainerTransfer) Save() revel.Result {
	db_session := c.DbSession

	var detail_form ContainerTransferDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	container_transfer := detail_form.ContainerTransfer

	var affected int64
	if detail_form.IsCreate() { 
		affected, err = db_session.Insert(&container_transfer)
		core.HandleError(err)
	} else { 
		affected, err = db_session.Id(container_transfer.Id).Update(&container_transfer)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpContainerTransfer) Delete() revel.Result {
	db_session := c.DbSession

	container_transfer_id_list := make([]int64, 0)
	c.Params.Bind(&container_transfer_id_list, "id_list")

	affected, err := db_session.In("id", container_transfer_id_list).Delete(new(models.ContainerTransferInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

