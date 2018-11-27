package controllers

import (
	"github.com/cnlh/httpMonitor/lib"
	"github.com/cnlh/httpMonitor/models"
)

type JobController struct {
	BaseController
}

func (self *JobController) List() {
	if self.Ctx.Request.Method == "POST" {
		start, limit := self.GetAjaxParams()
		filter := make(map[string]interface{})
		if group := self.GetIntNoErr("group", -1); group > 0 {
			filter["group__id"] = group
		}
		var jobd []*models.JobDetail
		cnt := models.GetTableList("mn_job_detail", filter, "", limit, start, &jobd)
		self.AjaxTable(jobd, cnt, len(jobd))
	} else {
		self.Data["jobg"] = models.GetAllGroup()
		self.SetInfo("任务列表")
		self.display()
	}
}
func (self *JobController) Add() {
	if self.Ctx.Request.Method == "POST" {
		job := models.JobDetail{
			Method:         self.GetString("Method"),
			Title:          self.GetString("Title"),
			NoticeTo:       self.GetString("NoticeTo"),
			Url:            self.GetString("Url"),
			Cron:           self.GetString("Cron"),
			CronType:       self.GetIntNoErr("CronType"),
			Header:         self.GetString("Header", ""),
			Data:           self.GetString("Data", ""),
			RegType:        self.GetIntNoErr("RegType"),
			Status:         1,
			RegVal:         self.GetString("RegVal"),
			Overtime:       self.GetIntNoErr("Overtime"),
			IsNotice:       self.GetBoolNoErr("IsNotice"),
			NoticeType:     self.GetIntNoErr("NoticeType"),
			ErrTimes:       self.GetIntNoErr("ErrTimes"),
			NoticeInterval: self.GetIntNoErr("NoticeInterval"),
		}
		isTest := self.GetBoolNoErr("test", false)
		if isTest {
			self.Data["json"] = lib.TestJob(&job)
			self.ServeJSON()
			self.StopRun()
		} else {
			job.Group = &models.JobGroup{Id: self.GetIntNoErr("Group")}
			if id, err := models.Insert(&job); err != nil {
				self.AjaxErr("任务添加时出错")
			} else {
				if err = lib.AddJobById(int(id)); err != nil {
					self.AjaxErr("任务初始化时出错，请及时修改")
				} else {
					self.AjaxOk("添加成功，已经开始执行")
				}
			}
		}
	} else {
		self.Data["jobg"] = models.GetAllGroup()
		self.SetInfo("添加任务")
		self.display()
	}
}
func (self *JobController) Change() {
	jobd := models.JobDetail{Id: self.GetIntNoErr("Id")}
	if models.Read(&jobd) == nil {
		jobd.Status = self.GetIntNoErr("Status")
		if _, err := models.Update(&jobd); err == nil {
			if jobd.Status == 1 { //重新开始
				lib.AddJobById(jobd.Id)
			} else {
				lib.DelJobById(jobd.Id)
			}
			self.AjaxOk("修改成功")
		}
	}
	self.AjaxErr("修改失败")
}
func (self *JobController) Edit() {
	if self.Ctx.Request.Method == "POST" {
		jobd := models.JobDetail{Id: self.GetIntNoErr("Id")}
		if models.Read(&jobd) == nil {
			jobd.Method = self.GetString("Method")
			jobd.Title = self.GetString("Title")
			jobd.NoticeTo = self.GetString("NoticeTo")
			jobd.Url = self.GetString("Url")
			jobd.Cron = self.GetString("Cron")
			jobd.CronType = self.GetIntNoErr("CronType")
			jobd.Header = self.GetString("Header", "")
			jobd.Data = self.GetString("Data", "")
			jobd.RegType = self.GetIntNoErr("RegType")
			jobd.RegVal = self.GetString("RegVal")
			jobd.Overtime = self.GetIntNoErr("Overtime")
			jobd.IsNotice = self.GetBoolNoErr("IsNotice")
			jobd.NoticeType = self.GetIntNoErr("NoticeType")
			jobd.ErrTimes = self.GetIntNoErr("ErrTimes")
			jobd.NoticeInterval = self.GetIntNoErr("NoticeInterval")
			jobd.Group = &models.JobGroup{Id: self.GetIntNoErr("Group")}
			if _, err := models.Update(&jobd); err == nil {
				//任务非暂停状态，任务重置
				if jobd.Status == 1 {
					err = lib.DelJobById(jobd.Id)
					if err == nil {
						err = lib.AddJobById(jobd.Id)
						if err == nil {
							self.AjaxOk("修改成功")
						}
					}
					self.AjaxErr("修改成功，任务重置失败！")
				}
				self.AjaxOk("修改成功")
			} else {
				self.AjaxErr("修改失败")
			}
		} else {
			self.AjaxErr("修改失败")
		}
	} else {
		jobd := models.JobDetail{Id: self.GetIntNoErr("id")}
		if models.Read(&jobd) != nil {
			self.error()
		} else {
			self.Data["jobd"] = jobd
			self.Data["jobg"] = models.GetAllGroup()
			self.SetInfo("查看修改任务")
			self.display()
		}
	}
}

func (self *JobController) Log() {
	if self.Ctx.Request.Method == "POST" {
		start, limit := self.GetAjaxParams()
		filter := make(map[string]interface{})
		filter["job_id"] = self.GetIntNoErr("id")
		var jobd []*models.JobRecord
		cnt := models.GetTableList("mn_job_record", filter, "-create_time", limit, start, &jobd)
		self.AjaxTable(jobd, cnt, len(jobd))
	} else {
		self.SetInfo("日志列表")
		self.display()
	}
}
func (self *JobController) Ldetail() {
	jobr := models.JobRecord{Id: self.GetIntNoErr("id")}
	if models.Read(&jobr) != nil {
		self.error()
	} else {
		self.Data["jobr"] = jobr
		self.SetInfo("日志详情")
		self.display()
	}
}
func (self *JobController) DelJob() {
	if _, err := models.Delete(&models.JobDetail{Id: self.GetIntNoErr("Id")}); err == nil {
		lib.DelJobById(self.GetIntNoErr("Id"))
		self.AjaxOk("删除成功")
	} else {
		self.AjaxErr("删除失败")
	}
}

//批量更改状态
func (self *JobController) ChangeAll() {
	var jobd []*models.JobDetail
	err := models.GetAndUpdate(&jobd, self.GetIntNoErr("group_id"), self.GetIntNoErr("status"))
	if err == nil {
		for _, v := range jobd {
			if v.Status == 1 { //重新开始
				lib.AddJobById(v.Id)
			} else {
				lib.DelJobById(v.Id)
			}
		}
		self.AjaxOk("修改成功")
	}
	self.AjaxErr("修改失败")
}

func (self *JobController) DelAllLog() {
	if err := models.DelLogById(self.GetIntNoErr("id")); err != nil {
		self.AjaxErr("删除失败")
	}
	self.AjaxOk("删除成功")
}
