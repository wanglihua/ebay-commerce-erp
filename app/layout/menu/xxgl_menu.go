package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/modules/auth/service/pass_service"
	"matrix/core"
	"matrix/modules/auth/pass_res"
)

func GetXxglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="xxgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-envelope-o"></i>
        <span class="menu-text">消息管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.XXGL_MESSAGE_ALL) {
		buffer.WriteString(`
        <li id="xxgl-all-message-menu" class="">
            <a href="/message">
                <i class="menu-icon fa fa-caret-right"></i>
                所有消息
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`
        <li id="xxgl-dcl-message-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                待处理消息
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`
        <li id="xxgl-wfp-menu" class="">
            <a href="#">
                <i class="menu-icon fa fa-caret-right"></i>
                未分配消息
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`
    </ul>
</li>
`)

	return buffer.String()
}
