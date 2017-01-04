package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetXsglMenu(db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xsgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-car"></i>
        <span class="menu-text">销售管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="xsgl-productimagetype-menu" class="">
            <a href="/erp/product/image/type">
                <i class="menu-icon fa fa-caret-right"></i>
                商品图片属性
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
