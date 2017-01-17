package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetXsbbMenu(db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xsbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-line-chart"></i>
        <span class="menu-text">销售报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="xsbb-diary-income-report-menu" class="">
            <a href="/erp/diary/income/report">
                <i class="menu-icon fa fa-caret-right"></i>
                日收入报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-diary-sales-report-menu" class="">
            <a href="/erp/diary/sales/report">
                <i class="menu-icon fa fa-caret-right"></i>
                日销售报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-month-sales-amount-report-menu" class="">
            <a href="/erp/month/sales/amount/report">
                <i class="menu-icon fa fa-caret-right"></i>
                月销售额报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-hot-sales-product-rank-menu" class="">
            <a href="/erp/hot/sales/product/rank">
                <i class="menu-icon fa fa-caret-right"></i>
                热销商品排名
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-hot-sales-store-rank-menu" class="">
            <a href="/erp/hot/sales/store/rank">
                <i class="menu-icon fa fa-caret-right"></i>
                热销店铺排名
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-area-sales-amount-report-menu" class="">
            <a href="/erp/area/sales/amount/report">
                <i class="menu-icon fa fa-caret-right"></i>
                按地区销售额报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-product-feed-report-menu" class="">
            <a href="/erp/product/feed/report">
                <i class="menu-icon fa fa-caret-right"></i>
                商品反馈报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-single-product-sales-query-menu" class="">
            <a href="/erp/single/product/sales/query">
                <i class="menu-icon fa fa-caret-right"></i>
                单商品销售查询
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-person-sales-report-menu" class="">
            <a href="/erp/person/sales/report">
                <i class="menu-icon fa fa-caret-right"></i>
                人员营业报表
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
