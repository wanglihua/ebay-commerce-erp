package menu

import (
	"bytes"

	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetXtpzMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	login_user := core.GetLoginUser(web_session)
	var is_admin = pass_service.Is_user_admin(login_user.UserId)

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="xtpz-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-cogs"></i>
        <span class="menu-text">系统配置</span>
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

	if check_pass_res_menu(pass_res.XTPZ_ECOMMERCE_PLATFORM) == true {
		buffer.WriteString(`
        <li id="xtpz-ecommerceplatform-menu" class="">
            <a href="/ecommerce/platform">
                <i class="menu-icon fa fa-caret-right"></i>
                电商平台
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.XTPZ_SHOP_ACCOUNT) == true {
		buffer.WriteString(`
        <li id="xtpz-shopaccount-menu" class="">
            <a href="/shop/account">
                <i class="menu-icon fa fa-caret-right"></i>
                店铺账号
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.XTPZ_PAYPAL_ACCOUNT) == true {
		buffer.WriteString(`
        <li id="xtpz-paypalaccount-menu" class="">
            <a href="/paypal/account">
                <i class="menu-icon fa fa-caret-right"></i>
                PayPal账号
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.XTPZ_EXCHANGE_RATE) == true {
		buffer.WriteString(`
        <li id="xtpz-exchangerate-menu" class="">
            <a href="/exchange/rate">
                <i class="menu-icon fa fa-caret-right"></i>
                汇率设置
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.XTPZ_LOGIN_LOG) == true {
		buffer.WriteString(`
        <li id="xtpz-loginlog-menu" class="">
            <a href="/login/log">
                <i class="menu-icon fa fa-caret-right"></i>
                系统日志
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.XTPZ_SCHEDULE_TASK) == true {
		buffer.WriteString(`
        <li id="xtpz-schedule-task-menu" class="">
            <a href="/schedule/task">
                <i class="menu-icon fa fa-caret-right"></i>
                定时任务
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.XTPZ_LITTLE_ICON) == true {
		buffer.WriteString(`
        <li id="xtpz-little-icon-menu" class="">
            <a href="/little/icon">
                <i class="menu-icon fa fa-caret-right"></i>
                图标管理
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.XTPZ_EMAIL_ADDRESS) == true {
		buffer.WriteString(`
        <li id="xtpz-emailaddress-menu" class="">
            <a href="/email/address">
                <i class="menu-icon fa fa-caret-right"></i>
                电子邮件地址
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
