package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type JobDetail struct {
	Id             int       `orm:"size(11);pk;auto"`
	Group          *JobGroup `orm:"rel(fk)"`
	Method         string    `orm:"size(10)"`                    //请求的方法
	Title          string    `orm:"size(100)"`                   //该任务的名称
	NoticeTo       string    `orm:"size(100)"`                   //该任务的名称
	Url            string    `orm:"size(255)"`                   //请求的url
	Cron           string    `orm:"size(20)"`                    //执行表达式分为表达式以及间隔XX秒执行一次
	CronType       int       `orm:"size(1);default(1)"`          //表达式类型1为间隔 2为表达式
	Header         string    `orm:"type(text)"`                  //请求头部信息
	Data           string    `orm:"type(text) "`                 //请求的数据
	RegType        int       `orm:"size(1) "`                    //检测匹配类型0、1、2
	RegVal         string    `orm:"size(255) "`                  //检测表达式
	Overtime       int       `orm:"size(11) "`                   //超时时间
	IsNotice       bool                                          //是否通知1
	NoticeType     int       `orm:"size(1)"`                     //通知类型：1邮件、短信
	LastNotice     time.Time `orm:"auto_now_add;type(datetime)"` //上次通知时间
	Status         int       `orm:"size(1);default(1)"`          //任务状态1为正常 0为暂停
	RunStatus      int       `orm:"size(1);default(1)"`          //任务执行状态1为正常 2为发生错误(并且未得到修正)
	ErrTimes       int       `orm:"size(1);default(3)"`          //连续出错X次通知次数
	NoticeInterval int       `orm:" size(11) "`                  //通知间隔 秒
	CreateTime     time.Time `orm:"auto_now_add;type(datetime)"` //创建时间
	UpdateTime     time.Time `orm:"auto_now;type(datetime)"`     //修改时间
}
type JobRecord struct {
	Id         int    `orm:"size(11);pk;auto"`
	JobId      int    `orm:"size(11);index"`                 //所属任务id
	StatusCode int    `orm:"size(10)"`                       //执行返回状态码
	ReqContent string `orm:"type(text)"`                     //请求返回内容
	RunRes     bool                                          //是否匹配成功
	Err        string    `orm:"size(255)"`                   //错误
	RunTime    string    `orm:"size(100)"`                   //执行耗费时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"` //创建时间
}

func init() {
	orm.RegisterModelWithPrefix("mn_", new(JobGroup), new(JobDetail), new(JobRecord))
}

//修改任务状态
func SetJobRunStatus(id, status int) {
	o := orm.NewOrm()
	job := JobDetail{Id: id}
	if o.Read(&job) == nil {
		job.RunStatus = status
		o.Update(&job)
	}
}

//修改上次通知时间
func UpdateJobLastNotice(id int) {
	o := orm.NewOrm()
	job := JobDetail{Id: id}
	if o.Read(&job) == nil {
		job.LastNotice = time.Now()
		o.Update(&job)
	}
}

//通过id获取job
func GetJobDetailById(id int) (*JobDetail, error) {
	o := orm.NewOrm()
	job := &JobDetail{Id: id}
	err := o.Read(job)
	if err != nil {
		return job, err
	}
	return job, nil
}

//根据分组id更改状态并返回列表
func GetAndUpdate(list interface{}, groupId int, status int) error {
	o := orm.NewOrm().QueryTable("mn_job_detail")
	if groupId > 0 {
		o = o.Filter("group_id", groupId)
	}
	_, err := o.Update(orm.Params{"status": status})
	if err != nil {
		return err
	}
	_, err = o.All(list)
	if err != nil {
		return err
	}
	return nil
}

//根据任务id删除日志
func DelLogById(jobId int) error {
	_, err := orm.NewOrm().QueryTable("mn_job_record").Filter("job_id", jobId).Delete()
	return err
}
