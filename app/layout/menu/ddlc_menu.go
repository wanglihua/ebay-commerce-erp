package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetDdlcMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="ddlc-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-random"></i>
        <span class="menu-text">订单流程</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.DDLC_ORDER_ALL) {
		buffer.WriteString(`<li id="ddlc-order-all-menu" class="">
            <a href="/erp/order/all">
                <i class="menu-icon fa fa-caret-right"></i>
                所有订单
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_YCJ) {
		buffer.WriteString(`<li id="ddlc-order-ycj-menu" class="">
            <a href="/erp/order/ycj">
                <i class="menu-icon fa fa-caret-right"></i>
                未支付
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_YZF) {
		buffer.WriteString(`<li id="ddlc-order-yzf-menu" class="">
            <a href="/erp/order/yzf">
                <i class="menu-icon fa fa-caret-right"></i>
                已支付
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_HOLD) {
		buffer.WriteString(`<li id="ddlc-order-hold-menu" class="">
            <a href="/erp/order/hold">
                <i class="menu-icon fa fa-caret-right"></i>
                Hold
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_DPH) {
		buffer.WriteString(`<li id="ddlc-order-dph-menu" class="">
            <a href="/erp/order/dph">
                <i class="menu-icon fa fa-caret-right"></i>
                待配货
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_PHZ) {
		buffer.WriteString(`<li id="ddlc-order-phz-menu" class="">
            <a href="/erp/order/phz">
                <i class="menu-icon fa fa-caret-right"></i>
                配货中
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_YFH) {
		buffer.WriteString(`<li id="ddlc-order-yfh-menu" class="">
            <a href="/erp/order/yfh">
                <i class="menu-icon fa fa-caret-right"></i>
                已发货
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_YWC) {
		buffer.WriteString(`<li id="ddlc-order-ywc-menu" class="">
            <a href="/erp/order/ywc">
                <i class="menu-icon fa fa-caret-right"></i>
                已完成
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_IMPORT) {
		buffer.WriteString(`<li id="ddlc-order-import-menu" class="">
            <a href="/erp/order/import">
                <i class="menu-icon fa fa-caret-right"></i>
                订单导入
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_LOGISTICS) {
		buffer.WriteString(`<li id="ddlc-orderlogistics-menu" class="">
            <a href="/erp/order/logistics">
                <i class="menu-icon fa fa-caret-right"></i>
                订单物流
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDLC_ORDER_LOGISTICS_FEE_DIFF) {
		buffer.WriteString(`<li id="ddlc-order-logistics-fee-diff-menu" class="">
            <a href="/erp/order/logistics/fee/diff">
                <i class="menu-icon fa fa-caret-right"></i>
                订单物流费用偏差
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
