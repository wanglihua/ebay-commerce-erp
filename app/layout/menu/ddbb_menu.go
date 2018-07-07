package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/core"
	"matrix/modules/auth/pass_res"
)

func GetDdbbMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="ddbb-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-area-chart"></i>
        <span class="menu-text">订单报表</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.DDBB_LOGISTICS_REPORT) {
		buffer.WriteString(`

        <li id="ddbb-logistics-report-menu" class="">
            <a href="/logistics/report">
                <i class="menu-icon fa fa-caret-right"></i>
                物流报表
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
