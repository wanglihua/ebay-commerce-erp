package layout

import (
	"bytes"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/app/layout/menu"
	"matrix/core"
	"matrix/core/pass"
	"matrix/modules/auth/models"
)

func GetHeader(title string, db_session *xorm.Session, web_session revel.Session) string {
	sys_name := GetSysName(db_session)

	login_user := core.GetLoginUser(web_session)
	login_user_login_name := ""
	login_user_nick_name := ""
	if login_user.UserId != 0 {
		login_user_login_name = login_user.LoginName
		login_user_nick_name = login_user.NickName
	}

	var header_buffer bytes.Buffer

	header_buffer.WriteString(GetSiteHeader(title, db_session, web_session))

	header_buffer.WriteString(`<div id="navbar" class="navbar navbar-default" xmlns="http://www.w3.org/1999/html">
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
		header_buffer.WriteString(fmt.Sprintf("工号：%s &nbsp;&nbsp;&nbsp;&nbsp;姓名：%s", login_user_login_name, login_user_nick_name))
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

	is_admin := login_user.IsAdmin

	var res_cate_id_map = make(map[int64]bool)
	var res_id_map = make(map[int64]bool)

	if is_admin == false {
		res_cate_id_list := pass.GetLoginUserPassResCateList(db_session, web_session)
		for _, res_cate_id := range res_cate_id_list {
			res_cate_id_map[res_cate_id] = true
		}

		res_id_list := pass.GetLoginUserPassResList(db_session, web_session)
		for _, res_id := range res_id_list {
			res_id_map[res_id] = true
		}
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_DDLC]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetDdlcMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_XXGL]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetXxglMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_PPGL]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetPpglMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_THGL]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetThglMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_DDPZ]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetDdpzMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_DDBB]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetDdbbMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_KCGL]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetKcglMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_SPPZ]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetSppzMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_SPBB]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetSpbbMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_CGGL]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetCgglMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_DBGL]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetDbglMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_XSGL]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetXsglMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_XSBB]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetXsbbMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_AUTH]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetAuthMenu(res_id_map, db_session, web_session))
	}

	if v, ok := res_cate_id_map[models.PASS_RES_CATE_XTPZ]; is_admin == true || (v == true && ok == true) {
		header_buffer.WriteString(menu.GetXtpzMenu(res_id_map, db_session, web_session))
	}

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

func GetSiteHeader(title string, db_session *xorm.Session, web_session revel.Session) string {
	sys_name := GetSysName(db_session)

	var site_header_buffer bytes.Buffer
	site_header_buffer.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
    <meta charset="utf-8"/>
    <title>`)

	site_header_buffer.WriteString(title)
	site_header_buffer.WriteString("-")
	site_header_buffer.WriteString(sys_name)

	site_header_buffer.WriteString(fmt.Sprintf(`</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0"/>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/all.css?%[1]s" class="ace-main-stylesheet" id="main-ace-style" />
    <!--[if lte IE 9]>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/ace-part2.min.css?%[1]s" class="ace-main-stylesheet"/>
    <![endif]-->
    <!--[if lte IE 9]>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/ace-ie.min.css?%[1]s"/>
    <![endif]-->
    <link type="text/css" rel="stylesheet" href="/static/site.css?%[1]s"/>`, core.GetStaticVersion()))

	if revel.DevMode {
		site_header_buffer.WriteString(`
    <style type="text/css">
        div.todo {
            display: block;
            margin: 10px 5px;
            padding: 10px;
            color: blue;
            background-color: yellow;
            line-height: 38px;
        }
    </style>`)
	}

	if revel.DevMode {
		site_header_buffer.WriteString(fmt.Sprintf(`<script type="text/javascript" src="/static/js/vue.js?%[1]s"></script>`, core.GetStaticVersion()))
	} else {
		site_header_buffer.WriteString(fmt.Sprintf(`<script type="text/javascript" src="/static/js/vue.min.js?%[1]s"></script>`, core.GetStaticVersion()))
	}

	site_header_buffer.WriteString(fmt.Sprintf(`<!--[if !IE]> -->
    <script type="text/javascript">
        window.jQuery || document.write("<script type='text/javascript' src='/static/ace/js/jquery.min.js?%[1]s'>" + "<" + "/script>");
    </script>
    <!-- <![endif]-->
    <!--[if IE]>
    <script type="text/javascript">
        window.jQuery || document.write("<script type='text/javascript' src='/static/ace/js/jquery1x.min.js?%[1]s'>" + "<" + "/script>");
    </script>
    <![endif]-->
    <script type="text/javascript">
        if ('ontouchstart' in document.documentElement) document.write("<script type='text/javascript' src='/static/ace/js/jquery.mobile.custom.min.js?%[1]s'>" + "<" + "/script>");
    </script>`, core.GetStaticVersion()))

	site_header_buffer.WriteString(fmt.Sprintf(`
        <script type="text/javascript" src="/static/ace/js/all1.js?%[1]s"></script>
        <!--[if lte IE 8]>
        <script type="text/javascript" src="/static/ace/js/html5shiv.min.js?%[1]s"></script>
        <script type="text/javascript" src="/static/ace/js/respond.min.js?%[1]s"></script>
        <![endif]-->
        <script type="text/javascript" src="/static/ace/js/all2.js?%[1]s"></script>
        <script type="text/javascript" src="/static/js/js.cookie.js?%[1]s"></script>
        <script type="text/javascript" src="/static/site.js?%[1]s"></script>`, core.GetStaticVersion()))

	site_header_buffer.WriteString(`
</head>

<body class="no-skin">
<div id="huaduCover" class="huaduCover" style="display: none;">
    <div class="huaduBg" id="huaduBg">&nbsp;</div>
    <div class="huaduProgressBar" id="huaduProgressBar"></div>
</div>
`)

	return site_header_buffer.String()
}

