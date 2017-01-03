package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetDdpzMenu(db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="ddpz-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-car"></i>
        <span class="menu-text">订单配置</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="ddpz-logisticsmode-menu" class="">
            <a href="/erp/logistics/mode">
                <i class="menu-icon fa fa-caret-right"></i>
                物流方式
            </a>
            <b class="arrow"></b>
        </li>

        <li id="ddpz-orderlog-menu" class="">
            <a href="/erp/order/log">
                <i class="menu-icon fa fa-caret-right"></i>
                订单修改日志
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
