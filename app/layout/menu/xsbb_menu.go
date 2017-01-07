package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetXsbbMenu(db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xsbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-car"></i>
        <span class="menu-text">销售报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="xsbb-a-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                日收入报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-b-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                日销售报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-c-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                月销售额报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-d-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                热销商品排名
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-e-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                热销店铺排名
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-f-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                按地区销售额报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-g-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                商品反馈报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-h-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                单商品销售查询
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-i-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                人员营业报表
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
