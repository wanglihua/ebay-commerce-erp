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

        <li id="xsgl-xmgl-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                项目管理
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsgl-ebay-zhgl-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                EBAY账号管理
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsgl-sku-xsbg-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                SKU销售报告
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsgl-cpgjz-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                产品关键字
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsgl-productimage-menu" class="">
            <a href="/erp/product/image">
                <i class="menu-icon fa fa-caret-right"></i>
                商品图片
            </a>
            <b class="arrow"></b>
        </li>

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
