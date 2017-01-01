package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetXtpzMenu(db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xtpz-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-cogs"></i>
        <span class="menu-text">系统配置</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="xtpz-shopaccount-menu" class="">
            <a href="/erp/shop/account">
                <i class="menu-icon fa fa-caret-right"></i>
                店铺账号
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xtpz-paypalaccount-menu" class="">
            <a href="/erp/paypal/account">
                <i class="menu-icon fa fa-caret-right"></i>
                Paypal账号
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xtpz-exchangerate-menu" class="">
            <a href="/erp/exchange/rate">
                <i class="menu-icon fa fa-caret-right"></i>
                汇率设置
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xtpz-eubgoodstype-menu" class="">
            <a href="/erp/eub/goods/type">
                <i class="menu-icon fa fa-caret-right"></i>
                EUB物品类型
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xtpz-eubinsurancetype-menu" class="">
            <a href="/erp/eub/insurance/type">
                <i class="menu-icon fa fa-caret-right"></i>
                EUB保险类型
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xtpz-loginlog-menu" class="">
            <a href="/erp/login/log">
                <i class="menu-icon fa fa-caret-right"></i>
                系统日志
            </a>
            <b class="arrow"></b>
        </li>
    </ul>
</li>
`)

	return buffer.String()
}
