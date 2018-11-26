package controllers

import (
	"github.com/astaxie/beego"
	"github.com/cnlh/httpMonitor/lib"
)

type LoginController struct {
	beego.Controller
}

func (self *LoginController) Index() {
	self.TplName = "login/index.html"
}
func (self *LoginController) Verify() {
	if lib.Str2md5(self.GetString("psd")) == lib.ConfigValue["sys::password"] {
		self.SetSession("auth", true)
		self.Data["json"] = map[string]interface{}{"status": 1, "msg": "验证成功"}
		self.ServeJSON()
	} else {
		self.Data["json"] = map[string]interface{}{"status": 0, "msg": "验证失败"}
		self.ServeJSON()
	}
}
func (self *LoginController) Out() {
	self.SetSession("auth", false)
	self.Redirect("/login/index", 302)
}
