package controllers

import (
	"github.com/cnlh/httpMonitor/lib"
	"github.com/cnlh/httpMonitor/models"
)

type IndexController struct {
	BaseController
}

func (self *IndexController) Index() {
	json := make(map[string]interface{})
	json["total"] = models.GetNum("mn_job_detail", map[string]interface{}{})
	json["start"] = models.GetNum("mn_job_detail", map[string]interface{}{"status": 1})
	json["run"] = models.GetNum("mn_job_detail", map[string]interface{}{"run_status": 1})
	json["stop"] = models.GetNum("mn_job_detail", map[string]interface{}{"run_status": 2})
	self.Data["num"] = json
	if log, err := lib.GetLogContentByLine(100); err == nil {
		self.Data["log"] = log
	}
	self.SetInfo("dashboard")
	self.display()
}

