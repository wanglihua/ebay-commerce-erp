package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetDbglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="dbgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-truck"></i>
        <span class="menu-text">调拨管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="dbgl-aaaa-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                调拨清单
            </a>
            <b class="arrow"></b>
        </li>

        <li id="dbgl-cccc-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                异常调拨单
            </a>
            <b class="arrow"></b>
        </li>
    </ul>
</li>
`)

	return buffer.String()
}
