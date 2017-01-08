package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetPpglMenu(db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="cggl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-car"></i>
        <span class="menu-text">PayPal管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="cggl-aaaa-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                PayPal对账
            </a>
            <b class="arrow"></b>
        </li>

        <li id="cggl-bbbb-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                退款管理
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
