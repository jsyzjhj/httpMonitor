package controllers

import (
	"github.com/cnlh/httpMonitor/models"
)

type GroupController struct {
	BaseController
}

func (self *GroupController) List() {
	if self.Ctx.Request.Method == "POST" {
		start, limit := self.GetAjaxParams()
		filter := make(map[string]interface{})
		var jobg []*models.JobGroup
		cnt := models.GetTableList("mn_job_group", filter, "", limit, start, &jobg)
		self.AjaxTable(jobg, cnt, len(jobg))
	} else {
		self.SetInfo("分组管理")
		self.display()
	}
}
func (self *GroupController) Add() {
	if self.Ctx.Request.Method == "POST" {
		jobg := models.JobGroup{
			Title:       self.GetString("Title"),
			Description: self.GetString("Description"),
		}
		if _, err := models.Insert(&jobg); err != nil {
			self.AjaxErr("分组添加时出错")
		} else {
			self.AjaxOk("添加成功")
		}
	} else {
		self.SetInfo("添加分组")
		self.display()
	}
}

func (self *GroupController) Edit() {
	if self.Ctx.Request.Method == "POST" {
		jobg := models.JobGroup{Id: self.GetIntNoErr("Id")}
		if models.Read(&jobg) == nil {
			jobg.Description = self.GetString("Method")
			jobg.Title = self.GetString("Title")
			if _, err := models.Update(&jobg); err == nil {
				self.AjaxOk("修改成功")
			} else {
				self.AjaxErr("修改失败")
			}
		} else {
			self.AjaxErr("修改失败")
		}
	} else {
		jobg := models.JobGroup{Id: self.GetIntNoErr("id")}
		if models.Read(&jobg) != nil {
			self.error()
		} else {
			self.Data["jobg"] = jobg
			self.SetInfo("修改分组")
			self.display()
		}
	}
}

func (self *GroupController) Del() {
	if _, err := models.Delete(&models.JobGroup{Id: self.GetIntNoErr("Id")}); err == nil {
		self.AjaxOk("删除成功")
	} else {
		self.AjaxErr("删除失败")
	}
}
