package controllers

import (
	"github.com/revel/revel"
	"fmt"
	"strings"
	"github.com/pborman/uuid"
	"path"
	"os"
	"io"
	"matrix/core"
)

type UEditor struct {
	*revel.Controller
}

func (c UEditor) Index() revel.Result {
	var action string
	c.Params.Bind(&action, "action")

	if action == "config" {
		return c.Config()
	} else if action == "uploadimage" {
		return c.UploadImage()
	} else if action == "uploadfile" {
		return c.UploadFile()
	}

	return c.RenderEmpty()
}

func (c UEditor) Test() revel.Result {
	return c.RenderTemplate("ueditor/ueditor_test.html")
}

func (c UEditor) Config() revel.Result {

	return c.RenderJson(map[string]interface{}{
		/* 上传图片配置 */
		"imageActionName": "uploadimage",
		"imageUrlPrefix":  "",
		"imagePath":       "/static/upload/ueditor/notice",
		"imageFieldName":  "upfile",
		"imageMaxSize":    2048,
		"imageAllowFiles": []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"},
		/* 上传文件配置 */
		"fileActionName": "uploadfile",                    /* controller里,执行上传视频的action名称 */
		"fileFieldName":  "upfile",                        /* 提交的文件表单名称 */
		"filePathFormat": "/static/upload/ueditor/notice", /* 上传保存路径,可以自定义保存路径和文件名格式 */
		"fileUrlPrefix":  "",                              /* 文件访问路径前缀 */
		"fileMaxSize":    51200000,                        /* 上传大小限制，单位B，默认50MB */
		"fileAllowFiles": []string{
			".png", ".jpg", ".jpeg", ".gif", ".bmp",
			".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
			".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
			".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
			".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
		},                                                 /* 上传文件格式显示 */
	})
}

func (c UEditor) UploadImage() revel.Result {

	file, header, err := c.Request.FormFile("upfile")
	core.HandleError(err)

	defer file.Close()

	filename := strings.Replace(uuid.NewUUID().String(), "-", "", -1) + path.Ext(header.Filename)

	err = os.MkdirAll(path.Join(revel.BasePath, "static", "upload", "ueditor", "notice"), 0775)
	core.HandleError(err)

	outFile, err := os.Create(path.Join(revel.BasePath, "static", "upload", "ueditor", "notice", filename))
	core.HandleError(err)

	defer outFile.Close()

	io.Copy(outFile, file)

	return c.RenderJson(map[string]interface{}{
		"url":      fmt.Sprintf("/static/upload/ueditor/notice/%s", filename), //保存后的文件路径
		"title":    "",                                                        //文件描述，对图片来说在前端会添加到title属性上
		"original": header.Filename,                                           //原始文件名
		"state":    "SUCCESS",                                                 //上传状态，成功时返回SUCCESS,其他任何值将原样返回至图片上传框中
	})
}

func (c UEditor) UploadFile() revel.Result {

	file, header, err := c.Request.FormFile("upfile")
	core.HandleError(err)

	defer file.Close()

	subdir := strings.Replace(uuid.NewUUID().String(), "-", "", -1)
	file_name := header.Filename

	err = os.MkdirAll(path.Join(revel.BasePath, "static", "upload", "ueditor", "notice", subdir), 0775)
	core.HandleError(err)

	outFile, err := os.Create(path.Join(revel.BasePath, "static", "upload", "ueditor", "notice", subdir, file_name))
	core.HandleError(err)

	defer outFile.Close()

	io.Copy(outFile, file)

	return c.RenderJson(map[string]interface{}{
		"url":      fmt.Sprintf("/static/upload/ueditor/notice/%s/%s", subdir, file_name), //保存后的文件路径
		"title":    "",                                                                    //文件描述，对图片来说在前端会添加到title属性上
		"original": header.Filename,                                                       //原始文件名
		"state":    "SUCCESS",                                                             //上传状态，成功时返回SUCCESS,其他任何值将原样返回至图片上传框中
	})
}
