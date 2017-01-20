package layout

import (
	"bytes"
	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/app/layout/menu"
	"matrix/core"
	"fmt"
)

func GetHeader(title string, db_session *xorm.Session, web_session revel.Session) string {
	sys_name := GetSysName(db_session)

	login_user := core.GetLoginUser(web_session)
	login_user_login_name := ""
	login_user_nick_name := ""
	if login_user != nil {
		login_user_login_name = login_user.LoginName
		login_user_nick_name = login_user.NickName
	}

	var header_buffer bytes.Buffer
	header_buffer.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
    <meta charset="utf-8"/>
    <title>`)

	header_buffer.WriteString(title)
	header_buffer.WriteString("-")
	header_buffer.WriteString(sys_name)

	header_buffer.WriteString(`</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0"/>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/all.css?201701202000" class="ace-main-stylesheet" id="main-ace-style" />
    <!--[if lte IE 9]>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/ace-part2.min.css?201701202000" class="ace-main-stylesheet"/>
    <![endif]-->
    <!--[if lte IE 9]>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/ace-ie.min.css?201701202000"/>
    <![endif]-->
    <link type="text/css" rel="stylesheet" href="/static/site.css?201701202000"/>`)

	if revel.DevMode {
		header_buffer.WriteString(`<script type="text/javascript" src="/static/js/vue.js?201701202000"></script>`)
	} else {
		header_buffer.WriteString(`<script type="text/javascript" src="/static/js/vue.min.js?201701202000"></script>`)
	}

	header_buffer.WriteString(`<!--[if !IE]> -->
    <script type="text/javascript">
        window.jQuery || document.write("<script type='text/javascript' src='/static/ace/js/jquery.min.js?201701202000'>" + "<" + "/script>");
    </script>
    <!-- <![endif]-->
    <!--[if IE]>
    <script type="text/javascript">
        window.jQuery || document.write("<script type='text/javascript' src='/static/ace/js/jquery1x.min.js?201701202000'>" + "<" + "/script>");
    </script>
    <![endif]-->
    <script type="text/javascript">
        if ('ontouchstart' in document.documentElement) document.write("<script type='text/javascript' src='/static/ace/js/jquery.mobile.custom.min.js?201701202000'>" + "<" + "/script>");
    </script>`)

	header_buffer.WriteString(`
        <script type="text/javascript" src="/static/ace/js/all1.js?201701202000"></script>
        <!--[if lte IE 8]>
        <script type="text/javascript" src="/static/ace/js/html5shiv.min.js?201701202000"></script>
        <script type="text/javascript" src="/static/ace/js/respond.min.js?201701202000"></script>
        <![endif]-->
        <script type="text/javascript" src="/static/ace/js/all2.js?201701202000"></script>
        <script type="text/javascript" src="/static/js/js.cookie.js?201701202000"></script>
        <script type="text/javascript" src="/static/site.js?201701202000"></script>`)

	header_buffer.WriteString(`
</head>

<body class="no-skin">
<div id="huaduCover" class="huaduCover" style="display: none;">
    <div class="huaduBg" id="huaduBg">&nbsp;</div>
    <div class="huaduProgressBar" id="huaduProgressBar"></div>
</div>
<div id="navbar" class="navbar navbar-default" xmlns="http://www.w3.org/1999/html">
    <script type="text/javascript">
        try {
            ace.settings.check('navbar', 'fixed')
        } catch (e) {
        }
    </script>
    <div class="navbar-container" id="navbar-container">
        <button type="button" class="navbar-toggle menu-toggler pull-left" id="menu-toggler" data-target="#sidebar">
            <span class="sr-only">Toggle sidebar</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
        </button>
        <div class="navbar-header pull-left">
            <a href="/" class="navbar-brand">
                <small>`)

	header_buffer.WriteString(sys_name)

	header_buffer.WriteString(`</small>
            </a>
            </div>
        <div class="navbar-buttons navbar-header pull-right" role="navigation">
            <ul class="nav ace-nav">
                <li class="light-blue">
                    <a data-toggle="dropdown" href="#" class="dropdown-toggle">
                                <span class="user-info">`)

	if login_user_login_name != "" {
		header_buffer.WriteString(fmt.Sprintf("工号：%s &nbsp;&nbsp;&nbsp;&nbsp;姓名：%s", login_user.LoginName, login_user_nick_name))
	} else {
		header_buffer.WriteString("未登录")
	}

	header_buffer.WriteString(`</span>
                        <i class="ace-icon fa fa-caret-down"></i>
                    </a>
                    <ul class="user-menu dropdown-menu-right dropdown-menu dropdown-yellow dropdown-caret dropdown-closer">
                        <li>
                            <a href="/auth/profile">
                                <i class="ace-icon fa fa-cog"></i>
                                个人设置
                            </a>
                        </li>
                        <li class="divider"></li>
                        <li>
                            <a href="javascript:logout()">
                                <i class="ace-icon fa fa-power-off"></i>
                                注销
                            </a>
                            <script type="text/javascript">
                                function logout() {
                                    $.msg.confirm('确认要退出系统吗？', function () {
                                        $.ajax({
                                                   url: '/logout',
                                                   type: "POST",
                                                   success: function (data) {
                                                       if (data.success == true) {
                                                           window.location.href = '/login';
                                                           //return;
                                                       }
                                                       else {
                                                           $.msg.error(data.message);
                                                           //return;
                                                       }
                                                   }
                                               }
                                        );
                                    });
                                }
                            </script>
                        </li>
                     </ul>
                 </li>
            </ul>
        </div>
     </div>
</div>
<div class="main-container" id="main-container">
    <script type="text/javascript">
        $(function () {
            try {
                var parentMenuId = $("div#breadcrumbs ul.breadcrumb li.parent").data().menu;
                var activeMenuId = $("div#breadcrumbs ul.breadcrumb li.active").data().menu;
                var parentMenu = $("#" + parentMenuId);
                parentMenu.addClass("open");
                parentMenu.find("#" + activeMenuId).addClass("active");

                var pparentBreadcrumb = $("div#breadcrumbs ul.breadcrumb li.pparent");
                if(pparentBreadcrumb.size() != 0) {
                    var pparentMenuId = pparentBreadcrumb.data().menu;
                    var pparentMenu = $("#" + pparentMenuId);
                    pparentMenu.addClass("open");
                }

                var ppparentBreadcrumb = $("div#breadcrumbs ul.breadcrumb li.ppparent");
                if(ppparentBreadcrumb.size() != 0) {
                    var ppparentMenuId = ppparentBreadcrumb.data().menu;
                    var ppparentMenu = $("#" + ppparentMenuId);
                    ppparentMenu.addClass("open");
                }
            }
            catch (e) {
            }
        });
    </script>
    <!--sidebar-->
    <div id="sidebar" class="sidebar responsive">
        <script type="text/javascript">
            try {
                ace.settings.check('sidebar', 'fixed')
            } catch (e) {
            }
        </script>
        <ul class="nav nav-list">
`)

	header_buffer.WriteString(menu.GetDdlcMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetPpglMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetThglMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetDdpzMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetDdbbMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetKcglMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetSppzMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetSpbbMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetCgglMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetDbglMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetXsglMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetXsbbMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetAuthMenu(db_session, web_session))
	header_buffer.WriteString(menu.GetXtpzMenu(db_session, web_session))

	header_buffer.WriteString(`
        <div class="sidebar-toggle sidebar-collapse" id="sidebar-collapse">
            <i class="ace-icon fa fa-angle-double-left" data-icon1="ace-icon fa fa-angle-double-left" data-icon2="ace-icon fa fa-angle-double-right"></i>
        </div>
        <script type="text/javascript">
            try {
                ace.settings.check('sidebar', 'collapsed')
            } catch (e) {
            }
        </script>
    </div>

    <div class="main-content">
`)
	return header_buffer.String()
}
