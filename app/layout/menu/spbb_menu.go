package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/auth/pass_res"
)

func GetSpbbMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="spbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-bar-chart-o"></i>
        <span class="menu-text">商品报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.SPBB_PRODUCT_SALE_RANK) {
		buffer.WriteString(`<li id="spbb-product-sale-rank-menu" class="">
            <a href="/product/sale/rank">
                <i class="menu-icon fa fa-caret-right"></i>
                商品销售排名
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPBB_PRODUCT_SALE_ANALYSIS) {
		buffer.WriteString(`<li id="spbb-product-sale-analysis-menu" class="">
            <a href="/product/sale/analysis">
                <i class="menu-icon fa fa-caret-right"></i>
                商品销售统计
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPBB_PRODUCT_SALE_DAY_ANALYSIS) {
		buffer.WriteString(`<li id="spbb-product-sale-day-analysis-menu" class="">
            <a href="/product/sale/day/analysis">
                <i class="menu-icon fa fa-caret-right"></i>
                商品销售日统计
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPBB_PRODUCT_SALE_ANALYSIS_SETTING) {
		buffer.WriteString(`<li id="spbb-product-sale-analysis-setting-menu" class="">
            <a href="/product/sale/analysis/setting">
                <i class="menu-icon fa fa-caret-right"></i>
                商品销售统计设置
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
