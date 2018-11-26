package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
}

//初始化参数
func (self *BaseController) Prepare() {
	controllerName, actionName := self.GetControllerAndAction()
	self.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	self.actionName = strings.ToLower(actionName)
	if self.GetSession("auth") != true {
		self.Redirect("/login/index", 302)
	}
}

//加载模板
func (self *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		tplname = self.controllerName + "/" + self.actionName + ".html"
	}
	self.Data["menu"] = self.controllerName
	self.Layout = "public/layout.html"
	self.TplName = tplname
}

//错误
func (self *BaseController) error() {
	self.Layout = "public/layout.html"
	self.TplName = "public/error.html"
}

//去掉没有err返回值的int
func (self *BaseController) GetIntNoErr(key string, def ...int) int {
	strv := self.Ctx.Input.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0]
	}
	val, _ := strconv.Atoi(strv)
	return val
}

//获取去掉错误的bool值
func (self *BaseController) GetBoolNoErr(key string, def ...bool) bool {
	strv := self.Ctx.Input.Query(key)
	if len(strv) == 0 && len(def) > 0 {
		return def[0]
	}
	val, _ := strconv.ParseBool(strv)
	return val
}

//ajax正确返回
func (self *BaseController) AjaxOk(str string) {
	self.Data["json"] = ajax(str, 1)
	self.ServeJSON()
	self.StopRun()
}

//ajax错误返回
func (self *BaseController) AjaxErr(str string) {
	self.Data["json"] = ajax(str, 0)
	self.ServeJSON()
	self.StopRun()
}

//组装ajax
func ajax(str string, status int) (map[string]interface{}) {
	json := make(map[string]interface{})
	json["status"] = status
	json["msg"] = str
	return json
}

//ajax table返回
func (self *BaseController) AjaxTable(list interface{}, cnt int64, recordsTotal int) {
	json := make(map[string]interface{})
	json["data"] = list
	json["draw"] = self.GetIntNoErr("draw")
	json["err"] = ""
	json["recordsTotal"] = recordsTotal
	json["recordsFiltered"] = cnt
	self.Data["json"] = json
	self.ServeJSON()
	self.StopRun()
}

//ajax table参数
func (self *BaseController) GetAjaxParams() (start, limit int) {
	self.Ctx.Input.Bind(&start, "start")
	self.Ctx.Input.Bind(&limit, "length")
	return
}

func (self *BaseController) SetInfo(name string) {
	self.Data["name"] = name
}
