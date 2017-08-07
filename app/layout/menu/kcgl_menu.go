package menu

import (
	"bytes"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetKcglMenu(res_id_map map[int64]bool, web_session revel.Session) string {
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
<li id="kcgl-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-building"></i>
        <span class="menu-text">库存管理</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.KCGL_PRODUCT) {
		buffer.WriteString(`<li id="kcgl-product-menu" class="">
            <a href="/product">
                <i class="menu-icon fa fa-caret-right"></i>
                商品列表
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.KCGL_PRODUCT_LWHW) {
		buffer.WriteString(`<li id="kcgl-product-weight-menu" class="">
            <a href="/product/weight">
                <i class="menu-icon fa fa-caret-right"></i>
                商品长宽高重量
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.KCGL_SPECIAL_SKU) {
		buffer.WriteString(`<li id="kcgl-special-product-menu" class="">
            <a href="/special/product">
                <i class="menu-icon fa fa-caret-right"></i>
                特殊SKU
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.KCGL_PRODUCT_BIND) {
		buffer.WriteString(`<li id="kcgl-productbind-menu" class="">
            <a href="/product/bind">
                <i class="menu-icon fa fa-caret-right"></i>
                绑定商品
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.KCGL_PRODUCT_IMPORT) || check_pass_res_menu(pass_res.KCGL_PRODUCT_EXPORT) {
		buffer.WriteString(`<li id="kcgl-product-import-menu" class="">
            <a href="/product/import">
                <i class="menu-icon fa fa-caret-right"></i>
                商品导入导出
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	/*
	if check_pass_res_menu(pass_res.KCGL_PRODUCT_EXPORT) {
		buffer.WriteString(`<li id="kcgl-product-export-menu" class="">
            <a href="/product/export">
                <i class="menu-icon fa fa-caret-right"></i>
                商品导出
            </a>
            <b class="arrow"></b>
        </li>`)
	}
	*/

	if check_pass_res_menu(pass_res.KCGL_INVENTORY) {
		buffer.WriteString(`<li id="kcgl-inventory-menu" class="">
            <a href="/inventory">
                <i class="menu-icon fa fa-caret-right"></i>
                商品库存
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.KCGL_INVENTORY_CHECK) {
		buffer.WriteString(`<li id="kcgl-inventorycheck-menu" class="">
            <a href="/inventory/check">
                <i class="menu-icon fa fa-caret-right"></i>
                库存盘点
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	/*
	if check_pass_res_menu(pass_res.KCGL_INVENTORY_CHECK_APPROVE) {
		buffer.WriteString(`<li id="kcgl-inventorycheck-approve-menu" class="">
            <a href="/inventory/check/approve">
                <i class="menu-icon fa fa-caret-right"></i>
                盘点审批
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.KCGL_STOCKBILL_TRANSFER) {
		buffer.WriteString(`<li id="kcgl-stockbill-transfer-menu" class="">
            <a href="/stock/bill/transfer">
                <i class="menu-icon fa fa-caret-right"></i>
                分仓调拨
            </a>
            <b class="arrow"></b>
        </li>`)
	}
	*/

	if check_pass_res_menu(pass_res.KCGL_STOCKBILL_CHECKIN) {
		buffer.WriteString(`<li id="kcgl-stockbill-checkin-menu" class="">
            <a href="/stock/bill/checkin">
                <i class="menu-icon fa fa-caret-right"></i>
                手工入库
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.KCGL_STOCKBILL_CHECKOUT) {
		buffer.WriteString(`<li id="kcgl-stockbill-checkout-menu" class="">
            <a href="/stock/bill/checkout">
                <i class="menu-icon fa fa-caret-right"></i>
                手工出库
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}
