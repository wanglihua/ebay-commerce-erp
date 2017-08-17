package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetSpzlMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	/*
	login_user := core.GetLoginUser(web_session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)

	var check_pass_res_menu = func(pass_res_id int64) bool {
		if v, ok := res_id_map[pass_res_id]; is_admin == true || (v == true && ok == true) {
			return true
		}
		return false
	}
	*/

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="spzl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-reply"></i>
        <span class="menu-text">商品质量反馈</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	buffer.WriteString(`<li id="spzl-product-feedback-menu" class="">
            <a href="/product/feedback">
                <i class="menu-icon fa fa-caret-right"></i>
                商品质量反馈
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`<li id="spzl-product-feedback-cate-menu" class="">
            <a href="/product/feedback/cate">
                <i class="menu-icon fa fa-caret-right"></i>
                商品质量反馈类别
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
