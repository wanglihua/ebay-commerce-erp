package app

import (
	"fmt"
	"html/template"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"matrix/app/layout"
)

func registerTags() {
	revel.TemplateFuncs["test"] = test_template_func
	revel.TemplateFuncs["header"] = header_template_func
	revel.TemplateFuncs["footer"] = footer_template_func
	revel.TemplateFuncs["site_header"] = site_header_template_func
	revel.TemplateFuncs["login_header"] = login_header_template_func
	revel.TemplateFuncs["site_footer"] = site_footer_template_func

	revel.TemplateFuncs["IntToStr"] = int_to_str_template_func
}

func test_template_func(renderArgs map[string]interface{}) template.HTML {
	session := renderArgs["session"].(revel.Session)
	return template.HTML("<span style='color:red;'>Session Id: " + session.Id() + "</span>")
}

type header_template_data struct {
	Title string
}

func header_template_func(title string, renderArgs map[string]interface{}) template.HTML {
	db_session := renderArgs["dbsession"].(*xorm.Session)
	web_session := renderArgs["session"].(revel.Session)
	return template.HTML(layout.GetHeader(title, db_session, web_session))
}

func footer_template_func() template.HTML {
	//just plain string
	return template.HTML(layout.GetFooter())
}

func site_header_template_func(title string, renderArgs map[string]interface{}) template.HTML {
	db_session := renderArgs["dbsession"].(*xorm.Session)
	web_session := renderArgs["session"].(revel.Session)
	return template.HTML(layout.GetSiteHeader(title, db_session, web_session))
}

func login_header_template_func(renderArgs map[string]interface{}) template.HTML {
	db_session := renderArgs["dbsession"].(*xorm.Session)
	web_session := renderArgs["session"].(revel.Session)
	return template.HTML(layout.GetLoginHeader(db_session, web_session))
}

func site_footer_template_func() template.HTML {
	//just plain string
	return template.HTML(layout.GetSiteFooter())
}

func int_to_str_template_func(val interface{}) string {
	return fmt.Sprintf("%d", val)
}
