package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetXxglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xxgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-envelope-o"></i>
        <span class="menu-text">消息管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="xxgl-all-message-menu" class="">
            <a href="/erp/message">
                <i class="menu-icon fa fa-caret-right"></i>
                所有消息
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xxgl-dcl-message-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                待处理消息
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xxgl-wfp-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                未分配消息
            </a>
            <b class="arrow"></b>
        </li>
    </ul>
</li>
`)

	return buffer.String()
}
