package controllers

import (
	"encoding/json"
	"strconv"
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
	"matrix/modules/erp/service/container_transfer_service"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
	"time"
	"matrix/modules/erp/service/container_transfer_log_service"
)

type ErpContainerTransfer struct {
	*revel.Controller
}

func (c ErpContainerTransfer) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	var can_to_zt_stock = pass_service.Is_login_user_has_pass(pass_res.HGBG_CONTAINER_TRANSFER_TO_ZT_STOCK, c.Session)
	var can_to_au_stock = pass_service.Is_login_user_has_pass(pass_res.HGBG_CONTAINER_TRANSFER_TO_AU_STOCK, c.Session)

	if can_to_zt_stock {
		c.RenderArgs["to_zt_stock_disabled"] = ""
	} else {
		c.RenderArgs["to_zt_stock_disabled"] = `disabled="disabled"`
	}

	if can_to_au_stock {
		c.RenderArgs["to_au_stock_disabled"] = ""
	} else {
		c.RenderArgs["to_au_stock_disabled"] = `disabled="disabled"`
	}

	return c.RenderTemplate("erp/hgbg/container_transfer/container_transfer_index.html")
}

func (c ErpContainerTransfer) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	type container_transfer_view_item_t struct {
		models.ContainerTransferInfo `xorm:"extends" json:"ct"`
		models.StatusInfo            `xorm:"extends" json:"s"`
	}

	query := db_session.
	Select("ct.*, s.*").
		Table("erp_container_transfer").Alias("ct").
		Join("inner", []string{"erp_status", "s"}, "ct.status_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	container_transfer_list := make([]container_transfer_view_item_t, 0, limit)
	err := data_query.Limit(limit, offset).Find(&container_transfer_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(container_transfer_view_item_t))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  container_transfer_list,
		Total: count,
	})
}

type ContainerTransferDetailForm struct {
	ContainerTransfer models.ContainerTransferInfo              `json:"transfer"`
	StartPortList     []models.ContainerTransferStartPortInfo   `json:"start_port_list"`
	DestPortList      []models.ContainerTransferDestPortInfo    `json:"dest_port_list"`
}

func (f *ContainerTransferDetailForm) IsCreate() bool {
	return f.ContainerTransfer.Id == 0
}

func (f *ContainerTransferDetailForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.ContainerTransfer.BatchCode).Message("批次号不能为空！")

	validation.Required(f.ContainerTransfer.StartPortId).Message("起始港口不能为空！")
	validation.Required(f.ContainerTransfer.DestPortId).Message("目的港口不能为空！")

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

	validation.Required(f.ContainerTransfer.PriceUpRate).Message("非征品上调不能为空！")

	validation.Required(f.ContainerTransfer.PriceDownRate).Message("征税品下调不能为空！")

	validation.Required(f.ContainerTransfer.CostUpRate).Message("成本上调不能为空！")

	validation.Required(f.ContainerTransfer.PackingTime).Message("装箱时间不能为空！")

	validation.Required(f.ContainerTransfer.DepartureTime).Message("到港时间不能为空！")

	validation.Required(f.ContainerTransfer.StockTime).Message("到仓时间不能为空！")

	return validation.HasErrors() == false
}

func (c ErpContainerTransfer) Detail(id int64) revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	var db_session = c.DbSession

	var transfer_id = id
	var container_transfer models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&container_transfer)
	core.HandleError(err)

	if container_transfer.Id == 0 {
		container_transfer.StatusId = models.CONTAINER_TRANSFER_STATUS_NEW
	}

	var status models.StatusInfo
	_, err = db_session.Id(container_transfer.StatusId).Get(&status)
	core.HandleError(err)

	c.RenderArgs["transfer_id"] = transfer_id
	c.RenderArgs["status_name"] = status.Name
	c.RenderArgs["can_modify"] = container_transfer_can_modify(container_transfer)
	c.RenderArgs["can_depart"] = container_transfer_can_depart(container_transfer)
	c.RenderArgs["can_stock"] = container_transfer_can_stock(container_transfer)
	c.RenderArgs["is_not_new"] = transfer_id != 0

	if id == 0 {
		c.RenderArgs["entries_display"] = "none"
	} else {
		c.RenderArgs["entries_display"] = "block"
	}

	return c.RenderTemplate("erp/hgbg/container_transfer/container_transfer_detail.html")
}

