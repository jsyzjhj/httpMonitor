package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type JobGroup struct {
	Id          int       `orm:"size(11);pk;auto"`
	Title       string    `orm:"size(100)"`                   //分组名称
	Description string    `orm:"size(255)"`                   //分组描述
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"` //创建时间
	UpdateTime  time.Time `orm:"auto_now;type(datetime)"`     //修改时间
}

func GetAllGroup() (jobg []*JobGroup) {
	orm.NewOrm().QueryTable("mn_job_group").All(&jobg)
	return
}
