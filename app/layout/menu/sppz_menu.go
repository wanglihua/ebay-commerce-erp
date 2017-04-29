package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
)

func GetSppzMenu(res_id_map map[int64]bool, db_session *xorm.Session, web_session revel.Session) string {
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
<li id="sppz-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-gear"></i>
        <span class="menu-text">商品配置</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">`)

	if check_pass_res_menu(pass_res.SPPZ_PRODUCT_BATCH) {
		buffer.WriteString(`<li id="sppz-product-batch-menu" class="">
            <a href="/erp/product/batch/edit">
                <i class="menu-icon fa fa-caret-right"></i>
                商品批量修改
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPPZ_PACKAGE_MATERIAL) {
		buffer.WriteString(`<li id="sppz-packagematerial-menu" class="">
            <a href="/erp/package/material">
                <i class="menu-icon fa fa-caret-right"></i>
                包材类型
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPPZ_STOCK) {
		buffer.WriteString(`<li id="sppz-stock-menu" class="">
            <a href="/erp/stock">
                <i class="menu-icon fa fa-caret-right"></i>
                仓库
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPPZ_STORAGE_LOC) {
		buffer.WriteString(`<li id="sppz-storageloc-menu" class="">
            <a href="/erp/storage/loc">
                <i class="menu-icon fa fa-caret-right"></i>
                仓位
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPPZ_PRODUCT_CATE) {
		buffer.WriteString(`<li id="sppz-productcate-menu" class="">
            <a href="/erp/product/cate">
                <i class="menu-icon fa fa-caret-right"></i>
                商品目录
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPPZ_PSL_BATCH) {
		buffer.WriteString(`<li id="sppz-psl-batch-menu" class="">
            <a href="/erp/product/storage/loc/batch/edit">
                <i class="menu-icon fa fa-caret-right"></i>
                批量修改仓位
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	if check_pass_res_menu(pass_res.SPPZ_SUPPLIER) {
		buffer.WriteString(`<li id="sppz-supplier-menu" class="">
            <a href="/erp/supplier">
                <i class="menu-icon fa fa-caret-right"></i>
                供应商
            </a>
            <b class="arrow"></b>
        </li>`)
	}

	buffer.WriteString(`<li id="sppz-product-status-menu" class="">
            <a href="/erp/product/status">
                <i class="menu-icon fa fa-caret-right"></i>
                商品状态
            </a>
            <b class="arrow"></b>
        </li>`)

	buffer.WriteString(`</ul>
</li>
`)

	return buffer.String()
}

/*
   <li id="sppz-brand-menu" class="">
       <a href="/erp/brand">
           <i class="menu-icon fa fa-caret-right"></i>
           品牌
       </a>
       <b class="arrow"></b>
   </li>
*/
