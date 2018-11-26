package lib

import (
	"github.com/qichengzx/qcloudsms_go"
	"net/smtp"
	"strconv"
	"strings"
)

func SendMail(toUser, title, content string) error {
	auth := smtp.PlainAuth("", ConfigValue["email::user"], ConfigValue["email::password"], ConfigValue["email::host"])
	to := []string{toUser}
	nickname := "system"
	user := ConfigValue["email::user"]
	content_type := "Content-Type: text/plain; charset=UTF-8"
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + title + "\r\n" + content_type + "\r\n\r\n" + content)
	err := smtp.SendMail(ConfigValue["email::addr"], auth, user, to, msg)
	return err
}

func SendMsg(tel string, title string, status int) error {
	opt := qcloudsms.NewOptions(ConfigValue["msg::appid"], ConfigValue["msg::appkey"], ConfigValue["msg::sign"])
	opt.Debug = true
	var client = qcloudsms.NewClient(opt)
	var sm = qcloudsms.SMSSingleReq{
		Type: 0,
		Msg:  title + "接口发生异常，返回状态码为" + strconv.Itoa(status) + "，请及时处理!",
		Tel:  qcloudsms.SMSTel{Nationcode: "86", Mobile: tel},
	}
	_, err := client.SendSMSSingle(sm)
	return err
}
