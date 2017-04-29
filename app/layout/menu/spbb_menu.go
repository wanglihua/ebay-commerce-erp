package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetSpbbMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="spbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-bar-chart-o"></i>
        <span class="menu-text">商品报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="spbb-brand-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                滞销商品报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="spbb-product-batch-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                热销商品排名
            </a>
            <b class="arrow"></b>
        </li>

        <li id="spbb-packagematerial-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                库存报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="spbb-stock-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                盘点查询
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
