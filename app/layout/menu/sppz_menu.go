package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetSppzMenu(res_id_map map[int64]bool, db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
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
    <ul class="submenu">

        <li id="sppz-product-batch-menu" class="">
            <a href="/erp/product/batch/edit">
                <i class="menu-icon fa fa-caret-right"></i>
                商品批量修改
            </a>
            <b class="arrow"></b>
        </li>

        <li id="sppz-packagematerial-menu" class="">
            <a href="/erp/package/material">
                <i class="menu-icon fa fa-caret-right"></i>
                包材类型
            </a>
            <b class="arrow"></b>
        </li>

        <li id="sppz-stock-menu" class="">
            <a href="/erp/stock">
                <i class="menu-icon fa fa-caret-right"></i>
                仓库
            </a>
            <b class="arrow"></b>
        </li>

        <li id="sppz-storageloc-menu" class="">
            <a href="/erp/storage/loc">
                <i class="menu-icon fa fa-caret-right"></i>
                仓位
            </a>
            <b class="arrow"></b>
        </li>

        <li id="sppz-productcate-menu" class="">
            <a href="/erp/product/cate">
                <i class="menu-icon fa fa-caret-right"></i>
                商品目录
            </a>
            <b class="arrow"></b>
        </li>

        <li id="sppz-psl-batch-menu" class="">
            <a href="/erp/product/storage/loc/batch/edit">
                <i class="menu-icon fa fa-caret-right"></i>
                批量修改仓位
            </a>
            <b class="arrow"></b>
        </li>

        <li id="sppz-supplier-menu" class="">
            <a href="/erp/supplier">
                <i class="menu-icon fa fa-caret-right"></i>
                供应商
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
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
