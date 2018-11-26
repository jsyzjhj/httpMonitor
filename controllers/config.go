package controllers

import (
	"github.com/cnlh/httpMonitor/lib"
)

type ConfigController struct {
	BaseController
}

func (self *ConfigController) Index() {
	config := make(map[string]string)
	config["user"] = lib.ConfigValue["email::user"]
	config["password"] = lib.ConfigValue["email::password"]
	config["host"] = lib.ConfigValue["email::host"]
	config["addr"] = lib.ConfigValue["email::addr"]
	config["appid"] = lib.ConfigValue["msg::appid"]
	config["appkey"] = lib.ConfigValue["msg::appkey"]
	config["sign"] = lib.ConfigValue["msg::sign"]
	self.Data["config"] = config
	self.SetInfo("配置管理")
	self.display()
}
func (self *ConfigController) Save() {
	if self.Ctx.Request.Method == "POST" {
		config := make(map[string]string)
		for k, _ := range lib.ConfigValue {
			lib.ConfigValue[k] = self.GetString(k)
			config[k] = self.GetString(k)
		}
		if lib.SetConf(config) != nil {
			self.AjaxErr("修改失败")
		} else {
			self.AjaxOk("修改成功")
		}
	}
}
