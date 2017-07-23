package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetHgbgMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="hgbg-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-university"></i>
        <span class="menu-text">海关报关</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.HGBG_PRODUCT_HS_CODE) {
		buffer.WriteString(`<li id="hgbg-product-hs-code-menu" class="">
            <a href="/product/hs/code">
                <i class="menu-icon fa fa-caret-right"></i>
                商品海关编码
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.HGBG_HS_CODE_CN_REPORT_ELEM) {
		buffer.WriteString(`<li id="hgbg-hs-code-cn-report-elem-menu" class="">
            <a href="/hs/code/cn">
                <i class="menu-icon fa fa-caret-right"></i>
                中国海关编码
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.HGBG_HS_CODE_AU_DUTY) {
		buffer.WriteString(`<li id="hgbg-hs-code-au-duty-menu" class="">
            <a href="/hs/code/au">
                <i class="menu-icon fa fa-caret-right"></i>
                澳洲海关编码
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.HGBG_CONTAINER_TRANSFER) {
		buffer.WriteString(`<li id="hgbg-container-transfer-menu" class="">
            <a href="/container/transfer">
                <i class="menu-icon fa fa-caret-right"></i>
                集装箱调拨
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.HGBG_CONTAINER_TRANSFER) {
		buffer.WriteString(`<li id="hgbg-product-invoice-menu" class="">
            <a href="/product/invoice">
                <i class="menu-icon fa fa-caret-right"></i>
                商品是否开票
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	/*
	buffer.WriteString(`<li id="hgbg-container-transfer-entry-menu" class="">
            <a href="/container/transfer/entry">
                <i class="menu-icon fa fa-caret-right"></i>
                集装箱调拨明细
            </a>
            <b class="arrow"></b>
        </li>`)
	 */

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