func (c ErpContainerTransfer) DetailData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

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
	} else {
		// 填入上次的
		var last_container_transfer models.ContainerTransferInfo
		_, err := db_session.Sql("SELECT ct.* FROM erp_container_transfer ct ORDER BY ct.id DESC LIMIT 1").Get(&last_container_transfer)
		core.HandleError(err)

		revel.TRACE.Println("last_container_transfer.Id", last_container_transfer.Id)

		container_transfer.StatusId = models.CONTAINER_TRANSFER_STATUS_NEW
		container_transfer.PriceDownRate = last_container_transfer.PriceDownRate
		container_transfer.CostUpRate = last_container_transfer.CostUpRate

		revel.TRACE.Println("container_transfer.PriceDownRate", container_transfer.PriceDownRate)
		revel.TRACE.Println("container_transfer.CostUpRate", container_transfer.CostUpRate)
	}

	var start_port_list = make([]models.ContainerTransferStartPortInfo, 0)
	core.HandleError(db_session.Find(&start_port_list))

	var dest_port_list = make([]models.ContainerTransferDestPortInfo, 0)
	core.HandleError(db_session.Find(&dest_port_list))

	var detail_form ContainerTransferDetailForm
	detail_form.ContainerTransfer = container_transfer
	detail_form.StartPortList = start_port_list
	detail_form.DestPortList = dest_port_list

	return c.RenderJson(detail_form)
}

func (c ErpContainerTransfer) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	var detail_form ContainerTransferDetailForm
	err := json.Unmarshal(c.GetRequestBody(), &detail_form)
	core.HandleError(err)

	if container_transfer_can_modify(detail_form.ContainerTransfer) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: "当前状态不能修改！"})
	}

	if detail_form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	var packing_time = time.Time(detail_form.ContainerTransfer.PackingTime)
	var departure_time = time.Time(detail_form.ContainerTransfer.DepartureTime)
	var stock_time = time.Time(detail_form.ContainerTransfer.StockTime)

	if packing_time.After(departure_time) {
		return c.RenderJson(core.JsonResult{Success: false, Message: "装箱时间 不能晚于 预计到港时间！"})
	}

	if departure_time.After(stock_time) {
		return c.RenderJson(core.JsonResult{Success: false, Message: "预计到港时间 不能晚于 预计到仓时间！"})
	}

	container_transfer := detail_form.ContainerTransfer

	var start_port models.ContainerTransferStartPortInfo
	_, err = db_session.Id(container_transfer.StartPortId).Get(&start_port)
	core.HandleError(err)
	container_transfer.StartPortName = start_port.PortName

	var dest_port models.ContainerTransferDestPortInfo
	_, err = db_session.Id(container_transfer.DestPortId).Get(&dest_port)
	core.HandleError(err)
	container_transfer.DestPortName = dest_port.PortName

	var is_create = detail_form.IsCreate()

	var affected int64
	if is_create == true {
		container_transfer.StatusId = models.CONTAINER_TRANSFER_STATUS_NEW
		affected, err = db_session.Insert(&container_transfer)
		core.HandleError(err)

		container_transfer_log_service.Insert_container_transfer_log(container_transfer, "新建调拨单", db_session, c.Session)
	} else {
		affected, err = db_session.Id(container_transfer.Id).Update(&container_transfer)
		core.HandleError(err)

		container_transfer_log_service.Insert_container_transfer_log(container_transfer, "修改调拨单", db_session, c.Session)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！" })
		}
	}

	container_transfer_service.CalculateAndSaveWeightVolumeAmount(container_transfer.Id, db_session)

	return c.RenderJson(map[string]interface{}{
		"success":     true,
		"is_create":   is_create,
		"transfer_id": container_transfer.Id,
		"message":     "数据保存成功",
	})
}

