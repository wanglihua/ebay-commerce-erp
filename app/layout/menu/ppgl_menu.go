package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetPpglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="ppgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-cc-paypal"></i>
        <span class="menu-text">PayPal管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="ppgl-paypal-payment-menu" class="">
            <a href="/paypal/payment">
                <i class="menu-icon fa fa-caret-right"></i>
                Paypal付款
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
