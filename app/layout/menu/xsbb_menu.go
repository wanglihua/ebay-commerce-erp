package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/core"
	"matrix/modules/auth/pass_res"
)

func GetXsbbMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="xsbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-line-chart"></i>
        <span class="menu-text">销售报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.XSBB_SALE_INCOME_REPORT_DAY) {
		buffer.WriteString(`
        <li id="xsbb-sale-income-report-day-menu" class="">
			<a href="/sale/income/report/day">
				<i class="menu-icon fa fa-caret-right"></i>
				日收入报表
			</a>
			<b class="arrow"></b>
		</li>`)
	}

	if check_pass_res_menu(pass_res.XSBB_SALE_AMOUNT_REPORT_DAY) {
		buffer.WriteString(`
		<li id="xsbb-sale-amount-report-day-menu" class="">
			<a href="/sale/amount/report/day">
				<i class="menu-icon fa fa-caret-right"></i>
				日销售额报表
			</a>
			<b class="arrow"></b>
		</li>`)
	}

	if check_pass_res_menu(pass_res.XSBB_SALE_AMOUNT_REPORT_MONTH) {
		buffer.WriteString(`
		<li id="xsbb-sale-amount-report-month-menu" class="">
			<a href="/sale/amount/report/month">
				<i class="menu-icon fa fa-caret-right"></i>
				月销售额报表
			</a>
			<b class="arrow"></b>
		</li>`)
	}

	if check_pass_res_menu(pass_res.XSBB_SALE_AMOUNT_REPORT_PRODUCT) {
		buffer.WriteString(`
        <li id="xsbb-hot-sales-product-rank-menu" class="">
            <a href="/hot/sales/product/rank">
                <i class="menu-icon fa fa-caret-right"></i>
                热销商品排名
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`
    </ul>
</li>
`)

	return buffer.String()
}

/*

<li id="xsbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-line-chart"></i>
        <span class="menu-text">销售报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">
		<li id="xsbb-sale-income-report-day-menu" class="">
			<a href="/sale/income/report/day">
				<i class="menu-icon fa fa-caret-right"></i>
				日收入报表
			</a>
			<b class="arrow"></b>
		</li>
		<li id="xsbb-sale-amount-report-day-menu" class="">
			<a href="/sale/amount/report/day">
				<i class="menu-icon fa fa-caret-right"></i>
				日销售额报表
			</a>
			<b class="arrow"></b>
		</li>

		<li id="xsbb-sale-amount-report-month-menu" class="">
			<a href="/sale/amount/report/month">
				<i class="menu-icon fa fa-caret-right"></i>
				月销售额报表
			</a>
			<b class="arrow"></b>
		</li>

        <li id="xsbb-hot-sales-product-rank-menu" class="">
            <a href="/hot/sales/product/rank">
                <i class="menu-icon fa fa-caret-right"></i>
                热销商品排名
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-hot-sales-store-rank-menu" class="">
            <a href="/hot/sales/store/rank">
                <i class="menu-icon fa fa-caret-right"></i>
                热销店铺排名
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-area-sales-amount-report-menu" class="">
            <a href="/area/sales/amount/report">
                <i class="menu-icon fa fa-caret-right"></i>
                按地区销售额报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-product-feed-report-menu" class="">
            <a href="/product/feed/report">
                <i class="menu-icon fa fa-caret-right"></i>
                商品反馈报表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-single-product-sales-query-menu" class="">
            <a href="/single/product/sales/query">
                <i class="menu-icon fa fa-caret-right"></i>
                单商品销售查询
            </a>
            <b class="arrow"></b>
        </li>

        <li id="xsbb-person-sales-report-menu" class="">
            <a href="/person/sales/report">
                <i class="menu-icon fa fa-caret-right"></i>
                人员营业报表
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>

 */
