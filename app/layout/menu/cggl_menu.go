package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetCgglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="cggl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-shopping-cart"></i>
        <span class="menu-text">采购管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="cggl-purchasebill-menu" class="">
            <a href="/erp/purchase/bill">
                <i class="menu-icon fa fa-caret-right"></i>
                采购管理
            </a>
            <b class="arrow"></b>
        </li>

        <li id="cggl-bbbb-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                入库管理
            </a>
            <b class="arrow"></b>
        </li>

        <li id="cggl-cccc-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                异常采购单
            </a>
            <b class="arrow"></b>
        </li>
    </ul>
</li>
`)

	return buffer.String()
}
