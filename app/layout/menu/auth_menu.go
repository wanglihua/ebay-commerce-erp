package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetAuthMenu(res_id_map map[int64]bool, db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("auth"); found == false {
		return ""
	}

	login_user := core.GetLoginUser(web_session)
	is_admin := false
	if login_user.UserId != 0 {
		is_admin = pass_service.Check_user_is_admin(login_user.UserId)
	}

	var check_pass_res_menu = func(pass_res_id int64) bool {
		if v, ok := res_id_map[pass_res_id]; is_admin == true || (v == true && ok == true) {
			return true
		}
		return false
	}

	var buffer bytes.Buffer

	if is_admin {
		buffer.WriteString(`
<li id="auth-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-users"></i>
        <span class="menu-text">权限管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

		if check_pass_res_menu(pass_res.AUTH_USER) {
			buffer.WriteString(`<li id="auth-user-menu" class="">
            <a href="/auth/user">
                <i class="menu-icon fa fa-caret-right"></i>
                用户管理
            </a>
            <b class="arrow"></b>
        </li>`)
		}

		if check_pass_res_menu(pass_res.AUTH_GROUP) {
			buffer.WriteString(`<li id="auth-group-menu" class="">
            <a href="/auth/group">
                <i class="menu-icon fa fa-caret-right"></i>
                群组管理
            </a>
            <b class="arrow"></b>
        </li>`)
		}

		if login_user.IsSuperAdmin == true {
			buffer.WriteString(`<li id="auth-admin-menu" class="">
            <a href="/auth/admin">
                <i class="menu-icon fa fa-caret-right"></i>
                管理员组
            </a>
            <b class="arrow"></b>
        </li>`)
		}

		buffer.WriteString(`</ul>
</li>
`)
	}
	return buffer.String()
}

/*

<li id="auth-passres-menu" class="">
    <a href="/auth/user/pass/res">
        <i class="menu-icon fa fa-caret-right"></i>
        用户权限
    </a>
    <b class="arrow"></b>
</li>

<li id="auth-group-menu" class="">
    <a href="/auth/group">
        <i class="menu-icon fa fa-caret-right"></i>
        群组管理
    </a>
    <b class="arrow"></b>
</li>

*/
