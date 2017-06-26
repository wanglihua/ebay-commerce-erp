package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetDdpzMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="ddpz-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-gear"></i>
        <span class="menu-text">订单配置</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.DDPZ_LOGISTICS_COMPANY) {
		buffer.WriteString(`<li id="ddpz-logisticscompany-menu" class="">
            <a href="/logistics/company">
                <i class="menu-icon fa fa-caret-right"></i>
                物流公司
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDPZ_LOGISTICS_MODE) {
		buffer.WriteString(`<li id="ddpz-logisticsmode-menu" class="">
            <a href="/logistics/mode">
                <i class="menu-icon fa fa-caret-right"></i>
                物流方式
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDPZ_LOGISTICS_FEE_RATE) {
		buffer.WriteString(`<li id="ddpz-logisticsfeerate-menu" class="">
            <a href="/logistics/fee/rate">
                <i class="menu-icon fa fa-caret-right"></i>
                物流费率
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDPZ_LOGISTICS_SETTING) {
		buffer.WriteString(`<li id="ddpz-logisticssetting-menu" class="">
            <a href="/logistics/setting">
                <i class="menu-icon fa fa-caret-right"></i>
                物流配置
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDPZ_ENVELOPE_FEE_RATE) {
		buffer.WriteString(`<li id="ddpz-envelopefeerate-menu" class="">
            <a href="/envelope/fee/rate">
                <i class="menu-icon fa fa-caret-right"></i>
                信封费率
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDPZ_PICK_FEE) {
		buffer.WriteString(`<li id="ddpz-pickfee-menu" class="">
             <a href="/pick/fee">
                 <i class="menu-icon fa fa-caret-right"></i>
                 拣货费用
             </a>
             <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDPZ_ORDER_HOLD_TYPE) {
		buffer.WriteString(`<li id="ddpz-orderholdtype-menu" class="">
            <a href="/order/hold/type">
                <i class="menu-icon fa fa-caret-right"></i>
                订单Hold类型
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.DDPZ_ORDER_TAG) {
		buffer.WriteString(`<li id="ddpz-ordertag-menu" class="">
            <a href="/erp/order/tag">
                <i class="menu-icon fa fa-caret-right"></i>
                订单标签
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	/*
                    <li id="ddpz-ordertag-menu" class="">
                        <a href="/erp/order/tag">
                            <i class="menu-icon fa fa-caret-right"></i>
                            订单标签
                        </a>
                        <b class="arrow"></b>
                    </li>
	 */

	/*
	buffer.WriteString(`<li id="ddpz-platformlogisticsmode-menu" class="">
            <a href="/platform/logistics/mode">
                <i class="menu-icon fa fa-caret-right"></i>
                电商平台物流方式
            </a>
            <b class="arrow"></b>
        </li>`)


	buffer.WriteString(`<li id="ddpz-ordercustomtype-menu" class="">
            <a href="/order/custom/type">
                <i class="menu-icon fa fa-caret-right"></i>
                订单自定义类型
            </a>
            <b class="arrow"></b>
        </li>`)
    */

	if check_pass_res_menu(pass_res.DDPZ_ORDER_LOG) {
		buffer.WriteString(`<li id="ddpz-orderlog-menu" class="">
            <a href="/order/log">
                <i class="menu-icon fa fa-caret-right"></i>
                订单修改日志
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
   <li id="ddpz-logisticsnumberbatch-menu" class="">
       <a href="/logistics/number/batch">
           <i class="menu-icon fa fa-caret-right"></i>
           订单运单号
       </a>
       <b class="arrow"></b>
   </li>

   <li id="ddpz-logisticsnumber-menu" class="">
       <a href="/logistics/number">
           <i class="menu-icon fa fa-caret-right"></i>
           运单号不规则号段
       </a>
       <b class="arrow"></b>
   </li>

   <li id="ddpz-logisticsnumberlog-menu" class="">
       <a href="/logistics/number/log">
           <i class="menu-icon fa fa-caret-right"></i>
           运单号操作日志
       </a>
       <b class="arrow"></b>
   </li>
*/
