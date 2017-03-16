package controllers

import (
	"github.com/revel/revel"
)

type Cpa_store struct {
	*revel.Controller
}

func (c Cpa_store) Index() revel.Result {
	return c.RenderTemplate("cpa/store/store_index.html")
}
