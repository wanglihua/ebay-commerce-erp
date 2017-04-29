package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetDdbbMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="ddbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-area-chart"></i>
        <span class="menu-text">订单报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="ddbb-logistics-dialy-report-menu" class="">
            <a href="/erp/logistics/dialy/report">
                <i class="menu-icon fa fa-caret-right"></i>
                日物流报表
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
