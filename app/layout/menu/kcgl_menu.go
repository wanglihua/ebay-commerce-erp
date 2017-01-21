package menu

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
)

func GetKcglMenu(res_id_map map[int64]bool, db_session *xorm.Session, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
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
    <ul class="submenu">

        <li id="kcgl-product-menu" class="">
            <a href="/erp/product">
                <i class="menu-icon fa fa-caret-right"></i>
                商品列表
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-productbind-menu" class="">
            <a href="/erp/product/bind">
                <i class="menu-icon fa fa-caret-right"></i>
                绑定商品
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-product-import-menu" class="">
            <a href="/erp/product/import">
                <i class="menu-icon fa fa-caret-right"></i>
                商品导入
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-product-export-menu" class="">
            <a href="/erp/product/export">
                <i class="menu-icon fa fa-caret-right"></i>
                商品导出
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-inventory-menu" class="">
            <a href="/erp/inventory">
                <i class="menu-icon fa fa-caret-right"></i>
                商品库存
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-inventorycheck-menu" class="">
            <a href="/erp/inventory/check">
                <i class="menu-icon fa fa-caret-right"></i>
                库存盘点
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-inventorycheck-approve-menu" class="">
            <a href="/erp/inventory/check/approve">
                <i class="menu-icon fa fa-caret-right"></i>
                盘点审批
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-stockbill-transfer-menu" class="">
            <a href="/erp/stock/bill/transfer">
                <i class="menu-icon fa fa-caret-right"></i>
                分仓调拨
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-product-quality-feedback-menu" class="">
            <a href="/erp/product/quality/feedback">
                <i class="menu-icon fa fa-caret-right"></i>
                商品质量反馈
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-stockbill-checkin-menu" class="">
            <a href="/erp/stock/bill/checkin">
                <i class="menu-icon fa fa-caret-right"></i>
                手工入库
            </a>
            <b class="arrow"></b>
        </li>

        <li id="kcgl-stockbill-checkout-menu" class="">
            <a href="/erp/stock/bill/checkout">
                <i class="menu-icon fa fa-caret-right"></i>
                手工出库
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
