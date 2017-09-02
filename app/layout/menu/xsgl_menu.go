package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

func GetXsglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	login_user := core.GetLoginUser(web_session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)

	var check_pass_res_menu = func(pass_res_id int64) bool {
		if v, ok := res_id_map[pass_res_id]; is_admin == true || (v == true && ok == true) {
			return true
		}
		return false
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xsgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-tags"></i>
        <span class="menu-text">销售管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.XSGL_PRODUCT_SALE_SETTING) {
		buffer.WriteString(`<li id="xsgl-product-sale-setting-menu" class="">
            <a href="/product/sale/setting">
                <i class="menu-icon fa fa-caret-right"></i>
                商品销售设置
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`<li id="xsgl-salesmanager-menu" class="">
            <a href="/sales/manager">
                <i class="menu-icon fa fa-caret-right"></i>
                销售负责人
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="xsgl-xmgl-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                项目管理
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="xsgl-sku-xsbg-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                SKU销售报告
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="xsgl-cpgjz-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                产品关键字
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="xsgl-productimage-menu" class="">
            <a href="/product/image">
                <i class="menu-icon fa fa-caret-right"></i>
                商品图片
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="xsgl-productimagetype-menu" class="">
            <a href="/product/image/type">
                <i class="menu-icon fa fa-caret-right"></i>
                商品图片属性
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
