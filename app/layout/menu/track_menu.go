package menu

import (
	"bytes"
	"github.com/revel/revel"
)

func GetTrackMenu(res_id_map map[int64]bool, web_session revel.Session) string {
	if _, found := revel.ModuleByName("erp"); found == false {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(`
<li id="track-menu" class="">
    <a href="#" class="dropdown-toggle">
        <i class="menu-icon fa fa-users"></i>
        <span class="menu-text">缺陷跟踪</span>
        <b class="arrow fa fa-angle-down"></b>
    </a>
    <b class="arrow"></b>
    <ul class="submenu">

        <li id="track-issue-menu" class="">
            <a href="/track/issue">
                <i class="menu-icon fa fa-caret-right"></i>
                软件缺陷
            </a>
            <b class="arrow"></b>
        </li>

        <li id="track-issuecate-menu" class="">
            <a href="/track/issue/cate">
                <i class="menu-icon fa fa-caret-right"></i>
                缺陷类别
            </a>
            <b class="arrow"></b>
        </li>

        <li id="track-issuestatus-menu" class="">
            <a href="/track/issue/status">
                <i class="menu-icon fa fa-caret-right"></i>
                缺陷状态
            </a>
            <b class="arrow"></b>
        </li>

        <li id="track-issuetype-menu" class="">
            <a href="/track/issue/type">
                <i class="menu-icon fa fa-caret-right"></i>
                缺陷类型
            </a>
            <b class="arrow"></b>
        </li>

    </ul>
</li>
`)

	return buffer.String()
}
