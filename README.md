# httpMonitor
http、api接口监控系统,可监控各种http请求，支持header、data设置，支持腾讯云短信和邮件警告通知。

[测试地址](http://123.206.77.88:8080) 管理密码123

## 界面
![image](https://raw.githubusercontent.com/cnlh/httpMonitor/master/img1.png)
![image](https://raw.githubusercontent.com/cnlh/httpMonitor/master/img2.png)
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

