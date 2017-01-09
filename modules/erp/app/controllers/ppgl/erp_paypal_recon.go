package controllers

import (
	"github.com/revel/revel"

	"matrix/core"
	"matrix/modules/erp/models"
)

type ErpPaypalRecon struct {
	*revel.Controller
}

func (c ErpPaypalRecon) Index() revel.Result {
	return c.RenderTemplate("erp/ppgl/paypal_recon/paypal_recon_index.html")
}

type ErpPaypalReconViewItem struct {
	models.PaypalReconInfo  `xorm:"extends" json:"r"`
	models.StatusInfo       `xorm:"extends" json:"s"`
}

func (c ErpPaypalRecon) ListData() revel.Result {
	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)

	/*
		SELECT r.*, s.*
		FROM erp_paypal_recon r
		    LEFT JOIN erp_status s ON r.result_id = s.id
	 */

	query := db_session.
	Select("r.*, s.*").
		Table("erp_paypal_recon").Alias("r").
		Join("LEFT", []string{"erp_status", "s"}, "r.result_id = s.id").
		Where(filter)

	//query extra filter here

	data_query := *query
	if order != "" {
		data_query = *data_query.OrderBy(order)
	} else {
		data_query = *data_query.Asc("id")
	}

	paypal_recon_list := make([]ErpPaypalReconViewItem, 0, limit)
	err := data_query.Limit(limit, offset).Find(&paypal_recon_list)
	core.HandleError(err)

	count_query := *query
	count, err := count_query.Count(new(ErpPaypalReconViewItem))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  paypal_recon_list,
		Total: count,
	})
}
