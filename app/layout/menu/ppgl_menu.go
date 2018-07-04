package menu

import (
	"bytes"
	"matrix/core"
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

func GetPpglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="ppgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-cc-paypal"></i>
        <span class="menu-text">付款退款</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.PPGL_ORDER_REFUND_APPROVE) {
		buffer.WriteString(`<li id="ppgl-refund-menu" class="">
            <a href="/refund">
                <i class="menu-icon fa fa-caret-right"></i>
                订单退款
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.PPGL_PAYPAL_PAYMENT) {
		buffer.WriteString(`<li id="ppgl-paypal-payment-menu" class="">
            <a href="/paypal/payment">
                <i class="menu-icon fa fa-caret-right"></i>
                Paypal付款
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`<li id="ppgl-refund-report-menu" class="">
            <a href="/refund/report">
                <i class="menu-icon fa fa-caret-right"></i>
                退款报表
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
