package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetDdlcMenu(res_id_map map[int64]bool, db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
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
    <ul class="submenu">

        <li id="ddlc-order-all-menu" class="">
            <a href="/erp/order/all">
                <i class="menu-icon fa fa-caret-right"></i>
                所有订单
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddlc-order-yzf-menu" class="">
            <a href="/erp/order/yzf">
                <i class="menu-icon fa fa-caret-right"></i>
                已支付
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddlc-order-hold-menu" class="">
            <a href="/erp/order/hold">
                <i class="menu-icon fa fa-caret-right"></i>
                Hold
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddlc-order-dph-menu" class="">
            <a href="/erp/order/dph">
                <i class="menu-icon fa fa-caret-right"></i>
                待配货
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddlc-order-phz-menu" class="">
            <a href="/erp/order/phz">
                <i class="menu-icon fa fa-caret-right"></i>
                配货中
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddlc-order-yfh-menu" class="">
            <a href="/erp/order/yfh">
                <i class="menu-icon fa fa-caret-right"></i>
                已发货
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddlc-order-ywc-menu" class="">
            <a href="/erp/order/ywc">
                <i class="menu-icon fa fa-caret-right"></i>
                已完成
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddlc-order-khb-menu" class="">
            <a href="/erp/order/khb">
                <i class="menu-icon fa fa-caret-right"></i>
                可合并的订单
            </a>
            <b class="arrow"></b>
        </li>
    </ul>
</li>
`)

	return buffer.String()
}
