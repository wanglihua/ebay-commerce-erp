package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetThglMenu(res_id_map map[int64]bool, db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="thgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-rotate-left"></i>
        <span class="menu-text">退货管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="thgl-order-return-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                退货列表
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
