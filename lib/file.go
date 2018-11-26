package lib

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego/config"
	"io"
	"os"
)

func GetLogContentByLine(line int) ([]string, error) {
	file, err := os.Open("./project.log")
	log := make([]string, 0)
	if err != nil {
		return log, err
	}
	b := bufio.NewReader(file)
	for {
		a, _, c := b.ReadLine()
		if c == io.EOF {
			break
		}
		log = append(log, string(a))
	}
	for i, j := 0, len(log)-1; i < j; i, j = i+1, j-1 {
		log[i], log[j] = log[j], log[i]
	}
	if len(log) < line {
		return log, nil
	}
	return log[len(log)-line:], nil
}

var ConfigValue = map[string]string{
	"email::user":     "",
	"email::password": "",
	"email::host":     "",
	"email::addr":     "",
	"msg::appid":      "",
	"msg::appkey":     "",
	"msg::sign":       "",
	"sys::password":   "",
}
//读取ini文件
func ReadConf() error {
	conf, err := config.NewConfig("ini", "./conf/config.conf")
	if err != nil {
		return err
	}
	for k, _ := range ConfigValue {
		ConfigValue[k] = conf.String(k)
	}
	return nil
}

func SetConf(c map[string]string) error {
	conf, err := config.NewConfig("ini", "./conf/config.conf")
	if err != nil {
		return err
	}
	for k, v := range c {
		err = conf.Set(k, v)
		if err != nil {
			return err
		}
	}
	return conf.SaveConfigFile("./conf/config.conf")
}

func Str2md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
