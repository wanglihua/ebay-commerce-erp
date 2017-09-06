package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/core"
	"matrix/modules/auth/pass_res"
)

func GetCgglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var login_user = core.GetLoginUser(web_session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)

	var check_pass_res_menu = func(pass_res_id int64) bool {
		if v, ok := res_id_map[pass_res_id]; is_admin == true || (v == true && ok == true) {
			return true
		}
		return false
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
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.CGGL_PRODUCT_SETTING) == true {
		buffer.WriteString(`<li id="cggl-product-purchase-setting-menu" class="">
            <a href="/product/purchase/setting">
                <i class="menu-icon fa fa-caret-right"></i>
                商品采购设置
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	/*
		<li id="cggl-intel-repl-menu" class="">
			<a href="/erp/intel/repl">
				<i class="menu-icon fa fa-caret-right"></i>
				智能补货
			</a>
			<b class="arrow"></b>
		</li>
	 */

	buffer.WriteString(`<li id="cggl-intel-repl-menu" class="">
            <a href="/intel/repl">
                <i class="menu-icon fa fa-caret-right"></i>
                智能补货
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="cggl-purchase-bill-menu" class="">
            <a href="/purchase/bill">
                <i class="menu-icon fa fa-caret-right"></i>
                采购管理
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="cggl-purchase-finance-menu" class="">
            <a href="/purchase/finance">
                <i class="menu-icon fa fa-caret-right"></i>
                财务管理
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="cggl-purchase-checkin-menu" class="">
            <a href="/purchase/checkin">
                <i class="menu-icon fa fa-caret-right"></i>
                入库管理
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="cggl-cccc-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                异常采购单
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="cggl-dddd-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                采购流水
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`
    </ul>
</li>
`)

	return buffer.String()
}

/*

智能补货
采购管理
财务管理
入库管理
异常采购管理
采购流水

 */
