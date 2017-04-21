package controllers

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	"github.com/revel/revel"
	"matrix/core"
	"matrix/modules/auth/pass_res"
	"matrix/modules/auth/service/pass_service"
	"matrix/modules/erp/models"
	"matrix/modules/erp/setting"
)

type ErpLittleIcon struct {
	*revel.Controller
}

func (c ErpLittleIcon) Index() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LITTLE_ICON, c.Session)

	c.RenderArgs["icon_path"] = setting.ICON_PATH
	return c.RenderTemplate("erp/xtpz/little_icon/little_icon_index.html")
}

func (c ErpLittleIcon) ListData() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LITTLE_ICON, c.Session)

	db_session := c.DbSession

	filter, order, offset, limit := core.GetGridRequestParam(c.Request)
	query := db_session.Where(filter)

	//query extra filter here

	dataQuery := *query
	if order != "" {
		dataQuery = *dataQuery.OrderBy(order)
	} else {
		dataQuery = *dataQuery.Asc("id")
	}

	little_icon_list := make([]models.LittleIconInfo, 0, limit)
	err := dataQuery.Limit(limit, offset).Find(&little_icon_list)
	core.HandleError(err)

	countQuery := *query
	count, err := countQuery.Count(new(models.LittleIconInfo))
	core.HandleError(err)

	return c.RenderJson(core.GridResult{
		Data:  little_icon_list,
		Total: count,
	})
}

type LittleIconInfoForm struct {
	LittleIcon models.LittleIconInfo
}

func (f *LittleIconInfoForm) IsCreate() bool {
	return f.LittleIcon.Id == 0
}

func (f *LittleIconInfoForm) Valid(validation *revel.Validation) bool {
	validation.Required(f.LittleIcon.IconFile).Message("图标不能为空！")

	return validation.HasErrors() == false
}

func (c ErpLittleIcon) Detail() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LITTLE_ICON, c.Session)

	db_session := c.DbSession

	var little_icon_id int64
	c.Params.Bind(&little_icon_id, "id")

	little_icon := new(models.LittleIconInfo)
	if little_icon_id != 0 {
		has, err := db_session.Id(little_icon_id).Get(little_icon)
		core.HandleError(err)
		if has == false {
			panic("指定的列表图标不存在！")
		}
	}

	form := new(LittleIconInfoForm)
	form.LittleIcon = *little_icon
	c.UnbindToRenderArgs(form, "form")
	c.RenderArgs["is_create"] = form.IsCreate()
	c.RenderArgs["icon_path"] = setting.ICON_PATH

	return c.RenderTemplate("erp/xtpz/little_icon/little_icon_detail.html")
}

func (c ErpLittleIcon) Save() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LITTLE_ICON, c.Session)

	db_session := c.DbSession

	form := new(LittleIconInfoForm)
	c.Params.Bind(form, "form")

	if form.Valid(c.Validation) == false {
		return c.RenderJson(core.JsonResult{Success: false, Message: c.GetValidationErrorMessage()})
	}

	little_icon := &form.LittleIcon

	var affected int64
	var err error
	if form.IsCreate() {
		affected, err = db_session.Insert(little_icon)
		core.HandleError(err)
	} else {
		affected, err = db_session.Id(little_icon.Id).Update(little_icon)
		core.HandleError(err)

		if affected == 0 {
			return c.RenderJson(core.JsonResult{Success: false, Message: "数据保存失败，请重试！"})
		}
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据保存成功!"})
}

func (c ErpLittleIcon) Delete() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LITTLE_ICON, c.Session)

	db_session := c.DbSession

	little_icon_id_list := make([]int64, 0)
	c.Params.Bind(&little_icon_id_list, "id_list")

	var little_icon_list = make([]models.LittleIconInfo, 0)
	err := db_session.In("id", little_icon_id_list).Find(&little_icon_list)
	core.HandleError(err)

	affected, err := db_session.In("id", little_icon_id_list).Delete(new(models.LittleIconInfo))
	core.HandleError(err)

	//同时删除文件
	for _, little_icon := range little_icon_list {
		os.Remove(filepath.Join(revel.BasePath, setting.ICON_PATH, little_icon.IconFile))
	}

	return c.RenderJson(core.JsonResult{Success: true, Message: strconv.FormatInt(affected, 10) + "条数据删除成功!"})
}

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

func (c ErpLittleIcon) Upload() revel.Result {
	pass_service.Check_login_user_pass(pass_res.XTPZ_LITTLE_ICON, c.Session)

	iconFileBytes := make([]byte, 0)
	c.Params.Bind(&iconFileBytes, "icon_file")

	c.Validation.Required(iconFileBytes).Message("指定文件内容未成功上传！")
	c.Validation.MinSize(iconFileBytes, 800).Message("上传的文件太小了！")
	c.Validation.MaxSize(iconFileBytes, 100*KB).Message("上传的文件太大了！")

	conf, format, err := image.DecodeConfig(bytes.NewReader(iconFileBytes))
	revel.TRACE.Println(err)
	noError := err == nil
	c.Validation.Required(noError).Message("解码错误，可能是上传的文件格式不对！")
	formatLimit := format == "jpeg" || format == "png"
	c.Validation.Required(formatLimit).Message("请上传jpg/jpeg或png格式图片！")

	minSize := conf.Height > 20 && conf.Width > 20
	c.Validation.Required(minSize).Message("不能上传尺寸小于20x20像素的图片！")
	maxSize := conf.Height < 100 && conf.Width < 100
	c.Validation.Required(maxSize).Message("不能上传尺寸大于100x100像素的图片！")
	sameWidthHeight := conf.Height == conf.Width
	c.Validation.Required(sameWidthHeight).Message("请上传长和宽相同的图片！")

	if c.Validation.HasErrors() {
		return c.RenderJson(map[string]interface{}{
			"success": false,
			"message": c.GetValidationErrorMessage(),
		})
	}

	//保存文件内容到指定文件夹
	//fileName := time.Now().Format("20060102150405")
	fileName := strings.Replace(uuid.NewUUID().String(), "-", "", -1)
	fileExt := ""
	if format == "jpeg" {
		fileExt = ".jpg"
	} else if format == "png" {
		fileExt = ".png"
	}

	cate_icon_path := filepath.Join(revel.BasePath, setting.ICON_PATH)

	if core.IsDirExists(cate_icon_path) == false {
		err = os.MkdirAll(cate_icon_path, 0777)
		core.HandleError(err)
	}

	imageFile, err := os.Create(filepath.Join(cate_icon_path, fileName+fileExt))
	defer imageFile.Close()

	imageFile.Write(iconFileBytes)

	return c.RenderJson(map[string]interface{}{
		"success":  true,
		"message":  "图标上传成功！",
		"fileName": fileName + fileExt,
	})
}