func (c ErpContainerTransfer) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER, c.Session)

	db_session := c.DbSession

	container_transfer_id_list := make([]int64, 0)
	c.Params.Bind(&container_transfer_id_list, "id_list")

	var container_transfer_list = make([]models.ContainerTransferInfo, 0)
	core.HandleError(db_session.In("id", container_transfer_id_list).Find(&container_transfer_list))

	for _, container_transfer := range container_transfer_list {
		if container_transfer_can_modify(container_transfer) == false {
			return c.RenderJson(core.JsonResult{Success: false, Message: "当前状态不能删除！"})
		}
	}

	for _, container_transfer := range container_transfer_list {
		container_transfer_log_service.Insert_container_transfer_log(container_transfer, "删除调拨单", db_session, c.Session)
	}

	affected, err := db_session.In("id", container_transfer_id_list).Delete(new(models.ContainerTransferInfo))
	core.HandleError(err)

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

func (c ErpContainerTransfer) Depart() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER_TO_ZT_STOCK, c.Session)

	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var ct models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&ct)
	core.HandleError(err)

	if ct.StatusId == models.CONTAINER_TRANSFER_STATUS_DEPART {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "当前集装箱调拨已经为已出发状态！",
		})
	}

	if container_transfer_can_depart(ct) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "当前集装箱调拨不能设置为已出发状态！",
		})
	}

	ct.StatusId = models.CONTAINER_TRANSFER_STATUS_DEPART

	_, err = db_session.Id(transfer_id).Cols("status_id").Update(&ct)
	core.HandleError(err)

	container_transfer_service.UpdateStockInventory(ct, c.Session, db_session)

	container_transfer_log_service.Insert_container_transfer_log(ct, "调拨至在途仓", db_session, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "设置集装箱调拨为 已出发状态 成功！",
	})
}

func (c ErpContainerTransfer) Stock() revel.Result {
	pass_service.Check_login_user_pass(pass_res.HGBG_CONTAINER_TRANSFER_TO_AU_STOCK, c.Session)

	db_session := c.DbSession

	var transfer_id int64
	c.Params.Bind(&transfer_id, "transfer_id")

	var ct models.ContainerTransferInfo
	_, err := db_session.Id(transfer_id).Get(&ct)
	core.HandleError(err)

	if ct.StatusId == models.CONTAINER_TRANSFER_STATUS_STOCK {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "当前集装箱调拨已经为已到仓状态！",
		})
	}

	if container_transfer_can_stock(ct) == false {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": "当前集装箱调拨不能设置为已到仓状态！",
		})
	}

	ct.StatusId = models.CONTAINER_TRANSFER_STATUS_STOCK

	_, err = db_session.Id(transfer_id).Cols("status_id").Update(&ct)
	core.HandleError(err)

	container_transfer_service.UpdateStockInventory(ct, c.Session, db_session)

	container_transfer_log_service.Insert_container_transfer_log(ct, "调拨至澳洲仓", db_session, c.Session)

	return c.RenderJson(map[string]interface{}{
		"success": true,
		"message": "设置集装箱调拨为 已到仓状态 成功！",
	})
}

func container_transfer_can_modify(container_transfer models.ContainerTransferInfo) bool {
	return container_transfer.StatusId == models.CONTAINER_TRANSFER_STATUS_NEW
}

func container_transfer_can_depart(container_transfer models.ContainerTransferInfo) bool {
	return container_transfer.StatusId == models.CONTAINER_TRANSFER_STATUS_NEW
}

func container_transfer_can_stock(container_transfer models.ContainerTransferInfo) bool {
	return container_transfer.StatusId == models.CONTAINER_TRANSFER_STATUS_DEPART
}