func GetLoginHeader(db_session *xorm.Session, web_session revel.Session) string {
	var site_header_buffer bytes.Buffer
	site_header_buffer.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
    <meta charset="utf-8"/>
    <title>`)

	site_header_buffer.WriteString("welcome to wollongzone!")

	site_header_buffer.WriteString(fmt.Sprintf(`</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0"/>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/all.css?%[1]s" class="ace-main-stylesheet" id="main-ace-style" />
    <!--[if lte IE 9]>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/ace-part2.min.css?%[1]s" class="ace-main-stylesheet"/>
    <![endif]-->
    <!--[if lte IE 9]>
    <link type="text/css" rel="stylesheet" href="/static/ace/css/ace-ie.min.css?%[1]s"/>
    <![endif]-->
    <link type="text/css" rel="stylesheet" href="/static/site.css?%[1]s"/>`, core.GetStaticVersion()))

	if revel.DevMode {
		site_header_buffer.WriteString(fmt.Sprintf(`<script type="text/javascript" src="/static/js/vue.js?%[1]s"></script>`, core.GetStaticVersion()))
	} else {
		site_header_buffer.WriteString(fmt.Sprintf(`<script type="text/javascript" src="/static/js/vue.min.js?%[1]s"></script>`, core.GetStaticVersion()))
	}

	site_header_buffer.WriteString(fmt.Sprintf(`<!--[if !IE]> -->
    <script type="text/javascript">
        window.jQuery || document.write("<script type='text/javascript' src='/static/ace/js/jquery.min.js?%[1]s'>" + "<" + "/script>");
    </script>
    <!-- <![endif]-->
    <!--[if IE]>
    <script type="text/javascript">
        window.jQuery || document.write("<script type='text/javascript' src='/static/ace/js/jquery1x.min.js?%[1]s'>" + "<" + "/script>");
    </script>
    <![endif]-->
    <script type="text/javascript">
        if ('ontouchstart' in document.documentElement) document.write("<script type='text/javascript' src='/static/ace/js/jquery.mobile.custom.min.js?%[1]s'>" + "<" + "/script>");
    </script>`, core.GetStaticVersion()))

	site_header_buffer.WriteString(fmt.Sprintf(`
        <script type="text/javascript" src="/static/ace/js/all1.js?%[1]s"></script>
        <!--[if lte IE 8]>
        <script type="text/javascript" src="/static/ace/js/html5shiv.min.js?%[1]s"></script>
        <script type="text/javascript" src="/static/ace/js/respond.min.js?%[1]s"></script>
        <![endif]-->
        <script type="text/javascript" src="/static/ace/js/all2.js?%[1]s"></script>
        <script type="text/javascript" src="/static/js/js.cookie.js?%[1]s"></script>
        <script type="text/javascript" src="/static/site.js?%[1]s"></script>`, core.GetStaticVersion()))

	site_header_buffer.WriteString(`
</head>

<body class="no-skin">
<div id="huaduCover" class="huaduCover" style="display: none;">
    <div class="huaduBg" id="huaduBg">&nbsp;</div>
    <div class="huaduProgressBar" id="huaduProgressBar"></div>
</div>
`)

	return site_header_buffer.String()
}
