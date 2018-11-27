package lib

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/cnlh/httpMonitor/cron"
	"github.com/cnlh/httpMonitor/models"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	mainCron *cron.Cron
	jobList  map[int]*Job
)

type Job struct {
	status     int       //任务
	lastNotice time.Time //上次通知时间
	errTimes   int
	index      int
}

func init() {
	mainCron = cron.New()
	mainCron.Start()
	jobList = make(map[int]*Job)
}

//添加单个job
func AddJob(job *models.JobDetail) error {
	if job.CronType == 1 {
		job.Cron = "0/" + job.Cron + " * * * * ?"
	}
	index, err := mainCron.AddFunc(job.Cron, NewJobFunc(job)) //添加出错可能性处理
	if err != nil {
		log.Println("任务添加出错", err)
		return err
	}
	s, _ := time.ParseDuration("-" + strconv.Itoa(job.NoticeInterval) + "s")
	jobList[job.Id] = &Job{1, time.Now().Add(s), 0, index}
	return nil
}

//通过id添加任务
func AddJobById(id int) error {
	if _, ok := jobList[id]; !ok {
		job, err := models.GetJobDetailById(id)
		if err != nil {
			return err
		}
		if err = AddJob(job); err != nil {
			return err
		}
	} else {
		return errors.New("任务已经存在")
	}
	return nil
}

//通过id删除任务
func DelJobById(id int) error {
	if v, ok := jobList[id]; ok {
		err := mainCron.DelJob(v.index)
		delete(jobList, id)
		if err != nil {
			return err
		}
	}
	return nil
}

//从数据库中读取，并添加
func InitJob() {
	var jobs []*models.JobDetail
	o := orm.NewOrm()
	_, err := o.QueryTable("mn_job_detail").Filter("status", 1).All(&jobs)
	if err != nil {
		log.Fatalf("任务读取错误，中断执行")
	}
	for _, v := range jobs {
		AddJob(v)
	}
}

//组装一个job方法，
func NewJobFunc(job *models.JobDetail) func() {
	return func() {
		start := time.Now()
		r := New(job.Method, job.Url, job.Data)
		if job.Header != "" {
			r.SetHeader(job.Header)
		}
		if job.Overtime != 0 {
			r.SetOvertime(job.Overtime)
		}
		statuCode, content, err := r.Do()
		cost := time.Since(start)
		//执行错误，添加日志
		if err != nil {
			Record(job, statuCode, content, false, cost.String(), err.Error())
			return
		}
		//验证是否执行成功,添加日志
		if ok, err := verify(job.RegType, statuCode, content, job.RegVal); err != nil {
			Record(job, statuCode, content, false, cost.String(), err.Error())
		} else if ok {
			Record(job, statuCode, content, true, cost.String(), "")
		} else {
			Record(job, statuCode, content, false, cost.String(), "验证失败")
		}
	}
}

func TestJob(job *models.JobDetail) map[string]interface{} {
	isPass := false
	start := time.Now()
	r := New(job.Method, job.Url, job.Data)
	if job.Header != "" {
		r.SetHeader(job.Header)
	}
	if job.Overtime != 0 {
		r.SetOvertime(job.Overtime)
	}
	statuCode, content, lerr := r.Do()
	if lerr != nil {
		isPass = false
	} else if ok, err := verify(job.RegType, statuCode, content, job.RegVal); ok {
		isPass = true
	} else {
		lerr = err
	}
	cost := time.Since(start)
	json := make(map[string]interface{})
	json["statuCode"] = statuCode
	json["content"] = content
	json["isPass"] = isPass
	json["cost"] = cost.String()
	if lerr != nil {
		json["err"] = lerr.Error()
	} else {
		json["err"] = ""
	}
	return json
}

//添加日志
func Record(job *models.JobDetail, status int, content string, res bool, runTime string, err string) {
	if !res {
		jobList[job.Id].errTimes++
		jobRecord := models.JobRecord{
			JobId:      job.Id,
			StatusCode: status,
			ReqContent: content,
			RunRes:     false,
			Err:        err,
			RunTime:    runTime,
		}
		if _, err := models.Insert(&jobRecord); err != nil {
			logs.Error("日志记录发生错误", err)
		}
	} else {
		models.SetJobRunStatus(job.Id, 1)
		jobList[job.Id].errTimes = 0
	}
	if !res && jobList[job.Id].errTimes == job.ErrTimes {
		models.SetJobRunStatus(job.Id, 2)
		jobList[job.Id].status = 2
		jobList[job.Id].errTimes = 0
		sendErrNotice(job, status, content, err)
	}
	if res && jobList[job.Id].status == 2 {
		models.SetJobRunStatus(job.Id, 1)
		jobList[job.Id].status = 1
		sendReNotice(job)
	}
}

//发送接口恢复信息
func sendReNotice(job *models.JobDetail) {
	if job.NoticeType == 1 {
		SendMail(job.NoticeTo, job.Title+"接口已经恢复正常", "已经恢复正常,时间："+time.Now().String())
	} else {

	}
}

//发送邮件或者短信通知
func sendErrNotice(job *models.JobDetail, status int, content string, err string) {
	if job.IsNotice == true && (int(time.Now().Unix()-jobList[job.Id].lastNotice.Unix()) > job.NoticeInterval) { //时间满足上一次间隔
		noticeUsers := strings.Split(job.NoticeTo, ",")
		if job.NoticeType == 1 {
			content := fmt.Sprintf("url:%s\n返回状态码:%s\n返回内容:%s\n错误提示:%s\n时间:%s", job.Url, strconv.Itoa(status), content, err, time.Now().String())
			for _, v := range noticeUsers {
				err := SendMail(v, job.Title+"接口发生异常，请及时进行处理", content)
				if err != nil { //可能是内容错误 再发送一次
					content := fmt.Sprintf("url:%s\n返回状态码:%s\n返回内容:%s\n错误提示:%s\n时间:%s", job.Url, strconv.Itoa(status), "暂时无法通过邮件发送", err, time.Now().String())
					err = SendMail(job.NoticeTo, job.Title+"接口发生异常，请及时进行处理", content)
					if err != nil {
						logs.Error("邮件发送失败！", err)
					}
				}
			}
		} else {
			for _, v := range noticeUsers {
				if err := SendMsg(v, job.Title, status); err != nil {
					logs.Error("短信发送出错", err.Error())
				}
			}
		}
		models.UpdateJobLastNotice(job.Id)
		jobList[job.Id].lastNotice = time.Now()
	}
}
