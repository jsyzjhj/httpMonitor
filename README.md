# httpMonitor
http、api接口监控系统,可监控各种http请求，支持header、data设置，支持腾讯云短信和邮件警告通知。
## 界面

## 安装

1. 安装本系统

```
go get github.com/cnlh/httpMonitor
```

2. 安装beego

```
go get github.com/astaxie/beego
```

3. 安装腾讯云短信支持

```
go get github.com/qichengzx/qcloudsms_go
```
4. 编译

```
go build
```

5. 初始化数据库
```
./httpMonitor -orm syncdb
```

6. 设置系统管理密码

```
./httpMonitor -psd=password
```

7.运行
```
./httpMonitor&
```

