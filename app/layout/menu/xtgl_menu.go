package menu

import (
	"bytes"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetXtglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	login_user := core.GetLoginUser(web_session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xtgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-legal"></i>
        <span class="menu-text">系统管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	var check_pass_res_menu = func(pass_res_id int64) bool {
		if v, ok := res_id_map[pass_res_id]; is_admin == true || (v == true && ok == true) {
			return true
		}
		return false
	}

	if check_pass_res_menu(pass_res.XTGL_NOTICE) == true {
		buffer.WriteString(`
        <li id="xtgl-notice-menu" class="">
            <a href="/notice">
                <i class="menu-icon fa fa-caret-right"></i>
                通知公告
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
   <li id="xtpz-eubgoodstype-menu" class="">
       <a href="/eub/goods/type">
           <i class="menu-icon fa fa-caret-right"></i>
           EUB物品类型
       </a>
       <b class="arrow"></b>
   </li>

   <li id="xtpz-eubinsurancetype-menu" class="">
       <a href="/eub/insurance/type">
           <i class="menu-icon fa fa-caret-right"></i>
           EUB保险类型
       </a>
       <b class="arrow"></b>
   </li>

*/
