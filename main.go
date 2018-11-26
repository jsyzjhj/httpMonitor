package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/cnlh/httpMonitor/lib"
	_ "github.com/cnlh/httpMonitor/models"
	_ "github.com/cnlh/httpMonitor/routers"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func init() {
	dbHost := beego.AppConfig.String("dbHost")
	dbPort := beego.AppConfig.String("dbPort")
	dbUser := beego.AppConfig.String("dbUser")
	dbPassword := beego.AppConfig.String("dbPassword")
	dbName := beego.AppConfig.String("dbName")
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&loc=Asia%2FShanghai"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dsn)
}
func main() {
	password := flag.String("psd", "", "password")
	flag.Parse()
	if *password != "" {
		json := make(map[string]string)
		json["sys::password"] = lib.Str2md5(*password)
		lib.SetConf(json)
		fmt.Println("密码设置成功")
		return
	}
	orm.RunCommand()
	lib.InitJob()
	beego.AddFuncMap("isColor", lib.IsColor)
	if err := lib.ReadConf(); err != nil {
		log.Fatalf("配置文件读取错误")
	}
	logs.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":5,"maxlines":0,"maxsize":0,"daily":false,"maxdays":10,"color":true}`)
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.Run()
}
